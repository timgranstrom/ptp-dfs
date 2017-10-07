package ptp

import (
	"github.com/timgranstrom/ptp-dfs/protoMessages"
	"github.com/golang/protobuf/proto"
	//"log"
	"log"
	"time"
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
Create a Wrapper Message with ping message
*/
func (protobufHandler *ProtobufHandler) CreateWrapperMessage_1(senderKademliaId *KademliaID, requestId int64, messageType protoMessages.MessageType, message *protoMessages.PingMessage, isReply bool) *protoMessages.WrapperMessage{
	wrapperMessage := &protoMessages.WrapperMessage{
		SenderKademliaId: senderKademliaId.String(),
		MessageType:messageType,
		RequestId:requestId,
		IsReply:isReply,
	}
	wrappedMsg := protoMessages.WrapperMessage_Msg_1{message}
	wrapperMessage.Messages = &wrappedMsg

	return wrapperMessage
}

/*
Create a Wrapper Message with lookup contact message
*/
func (protobufHandler *ProtobufHandler) CreateWrapperMessage_2(senderKademliaId *KademliaID, requestId int64, messageType protoMessages.MessageType, message *protoMessages.LookupContactMessage, isReply bool) *protoMessages.WrapperMessage{
	wrapperMessage := &protoMessages.WrapperMessage{
		SenderKademliaId: senderKademliaId.String(),
		MessageType:messageType,
		RequestId:requestId,
		IsReply:isReply,
	}
	wrappedMsg := protoMessages.WrapperMessage_Msg_2{message}
	wrapperMessage.Messages = &wrappedMsg

	return wrapperMessage
}

/*
Create a Wrapper Message with lookup data message
*/
func (protobufHandler *ProtobufHandler) CreateWrapperMessage_3(senderKademliaId *KademliaID, requestId int64, messageType protoMessages.MessageType, message *protoMessages.LookupDataMessage, isReply bool) *protoMessages.WrapperMessage{
	wrapperMessage := &protoMessages.WrapperMessage{
		SenderKademliaId: senderKademliaId.String(),
		MessageType:messageType,
		RequestId:requestId,
		IsReply:isReply,
	}
	wrappedMsg := protoMessages.WrapperMessage_Msg_3{message}
	wrapperMessage.Messages = &wrappedMsg

	return wrapperMessage
}

/*
Create a Wrapper Message with store message
*/
func (protobufHandler *ProtobufHandler) CreateWrapperMessage_4(senderKademliaId *KademliaID, requestId int64, messageType protoMessages.MessageType, message *protoMessages.StoreMessage, isReply bool) *protoMessages.WrapperMessage{
	wrapperMessage := &protoMessages.WrapperMessage{
		SenderKademliaId: senderKademliaId.String(),
		MessageType:messageType,
		RequestId:requestId,
		IsReply:isReply,
	}
	wrappedMsg := protoMessages.WrapperMessage_Msg_4{message}
	wrapperMessage.Messages = &wrappedMsg

	return wrapperMessage
}

/*
Create a Lookup Contact Message
 */
func (protobufHandler *ProtobufHandler) CreateLookupContactMessage(KademliaTargetId *KademliaID) *protoMessages.LookupContactMessage{
	lookupContactMessage := &protoMessages.LookupContactMessage{
		KademliaTargetId: KademliaTargetId.String(),
	}
	return lookupContactMessage
}

/*
Create a proto contact
 */
func (protobufHandler *ProtobufHandler) CreateContactMessage(kademliaId *KademliaID, address string) *protoMessages.ProtoContact{
	contactMessage := &protoMessages.ProtoContact{
		KademliaId: kademliaId.String(),
		Address: address,
	}
	return contactMessage
}

/*
Create a list of proto contacts
 */
func (protobufHandler *ProtobufHandler) CreateContactMessages(contacts []Contact) []*protoMessages.ProtoContact{
	protoContacts := []*protoMessages.ProtoContact{}
	for _,elem := range contacts{
	contactMessage := &protoMessages.ProtoContact{
		KademliaId: elem.ID.String(),
		Address: elem.Address,
	}
		protoContacts = append(protoContacts, contactMessage)
	}
	return protoContacts
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
func (protobufHandler *ProtobufHandler) CreateStoreMessage(key []byte, data []byte,lifeTime time.Duration) *protoMessages.StoreMessage{
	storeMessage := &protoMessages.StoreMessage{
		KeyStore: string(key), //Set kademlia id as key
		ValueStore: string(data), //Set ip address as stored value
	}
	return storeMessage
}

/*
Create a Lookup Data Message
 */
func (protobufHandler *ProtobufHandler) CreateLookupDataMessage(kademliaId *KademliaID, foundFile bool, data []byte, contacts []Contact) *protoMessages.LookupDataMessage{
	lookupDataMessage := &protoMessages.LookupDataMessage{
		KademliaTargetId: kademliaId.String(),
		FoundFile: foundFile,
	}
	if foundFile {
		lookupDataMessage.FileData = string(data)
	} else {
		lookupDataMessage.Contacts = protobufHandler.CreateContactMessages(contacts)
	}

	return lookupDataMessage
}

func ConvertProtobufContacts(protoContacts []*protoMessages.ProtoContact, me Contact) []Contact {
	contacts := []Contact{}
	for _,protoContact := range protoContacts {
		protoKademliaID := NewKademliaID(protoContact.GetKademliaId())
		if me.ID != protoKademliaID {
			replyContact := NewContact(protoKademliaID, protoContact.GetAddress())
			contacts = append(contacts, replyContact)
		} else {
			log.Println(me.Address,": Filtered from sending LOOKUP_DATA message to self")
		}
	}
	return contacts
}

