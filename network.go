package ptp
import (
	"github.com/timgranstrom/ptp-dfs/protoMessages"
	"fmt"
	"os"
	"net"
	"log"
	"strconv"
)

// A buffered channel that we can send work requests on.

//var WorkQueue = make(chan WorkRequest, 100)

type Network struct {
	protobufhandler *ProtobufHandler
	routingTable RoutingTable
	WorkerQueue chan Worker
	WorkQueue chan WorkRequest
	SendQueue chan Sender
	listenerActive bool
}

type Sender struct {
	address string
	data []byte
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
			SendQueue:make(chan Sender,500)}
	return network
}



/* A Simple function to verify error */
func CheckError(err error) {
	if err  != nil {
		log.Fatal("Error: " , err)
		os.Exit(0)
	}
}

func (network *Network)SetupUDPListener(address string) *net.UDPConn{
	/* Lets prepare a address at any address at :8000*/
	ServerAddr,err := net.ResolveUDPAddr("udp",address)
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	log.Println(address+": Connection Established at "+ServerAddr.IP.String()+":"+strconv.Itoa(ServerAddr.Port))
	CheckError(err)
	network.listenerActive = true
	return ServerConn
}

func (network *Network) Sender() {
	for{
		select {
			case sender := <- network.SendQueue:
				ServerAddr,err := net.ResolveUDPAddr("udp",sender.address)
				CheckError(err)

				LocalAddr, err := net.ResolveUDPAddr("udp", network.routingTable.me.Address)
				CheckError(err)

				Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
				CheckError(err)

				_,error := Conn.Write(sender.data)
				if error != nil {
					log.Fatal(err)
				}
				Conn.Close()
		}

	}
}

func (network *Network) Listen() {
	ServerConn := network.SetupUDPListener(network.routingTable.me.Address)
	defer ServerConn.Close()

	buffer := make([]byte, 1024)

	for { //Infinite for-loop to check for incomming messages
		_,addr,err := ServerConn.ReadFromUDP(buffer)
		go network.HandleRecievedMessage(buffer,addr,err)
	}
}

func (network *Network) HandleRecievedMessage(bufferMsg []byte,addr *net.UDPAddr,err error){
	msg := network.protobufhandler.UnMarshalWrapperMessage(bufferMsg)
	log.Println(network.routingTable.me.Address+" :Received ",msg.MessageType.String(), " from ",addr)

	if err != nil {
		fmt.Println("Error: ",err)
	}
	//take the message and make a WorkRequest out of them.
	work := WorkRequest{*msg.RequestId,*msg,addr.String()}

	//When handling work, make sure to always add the message sender as a contact
	kadId := NewKademliaID(*work.message.SenderKademliaId)
	requestContact := NewContact(kadId,work.senderAddress)
	network.routingTable.AddContact(requestContact)


	//Push the work onto the queue.
	network.WorkQueue <- work



	log.Println(network.routingTable.me.Address+": Work request queued")
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(targetContact *Contact,contact *Contact, requestId int64) {
	// TODO
	lookupContactMessage := network.protobufhandler.CreateLookupContactMessage(targetContact.ID) //Create a lookupContact message for the target contact

	//Create wrapper message for the request
	wrapperMessage := network.protobufhandler.CreateWrapperMessage_2(network.routingTable.me.ID,requestId,protoMessages.MessageType_FIND_CONTACT,lookupContactMessage,false)

	data := network.protobufhandler.MarshalMessage(wrapperMessage) //Marshal the message for network transport

	sender := *NewSender(contact.Address,&data) //Create sender to put on sender queue
	network.SendQueue <- sender //put sender on sender queue

	log.Println(network.routingTable.me.Address+" :Sent Find Contact Message ("+targetContact.Address+")"+ "to "+contact.Address)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

//RENAME "PROTOCONTACT" TO PROTOCONTACT INSTEAD OF CONTACT
func (network *Network) RecieveFindContactMessage(workRequest *WorkRequest) {

	log.Println(network.routingTable.me.Address,": Recieved find contact request from ",workRequest.senderAddress)

	targetKadId := NewKademliaID(*workRequest.message.GetMsg_2().KademliaTargetId)
	contacts := network.routingTable.FindClosestContacts(targetKadId,3)
	lookupContactMsg := network.protobufhandler.CreateLookupContactMessage(targetKadId)

	log.Println(network.routingTable.me.Address,": Contacts to return:")
	log.Println(network.routingTable.me.Address,": <LIST START>")

	for _,elem := range contacts{
		log.Println(network.routingTable.me.Address+":",elem.Address)
		protoContact := network.protobufhandler.CreateContactMessage(elem.ID,elem.Address)
		lookupContactMsg.Contacts = append(lookupContactMsg.Contacts, protoContact)
	}
	log.Println(network.routingTable.me.Address,": <LIST END>")


	wrapperMsg := network.protobufhandler.CreateWrapperMessage_2(network.routingTable.me.ID,workRequest.id,protoMessages.MessageType_FIND_CONTACT, lookupContactMsg,true)
	marshaledMsg := network.protobufhandler.MarshalMessage(wrapperMsg)
	log.Println(network.routingTable.me.Address+" :Sending response to "+workRequest.senderAddress)
	sender := *NewSender(workRequest.senderAddress,&marshaledMsg) //Create sender to put on sender queue
	network.SendQueue <- sender //put sender on sender queue
	//network.Send(workRequest.senderAddress,marshaledMsg)
}
