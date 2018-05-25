package socks

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import net "v2ray.com/core/common/net"
import protocol "v2ray.com/core/common/protocol"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AuthType int32

const (
	AuthType_NO_AUTH  AuthType = 0
	AuthType_PASSWORD AuthType = 1
)

var AuthType_name = map[int32]string{
	0: "NO_AUTH",
	1: "PASSWORD",
}
var AuthType_value = map[string]int32{
	"NO_AUTH":  0,
	"PASSWORD": 1,
}

func (x AuthType) String() string {
	return proto.EnumName(AuthType_name, int32(x))
}
func (AuthType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_45beb3f2ac36b2a3, []int{0}
}

type Account struct {
	Username             string   `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_45beb3f2ac36b2a3, []int{0}
}
func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (dst *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(dst, src)
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

type ServerConfig struct {
	AuthType             AuthType          `protobuf:"varint,1,opt,name=auth_type,json=authType,enum=v2ray.core.proxy.socks.AuthType" json:"auth_type,omitempty"`
	Accounts             map[string]string `protobuf:"bytes,2,rep,name=accounts" json:"accounts,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Address              *net.IPOrDomain   `protobuf:"bytes,3,opt,name=address" json:"address,omitempty"`
	UdpEnabled           bool              `protobuf:"varint,4,opt,name=udp_enabled,json=udpEnabled" json:"udp_enabled,omitempty"`
	Timeout              uint32            `protobuf:"varint,5,opt,name=timeout" json:"timeout,omitempty"` // Deprecated: Do not use.
	UserLevel            uint32            `protobuf:"varint,6,opt,name=user_level,json=userLevel" json:"user_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_45beb3f2ac36b2a3, []int{1}
}
func (m *ServerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerConfig.Unmarshal(m, b)
}
func (m *ServerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerConfig.Marshal(b, m, deterministic)
}
func (dst *ServerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerConfig.Merge(dst, src)
}
func (m *ServerConfig) XXX_Size() int {
	return xxx_messageInfo_ServerConfig.Size(m)
}
func (m *ServerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ServerConfig proto.InternalMessageInfo

func (m *ServerConfig) GetAuthType() AuthType {
	if m != nil {
		return m.AuthType
	}
	return AuthType_NO_AUTH
}

func (m *ServerConfig) GetAccounts() map[string]string {
	if m != nil {
		return m.Accounts
	}
	return nil
}

func (m *ServerConfig) GetAddress() *net.IPOrDomain {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *ServerConfig) GetUdpEnabled() bool {
	if m != nil {
		return m.UdpEnabled
	}
	return false
}

// Deprecated: Do not use.
func (m *ServerConfig) GetTimeout() uint32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *ServerConfig) GetUserLevel() uint32 {
	if m != nil {
		return m.UserLevel
	}
	return 0
}

type ClientConfig struct {
	Server               []*protocol.ServerEndpoint `protobuf:"bytes,1,rep,name=server" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ClientConfig) Reset()         { *m = ClientConfig{} }
func (m *ClientConfig) String() string { return proto.CompactTextString(m) }
func (*ClientConfig) ProtoMessage()    {}
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_45beb3f2ac36b2a3, []int{2}
}
func (m *ClientConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientConfig.Unmarshal(m, b)
}
func (m *ClientConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientConfig.Marshal(b, m, deterministic)
}
func (dst *ClientConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientConfig.Merge(dst, src)
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
	proto.RegisterType((*Account)(nil), "v2ray.core.proxy.socks.Account")
	proto.RegisterType((*ServerConfig)(nil), "v2ray.core.proxy.socks.ServerConfig")
	proto.RegisterMapType((map[string]string)(nil), "v2ray.core.proxy.socks.ServerConfig.AccountsEntry")
	proto.RegisterType((*ClientConfig)(nil), "v2ray.core.proxy.socks.ClientConfig")
	proto.RegisterEnum("v2ray.core.proxy.socks.AuthType", AuthType_name, AuthType_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/socks/config.proto", fileDescriptor_config_45beb3f2ac36b2a3)
}

var fileDescriptor_config_45beb3f2ac36b2a3 = []byte{
	// 470 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0x75, 0xb2, 0xb6, 0x4d, 0x6f, 0xbb, 0x52, 0x06, 0x59, 0x42, 0x51, 0x8c, 0x05, 0xb1, 0xec,
	0xc3, 0x44, 0xe2, 0x8b, 0xb8, 0x28, 0xb4, 0xdd, 0x82, 0x82, 0x6c, 0xcb, 0x74, 0x55, 0xf0, 0x25,
	0xcc, 0x4e, 0x46, 0x37, 0x6c, 0x32, 0x13, 0x66, 0x26, 0xd5, 0xfc, 0x25, 0xff, 0x9f, 0xef, 0x92,
	0xaf, 0x65, 0x95, 0xee, 0xdb, 0xfd, 0x38, 0xf7, 0xcc, 0x3d, 0xe7, 0x0e, 0xbc, 0xdc, 0x87, 0x9a,
	0x95, 0x84, 0xab, 0x2c, 0xe0, 0x4a, 0x8b, 0x20, 0xd7, 0xea, 0x57, 0x19, 0x18, 0xc5, 0x6f, 0x4c,
	0xc0, 0x95, 0xfc, 0x9e, 0xfc, 0x20, 0xb9, 0x56, 0x56, 0xe1, 0x93, 0x0e, 0xa8, 0x05, 0xa9, 0x41,
	0xa4, 0x06, 0x4d, 0xff, 0x27, 0xe0, 0x2a, 0xcb, 0x94, 0x0c, 0xa4, 0xb0, 0x01, 0x8b, 0x63, 0x2d,
	0x8c, 0x69, 0x08, 0xa6, 0xaf, 0x0e, 0x03, 0xeb, 0x26, 0x57, 0x69, 0x60, 0x84, 0xde, 0x0b, 0x1d,
	0x99, 0x5c, 0xf0, 0x66, 0x62, 0xb6, 0x80, 0xc1, 0x82, 0x73, 0x55, 0x48, 0x8b, 0xa7, 0xe0, 0x16,
	0x46, 0x68, 0xc9, 0x32, 0xe1, 0x21, 0x1f, 0xcd, 0x87, 0xf4, 0x36, 0xaf, 0x7a, 0x39, 0x33, 0xe6,
	0xa7, 0xd2, 0xb1, 0xe7, 0x34, 0xbd, 0x2e, 0x9f, 0xfd, 0x71, 0x60, 0xbc, 0xab, 0x89, 0x57, 0xb5,
	0x18, 0xfc, 0x0e, 0x86, 0xac, 0xb0, 0xd7, 0x91, 0x2d, 0xf3, 0x86, 0xe9, 0x51, 0xe8, 0x93, 0xc3,
	0xd2, 0xc8, 0xa2, 0xb0, 0xd7, 0x97, 0x65, 0x2e, 0xa8, 0xcb, 0xda, 0x08, 0x5f, 0x80, 0xcb, 0x9a,
	0x95, 0x8c, 0xe7, 0xf8, 0x47, 0xf3, 0x51, 0x18, 0xde, 0x37, 0x7d, 0xf7, 0x59, 0xd2, 0xea, 0x30,
	0x6b, 0x69, 0x75, 0x49, 0x6f, 0x39, 0xf0, 0x19, 0x0c, 0x5a, 0x97, 0xbc, 0x23, 0x1f, 0xcd, 0x47,
	0xe1, 0xf3, 0xbb, 0x74, 0x8d, 0x45, 0x44, 0x0a, 0x4b, 0x3e, 0x6e, 0x37, 0xfa, 0x5c, 0x65, 0x2c,
	0x91, 0xb4, 0x9b, 0xc0, 0xcf, 0x60, 0x54, 0xc4, 0x79, 0x24, 0x24, 0xbb, 0x4a, 0x45, 0xec, 0x3d,
	0xf4, 0xd1, 0xdc, 0xa5, 0x50, 0xc4, 0xf9, 0xba, 0xa9, 0xe0, 0x27, 0x30, 0xb0, 0x49, 0x26, 0x54,
	0x61, 0xbd, 0x9e, 0x8f, 0xe6, 0xc7, 0x4b, 0xc7, 0x43, 0xb4, 0x2b, 0xe1, 0xa7, 0x00, 0x95, 0x87,
	0x51, 0x2a, 0xf6, 0x22, 0xf5, 0xfa, 0x15, 0x80, 0x0e, 0xab, 0xca, 0xa7, 0xaa, 0x30, 0x3d, 0x83,
	0xe3, 0x7f, 0xb6, 0xc6, 0x13, 0x38, 0xba, 0x11, 0x65, 0x6b, 0x7f, 0x15, 0xe2, 0xc7, 0xd0, 0xdb,
	0xb3, 0xb4, 0x10, 0xad, 0xed, 0x4d, 0xf2, 0xd6, 0x79, 0x83, 0x66, 0x14, 0xc6, 0xab, 0x34, 0x11,
	0xd2, 0xb6, 0xb6, 0x2f, 0xa1, 0xdf, 0xdc, 0xd7, 0x43, 0xb5, 0x6b, 0xa7, 0x07, 0x64, 0x76, 0x3f,
	0xa1, 0x75, 0x6e, 0x2d, 0xe3, 0x5c, 0x25, 0xd2, 0xd2, 0x76, 0xf2, 0xf4, 0x05, 0xb8, 0xdd, 0x45,
	0xf0, 0x08, 0x06, 0x17, 0x9b, 0x68, 0xf1, 0xf9, 0xf2, 0xc3, 0xe4, 0x01, 0x1e, 0x83, 0xbb, 0x5d,
	0xec, 0x76, 0x5f, 0x37, 0xf4, 0x7c, 0x82, 0x96, 0xef, 0x61, 0xca, 0x55, 0x76, 0xcf, 0x55, 0xb6,
	0xe8, 0x5b, 0xaf, 0x0e, 0x7e, 0x3b, 0x27, 0x5f, 0x42, 0xca, 0x4a, 0xb2, 0xaa, 0x10, 0xdb, 0x1a,
	0xb1, 0xab, 0x1a, 0x57, 0xfd, 0x7a, 0x8f, 0xd7, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x23, 0xac,
	0x72, 0x71, 0x1a, 0x03, 0x00, 0x00,
}
