package http

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
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_e66c3db3a635d8e4, []int{0}
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

func (m *Account) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Account) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

// Config for HTTP proxy server.
type ServerConfig struct {
	Timeout              uint32            `protobuf:"varint,1,opt,name=timeout,proto3" json:"timeout,omitempty"` // Deprecated: Do not use.
	Accounts             map[string]string `protobuf:"bytes,2,rep,name=accounts,proto3" json:"accounts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	AllowTransparent     bool              `protobuf:"varint,3,opt,name=allow_transparent,json=allowTransparent,proto3" json:"allow_transparent,omitempty"`
	UserLevel            uint32            `protobuf:"varint,4,opt,name=user_level,json=userLevel,proto3" json:"user_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_e66c3db3a635d8e4, []int{1}
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

// Deprecated: Do not use.
func (m *ServerConfig) GetTimeout() uint32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *ServerConfig) GetAccounts() map[string]string {
	if m != nil {
		return m.Accounts
	}
	return nil
}

func (m *ServerConfig) GetAllowTransparent() bool {
	if m != nil {
		return m.AllowTransparent
	}
	return false
}

func (m *ServerConfig) GetUserLevel() uint32 {
	if m != nil {
		return m.UserLevel
	}
	return 0
}

// ClientConfig is the protobuf config for HTTP proxy client.
type ClientConfig struct {
	// Sever is a list of HTTP server addresses.
	Server               []*protocol.ServerEndpoint `protobuf:"bytes,1,rep,name=server,proto3" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ClientConfig) Reset()         { *m = ClientConfig{} }
func (m *ClientConfig) String() string { return proto.CompactTextString(m) }
func (*ClientConfig) ProtoMessage()    {}
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_e66c3db3a635d8e4, []int{2}
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

func (m *ClientConfig) GetServer() []*protocol.ServerEndpoint {
	if m != nil {
		return m.Server
	}
	return nil
}

func init() {
	proto.RegisterType((*Account)(nil), "v2ray.core.proxy.http.Account")
	proto.RegisterType((*ServerConfig)(nil), "v2ray.core.proxy.http.ServerConfig")
	proto.RegisterMapType((map[string]string)(nil), "v2ray.core.proxy.http.ServerConfig.AccountsEntry")
	proto.RegisterType((*ClientConfig)(nil), "v2ray.core.proxy.http.ClientConfig")
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/http/config.proto", fileDescriptor_e66c3db3a635d8e4)
}

var fileDescriptor_e66c3db3a635d8e4 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0x4d, 0x6b, 0xe3, 0x30,
	0x10, 0xc5, 0x4e, 0x36, 0x1f, 0xda, 0x04, 0xb2, 0x62, 0x03, 0xde, 0xb0, 0x0b, 0x21, 0x87, 0x25,
	0xb4, 0x20, 0xb7, 0xe9, 0xa5, 0x34, 0xa7, 0x24, 0x04, 0x7a, 0x68, 0x21, 0xa8, 0xa5, 0x87, 0x5e,
	0x82, 0xaa, 0xa8, 0xad, 0xa9, 0x2d, 0x09, 0x49, 0x76, 0xea, 0x7b, 0x7f, 0x4d, 0x7f, 0x65, 0x91,
	0x6c, 0xa7, 0x69, 0xc9, 0xc9, 0x9e, 0xf7, 0x66, 0x9e, 0xe6, 0xbd, 0x01, 0xff, 0xb3, 0x89, 0x22,
	0x39, 0xa2, 0x22, 0x09, 0xa9, 0x50, 0x2c, 0x94, 0x4a, 0xbc, 0xe6, 0xe1, 0xb3, 0x31, 0x32, 0xa4,
	0x82, 0x3f, 0x46, 0x4f, 0x48, 0x2a, 0x61, 0x04, 0xec, 0x57, 0x7d, 0x8a, 0x21, 0xd7, 0x83, 0x6c,
	0xcf, 0xe0, 0xe4, 0xdb, 0x38, 0x15, 0x49, 0x22, 0x78, 0xe8, 0x66, 0xa8, 0x88, 0x43, 0xcd, 0x54,
	0xc6, 0xd4, 0x5a, 0x4b, 0x46, 0x0b, 0xa1, 0xd1, 0x0c, 0x34, 0x67, 0x94, 0x8a, 0x94, 0x1b, 0x38,
	0x00, 0xad, 0x54, 0x33, 0xc5, 0x49, 0xc2, 0x02, 0x6f, 0xe8, 0x8d, 0xdb, 0x78, 0x57, 0x5b, 0x4e,
	0x12, 0xad, 0xb7, 0x42, 0x6d, 0x02, 0xbf, 0xe0, 0xaa, 0x7a, 0xf4, 0xe6, 0x83, 0xce, 0x8d, 0x13,
	0x5e, 0xb8, 0x15, 0xe1, 0x5f, 0xd0, 0x34, 0x51, 0xc2, 0x44, 0x6a, 0x9c, 0x4e, 0x77, 0xee, 0x07,
	0x1e, 0xae, 0x20, 0x78, 0x0d, 0x5a, 0xa4, 0x78, 0x51, 0x07, 0xfe, 0xb0, 0x36, 0xfe, 0x39, 0x39,
	0x45, 0x07, 0xdd, 0xa0, 0x7d, 0x51, 0x54, 0x6e, 0xa9, 0x97, 0xdc, 0xa8, 0x1c, 0xef, 0x24, 0xe0,
	0x31, 0xf8, 0x45, 0xe2, 0x58, 0x6c, 0xd7, 0x46, 0x11, 0xae, 0x25, 0x51, 0x8c, 0x9b, 0xa0, 0x36,
	0xf4, 0xc6, 0x2d, 0xdc, 0x73, 0xc4, 0xed, 0x27, 0x0e, 0xff, 0x01, 0x60, 0x2d, 0xad, 0x63, 0x96,
	0xb1, 0x38, 0xa8, 0xdb, 0xe5, 0x70, 0xdb, 0x22, 0x57, 0x16, 0x18, 0x4c, 0x41, 0xf7, 0xcb, 0x33,
	0xb0, 0x07, 0x6a, 0x2f, 0x2c, 0x2f, 0xd3, 0xb0, 0xbf, 0xf0, 0x37, 0xf8, 0x91, 0x91, 0x38, 0x65,
	0x65, 0x0a, 0x45, 0x71, 0xe1, 0x9f, 0x7b, 0x23, 0x0c, 0x3a, 0x8b, 0x38, 0x62, 0xdc, 0x94, 0x29,
	0xcc, 0x41, 0xa3, 0x88, 0x3b, 0xf0, 0x9c, 0xcb, 0xa3, 0x7d, 0x97, 0xc5, 0x61, 0x50, 0x75, 0x98,
	0xd2, 0xea, 0x92, 0x6f, 0xa4, 0x88, 0xb8, 0xc1, 0xe5, 0xe4, 0x7c, 0x0a, 0xfe, 0x50, 0x91, 0x1c,
	0x8e, 0x67, 0xe5, 0xdd, 0xd7, 0xed, 0xf7, 0xdd, 0xef, 0xdf, 0x4d, 0x30, 0xc9, 0xd1, 0xc2, 0xf2,
	0x2b, 0xc7, 0x5f, 0x1a, 0x23, 0x1f, 0x1a, 0x4e, 0xfd, 0xec, 0x23, 0x00, 0x00, 0xff, 0xff, 0xac,
	0x7a, 0x67, 0x04, 0x54, 0x02, 0x00, 0x00,
}
