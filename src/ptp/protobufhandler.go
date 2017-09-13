package ptp

import (
	"ptp/proto"
	"github.com/golang/protobuf/proto"
	"log"
	"fmt"
)

type ProtobufHandler struct{}

func (protobufHandler *ProtobufHandler) UnMarshalWrapperMessage(message []byte) {
	unwrappedMessage := &protoMessages.WrapperMessage{} //Create the variable to store the unwrapped result
	proto.Unmarshal(message,unwrappedMessage)

	log.Println("SUCCESS! ")
}

func (protobufHandler *ProtobufHandler) CreateLookupContactMessage() *protoMessages.FindContactMessage{
	kadID:=NewKademliaID("hej123")
	log.Println("TEST: ")
	kadString := kadID.String()
	fmt.Println(kadString) //IT WORKS, PROBLEM IS THAT IT CANNOT PRINT THIS LARGE STRING MAYBE!?

	lookupContactMessage := &protoMessages.FindContactMessage{
		KademliaTargetId: proto.String("hej"),
	}
	return lookupContactMessage
}
