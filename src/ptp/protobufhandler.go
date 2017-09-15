package ptp

import (
	"ptp/proto"
	"github.com/golang/protobuf/proto"
	//"log"
	"log"
)

type ProtobufHandler struct{}

/**
********************************************************
************* MARSHAL MESSAGE PART *********************
********************************************************
 */

/**
Marshal a Message and return the marshaled message data
*/
func (protobufHandler *ProtobufHandler) MarshalMessage(message proto.Message) []byte {
	data, err := proto.Marshal(message) //Marshal the message to byte data and store error if any
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	return data
}

/**
********************************************************
************* UNMARSHAL MESSAGES PART *********************
********************************************************
 */

/**
Unmarshal a Ping message and return the unmarshaled message
*/
func (protobufHandler *ProtobufHandler) UnMarshalPingMessage(message []byte) *protoMessages.PingMessage {
	unmarshaledMessage := &protoMessages.PingMessage{} //Create the variable to store the unmarshaled result
	proto.Unmarshal(message,unmarshaledMessage)
	return unmarshaledMessage
}

/**
Unmarshal a Store message and return the unmarshaled message
*/
func (protobufHandler *ProtobufHandler) UnMarshalStoreMessage(message []byte) *protoMessages.StoreMessage {
	unmarshaledMessage := &protoMessages.StoreMessage{} //Create the variable to store the unmarshaled result
	proto.Unmarshal(message,unmarshaledMessage)
	return unmarshaledMessage
}

/**
Unmarshal a Lookup Data message and return the unmarshaled message
*/
func (protobufHandler *ProtobufHandler) UnMarshalLookupDataMessage(message []byte) *protoMessages.LookupDataMessage {
	unmarshaledMessage := &protoMessages.LookupDataMessage{} //Create the variable to store the unmarshaled result
	proto.Unmarshal(message,unmarshaledMessage)
	return unmarshaledMessage
}

/**
Unmarshal a Lookup Contact message and return the unmarshaled message
*/
func (protobufHandler *ProtobufHandler) UnMarshalLookupContactMessage(message []byte) *protoMessages.LookupContactMessage {
	unmarshaledMessage := &protoMessages.LookupContactMessage{} //Create the variable to store the unmarshaled result
	proto.Unmarshal(message,unmarshaledMessage)
	return unmarshaledMessage
}

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
Create a Wrapper Message
User still needs to add the sub-message manually after creation.
*/
func (protobufHandler *ProtobufHandler) CreateWrapperMessage(kademliaId *KademliaID, requestId int64, messageType protoMessages.MessageType, message proto.Message) *protoMessages.WrapperMessage{
	wrapperMessage := &protoMessages.WrapperMessage{
		SenderKademliaId: proto.String(kademliaId.String()),
		MessageType:&messageType,
		RequestId:proto.Int64(requestId),
	}

	switch messageType{
	case protoMessages.MessageType_FIND_CONTACT:
		test := protoMessages.LookupContactMessage(message)
		wrappedMsg := protoMessages.WrapperMessage_Msg_2{&test}
		wrapperMessage.Messages = &wrappedMsg
		println(messageType.String())
	}
	return wrapperMessage
}

/*
Create a Lookup Contact Message
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

