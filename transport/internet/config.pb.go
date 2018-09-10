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
	return fileDescriptor_config_d7068fe3707ce485, []int{0}
}

type SocketConfig_TCPFastOpenState int32

const (
	// AsIs is to leave the current TFO state as is, unmodified.
	SocketConfig_AsIs SocketConfig_TCPFastOpenState = 0
	// Enable is for enabling TFO explictly.
	SocketConfig_Enable SocketConfig_TCPFastOpenState = 1
	// Disable is for disabling TFO explictly.
	SocketConfig_Disable SocketConfig_TCPFastOpenState = 2
)

var SocketConfig_TCPFastOpenState_name = map[int32]string{
	0: "AsIs",
	1: "Enable",
	2: "Disable",
}
var SocketConfig_TCPFastOpenState_value = map[string]int32{
	"AsIs":    0,
	"Enable":  1,
	"Disable": 2,
}

func (x SocketConfig_TCPFastOpenState) String() string {
	return proto.EnumName(SocketConfig_TCPFastOpenState_name, int32(x))
}
func (SocketConfig_TCPFastOpenState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_d7068fe3707ce485, []int{3, 0}
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
	return fileDescriptor_config_d7068fe3707ce485, []int{0}
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
	return fileDescriptor_config_d7068fe3707ce485, []int{1}
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
	return fileDescriptor_config_d7068fe3707ce485, []int{2}
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

// SocketConfig is options to be applied on network sockets.
type SocketConfig struct {
	// Mark of the connection. If non-zero, the value will be set to SO_MARK.
	Mark int32 `protobuf:"varint,1,opt,name=mark,proto3" json:"mark,omitempty"`
	// TFO is the state of TFO settings.
	Tfo                  SocketConfig_TCPFastOpenState `protobuf:"varint,2,opt,name=tfo,proto3,enum=v2ray.core.transport.internet.SocketConfig_TCPFastOpenState" json:"tfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SocketConfig) Reset()         { *m = SocketConfig{} }
func (m *SocketConfig) String() string { return proto.CompactTextString(m) }
func (*SocketConfig) ProtoMessage()    {}
func (*SocketConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_d7068fe3707ce485, []int{3}
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

func (m *SocketConfig) GetTfo() SocketConfig_TCPFastOpenState {
	if m != nil {
		return m.Tfo
	}
	return SocketConfig_AsIs
}

func init() {
	proto.RegisterType((*TransportConfig)(nil), "v2ray.core.transport.internet.TransportConfig")
	proto.RegisterType((*StreamConfig)(nil), "v2ray.core.transport.internet.StreamConfig")
	proto.RegisterType((*ProxyConfig)(nil), "v2ray.core.transport.internet.ProxyConfig")
	proto.RegisterType((*SocketConfig)(nil), "v2ray.core.transport.internet.SocketConfig")
	proto.RegisterEnum("v2ray.core.transport.internet.TransportProtocol", TransportProtocol_name, TransportProtocol_value)
	proto.RegisterEnum("v2ray.core.transport.internet.SocketConfig_TCPFastOpenState", SocketConfig_TCPFastOpenState_name, SocketConfig_TCPFastOpenState_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_config_d7068fe3707ce485)
}

var fileDescriptor_config_d7068fe3707ce485 = []byte{
	// 518 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xd1, 0x8a, 0x13, 0x31,
	0x14, 0x86, 0x77, 0x3a, 0x6d, 0xb7, 0x3d, 0xed, 0x76, 0xd3, 0x5c, 0x15, 0x61, 0xb1, 0x56, 0x90,
	0xa2, 0x90, 0x59, 0x46, 0xbc, 0xf3, 0xc6, 0xb6, 0x8a, 0x8b, 0xee, 0xee, 0x30, 0x1d, 0x15, 0x16,
	0xa4, 0xa4, 0x63, 0xb6, 0x0c, 0xdb, 0x99, 0x94, 0x24, 0x8a, 0x7d, 0x24, 0xbd, 0xf6, 0x21, 0x7c,
	0x2c, 0x49, 0x66, 0x12, 0xca, 0x2a, 0x75, 0xc5, 0xbb, 0x33, 0x73, 0xfe, 0xfc, 0xe7, 0xff, 0x72,
	0x02, 0xe4, 0x4b, 0x28, 0xe8, 0x96, 0xa4, 0x3c, 0x0f, 0x52, 0x2e, 0x58, 0xa0, 0x04, 0x2d, 0xe4,
	0x86, 0x0b, 0x15, 0x64, 0x85, 0x62, 0xa2, 0x60, 0x2a, 0x48, 0x79, 0x71, 0x9d, 0xad, 0xc8, 0x46,
	0x70, 0xc5, 0xf1, 0x89, 0xd5, 0x0b, 0x46, 0x9c, 0x96, 0x58, 0xed, 0xbd, 0xd3, 0x5b, 0x76, 0x29,
	0xcf, 0x73, 0x5e, 0x04, 0x92, 0x89, 0x8c, 0xae, 0x03, 0xb5, 0xdd, 0xb0, 0x4f, 0x8b, 0x9c, 0x49,
	0x49, 0x57, 0xac, 0x34, 0x1c, 0xfd, 0xf4, 0xe0, 0x38, 0xb1, 0x46, 0x53, 0x33, 0x0a, 0xbf, 0x85,
	0x96, 0x69, 0xa6, 0x7c, 0x3d, 0xf0, 0x86, 0xde, 0xb8, 0x17, 0x9e, 0x92, 0xbd, 0x73, 0x89, 0x73,
	0x88, 0xaa, 0x73, 0xb1, 0x73, 0xc0, 0x0f, 0xe1, 0xc8, 0xd6, 0x8b, 0x82, 0xe6, 0x6c, 0xe0, 0x0f,
	0xbd, 0x71, 0x3b, 0xee, 0xda, 0x9f, 0x17, 0x34, 0x67, 0x78, 0x02, 0x2d, 0xc9, 0x94, 0xca, 0x8a,
	0x95, 0x1c, 0xd4, 0x86, 0xde, 0xb8, 0x13, 0x3e, 0xda, 0x1d, 0x59, 0x72, 0x90, 0x92, 0x83, 0x24,
	0x9a, 0xe3, 0xbc, 0xc4, 0x88, 0xdd, 0xb9, 0xd1, 0x0f, 0x1f, 0xba, 0x73, 0x25, 0x18, 0xcd, 0x2b,
	0x8e, 0xe8, 0xff, 0x39, 0x26, 0xb5, 0x81, 0xb7, 0x8f, 0xa5, 0xf1, 0x07, 0x96, 0x8f, 0x80, 0x9d,
	0xf5, 0x62, 0x87, 0xca, 0x1f, 0x77, 0x42, 0x72, 0xd7, 0x00, 0x25, 0x42, 0xdc, 0x77, 0x9a, 0x79,
	0x65, 0xa4, 0x33, 0x48, 0x96, 0x7e, 0x16, 0x99, 0xda, 0x2e, 0xf4, 0x46, 0xed, 0x7d, 0xda, 0x9f,
	0xfa, 0x76, 0xf0, 0x1c, 0xfa, 0x4e, 0xe4, 0x22, 0xd4, 0x4d, 0x84, 0xbb, 0x5e, 0x2c, 0xb2, 0x06,
	0x6e, 0x72, 0x02, 0xc7, 0x92, 0xa7, 0x37, 0x6c, 0x87, 0xaa, 0x69, 0x76, 0xf5, 0xe4, 0x2f, 0x54,
	0x73, 0x73, 0xaa, 0x42, 0xea, 0x95, 0x1e, 0xd6, 0x75, 0x74, 0x1f, 0x3a, 0x91, 0xe0, 0x5f, 0xb7,
	0xd5, 0xd2, 0x10, 0xf8, 0x8a, 0xae, 0xcc, 0xbe, 0xda, 0xb1, 0x2e, 0x47, 0xdf, 0x3c, 0xe8, 0xee,
	0x3a, 0x60, 0x0c, 0xf5, 0x9c, 0x8a, 0x1b, 0xa3, 0x69, 0xc4, 0xa6, 0xc6, 0x17, 0xe0, 0xab, 0x6b,
	0x6e, 0xde, 0x4e, 0x2f, 0x7c, 0xfe, 0x0f, 0x79, 0x48, 0x32, 0x8d, 0x5e, 0x51, 0xa9, 0x2e, 0x37,
	0xac, 0x98, 0x2b, 0xaa, 0x58, 0xac, 0x8d, 0x46, 0xcf, 0x00, 0xdd, 0x6e, 0xe0, 0x16, 0xd4, 0x5f,
	0xc8, 0x33, 0x89, 0x0e, 0x30, 0x40, 0xf3, 0x65, 0x41, 0x97, 0x6b, 0x86, 0x3c, 0xdc, 0x81, 0xc3,
	0x59, 0x26, 0xcd, 0x47, 0xed, 0xf1, 0x15, 0xf4, 0x7f, 0x7b, 0x43, 0xf8, 0x10, 0xfc, 0x64, 0x1a,
	0xa1, 0x03, 0x5d, 0xbc, 0x9b, 0x45, 0xc8, 0xd3, 0x4e, 0xe7, 0x6f, 0xa6, 0x11, 0xaa, 0xe1, 0x23,
	0x68, 0x7f, 0x60, 0xcb, 0x32, 0x10, 0xf2, 0x75, 0xe3, 0x75, 0x92, 0x44, 0xa8, 0x8e, 0x11, 0x74,
	0x67, 0x3c, 0xa7, 0x59, 0x51, 0xf5, 0x1a, 0x93, 0x4b, 0x78, 0x90, 0xf2, 0x7c, 0x3f, 0x5a, 0xe4,
	0x5d, 0xb5, 0x6c, 0xfd, 0xbd, 0x76, 0xf2, 0x3e, 0x8c, 0xe9, 0x96, 0x4c, 0xb5, 0xd6, 0xc5, 0x22,
	0x67, 0x55, 0x7f, 0xd9, 0x34, 0xcf, 0xf6, 0xe9, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf3, 0x8d,
	0xe5, 0xd8, 0x85, 0x04, 0x00, 0x00,
}
