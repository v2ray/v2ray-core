package internet

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import serial "v2ray.com/core/common/serial"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TransportProtocol int32

const (
	TransportProtocol_TCP          TransportProtocol = 0
	TransportProtocol_UDP          TransportProtocol = 1
	TransportProtocol_MKCP         TransportProtocol = 2
	TransportProtocol_WebSocket    TransportProtocol = 3
	TransportProtocol_HTTP         TransportProtocol = 4
	TransportProtocol_DomainSocket TransportProtocol = 5
)

var TransportProtocol_name = map[int32]string{
	0: "TCP",
	1: "UDP",
	2: "MKCP",
	3: "WebSocket",
	4: "HTTP",
	5: "DomainSocket",
}
var TransportProtocol_value = map[string]int32{
	"TCP":          0,
	"UDP":          1,
	"MKCP":         2,
	"WebSocket":    3,
	"HTTP":         4,
	"DomainSocket": 5,
}

func (x TransportProtocol) String() string {
	return proto.EnumName(TransportProtocol_name, int32(x))
}
func (TransportProtocol) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_807c5df32db81c88, []int{0}
}

type TransportConfig struct {
	// Type of network that this settings supports.
	// Deprecated. Use the string form below.
	Protocol TransportProtocol `protobuf:"varint,1,opt,name=protocol,proto3,enum=v2ray.core.transport.internet.TransportProtocol" json:"protocol,omitempty"`
	// Type of network that this settings supports.
	ProtocolName string `protobuf:"bytes,3,opt,name=protocol_name,json=protocolName,proto3" json:"protocol_name,omitempty"`
	// Specific settings. Must be of the transports.
	Settings             *serial.TypedMessage `protobuf:"bytes,2,opt,name=settings,proto3" json:"settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TransportConfig) Reset()         { *m = TransportConfig{} }
func (m *TransportConfig) String() string { return proto.CompactTextString(m) }
func (*TransportConfig) ProtoMessage()    {}
func (*TransportConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_807c5df32db81c88, []int{0}
}
func (m *TransportConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransportConfig.Unmarshal(m, b)
}
func (m *TransportConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransportConfig.Marshal(b, m, deterministic)
}
func (dst *TransportConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransportConfig.Merge(dst, src)
}
func (m *TransportConfig) XXX_Size() int {
	return xxx_messageInfo_TransportConfig.Size(m)
}
func (m *TransportConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_TransportConfig.DiscardUnknown(m)
}

var xxx_messageInfo_TransportConfig proto.InternalMessageInfo

func (m *TransportConfig) GetProtocol() TransportProtocol {
	if m != nil {
		return m.Protocol
	}
	return TransportProtocol_TCP
}

func (m *TransportConfig) GetProtocolName() string {
	if m != nil {
		return m.ProtocolName
	}
	return ""
}

func (m *TransportConfig) GetSettings() *serial.TypedMessage {
	if m != nil {
		return m.Settings
	}
	return nil
}

type StreamConfig struct {
	// Effective network. Deprecated. Use the string form below.
	Protocol TransportProtocol `protobuf:"varint,1,opt,name=protocol,proto3,enum=v2ray.core.transport.internet.TransportProtocol" json:"protocol,omitempty"` // Deprecated: Do not use.
	// Effective network.
	ProtocolName      string             `protobuf:"bytes,5,opt,name=protocol_name,json=protocolName,proto3" json:"protocol_name,omitempty"`
	TransportSettings []*TransportConfig `protobuf:"bytes,2,rep,name=transport_settings,json=transportSettings,proto3" json:"transport_settings,omitempty"`
	// Type of security. Must be a message name of the settings proto.
	SecurityType string `protobuf:"bytes,3,opt,name=security_type,json=securityType,proto3" json:"security_type,omitempty"`
	// Settings for transport security. For now the only choice is TLS.
	SecuritySettings     []*serial.TypedMessage `protobuf:"bytes,4,rep,name=security_settings,json=securitySettings,proto3" json:"security_settings,omitempty"`
	SocketSettings       *SocketConfig          `protobuf:"bytes,6,opt,name=socket_settings,json=socketSettings,proto3" json:"socket_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *StreamConfig) Reset()         { *m = StreamConfig{} }
func (m *StreamConfig) String() string { return proto.CompactTextString(m) }
func (*StreamConfig) ProtoMessage()    {}
func (*StreamConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_807c5df32db81c88, []int{1}
}
func (m *StreamConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamConfig.Unmarshal(m, b)
}
func (m *StreamConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamConfig.Marshal(b, m, deterministic)
}
func (dst *StreamConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamConfig.Merge(dst, src)
}
func (m *StreamConfig) XXX_Size() int {
	return xxx_messageInfo_StreamConfig.Size(m)
}
func (m *StreamConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamConfig.DiscardUnknown(m)
}

var xxx_messageInfo_StreamConfig proto.InternalMessageInfo

// Deprecated: Do not use.
func (m *StreamConfig) GetProtocol() TransportProtocol {
	if m != nil {
		return m.Protocol
	}
	return TransportProtocol_TCP
}

func (m *StreamConfig) GetProtocolName() string {
	if m != nil {
		return m.ProtocolName
	}
	return ""
}

func (m *StreamConfig) GetTransportSettings() []*TransportConfig {
	if m != nil {
		return m.TransportSettings
	}
	return nil
}

func (m *StreamConfig) GetSecurityType() string {
	if m != nil {
		return m.SecurityType
	}
	return ""
}

func (m *StreamConfig) GetSecuritySettings() []*serial.TypedMessage {
	if m != nil {
		return m.SecuritySettings
	}
	return nil
}

func (m *StreamConfig) GetSocketSettings() *SocketConfig {
	if m != nil {
		return m.SocketSettings
	}
	return nil
}

type ProxyConfig struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyConfig) Reset()         { *m = ProxyConfig{} }
func (m *ProxyConfig) String() string { return proto.CompactTextString(m) }
func (*ProxyConfig) ProtoMessage()    {}
func (*ProxyConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_807c5df32db81c88, []int{2}
}
func (m *ProxyConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyConfig.Unmarshal(m, b)
}
func (m *ProxyConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyConfig.Marshal(b, m, deterministic)
}
func (dst *ProxyConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyConfig.Merge(dst, src)
}
func (m *ProxyConfig) XXX_Size() int {
	return xxx_messageInfo_ProxyConfig.Size(m)
}
func (m *ProxyConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyConfig proto.InternalMessageInfo

func (m *ProxyConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type SocketConfig struct {
	Mark                 int32    `protobuf:"varint,1,opt,name=mark,proto3" json:"mark,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SocketConfig) Reset()         { *m = SocketConfig{} }
func (m *SocketConfig) String() string { return proto.CompactTextString(m) }
func (*SocketConfig) ProtoMessage()    {}
func (*SocketConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_807c5df32db81c88, []int{3}
}
func (m *SocketConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SocketConfig.Unmarshal(m, b)
}
func (m *SocketConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SocketConfig.Marshal(b, m, deterministic)
}
func (dst *SocketConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SocketConfig.Merge(dst, src)
}
func (m *SocketConfig) XXX_Size() int {
	return xxx_messageInfo_SocketConfig.Size(m)
}
func (m *SocketConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SocketConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SocketConfig proto.InternalMessageInfo

func (m *SocketConfig) GetMark() int32 {
	if m != nil {
		return m.Mark
	}
	return 0
}

func init() {
	proto.RegisterType((*TransportConfig)(nil), "v2ray.core.transport.internet.TransportConfig")
	proto.RegisterType((*StreamConfig)(nil), "v2ray.core.transport.internet.StreamConfig")
	proto.RegisterType((*ProxyConfig)(nil), "v2ray.core.transport.internet.ProxyConfig")
	proto.RegisterType((*SocketConfig)(nil), "v2ray.core.transport.internet.SocketConfig")
	proto.RegisterEnum("v2ray.core.transport.internet.TransportProtocol", TransportProtocol_name, TransportProtocol_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_config_807c5df32db81c88)
}

var fileDescriptor_config_807c5df32db81c88 = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xe1, 0x8a, 0x13, 0x31,
	0x10, 0xc7, 0xdd, 0x6e, 0x7b, 0xb6, 0xd3, 0xde, 0xdd, 0x36, 0x9f, 0x8a, 0x70, 0x58, 0x57, 0x90,
	0xa2, 0x90, 0x3d, 0xd6, 0x37, 0x68, 0xef, 0x83, 0xa2, 0xa7, 0xcb, 0x76, 0x55, 0x38, 0x90, 0x92,
	0x8b, 0xb1, 0x2c, 0x77, 0x49, 0x4a, 0x12, 0xc5, 0x7d, 0x25, 0x3f, 0xfb, 0x10, 0x3e, 0x96, 0x24,
	0xbb, 0x09, 0x45, 0xa5, 0x1e, 0xf8, 0x6d, 0xc8, 0xfc, 0xe7, 0x3f, 0xf3, 0x9b, 0x09, 0xe0, 0xaf,
	0xb9, 0x22, 0x0d, 0xa6, 0x92, 0x67, 0x54, 0x2a, 0x96, 0x19, 0x45, 0x84, 0xde, 0x49, 0x65, 0xb2,
	0x5a, 0x18, 0xa6, 0x04, 0x33, 0x19, 0x95, 0xe2, 0x73, 0xbd, 0xc5, 0x3b, 0x25, 0x8d, 0x44, 0x67,
	0x5e, 0xaf, 0x18, 0x0e, 0x5a, 0xec, 0xb5, 0x0f, 0xce, 0x7f, 0xb3, 0xa3, 0x92, 0x73, 0x29, 0x32,
	0xcd, 0x54, 0x4d, 0x6e, 0x33, 0xd3, 0xec, 0xd8, 0xa7, 0x0d, 0x67, 0x5a, 0x93, 0x2d, 0x6b, 0x0d,
	0xd3, 0x9f, 0x11, 0x9c, 0x56, 0xde, 0x68, 0xe5, 0x5a, 0xa1, 0xd7, 0x30, 0x74, 0x49, 0x2a, 0x6f,
	0x67, 0xd1, 0x3c, 0x5a, 0x9c, 0xe4, 0xe7, 0xf8, 0x60, 0x5f, 0x1c, 0x1c, 0x8a, 0xae, 0xae, 0x0c,
	0x0e, 0xe8, 0x31, 0x1c, 0xfb, 0x78, 0x23, 0x08, 0x67, 0xb3, 0x78, 0x1e, 0x2d, 0x46, 0xe5, 0xc4,
	0x3f, 0xbe, 0x21, 0x9c, 0xa1, 0x25, 0x0c, 0x35, 0x33, 0xa6, 0x16, 0x5b, 0x3d, 0xeb, 0xcd, 0xa3,
	0xc5, 0x38, 0x7f, 0xb2, 0xdf, 0xb2, 0xe5, 0xc0, 0x2d, 0x07, 0xae, 0x2c, 0xc7, 0x65, 0x8b, 0x51,
	0x86, 0xba, 0xf4, 0x47, 0x0c, 0x93, 0xb5, 0x51, 0x8c, 0xf0, 0x8e, 0xa3, 0xf8, 0x7f, 0x8e, 0x65,
	0x6f, 0x16, 0x1d, 0x62, 0x19, 0xfc, 0x85, 0xe5, 0x23, 0xa0, 0x60, 0xbd, 0xd9, 0xa3, 0x8a, 0x17,
	0xe3, 0x1c, 0xdf, 0x75, 0x80, 0x16, 0xa1, 0x9c, 0x06, 0xcd, 0xba, 0x33, 0xb2, 0x33, 0x68, 0x46,
	0xbf, 0xa8, 0xda, 0x34, 0x1b, 0x7b, 0x51, 0xbf, 0x4f, 0xff, 0x68, 0xb7, 0x83, 0xd6, 0x30, 0x0d,
	0xa2, 0x30, 0x42, 0xdf, 0x8d, 0x70, 0xd7, 0xc5, 0x26, 0xde, 0x20, 0x74, 0xae, 0xe0, 0x54, 0x4b,
	0x7a, 0xc3, 0xf6, 0xa8, 0x8e, 0xdc, 0xad, 0x9e, 0xfd, 0x83, 0x6a, 0xed, 0xaa, 0x3a, 0xa4, 0x93,
	0xd6, 0xc3, 0xbb, 0xa6, 0x0f, 0x61, 0x5c, 0x28, 0xf9, 0xad, 0xe9, 0x8e, 0x96, 0x40, 0x6c, 0xc8,
	0xd6, 0xdd, 0x6b, 0x54, 0xda, 0x30, 0x4d, 0x61, 0xb2, 0x6f, 0x80, 0x10, 0xf4, 0x39, 0x51, 0x37,
	0x4e, 0x32, 0x28, 0x5d, 0xfc, 0xf4, 0x0a, 0xa6, 0x7f, 0xdc, 0x0e, 0xdd, 0x87, 0xb8, 0x5a, 0x15,
	0xc9, 0x3d, 0x1b, 0xbc, 0xbb, 0x28, 0x92, 0x08, 0x0d, 0xa1, 0x7f, 0xf9, 0x6a, 0x55, 0x24, 0x3d,
	0x74, 0x0c, 0xa3, 0x0f, 0xec, 0xba, 0xf5, 0x4d, 0x62, 0x9b, 0x78, 0x51, 0x55, 0x45, 0xd2, 0x47,
	0x09, 0x4c, 0x2e, 0x24, 0x27, 0xb5, 0xe8, 0x72, 0x83, 0xe5, 0x5b, 0x78, 0x44, 0x25, 0x3f, 0x8c,
	0x58, 0x44, 0x57, 0x43, 0x1f, 0x7f, 0xef, 0x9d, 0xbd, 0xcf, 0x4b, 0xd2, 0xe0, 0x95, 0xd5, 0x86,
	0xb1, 0xf0, 0xcb, 0x2e, 0x7f, 0x7d, 0xe4, 0xbe, 0xcb, 0xf3, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xba, 0x0d, 0x33, 0xbb, 0xfd, 0x03, 0x00, 0x00,
}
