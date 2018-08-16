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
	return fileDescriptor_config_59931ebeb80dc13e, []int{0}
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
	return fileDescriptor_config_59931ebeb80dc13e, []int{0}
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
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *StreamConfig) Reset()         { *m = StreamConfig{} }
func (m *StreamConfig) String() string { return proto.CompactTextString(m) }
func (*StreamConfig) ProtoMessage()    {}
func (*StreamConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_59931ebeb80dc13e, []int{1}
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
	return fileDescriptor_config_59931ebeb80dc13e, []int{2}
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

func init() {
	proto.RegisterType((*TransportConfig)(nil), "v2ray.core.transport.internet.TransportConfig")
	proto.RegisterType((*StreamConfig)(nil), "v2ray.core.transport.internet.StreamConfig")
	proto.RegisterType((*ProxyConfig)(nil), "v2ray.core.transport.internet.ProxyConfig")
	proto.RegisterEnum("v2ray.core.transport.internet.TransportProtocol", TransportProtocol_name, TransportProtocol_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_config_59931ebeb80dc13e)
}

var fileDescriptor_config_59931ebeb80dc13e = []byte{
	// 419 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0xd1, 0x8a, 0x13, 0x31,
	0x14, 0x86, 0xcd, 0x4c, 0x57, 0xdb, 0xd3, 0xae, 0xa6, 0xb9, 0x2a, 0xc2, 0x62, 0xad, 0x20, 0xc5,
	0x8b, 0xcc, 0x32, 0xbe, 0x41, 0xbb, 0x17, 0x8a, 0xae, 0x86, 0xe9, 0xa8, 0xb0, 0x20, 0x25, 0x1b,
	0x63, 0x19, 0xdc, 0x24, 0x25, 0x13, 0xc5, 0x79, 0x25, 0x9f, 0xc4, 0xa7, 0xf0, 0x59, 0x24, 0x33,
	0x93, 0xb0, 0xa8, 0x94, 0x82, 0x77, 0x87, 0x9c, 0x3f, 0xff, 0x39, 0xdf, 0x7f, 0x80, 0x7e, 0xcb,
	0x2d, 0x6f, 0xa8, 0x30, 0x2a, 0x13, 0xc6, 0xca, 0xcc, 0x59, 0xae, 0xeb, 0xbd, 0xb1, 0x2e, 0xab,
	0xb4, 0x93, 0x56, 0x4b, 0x97, 0x09, 0xa3, 0x3f, 0x57, 0x3b, 0xba, 0xb7, 0xc6, 0x19, 0x72, 0x16,
	0xf4, 0x56, 0xd2, 0xa8, 0xa5, 0x41, 0xfb, 0xf0, 0xfc, 0x0f, 0x3b, 0x61, 0x94, 0x32, 0x3a, 0xab,
	0xa5, 0xad, 0xf8, 0x4d, 0xe6, 0x9a, 0xbd, 0xfc, 0xb4, 0x55, 0xb2, 0xae, 0xf9, 0x4e, 0x76, 0x86,
	0x8b, 0x9f, 0x08, 0x1e, 0x94, 0xc1, 0x68, 0xdd, 0x8e, 0x22, 0xaf, 0x61, 0xd8, 0x36, 0x85, 0xb9,
	0x99, 0xa1, 0x39, 0x5a, 0xde, 0xcf, 0xcf, 0xe9, 0xc1, 0xb9, 0x34, 0x3a, 0xb0, 0xfe, 0x5f, 0x11,
	0x1d, 0xc8, 0x13, 0x38, 0x0d, 0xf5, 0x56, 0x73, 0x25, 0x67, 0xe9, 0x1c, 0x2d, 0x47, 0xc5, 0x24,
	0x3c, 0xbe, 0xe1, 0x4a, 0x92, 0x15, 0x0c, 0x6b, 0xe9, 0x5c, 0xa5, 0x77, 0xf5, 0x2c, 0x99, 0xa3,
	0xe5, 0x38, 0x7f, 0x7a, 0x7b, 0x64, 0xc7, 0x41, 0x3b, 0x0e, 0x5a, 0x7a, 0x8e, 0xcb, 0x0e, 0xa3,
	0x88, 0xff, 0x16, 0xbf, 0x12, 0x98, 0x6c, 0x9c, 0x95, 0x5c, 0xf5, 0x1c, 0xec, 0xff, 0x39, 0x56,
	0xc9, 0x0c, 0x1d, 0x62, 0x39, 0xf9, 0x07, 0xcb, 0x47, 0x20, 0xd1, 0x7a, 0x7b, 0x8b, 0x2a, 0x5d,
	0x8e, 0x73, 0x7a, 0xec, 0x02, 0x1d, 0x42, 0x31, 0x8d, 0x9a, 0x4d, 0x6f, 0xe4, 0x77, 0xa8, 0xa5,
	0xf8, 0x6a, 0x2b, 0xd7, 0x6c, 0xfd, 0x45, 0x43, 0x9e, 0xe1, 0xd1, 0xa7, 0x43, 0x36, 0x30, 0x8d,
	0xa2, 0xb8, 0xc2, 0xa0, 0x5d, 0xe1, 0xd8, 0x60, 0x71, 0x30, 0x08, 0x93, 0x17, 0x8f, 0x60, 0xcc,
	0xac, 0xf9, 0xde, 0xf4, 0xf1, 0x62, 0x48, 0x1d, 0xdf, 0xb5, 0xc9, 0x8e, 0x0a, 0x5f, 0x3e, 0xbb,
	0x82, 0xe9, 0x5f, 0x09, 0x92, 0x7b, 0x90, 0x96, 0x6b, 0x86, 0xef, 0xf8, 0xe2, 0xdd, 0x05, 0xc3,
	0x88, 0x0c, 0x61, 0x70, 0xf9, 0x6a, 0xcd, 0x70, 0x42, 0x4e, 0x61, 0xf4, 0x41, 0x5e, 0x6f, 0x8c,
	0xf8, 0x22, 0x1d, 0x4e, 0x7d, 0xe3, 0x45, 0x59, 0x32, 0x3c, 0x20, 0x18, 0x26, 0x17, 0x46, 0xf1,
	0x4a, 0xf7, 0xbd, 0x93, 0xd5, 0x5b, 0x78, 0x2c, 0x8c, 0x3a, 0x1c, 0x1f, 0x43, 0x57, 0xc3, 0x50,
	0xff, 0x48, 0xce, 0xde, 0xe7, 0x05, 0x6f, 0xe8, 0xda, 0x6b, 0xe3, 0x5a, 0xf4, 0x65, 0xdf, 0xbf,
	0xbe, 0xdb, 0x1e, 0xed, 0xf9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1c, 0x47, 0x39, 0x56, 0x83,
	0x03, 0x00, 0x00,
}
