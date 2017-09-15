package ptp
import (
	//"log"
	//"github.com/golang/protobuf/proto"
	//"net"
	"ptp/proto"
	//"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/proto"
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

	//wrapMsg := protoMessages.WrapperMessage_Msg_2{lookupContactMessage}

	wrapperMessage := probufHandler.CreateWrapperMessage(contact.ID,45,protoMessages.MessageType_FIND_CONTACT,proto.Message(lookupContactMessage))
	//wrapperMessage.Messages = &wrapMsg
	//data,_ := proto.Marshal(wrapperMessage) //Marshal the wrapper message
	data := probufHandler.MarshalMessage(wrapperMessage)
	unwrappedMsg := probufHandler.UnMarshalWrapperMessage(data)

	print(unwrappedMsg.SenderKademliaId)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
