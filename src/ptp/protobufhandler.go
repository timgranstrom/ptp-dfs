package ptp

import (
	"ptp/proto"
	"github.com/golang/protobuf/proto"
	//"log"
)

type ProtobufHandler struct{}

/**
Unmarshal a wrapper message and return the unmarshaled message
 */
func (protobufHandler *ProtobufHandler) UnMarshalWrapperMessage(message []byte) *protoMessages.WrapperMessage {
	unwrappedMessage := &protoMessages.WrapperMessage{} //Create the variable to store the unwrapped result
	proto.Unmarshal(message,unwrappedMessage)
	return unwrappedMessage
}
/*
Create a LookupContactMessage
 */
func (protobufHandler *ProtobufHandler) CreateLookupContactMessage(kademliaId *KademliaID) *protoMessages.LookupContactMessage{
	lookupContactMessage := &protoMessages.LookupContactMessage{
		KademliaTargetId: proto.String(kademliaId.String()),
	}
	return lookupContactMessage
}
/*
Create a contact in message form
 */
func (protobufHandler *ProtobufHandler) CreateContactMessage(kademliaId *KademliaID, address string) *protoMessages.Contact{
	contactMessage := &protoMessages.Contact{
		KademliaId: proto.String(kademliaId.String()),
		Address: proto.String(address),
	}
	return contactMessage
}


