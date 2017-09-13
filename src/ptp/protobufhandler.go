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


/**
********************************************************
************* CREATE MESSAGES PART *********************
********************************************************
 */

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

/*
Create a Ping Message
 */
func (protobufHandler *ProtobufHandler) CreatePingMessage() *protoMessages.PingMessage{
	pingMessage := &protoMessages.PingMessage{
	}
	return pingMessage
}

/*
Create a Store Message
 */
func (protobufHandler *ProtobufHandler) CreateStoreMessage(kademliaId *KademliaID, address string) *protoMessages.StoreMessage{
	storeMessage := &protoMessages.StoreMessage{
		KeyStore: proto.String(kademliaId.String()), //Set kademlia id as key
		ValueStore: proto.String(address), //Set ip address as stored value
	}
	return storeMessage
}

/*
Create a Lookup Data Message
 */
func (protobufHandler *ProtobufHandler) CreateLookupDataMessage(kademliaId *KademliaID) *protoMessages.LookupDataMessage{
	lookupDataMessage := &protoMessages.LookupDataMessage{
		KademliaTargetId: proto.String(kademliaId.String()),
	}
	return lookupDataMessage
}

