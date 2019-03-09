package mtproto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	protocol "v2ray.com/core/common/protocol"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Account struct {
	Secret               []byte   `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_64514e21c693811b, []int{0}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetSecret() []byte {
	if m != nil {
		return m.Secret
	}
	return nil
}

type ServerConfig struct {
	// User is a list of users that allowed to connect to this inbound.
	// Although this is a repeated field, only the first user is effective for now.
	User                 []*protocol.User `protobuf:"bytes,1,rep,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_64514e21c693811b, []int{1}
}

func (m *ServerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerConfig.Unmarshal(m, b)
}
func (m *ServerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerConfig.Marshal(b, m, deterministic)
}
func (m *ServerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerConfig.Merge(m, src)
}
func (m *ServerConfig) XXX_Size() int {
	return xxx_messageInfo_ServerConfig.Size(m)
}
func (m *ServerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ServerConfig proto.InternalMessageInfo

func (m *ServerConfig) GetUser() []*protocol.User {
	if m != nil {
		return m.User
	}
	return nil
}

type ClientConfig struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClientConfig) Reset()         { *m = ClientConfig{} }
func (m *ClientConfig) String() string { return proto.CompactTextString(m) }
func (*ClientConfig) ProtoMessage()    {}
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_64514e21c693811b, []int{2}
}

func (m *ClientConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientConfig.Unmarshal(m, b)
}
func (m *ClientConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientConfig.Marshal(b, m, deterministic)
}
func (m *ClientConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientConfig.Merge(m, src)
}
func (m *ClientConfig) XXX_Size() int {
	return xxx_messageInfo_ClientConfig.Size(m)
}
func (m *ClientConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ClientConfig proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Account)(nil), "v2ray.core.proxy.mtproto.Account")
	proto.RegisterType((*ServerConfig)(nil), "v2ray.core.proxy.mtproto.ServerConfig")
	proto.RegisterType((*ClientConfig)(nil), "v2ray.core.proxy.mtproto.ClientConfig")
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/mtproto/config.proto", fileDescriptor_64514e21c693811b)
}

var fileDescriptor_64514e21c693811b = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x89, 0xca, 0x2e, 0xc4, 0xe2, 0xa1, 0x07, 0x09, 0xe2, 0xa1, 0xf6, 0xb4, 0x5e, 0x26,
	0x50, 0x7d, 0x01, 0xed, 0x5e, 0x85, 0xa5, 0xa2, 0x07, 0x6f, 0xeb, 0x30, 0xca, 0xc2, 0x26, 0x53,
	0xa6, 0x69, 0xb1, 0xaf, 0xe4, 0x53, 0x4a, 0x93, 0x16, 0x44, 0xf0, 0x94, 0xfc, 0xfc, 0x1f, 0xdf,
	0x9f, 0xe8, 0xdb, 0xa1, 0x92, 0xfd, 0x08, 0xc8, 0xce, 0x22, 0x0b, 0xd9, 0x56, 0xf8, 0x6b, 0xb4,
	0x2e, 0xb4, 0xc2, 0x81, 0x2d, 0xb2, 0xff, 0x38, 0x7c, 0x42, 0x0c, 0xb9, 0x59, 0x50, 0x21, 0x88,
	0x18, 0xcc, 0xd8, 0xd5, 0x5f, 0x09, 0xb2, 0x73, 0xec, 0x6d, 0x2c, 0x91, 0x8f, 0xb6, 0xef, 0x48,
	0x92, 0xa4, 0xbc, 0xd1, 0xeb, 0x07, 0x44, 0xee, 0x7d, 0xc8, 0x2f, 0xf5, 0xaa, 0x23, 0x14, 0x0a,
	0x46, 0x15, 0x6a, 0x93, 0x35, 0x73, 0x2a, 0xb7, 0x3a, 0x7b, 0x26, 0x19, 0x48, 0xea, 0xb8, 0x9e,
	0xdf, 0xeb, 0xb3, 0x49, 0x60, 0x54, 0x71, 0xba, 0x39, 0xaf, 0x0a, 0xf8, 0xf5, 0x8c, 0x34, 0x04,
	0xcb, 0x10, 0xbc, 0x74, 0x24, 0x4d, 0xa4, 0xcb, 0x0b, 0x9d, 0xd5, 0xc7, 0x03, 0xf9, 0x90, 0x2c,
	0x8f, 0x5b, 0x7d, 0x8d, 0xec, 0xe0, 0xbf, 0x3f, 0xec, 0xd4, 0xdb, 0x7a, 0xbe, 0x7e, 0x9f, 0x98,
	0xd7, 0xaa, 0xd9, 0x8f, 0x50, 0x4f, 0xd4, 0x2e, 0x52, 0x4f, 0xa9, 0x7a, 0x5f, 0xc5, 0xe3, 0xee,
	0x27, 0x00, 0x00, 0xff, 0xff, 0x54, 0x23, 0xa0, 0xae, 0x37, 0x01, 0x00, 0x00,
}
