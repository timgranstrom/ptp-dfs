package ptp
import (
	//"log"
	//"net"
	"ptp/proto"
)

type Network struct {
	protobufhandler *ProtobufHandler
}

func NewNetwork() *Network{
	network := &Network{}
	return network
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
	lookupContactMessage := network.protobufhandler.CreateLookupContactMessage(contact.ID)

	wrapperMessage := network.protobufhandler.CreateWrapperMessage_2(contact.ID,45,protoMessages.MessageType_FIND_CONTACT,lookupContactMessage)

	data := network.protobufhandler.MarshalMessage(wrapperMessage)
	unwrappedMsg := network.protobufhandler.UnMarshalWrapperMessage(data)

	println(unwrappedMsg.SenderKademliaId)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
