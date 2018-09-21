package shadowsocks

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	net "v2ray.com/core/common/net"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CipherType int32

const (
	CipherType_UNKNOWN           CipherType = 0
	CipherType_AES_128_CFB       CipherType = 1
	CipherType_AES_256_CFB       CipherType = 2
	CipherType_CHACHA20          CipherType = 3
	CipherType_CHACHA20_IETF     CipherType = 4
	CipherType_AES_128_GCM       CipherType = 5
	CipherType_AES_256_GCM       CipherType = 6
	CipherType_CHACHA20_POLY1305 CipherType = 7
	CipherType_NONE              CipherType = 8
)

var CipherType_name = map[int32]string{
	0: "UNKNOWN",
	1: "AES_128_CFB",
	2: "AES_256_CFB",
	3: "CHACHA20",
	4: "CHACHA20_IETF",
	5: "AES_128_GCM",
	6: "AES_256_GCM",
	7: "CHACHA20_POLY1305",
	8: "NONE",
}

var CipherType_value = map[string]int32{
	"UNKNOWN":           0,
	"AES_128_CFB":       1,
	"AES_256_CFB":       2,
	"CHACHA20":          3,
	"CHACHA20_IETF":     4,
	"AES_128_GCM":       5,
	"AES_256_GCM":       6,
	"CHACHA20_POLY1305": 7,
	"NONE":              8,
}

func (x CipherType) String() string {
	return proto.EnumName(CipherType_name, int32(x))
}

func (CipherType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8d089a30c2106007, []int{0}
}

type Account_OneTimeAuth int32

const (
	Account_Auto     Account_OneTimeAuth = 0
	Account_Disabled Account_OneTimeAuth = 1
	Account_Enabled  Account_OneTimeAuth = 2
)

var Account_OneTimeAuth_name = map[int32]string{
	0: "Auto",
	1: "Disabled",
	2: "Enabled",
}

var Account_OneTimeAuth_value = map[string]int32{
	"Auto":     0,
	"Disabled": 1,
	"Enabled":  2,
}

func (x Account_OneTimeAuth) String() string {
	return proto.EnumName(Account_OneTimeAuth_name, int32(x))
}

func (Account_OneTimeAuth) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8d089a30c2106007, []int{0, 0}
}

type Account struct {
	Password             string              `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
	CipherType           CipherType          `protobuf:"varint,2,opt,name=cipher_type,json=cipherType,proto3,enum=v2ray.core.proxy.shadowsocks.CipherType" json:"cipher_type,omitempty"`
	Ota                  Account_OneTimeAuth `protobuf:"varint,3,opt,name=ota,proto3,enum=v2ray.core.proxy.shadowsocks.Account_OneTimeAuth" json:"ota,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d089a30c2106007, []int{0}
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

func (m *Account) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *Account) GetCipherType() CipherType {
	if m != nil {
		return m.CipherType
	}
	return CipherType_UNKNOWN
}

func (m *Account) GetOta() Account_OneTimeAuth {
	if m != nil {
		return m.Ota
	}
	return Account_Auto
}

type ServerConfig struct {
	// UdpEnabled specified whether or not to enable UDP for Shadowsocks.
	// Deprecated. Use 'network' field.
	UdpEnabled           bool           `protobuf:"varint,1,opt,name=udp_enabled,json=udpEnabled,proto3" json:"udp_enabled,omitempty"` // Deprecated: Do not use.
	User                 *protocol.User `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Network              []net.Network  `protobuf:"varint,3,rep,packed,name=network,proto3,enum=v2ray.core.common.net.Network" json:"network,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ServerConfig) Reset()         { *m = ServerConfig{} }
func (m *ServerConfig) String() string { return proto.CompactTextString(m) }
func (*ServerConfig) ProtoMessage()    {}
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d089a30c2106007, []int{1}
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
func (m *ServerConfig) GetUdpEnabled() bool {
	if m != nil {
		return m.UdpEnabled
	}
	return false
}

func (m *ServerConfig) GetUser() *protocol.User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *ServerConfig) GetNetwork() []net.Network {
	if m != nil {
		return m.Network
	}
	return nil
}

type ClientConfig struct {
	Server               []*protocol.ServerEndpoint `protobuf:"bytes,1,rep,name=server,proto3" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ClientConfig) Reset()         { *m = ClientConfig{} }
func (m *ClientConfig) String() string { return proto.CompactTextString(m) }
func (*ClientConfig) ProtoMessage()    {}
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d089a30c2106007, []int{2}
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
	proto.RegisterType((*Account)(nil), "v2ray.core.proxy.shadowsocks.Account")
	proto.RegisterType((*ServerConfig)(nil), "v2ray.core.proxy.shadowsocks.ServerConfig")
	proto.RegisterType((*ClientConfig)(nil), "v2ray.core.proxy.shadowsocks.ClientConfig")
	proto.RegisterEnum("v2ray.core.proxy.shadowsocks.CipherType", CipherType_name, CipherType_value)
	proto.RegisterEnum("v2ray.core.proxy.shadowsocks.Account_OneTimeAuth", Account_OneTimeAuth_name, Account_OneTimeAuth_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/shadowsocks/config.proto", fileDescriptor_8d089a30c2106007)
}

var fileDescriptor_8d089a30c2106007 = []byte{
	// 522 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x6e, 0xd3, 0x4e,
	0x10, 0xc6, 0xbb, 0x71, 0xff, 0x4d, 0xfe, 0xe3, 0x50, 0xdc, 0x95, 0x90, 0xac, 0xa8, 0x42, 0x56,
	0x38, 0x10, 0x2a, 0xb1, 0x4e, 0x5c, 0x8a, 0x7a, 0x75, 0x4c, 0x4a, 0x2b, 0xc0, 0x89, 0x9c, 0x14,
	0x04, 0x17, 0xcb, 0x5d, 0x2f, 0xc4, 0x6a, 0xe2, 0xb5, 0xd6, 0x76, 0x43, 0x9e, 0x86, 0x03, 0x37,
	0xde, 0x8c, 0xb7, 0x40, 0x5e, 0x3b, 0xa9, 0x85, 0xaa, 0x70, 0x88, 0x94, 0x99, 0xfd, 0x7d, 0x9f,
	0x66, 0xbe, 0x31, 0xbc, 0xbc, 0xb3, 0x44, 0xb0, 0x26, 0x94, 0x2f, 0x4d, 0xca, 0x05, 0x33, 0x13,
	0xc1, 0xbf, 0xaf, 0xcd, 0x74, 0x1e, 0x84, 0x7c, 0x95, 0x72, 0x7a, 0x9b, 0x9a, 0x94, 0xc7, 0x5f,
	0xa3, 0x6f, 0x24, 0x11, 0x3c, 0xe3, 0xf8, 0x78, 0x83, 0x0b, 0x46, 0x24, 0x4a, 0x6a, 0x68, 0xe7,
	0xf9, 0x5f, 0x66, 0x94, 0x2f, 0x97, 0x3c, 0x36, 0x63, 0x96, 0x15, 0xbf, 0x15, 0x17, 0xb7, 0xa5,
	0x4d, 0xe7, 0xc5, 0xc3, 0xa0, 0x7c, 0xa4, 0x7c, 0x61, 0xe6, 0x29, 0x13, 0x15, 0xda, 0xff, 0x07,
	0x9a, 0x32, 0x71, 0xc7, 0x84, 0x9f, 0x26, 0x8c, 0x96, 0x8a, 0xee, 0x6f, 0x04, 0x4d, 0x9b, 0x52,
	0x9e, 0xc7, 0x19, 0xee, 0x40, 0x2b, 0x09, 0xd2, 0x74, 0xc5, 0x45, 0xa8, 0x23, 0x03, 0xf5, 0xfe,
	0xf7, 0xb6, 0x35, 0xbe, 0x02, 0x95, 0x46, 0xc9, 0x9c, 0x09, 0x3f, 0x5b, 0x27, 0x4c, 0x6f, 0x18,
	0xa8, 0x77, 0x68, 0xf5, 0xc8, 0xae, 0x0d, 0x89, 0x23, 0x05, 0xb3, 0x75, 0xc2, 0x3c, 0xa0, 0xdb,
	0xff, 0xd8, 0x01, 0x85, 0x67, 0x81, 0xae, 0x48, 0x8b, 0xc1, 0x6e, 0x8b, 0x6a, 0x34, 0x32, 0x8e,
	0xd9, 0x2c, 0x5a, 0x32, 0x3b, 0xcf, 0xe6, 0x5e, 0xa1, 0xee, 0x5a, 0xa0, 0xd6, 0x7a, 0xb8, 0x05,
	0xfb, 0x76, 0x9e, 0x71, 0x6d, 0x0f, 0xb7, 0xa1, 0xf5, 0x26, 0x4a, 0x83, 0x9b, 0x05, 0x0b, 0x35,
	0x84, 0x55, 0x68, 0x8e, 0xe2, 0xb2, 0x68, 0x74, 0x7f, 0x22, 0x68, 0x4f, 0x65, 0x02, 0x8e, 0x3c,
	0x13, 0x7e, 0x06, 0x6a, 0x1e, 0x26, 0x3e, 0x2b, 0x09, 0xb9, 0x73, 0x6b, 0xd8, 0xd0, 0x91, 0x07,
	0x79, 0x98, 0x54, 0x3a, 0xfc, 0x0a, 0xf6, 0x8b, 0x84, 0xe5, 0xca, 0xaa, 0x65, 0xd4, 0xe7, 0x2d,
	0xe3, 0x25, 0x9b, 0x78, 0xc9, 0x75, 0xca, 0x84, 0x27, 0x69, 0x7c, 0x0e, 0xcd, 0xea, 0x8a, 0xba,
	0x62, 0x28, 0xbd, 0x43, 0xeb, 0xe9, 0x03, 0xc2, 0x98, 0x65, 0xc4, 0x2d, 0x29, 0x6f, 0x83, 0x77,
	0x3d, 0x68, 0x3b, 0x8b, 0x88, 0xc5, 0x59, 0x35, 0xe4, 0x10, 0x0e, 0xca, 0xb3, 0xe9, 0xc8, 0x50,
	0x7a, 0xaa, 0x75, 0xb2, 0x6b, 0x82, 0x72, 0xbd, 0x51, 0x1c, 0x26, 0x3c, 0x8a, 0x33, 0xaf, 0x52,
	0x9e, 0xfc, 0x40, 0x00, 0xf7, 0xd7, 0x28, 0x52, 0xb9, 0x76, 0xdf, 0xb9, 0xe3, 0x4f, 0xae, 0xb6,
	0x87, 0x1f, 0x83, 0x6a, 0x8f, 0xa6, 0xfe, 0xc0, 0x3a, 0xf7, 0x9d, 0x8b, 0xa1, 0x86, 0x36, 0x0d,
	0xeb, 0xec, 0xb5, 0x6c, 0x34, 0x8a, 0x48, 0x9d, 0x4b, 0xdb, 0xb9, 0xb4, 0xad, 0xbe, 0xa6, 0xe0,
	0x23, 0x78, 0xb4, 0xa9, 0xfc, 0xab, 0xd1, 0xec, 0x42, 0xdb, 0xaf, 0x5b, 0xbc, 0x75, 0x3e, 0x68,
	0xff, 0xd5, 0x2d, 0x8a, 0xc6, 0x01, 0x7e, 0x02, 0x47, 0x5b, 0xd1, 0x64, 0xfc, 0xfe, 0xf3, 0xe0,
	0xb4, 0x7f, 0xa6, 0x35, 0x8b, 0xb3, 0xb9, 0x63, 0x77, 0xa4, 0xb5, 0x86, 0x13, 0x30, 0x28, 0x5f,
	0xee, 0xfc, 0x18, 0x26, 0xe8, 0x8b, 0x5a, 0x2b, 0x7f, 0x35, 0x8e, 0x3f, 0x5a, 0x5e, 0xb0, 0x26,
	0x4e, 0x41, 0x4f, 0x24, 0x3d, 0xbd, 0x7f, 0xbe, 0x39, 0x90, 0xa1, 0x9c, 0xfe, 0x09, 0x00, 0x00,
	0xff, 0xff, 0xdc, 0x7e, 0x6b, 0x61, 0xb5, 0x03, 0x00, 0x00,
}
