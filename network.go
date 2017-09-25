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
}

func NewNetwork(routingTable RoutingTable) *Network{
	network := &Network{routingTable:routingTable,
			WorkerQueue: make(chan Worker,100),
			WorkQueue:make(chan WorkRequest,100)}
	return network
}



/* A Simple function to verify error */
func CheckError(err error) {
	if err  != nil {
		log.Fatal("Error: " , err)
		os.Exit(0)
	}
}

func SetupUDPListener(address string) *net.UDPConn{
	/* Lets prepare a address at any address at :8000*/
	ServerAddr,err := net.ResolveUDPAddr("udp",address)
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	log.Println("Connection Established at "+ServerAddr.IP.String()+":"+strconv.Itoa(ServerAddr.Port))
	CheckError(err)
	return ServerConn
}

func (network *Network) Send(address string, buffMsg []byte) {

	ServerAddr,err := net.ResolveUDPAddr("udp",address)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", network.routingTable.me.Address)
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)
	defer Conn.Close()

	_,error := Conn.Write(buffMsg)
	if error != nil {
		log.Fatal(err)
	}
}

func (network *Network) Listen() {
	ServerConn := SetupUDPListener(network.routingTable.me.Address)
	defer ServerConn.Close()

	buffer := make([]byte, 1024)

	for { //Infinite for-loop to check for incomming messages
		_,addr,err := ServerConn.ReadFromUDP(buffer)
		go network.HandleRecievedMessage(buffer,addr,err)
	}
}

func (network *Network) HandleRecievedMessage(bufferMsg []byte,addr *net.UDPAddr,err error){
	msg := network.protobufhandler.UnMarshalWrapperMessage(bufferMsg)
	log.Println("Received ",msg.MessageType.String(), " from ",addr)

	if err != nil {
		fmt.Println("Error: ",err)
	}
	//Now, we take the delay, and the person's name, and make a WorkRequest out of them.
	work := WorkRequest{*msg.RequestId,*msg}

	//Push the work onto the queue.
	network.WorkQueue <- work
	log.Println("Work request queued")
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact, requestId int64) {
	// TODO
	lookupContactMessage := network.protobufhandler.CreateLookupContactMessage(contact.ID) //Create a lookupContact message for the target contact

	//Create wrapper message for the request
	wrapperMessage := network.protobufhandler.CreateWrapperMessage_2(network.routingTable.me.ID,requestId,protoMessages.MessageType_FIND_CONTACT,lookupContactMessage,false)

	data := network.protobufhandler.MarshalMessage(wrapperMessage) //Marshal the message for network transport

	network.Send(contact.Address,data) //Send Message to target address
	log.Println("Sent Find Contact Message to "+contact.Address)

	//unwrappedMsg := network.protobufhandler.UnMarshalWrapperMessage(data)

	//println(unwrappedMsg.SenderKademliaId)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

//TODO FIX THIS
func (network *Network) RecieveFindContactMessage(workRequest *WorkRequest) {
	log.Println("RECIEVED FIND CONTACT MESSAGE REQUEST")
	//targetKadId := NewKademliaID(*workRequest.message.GetMsg_2().KademliaTargetId)
	//contacts := network.routingTable.FindClosestContacts(targetKadId,3)
	//lookupContactMsg := network.protobufhandler.CreateLookupContactMessage(targetKadId)
	/*for _,elem := range contacts{
		lookupContactMsg.Contacts = append(lookupContactMsg.Contacts, elem)
	}*/
	//lookupContactMsg.Contacts = contacts
}
