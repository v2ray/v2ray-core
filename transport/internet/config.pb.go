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
	return fileDescriptor_config_6493eeef2ca10012, []int{0}
}

type TransportConfig struct {
	// Type of network that this settings supports.
	Protocol TransportProtocol `protobuf:"varint,1,opt,name=protocol,proto3,enum=v2ray.core.transport.internet.TransportProtocol" json:"protocol,omitempty"`
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
	return fileDescriptor_config_6493eeef2ca10012, []int{0}
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

func (m *TransportConfig) GetSettings() *serial.TypedMessage {
	if m != nil {
		return m.Settings
	}
	return nil
}

type StreamConfig struct {
	// Effective network.
	Protocol          TransportProtocol  `protobuf:"varint,1,opt,name=protocol,proto3,enum=v2ray.core.transport.internet.TransportProtocol" json:"protocol,omitempty"`
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
	return fileDescriptor_config_6493eeef2ca10012, []int{1}
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

func (m *StreamConfig) GetProtocol() TransportProtocol {
	if m != nil {
		return m.Protocol
	}
	return TransportProtocol_TCP
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
	return fileDescriptor_config_6493eeef2ca10012, []int{2}
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
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_config_6493eeef2ca10012)
}

var fileDescriptor_config_6493eeef2ca10012 = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x91, 0xcf, 0x6a, 0xdb, 0x40,
	0x10, 0x87, 0x2b, 0xc9, 0x6d, 0xe5, 0xb1, 0xdd, 0xae, 0xf7, 0x64, 0x0a, 0xa6, 0xae, 0x0b, 0x45,
	0xf4, 0xb0, 0x32, 0xea, 0x1b, 0x58, 0x3e, 0xb4, 0xb4, 0xa6, 0x42, 0x52, 0x5b, 0x30, 0x14, 0xb3,
	0xde, 0x6c, 0x84, 0x88, 0xa5, 0x35, 0xab, 0x4d, 0x88, 0x9e, 0x27, 0xb7, 0xdc, 0xf3, 0x7e, 0x41,
	0xff, 0x16, 0x93, 0x80, 0xf1, 0x25, 0xb7, 0x41, 0xf3, 0x9b, 0x6f, 0x3e, 0xcd, 0x02, 0xb9, 0xf1,
	0x24, 0x2d, 0x09, 0x13, 0x99, 0xcb, 0x84, 0xe4, 0xae, 0x92, 0x34, 0x2f, 0x0e, 0x42, 0x2a, 0x37,
	0xcd, 0x15, 0x97, 0x39, 0x57, 0x2e, 0x13, 0xf9, 0x65, 0x9a, 0x90, 0x83, 0x14, 0x4a, 0xe0, 0x69,
	0x97, 0x97, 0x9c, 0xe8, 0x2c, 0xe9, 0xb2, 0x1f, 0x16, 0x4f, 0x70, 0x4c, 0x64, 0x99, 0xc8, 0xdd,
	0x82, 0xcb, 0x94, 0xee, 0x5d, 0x55, 0x1e, 0xf8, 0xc5, 0x36, 0xe3, 0x45, 0x41, 0x13, 0xde, 0x00,
	0xe7, 0x77, 0x06, 0xbc, 0x8f, 0x3b, 0x90, 0x5f, 0xaf, 0xc2, 0xbf, 0xc0, 0xae, 0x9b, 0x4c, 0xec,
	0x27, 0xc6, 0xcc, 0x70, 0xde, 0x79, 0x0b, 0x72, 0x72, 0x2f, 0xd1, 0x84, 0xa0, 0x9d, 0x0b, 0x35,
	0x01, 0x2f, 0xc1, 0x2e, 0xb8, 0x52, 0x69, 0x9e, 0x14, 0x13, 0x73, 0x66, 0x38, 0x03, 0xef, 0xcb,
	0x31, 0xad, 0x51, 0x24, 0x8d, 0x22, 0x89, 0x2b, 0xc5, 0x75, 0x63, 0x18, 0xea, 0xb9, 0xf9, 0x83,
	0x09, 0xc3, 0x48, 0x49, 0x4e, 0xb3, 0x17, 0x51, 0xfc, 0x0f, 0x58, 0x4f, 0x6c, 0x8f, 0x64, 0x2d,
	0x67, 0xe0, 0x91, 0x73, 0xb9, 0x8d, 0x59, 0x38, 0xd6, 0x99, 0xa8, 0x05, 0xe1, 0xcf, 0x30, 0x2a,
	0x38, 0xbb, 0x96, 0xa9, 0x2a, 0xb7, 0xd5, 0x1b, 0x4c, 0xac, 0x99, 0xe1, 0xf4, 0xc3, 0x61, 0xf7,
	0xb1, 0xfa, 0x69, 0x1c, 0xc1, 0x58, 0x87, 0xb4, 0x42, 0xaf, 0x56, 0x38, 0xf7, 0x5e, 0xa8, 0x03,
	0x74, 0x9b, 0xe7, 0x1f, 0x61, 0x10, 0x48, 0x71, 0x5b, 0xb6, 0x57, 0x43, 0x60, 0x29, 0x9a, 0xd4,
	0x07, 0xeb, 0x87, 0x55, 0xf9, 0x75, 0x03, 0xe3, 0x67, 0x87, 0xc1, 0x6f, 0xc1, 0x8a, 0xfd, 0x00,
	0xbd, 0xaa, 0x8a, 0x3f, 0xab, 0x00, 0x19, 0xd8, 0x86, 0xde, 0xfa, 0xa7, 0x1f, 0x20, 0x13, 0x8f,
	0xa0, 0xff, 0x8f, 0xef, 0x22, 0xc1, 0xae, 0xb8, 0x42, 0x56, 0xd5, 0xf8, 0x1e, 0xc7, 0x01, 0xea,
	0x61, 0x04, 0xc3, 0x95, 0xc8, 0x68, 0x9a, 0xb7, 0xbd, 0xd7, 0xcb, 0xdf, 0xf0, 0x89, 0x89, 0xec,
	0xf4, 0xf9, 0x02, 0x63, 0x63, 0x77, 0xf5, 0xbd, 0x39, 0xfd, 0xeb, 0x85, 0xb4, 0x24, 0x7e, 0x95,
	0xd5, 0x5a, 0xe4, 0x47, 0xdb, 0xdf, 0xbd, 0xa9, 0x1f, 0xec, 0xdb, 0x63, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xa2, 0xdf, 0xde, 0xa4, 0x35, 0x03, 0x00, 0x00,
}
