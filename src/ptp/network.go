package ptp
import (
	//"log"
	"github.com/golang/protobuf/proto"
	//"net"
	"ptp/proto"
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
	kadID := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	lookupContactMessage := probufHandler.CreateLookupContactMessage(kadID)

	wrapMsg := protoMessages.WrapperMessage_Msg_2{lookupContactMessage}

	wrapperMessage := &protoMessages.WrapperMessage{
		RequestId: proto.Int64(11),
		MessageType: protoMessages.MessageType.Enum(1),
		Messages:&wrapMsg,
	}

	data,_ := proto.Marshal(wrapperMessage) //Marshal the wrapper message

	probufHandler.UnMarshalWrapperMessage(data)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
