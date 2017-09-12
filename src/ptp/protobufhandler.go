package ptp

import (
	"ptp/proto"
	"github.com/golang/protobuf/proto"
	"log"
)

type ProtobufHandler struct{}

func (protobufHandler *ProtobufHandler) UnMarshalWrapperMessage(message []byte) {
	unwrappedMessage := &protoMessages.WrapperMessage{} //Create the variable to store the unwrapped result
	proto.Unmarshal(message,unwrappedMessage)

	log.Println("SUCCESS! ")
}

func (protobufHandler *ProtobufHandler)
