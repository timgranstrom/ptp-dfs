package ptp
import (
	"github.com/timgranstrom/ptp-dfs/protoMessages"
	"fmt"
	"os"
	"net"
	"log"
	"strconv"
	"time"
	"encoding/hex"
	"go/types"
)

// A buffered channel that we can send work requests on.

//var WorkQueue = make(chan WorkRequest, 100)

type Network struct {
	protobufhandler *ProtobufHandler
	routingTable RoutingTable
	WorkerQueue chan Worker
	WorkQueue chan WorkRequest
	SendQueue chan Sender
	PingQueue chan Ping
	listenerActive bool
	store Store
}

type Sender struct {
	address string //Send to this address
	data []byte
}

//Contact to ping and channel to send the response through
type Ping struct {
	reply chan bool
	target *Contact
}

func NewSender(address string, data *[]byte) *Sender{
	sender := &Sender{
		address:address,
		data:*data,
	}
	return sender
}

func NewNetwork(routingTable RoutingTable) *Network{
	network := &Network{routingTable:routingTable,
			WorkerQueue: make(chan Worker,100),
			WorkQueue:make(chan WorkRequest,100),
			SendQueue:make(chan Sender,500),
			PingQueue:make(chan Ping),
			store: *MakeStore()}
	return network
}



/* A Simple function to verify error */
func CheckError(err error) {
	if err  != nil {
		log.Fatal("Error: " , err)
		os.Exit(0)
	}
}


func (network *Network) Sender() {
	log.Println(network.routingTable.me.Address,": SENDER STARTED")

	for{
		select {
			case sender := <- network.SendQueue:

				ServerAddr,err := net.ResolveUDPAddr("udp",sender.address)
				CheckError(err)

				LocalAddr, err := net.ResolveUDPAddr("udp", network.routingTable.me.Address)
				//CheckError(err)
				Conn, err := net.DialUDP("udp", LocalAddr,ServerAddr)
				CheckError(err)
			//	log.Println(network.routingTable.me.Address+" :Sent Message to ",sender.address)
				log.Println(network.routingTable.me.Address,": SENT MESSAGE TO OF SIZE ",len(sender.data),"FROM ",LocalAddr.String(),"TO",ServerAddr.String())
				_,err = Conn.Write(sender.data)
				CheckError(err)

				Conn.Close()
		}
	}
}


func (network *Network)SetupUDPListener(address string) *net.UDPConn{
	/* Lets prepare a address at any address at :8000*/
	ServerAddr,err := net.ResolveUDPAddr("udp",address)
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	//ServerConn.SetReadBuffer(4096)
	ServerConn.SetWriteBuffer(1048576)
	ServerConn.SetReadBuffer(1048576)
	log.Println(address+": Connection Established at "+ServerAddr.IP.String()+":"+strconv.Itoa(ServerAddr.Port))
	CheckError(err)
	network.listenerActive = true
	return ServerConn
}

func (network *Network) Listen() {
	log.Println(network.routingTable.me.Address,": LISTENER")
	ServerConn := network.SetupUDPListener(network.routingTable.me.Address)
	defer ServerConn.Close()

	buffer := make([]byte, 1024)

	for { //Infinite for-loop to check for incomming messages
		i,addr,err := ServerConn.ReadFromUDP(buffer)
		log.Println(network.routingTable.me.Address, ": RECEIVED MESSAGE OF SIZE", i, " FROM ", addr)
		go network.HandleRecievedMessage(buffer, addr, err)
	}
}

func (network *Network) HandleRecievedMessage(bufferMsg []byte,addr *net.UDPAddr,err error){
	msg := network.protobufhandler.UnMarshalWrapperMessage(bufferMsg)
	//log.Println(network.routingTable.me.Address+" :Received ",msg.MessageType.String(), " from ",addr)

	/*if *msg.MessageType == protoMessages.MessageType_FIND_CONTACT{
		log.Println(network.routingTable.me.Address," TARGET KAD ID: ",*msg.GetMsg_2().KademliaTargetId)
	}*/

	if err != nil {
		fmt.Println("Error: ",err)
	}
	//take the message and make a WorkRequest out of them.
	work := WorkRequest{msg.RequestId,*msg,addr.String()}

	//When handling work, make sure to always add the message sender as a contact
	kadId := NewKademliaID(work.message.SenderKademliaId)
	requestContact := NewContact(kadId,work.senderAddress)

	//Add the contact and see if there is a contact to ping
	pingContact := network.routingTable.AddContact(requestContact)
	if pingContact != nil {
		ping := Ping{ make(chan bool), pingContact }
		network.PingQueue <- ping //Tell self to start a worker that pings this node
		select {
			case reply := <- ping.reply: //Wait for response
				//If there is no reply (timeout), replace the pinged contact and move it to the front (using the addcontact func)
				if !reply {
					pingContact = &requestContact
					network.routingTable.AddContact(requestContact)
				}
		}
	}
	//log.Println(network.routingTable.me.Address+" :Received message from ",addr,": ADDED AS CONTACT")


	//Push the work onto the queue.
	network.WorkQueue <- work



	//log.Println(network.routingTable.me.Address+": Work request queued")
}

func (network *Network) SendPingMessage(targetAddress string, requestId int64, isReply bool) {
	//Create the ping message and marshal it
	pingContactMessage := network.protobufhandler.CreatePingMessage()
	wrapperMessage := network.protobufhandler.CreateWrapperMessage_1(network.routingTable.me.ID, requestId, protoMessages.MessageType_PING, pingContactMessage, isReply)
	data := network.protobufhandler.MarshalMessage(wrapperMessage)

	//Send the message using the queue
	sender := *NewSender(targetAddress, &data)
	network.SendQueue <- sender
}

func (network *Network) SendFindContactMessage(targetContact *Contact, sendToContact *Contact, requestId int64, responseContacts []Contact,isReply bool) {
	// TODO
	lookupContactMessage := network.protobufhandler.CreateLookupContactMessage(targetContact.ID) //Create a lookupContact message for the target sendToContact
	responseProtoContacts := network.protobufhandler.CreateContactMessages(responseContacts)
	lookupContactMessage.Contacts = responseProtoContacts
	//Create wrapper message for the request
	wrapperMessage := network.protobufhandler.CreateWrapperMessage_2(network.routingTable.me.ID,requestId,protoMessages.MessageType_FIND_CONTACT,lookupContactMessage,isReply)



	/*log.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	log.Println(network.routingTable.me.Address+" : Sent Find Contact Message ("+targetContact.Address+")"+ "to "+ sendToContact.Address)
	log.Println(network.routingTable.me.Address+"Contacts sent: ",len(wrapperMessage.GetMsg_2().Contacts))
	log.Println(network.routingTable.me.Address+" : KAD ID: ("+*wrapperMessage.GetMsg_2().KademliaTargetId+")"+ "to "+ sendToContact.Address)
	log.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")*/

	data := network.protobufhandler.MarshalMessage(wrapperMessage) //Marshal the message for network transport
	//unmarshaledMsg := network.protobufhandler.UnMarshalWrapperMessage(data)

	//log.Println(network.routingTable.me.Address+" :MARSHALED TARGET KAD ID: "+unmarshaledMsg.GetMsg_2().KademliaTargetId)

	sender := *NewSender(sendToContact.Address,&data) //Create sender to put on sender queue
	network.SendQueue <- sender                       //put sender on sender queue

}

func (network *Network) SendFindDataMessage(targetId *KademliaID, sendToContact Contact, requestID int64, isReply bool, foundFile bool, data []byte, responseContacts []Contact) {
	//Create the message and marshal it
	findDataMessage := network.protobufhandler.CreateLookupDataMessage(targetId, foundFile, data, responseContacts)
	wrapperMsg := network.protobufhandler.CreateWrapperMessage_3(targetId, requestID, protoMessages.MessageType_FIND_DATA, findDataMessage, isReply)
	marshaledMsg := network.protobufhandler.MarshalMessage(wrapperMsg)

	//Send the message with the queue
	sender := *NewSender(sendToContact.Address, &marshaledMsg)
	network.SendQueue <- sender
}

func (network *Network) SendStoreMessage(sendToContact *Contact, key []byte, data []byte, lifeTime time.Duration, requestId int64, isReply bool) {
	storeMsg := network.protobufhandler.CreateStoreMessage(key,data,lifeTime) //Create a store message
	//Create wrapper message for the store message
	wrapperMsg := network.protobufhandler.CreateWrapperMessage_4(network.routingTable.me.ID,requestId,protoMessages.MessageType_SEND_STORE,storeMsg,isReply)

	marshaledMsg := network.protobufhandler.MarshalMessage(wrapperMsg) //Marshal the message for network transport

	sender := *NewSender(sendToContact.Address,&marshaledMsg) //Create sender to put on sender queue
	network.SendQueue <- sender                       //put sender on sender queue
}

//RENAME "PROTOCONTACT" TO PROTOCONTACT INSTEAD OF CONTACT
func (network *Network) RecieveFindContactMessage(workRequest *WorkRequest) {

	//log.Println(network.routingTable.me.Address,": Recieved [Find Contact Request] from ",workRequest.senderAddress)
	//log.Println(network.routingTable.me.Address,": TARGET KAD ID: ",*workRequest.message.GetMsg_2().KademliaTargetId)

	targetKadId := NewKademliaID(workRequest.message.GetMsg_2().KademliaTargetId)
	targetContact := NewContact(targetKadId,"") //Ignore address, we only care about the target kademlia ID here
	contacts := network.routingTable.FindClosestContacts(targetKadId,3)
	sendContact := NewContact(nil,workRequest.senderAddress) //Ignore kad id, we only care about the address to send the response
	//lookupContactMsg := network.protobufhandler.CreateLookupContactMessage(targetKadId)

	/*log.Println(network.routingTable.me.Address,": Contacts to return:")
	log.Println(network.routingTable.me.Address,": <LIST START>")

	for _,elem := range contacts{
		log.Println(network.routingTable.me.Address+":",elem.Address)
		protoContact := network.protobufhandler.CreateContactMessage(elem.ID,elem.Address)
		lookupContactMsg.Contacts = append(lookupContactMsg.Contacts, protoContact)
	}
	log.Println(network.routingTable.me.Address,": <LIST END>")*/


	//wrapperMsg := network.protobufhandler.CreateWrapperMessage_2(network.routingTable.me.ID,workRequest.id,protoMessages.MessageType_FIND_CONTACT, lookupContactMsg,true)
	//log.Println(network.routingTable.me.Address,": <LIST OF FOUND CONTACTS START>")
	//for _,contact := range contacts{
	//	log.Println(network.routingTable.me.Address,": -->",contact.Address)

	//}
	//log.Println(network.routingTable.me.Address,": </LIST OF FOUND CONTACTS END>")

	network.SendFindContactMessage(&targetContact,&sendContact,workRequest.id,contacts,true)
	//marshaledMsg := network.protobufhandler.MarshalMessage(wrapperMsg)
	//unmarshaledMsg := network.protobufhandler.UnMarshalWrapperMessage(marshaledMsg)

	//log.Println(network.routingTable.me.Address+" :(SEND FROM RECEIVED) MARSHALED TARGET KAD ID: "+sendContact.Address)
	//sender := *NewSender(workRequest.senderAddress,&marshaledMsg) //Create sender to put on sender queue
	//network.SendQueue <- sender //put sender on sender queue
	//network.Send(workRequest.senderAddress,marshaledMsg)
}

func (network *Network) ReceivePingContactMessage(request *WorkRequest) {
	network.SendPingMessage(request.senderAddress, request.id, true) //Just need to ping right back
}

func (network *Network) ReceiveFindDataMessage(request *WorkRequest) {
	contacts := []Contact{}
	targetId := NewKademliaID(request.message.GetMsg_3().KademliaTargetId)
	targetIdBytes,err := hex.DecodeString(request.message.GetMsg_3().KademliaTargetId)
	CheckError(err)
	data, isFound := network.store.RetrieveData(targetIdBytes)
	if !isFound {
		contacts = network.routingTable.FindClosestContacts(targetId, 3)
	}
	network.SendFindDataMessage(targetId, NewContact(nil, request.senderAddress), request.id, true, isFound, data, contacts)
}

/*Receive store message.
Stores data with a specific key from the message in the key-value store.
 */
func (network *Network) RecieveStoreMessage(workRequest *WorkRequest) {

	key := workRequest.message.GetMsg_4().KeyStore
	data := workRequest.message.GetMsg_4().ValueStore
	lifeTime,error := time.ParseDuration(workRequest.message.GetMsg_4().LifeTime)
	if error != nil{
		log.Fatal("COULDN'T PARSE LIFETIME!")
	}
	republishTime := time.Now().Add(time.Minute) //TODO set to be 24 hours after testing
	network.store.StoreData([]byte(key),[]byte(data),lifeTime,republishTime,false)
}
