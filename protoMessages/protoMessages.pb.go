// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protoMessages.proto

/*
Package protoMessages is a generated protocol buffer package.

It is generated from these files:
	protoMessages.proto

It has these top-level messages:
	WrapperMessage
	PingMessage
	LookupContactMessage
	ProtoContact
	LookupDataMessage
	StoreMessage
*/
package protoMessages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Enum For Message Types
type MessageType int32

const (
	MessageType_PING         MessageType = 0
	MessageType_FIND_CONTACT MessageType = 1
	MessageType_FIND_DATA    MessageType = 2
	MessageType_SEND_STORE   MessageType = 3
)

var MessageType_name = map[int32]string{
	0: "PING",
	1: "FIND_CONTACT",
	2: "FIND_DATA",
	3: "SEND_STORE",
}
var MessageType_value = map[string]int32{
	"PING":         0,
	"FIND_CONTACT": 1,
	"FIND_DATA":    2,
	"SEND_STORE":   3,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (MessageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Common wrapper for all 4 types of messages
type WrapperMessage struct {
	MessageType      MessageType `protobuf:"varint,1,opt,name=message_type,json=messageType,enum=protoMessages.MessageType" json:"message_type,omitempty"`
	RequestId        int64       `protobuf:"varint,2,opt,name=request_id,json=requestId" json:"request_id,omitempty"`
	SenderKademliaId string      `protobuf:"bytes,3,opt,name=sender_kademlia_id,json=senderKademliaId" json:"sender_kademlia_id,omitempty"`
	IsReply          bool        `protobuf:"varint,4,opt,name=is_reply,json=isReply" json:"is_reply,omitempty"`
	// messages { //Type-of-message-specific header, only one can be used
	Msg_1 *PingMessage          `protobuf:"bytes,5,opt,name=msg_1,json=msg1" json:"msg_1,omitempty"`
	Msg_2 *LookupContactMessage `protobuf:"bytes,6,opt,name=msg_2,json=msg2" json:"msg_2,omitempty"`
	Msg_3 *LookupDataMessage    `protobuf:"bytes,7,opt,name=msg_3,json=msg3" json:"msg_3,omitempty"`
	Msg_4 *StoreMessage         `protobuf:"bytes,8,opt,name=msg_4,json=msg4" json:"msg_4,omitempty"`
}

func (m *WrapperMessage) Reset()                    { *m = WrapperMessage{} }
func (m *WrapperMessage) String() string            { return proto.CompactTextString(m) }
func (*WrapperMessage) ProtoMessage()               {}
func (*WrapperMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *WrapperMessage) GetMessageType() MessageType {
	if m != nil {
		return m.MessageType
	}
	return MessageType_PING
}

func (m *WrapperMessage) GetRequestId() int64 {
	if m != nil {
		return m.RequestId
	}
	return 0
}

func (m *WrapperMessage) GetSenderKademliaId() string {
	if m != nil {
		return m.SenderKademliaId
	}
	return ""
}

func (m *WrapperMessage) GetIsReply() bool {
	if m != nil {
		return m.IsReply
	}
	return false
}

func (m *WrapperMessage) GetMsg_1() *PingMessage {
	if m != nil {
		return m.Msg_1
	}
	return nil
}

func (m *WrapperMessage) GetMsg_2() *LookupContactMessage {
	if m != nil {
		return m.Msg_2
	}
	return nil
}

func (m *WrapperMessage) GetMsg_3() *LookupDataMessage {
	if m != nil {
		return m.Msg_3
	}
	return nil
}

func (m *WrapperMessage) GetMsg_4() *StoreMessage {
	if m != nil {
		return m.Msg_4
	}
	return nil
}

// Sub-header to ping a contact
type PingMessage struct {
}

func (m *PingMessage) Reset()                    { *m = PingMessage{} }
func (m *PingMessage) String() string            { return proto.CompactTextString(m) }
func (*PingMessage) ProtoMessage()               {}
func (*PingMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// Sub-header to find contacts, returns a "list" of contacts
type LookupContactMessage struct {
	KademliaTargetId string          `protobuf:"bytes,9,opt,name=kademlia_target_id,json=kademliaTargetId" json:"kademlia_target_id,omitempty"`
	Contacts         []*ProtoContact `protobuf:"bytes,10,rep,name=contacts" json:"contacts,omitempty"`
}

func (m *LookupContactMessage) Reset()                    { *m = LookupContactMessage{} }
func (m *LookupContactMessage) String() string            { return proto.CompactTextString(m) }
func (*LookupContactMessage) ProtoMessage()               {}
func (*LookupContactMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *LookupContactMessage) GetKademliaTargetId() string {
	if m != nil {
		return m.KademliaTargetId
	}
	return ""
}

func (m *LookupContactMessage) GetContacts() []*ProtoContact {
	if m != nil {
		return m.Contacts
	}
	return nil
}

// Protobuf representation of Contact.go
type ProtoContact struct {
	KademliaId string `protobuf:"bytes,11,opt,name=kademlia_id,json=kademliaId" json:"kademlia_id,omitempty"`
	Address    string `protobuf:"bytes,12,opt,name=address" json:"address,omitempty"`
}

func (m *ProtoContact) Reset()                    { *m = ProtoContact{} }
func (m *ProtoContact) String() string            { return proto.CompactTextString(m) }
func (*ProtoContact) ProtoMessage()               {}
func (*ProtoContact) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ProtoContact) GetKademliaId() string {
	if m != nil {
		return m.KademliaId
	}
	return ""
}

func (m *ProtoContact) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// Sub-header to find data with a certain id, returns closer contacts if it wasn't found
type LookupDataMessage struct {
	KademliaTargetId string          `protobuf:"bytes,13,opt,name=kademlia_target_id,json=kademliaTargetId" json:"kademlia_target_id,omitempty"`
	FoundFile        bool            `protobuf:"varint,14,opt,name=found_file,json=foundFile" json:"found_file,omitempty"`
	FileData         string          `protobuf:"bytes,15,opt,name=file_data,json=fileData" json:"file_data,omitempty"`
	Contacts         []*ProtoContact `protobuf:"bytes,16,rep,name=contacts" json:"contacts,omitempty"`
}

func (m *LookupDataMessage) Reset()                    { *m = LookupDataMessage{} }
func (m *LookupDataMessage) String() string            { return proto.CompactTextString(m) }
func (*LookupDataMessage) ProtoMessage()               {}
func (*LookupDataMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LookupDataMessage) GetKademliaTargetId() string {
	if m != nil {
		return m.KademliaTargetId
	}
	return ""
}

func (m *LookupDataMessage) GetFoundFile() bool {
	if m != nil {
		return m.FoundFile
	}
	return false
}

func (m *LookupDataMessage) GetFileData() string {
	if m != nil {
		return m.FileData
	}
	return ""
}

func (m *LookupDataMessage) GetContacts() []*ProtoContact {
	if m != nil {
		return m.Contacts
	}
	return nil
}

// Sub-header to share where a nodes shared data can be found
type StoreMessage struct {
	KeyStore   string `protobuf:"bytes,17,opt,name=key_store,json=keyStore" json:"key_store,omitempty"`
	ValueStore string `protobuf:"bytes,18,opt,name=value_store,json=valueStore" json:"value_store,omitempty"`
}

func (m *StoreMessage) Reset()                    { *m = StoreMessage{} }
func (m *StoreMessage) String() string            { return proto.CompactTextString(m) }
func (*StoreMessage) ProtoMessage()               {}
func (*StoreMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *StoreMessage) GetKeyStore() string {
	if m != nil {
		return m.KeyStore
	}
	return ""
}

func (m *StoreMessage) GetValueStore() string {
	if m != nil {
		return m.ValueStore
	}
	return ""
}

func init() {
	proto.RegisterType((*WrapperMessage)(nil), "protoMessages.WrapperMessage")
	proto.RegisterType((*PingMessage)(nil), "protoMessages.PingMessage")
	proto.RegisterType((*LookupContactMessage)(nil), "protoMessages.LookupContactMessage")
	proto.RegisterType((*ProtoContact)(nil), "protoMessages.ProtoContact")
	proto.RegisterType((*LookupDataMessage)(nil), "protoMessages.LookupDataMessage")
	proto.RegisterType((*StoreMessage)(nil), "protoMessages.StoreMessage")
	proto.RegisterEnum("protoMessages.MessageType", MessageType_name, MessageType_value)
}

func init() { proto.RegisterFile("protoMessages.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 503 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0xc9, 0xda, 0xad, 0xc9, 0x4b, 0x5b, 0x32, 0xc3, 0xc1, 0x30, 0x4d, 0x44, 0xe1, 0x12,
	0x21, 0x34, 0x58, 0x3b, 0x04, 0x17, 0x0e, 0x55, 0xbb, 0x41, 0xc4, 0xe8, 0xaa, 0x34, 0x12, 0x47,
	0xcb, 0xcc, 0x5e, 0x14, 0x25, 0x69, 0x42, 0x9c, 0x22, 0xe5, 0xc0, 0x27, 0xe0, 0xfb, 0xf0, 0xf9,
	0x90, 0x9d, 0x6c, 0x4b, 0xb7, 0xee, 0xb0, 0x5b, 0xde, 0xff, 0xff, 0x7e, 0x7e, 0xcf, 0x7f, 0xb7,
	0xf0, 0x2c, 0x2f, 0xb2, 0x32, 0xfb, 0xce, 0x85, 0xa0, 0x21, 0x17, 0x47, 0xaa, 0x42, 0x83, 0x0d,
	0xd1, 0xf9, 0xdb, 0x81, 0xe1, 0x8f, 0x82, 0xe6, 0x39, 0x2f, 0x1a, 0x0d, 0x7d, 0x86, 0x7e, 0x5a,
	0x7f, 0x92, 0xb2, 0xca, 0x39, 0xd6, 0x6c, 0xcd, 0x1d, 0x8e, 0x5e, 0x1e, 0x6d, 0x9e, 0xd6, 0x7c,
	0x04, 0x55, 0xce, 0x7d, 0x33, 0xbd, 0x2d, 0xd0, 0x21, 0x40, 0xc1, 0x7f, 0xad, 0xb9, 0x28, 0x49,
	0xc4, 0xf0, 0x8e, 0xad, 0xb9, 0x1d, 0xdf, 0x68, 0x14, 0x8f, 0xa1, 0xb7, 0x80, 0x04, 0x5f, 0x31,
	0x5e, 0x90, 0x98, 0x32, 0x9e, 0x26, 0x11, 0x95, 0x6d, 0x1d, 0x5b, 0x73, 0x0d, 0xdf, 0xaa, 0x9d,
	0x6f, 0x8d, 0xe1, 0x31, 0xf4, 0x02, 0xf4, 0x48, 0x90, 0x82, 0xe7, 0x49, 0x85, 0xbb, 0xb6, 0xe6,
	0xea, 0x7e, 0x2f, 0x12, 0xbe, 0x2c, 0xd1, 0x3b, 0xd8, 0x4d, 0x45, 0x48, 0x8e, 0xf1, 0xae, 0xad,
	0xb9, 0xe6, 0xbd, 0xfd, 0x16, 0xd1, 0x2a, 0x6c, 0x0a, 0xbf, 0x9b, 0x8a, 0xf0, 0x18, 0x7d, 0xaa,
	0x81, 0x11, 0xde, 0x53, 0xc0, 0xeb, 0x3b, 0xc0, 0x79, 0x96, 0xc5, 0xeb, 0x7c, 0x9a, 0xad, 0x4a,
	0x7a, 0x59, 0xb6, 0xc9, 0x11, 0xfa, 0x50, 0x93, 0x63, 0xdc, 0x53, 0xa4, 0xbd, 0x95, 0x9c, 0xd1,
	0x92, 0xb6, 0xb1, 0x31, 0x7a, 0x5f, 0x63, 0x27, 0x58, 0x57, 0xd8, 0xc1, 0x1d, 0x6c, 0x59, 0x66,
	0x05, 0x6f, 0x13, 0x27, 0xce, 0x00, 0xcc, 0xd6, 0xde, 0xce, 0x1f, 0x78, 0xbe, 0x6d, 0x2b, 0x99,
	0xe1, 0x4d, 0x78, 0x25, 0x2d, 0x42, 0xae, 0xa2, 0x36, 0xea, 0x0c, 0xaf, 0x9d, 0x40, 0x19, 0x1e,
	0x43, 0x1f, 0x41, 0xbf, 0xac, 0x79, 0x81, 0xc1, 0xee, 0x6c, 0xd9, 0x64, 0x21, 0xab, 0x66, 0x86,
	0x7f, 0xd3, 0xec, 0x78, 0xd0, 0x6f, 0x3b, 0xe8, 0x15, 0x98, 0xed, 0x37, 0x33, 0xd5, 0x3c, 0x88,
	0x6f, 0x5f, 0x0b, 0x43, 0x8f, 0x32, 0x56, 0x70, 0x21, 0x70, 0x5f, 0x99, 0xd7, 0xa5, 0xf3, 0x4f,
	0x83, 0xfd, 0x7b, 0x31, 0x3d, 0x70, 0x8f, 0xc1, 0x03, 0xf7, 0x38, 0x04, 0xb8, 0xca, 0xd6, 0x2b,
	0x46, 0xae, 0xa2, 0x84, 0xe3, 0xa1, 0xfa, 0x35, 0x18, 0x4a, 0x39, 0x8b, 0x12, 0x8e, 0x0e, 0xc0,
	0x90, 0x06, 0x61, 0xb4, 0xa4, 0xf8, 0xa9, 0x3a, 0x43, 0x97, 0x82, 0x1c, 0xb8, 0x91, 0x81, 0xf5,
	0x98, 0x0c, 0xce, 0xa1, 0xdf, 0x7e, 0x27, 0x39, 0x25, 0xe6, 0x15, 0x11, 0x52, 0xc3, 0xfb, 0xf5,
	0x94, 0x98, 0x57, 0xaa, 0x47, 0x06, 0xf4, 0x9b, 0x26, 0x6b, 0xde, 0xd8, 0xa8, 0x0e, 0x48, 0x49,
	0xaa, 0xe1, 0xcd, 0x57, 0x30, 0x5b, 0xff, 0x1b, 0xa4, 0x43, 0x77, 0xe1, 0xcd, 0xbf, 0x58, 0x4f,
	0x90, 0x05, 0xfd, 0x33, 0x6f, 0x3e, 0x23, 0xd3, 0x8b, 0x79, 0x30, 0x99, 0x06, 0x96, 0x86, 0x06,
	0x60, 0x28, 0x65, 0x36, 0x09, 0x26, 0xd6, 0x0e, 0x1a, 0x02, 0x2c, 0x4f, 0xe7, 0x33, 0xb2, 0x0c,
	0x2e, 0xfc, 0x53, 0xab, 0xf3, 0x73, 0x4f, 0x6d, 0x3f, 0xfe, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x8c,
	0x62, 0x65, 0x2f, 0xe4, 0x03, 0x00, 0x00,
}
