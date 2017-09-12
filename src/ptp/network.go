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
	testMessage := &protoMessages.KademliaMessage{
		KademliaID: proto.String("FuckingWork"),
		MethodType: protoMessages.MessageType.Enum(1),
	}

	data, err := proto.Marshal(testMessage)
	if(err != nil){
		log.Fatal("MARSHALING ERROR: ",err)
	}else{
		newMessage := &protoMessages.KademliaMessage{}
		proto.Unmarshal(data,newMessage)

		log.Println("SUCCESS! "+newMessage.GetKademliaID())
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
