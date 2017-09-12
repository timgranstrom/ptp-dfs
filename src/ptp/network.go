package ptp
import (
	"log"
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

	lookupContactMessage := &protoMessages.FindContactMessage{
		KademliaTargetId: proto.String("KadTestId"),

	}

	wrapperMessage := &protoMessages.WrapperMessage{
		MessageId: proto.Int64(11),
		MessageType: protoMessages.MessageType.Enum(1),
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
