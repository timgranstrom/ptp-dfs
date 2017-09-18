package ptp
import (
	//"log"
	//"github.com/golang/protobuf/proto"
	//"net"
	"ptp/proto"
	//"github.com/golang/protobuf/proto"
	//"github.com/golang/protobuf/proto"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
	probufHandler := ProtobufHandler{}
	lookupContactMessage := probufHandler.CreateLookupContactMessage(contact.ID)

	wrapperMessage := probufHandler.CreateWrapperMessage_2(contact.ID,45,protoMessages.MessageType_FIND_CONTACT,lookupContactMessage)

	data := probufHandler.MarshalMessage(wrapperMessage)
	unwrappedMsg := probufHandler.UnMarshalWrapperMessage(data)

	println(unwrappedMsg.SenderKademliaId)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
