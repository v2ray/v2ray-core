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
	return fileDescriptor_91dbc815c3d97a05, []int{0}
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
	return fileDescriptor_91dbc815c3d97a05, []int{3, 0}
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
	return fileDescriptor_91dbc815c3d97a05, []int{0}
}
func (m *TransportConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransportConfig.Unmarshal(m, b)
}
func (m *TransportConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransportConfig.Marshal(b, m, deterministic)
}
func (m *TransportConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransportConfig.Merge(m, src)
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
	return fileDescriptor_91dbc815c3d97a05, []int{1}
}
func (m *StreamConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamConfig.Unmarshal(m, b)
}
func (m *StreamConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamConfig.Marshal(b, m, deterministic)
}
func (m *StreamConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamConfig.Merge(m, src)
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
	return fileDescriptor_91dbc815c3d97a05, []int{2}
}
func (m *ProxyConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyConfig.Unmarshal(m, b)
}
func (m *ProxyConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyConfig.Marshal(b, m, deterministic)
}
func (m *ProxyConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyConfig.Merge(m, src)
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
	Tfo SocketConfig_TCPFastOpenState `protobuf:"varint,2,opt,name=tfo,proto3,enum=v2ray.core.transport.internet.SocketConfig_TCPFastOpenState" json:"tfo,omitempty"`
	// TProxy is for enabling TProxy socket option.
	Tproxy               bool     `protobuf:"varint,3,opt,name=tproxy,proto3" json:"tproxy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SocketConfig) Reset()         { *m = SocketConfig{} }
func (m *SocketConfig) String() string { return proto.CompactTextString(m) }
func (*SocketConfig) ProtoMessage()    {}
func (*SocketConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_91dbc815c3d97a05, []int{3}
}
func (m *SocketConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SocketConfig.Unmarshal(m, b)
}
func (m *SocketConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SocketConfig.Marshal(b, m, deterministic)
}
func (m *SocketConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SocketConfig.Merge(m, src)
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

func (m *SocketConfig) GetTproxy() bool {
	if m != nil {
		return m.Tproxy
	}
	return false
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
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_91dbc815c3d97a05)
}

var fileDescriptor_91dbc815c3d97a05 = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xd1, 0x6b, 0x13, 0x4f,
	0x10, 0xc7, 0x7b, 0xb9, 0x24, 0x4d, 0x26, 0x69, 0xba, 0xd9, 0x87, 0x1f, 0xe1, 0x07, 0xc5, 0x18,
	0x41, 0x82, 0xc2, 0x5e, 0x39, 0xf1, 0xcd, 0x17, 0x93, 0x28, 0x16, 0x6d, 0x7b, 0x5c, 0x4e, 0x85,
	0x82, 0x84, 0xcd, 0xb9, 0x0d, 0x47, 0x73, 0xb7, 0x61, 0x77, 0x15, 0xf3, 0x2f, 0xf9, 0xec, 0x3f,
	0xe0, 0x9b, 0x7f, 0x96, 0xec, 0xde, 0xed, 0x12, 0xaa, 0xc4, 0x8a, 0x6f, 0x73, 0x37, 0x33, 0xdf,
	0xf9, 0x7e, 0x66, 0x16, 0xc8, 0xe7, 0x50, 0xd0, 0x2d, 0x49, 0x79, 0x1e, 0xa4, 0x5c, 0xb0, 0x40,
	0x09, 0x5a, 0xc8, 0x0d, 0x17, 0x2a, 0xc8, 0x0a, 0xc5, 0x44, 0xc1, 0x54, 0x90, 0xf2, 0xe2, 0x3a,
	0x5b, 0x91, 0x8d, 0xe0, 0x8a, 0xe3, 0x13, 0x5b, 0x2f, 0x18, 0x71, 0xb5, 0xc4, 0xd6, 0xfe, 0x7f,
	0x7a, 0x4b, 0x2e, 0xe5, 0x79, 0xce, 0x8b, 0x40, 0x32, 0x91, 0xd1, 0x75, 0xa0, 0xb6, 0x1b, 0xf6,
	0x71, 0x91, 0x33, 0x29, 0xe9, 0x8a, 0x95, 0x82, 0xa3, 0x1f, 0x1e, 0x1c, 0x27, 0x56, 0x68, 0x6a,
	0x46, 0xe1, 0x37, 0xd0, 0x32, 0xc9, 0x94, 0xaf, 0x07, 0xde, 0xd0, 0x1b, 0xf7, 0xc2, 0x53, 0xb2,
	0x77, 0x2e, 0x71, 0x0a, 0x51, 0xd5, 0x17, 0x3b, 0x05, 0xfc, 0x00, 0x8e, 0x6c, 0xbc, 0x28, 0x68,
	0xce, 0x06, 0xfe, 0xd0, 0x1b, 0xb7, 0xe3, 0xae, 0xfd, 0x79, 0x41, 0x73, 0x86, 0x27, 0xd0, 0x92,
	0x4c, 0xa9, 0xac, 0x58, 0xc9, 0x41, 0x6d, 0xe8, 0x8d, 0x3b, 0xe1, 0xc3, 0xdd, 0x91, 0x25, 0x07,
	0x29, 0x39, 0x48, 0xa2, 0x39, 0xce, 0x4b, 0x8c, 0xd8, 0xf5, 0x8d, 0xbe, 0xf9, 0xd0, 0x9d, 0x2b,
	0xc1, 0x68, 0x5e, 0x71, 0x44, 0xff, 0xce, 0x31, 0xa9, 0x0d, 0xbc, 0x7d, 0x2c, 0x8d, 0xdf, 0xb0,
	0x7c, 0x00, 0xec, 0xa4, 0x17, 0x3b, 0x54, 0xfe, 0xb8, 0x13, 0x92, 0xbb, 0x1a, 0x28, 0x11, 0xe2,
	0xbe, 0xab, 0x99, 0x57, 0x42, 0xda, 0x83, 0x64, 0xe9, 0x27, 0x91, 0xa9, 0xed, 0x42, 0x5f, 0xd4,
	0xee, 0xd3, 0xfe, 0xd4, 0xdb, 0xc1, 0x73, 0xe8, 0xbb, 0x22, 0x67, 0xa1, 0x6e, 0x2c, 0xdc, 0x75,
	0xb1, 0xc8, 0x0a, 0xb8, 0xc9, 0x09, 0x1c, 0x4b, 0x9e, 0xde, 0xb0, 0x1d, 0xaa, 0xa6, 0xb9, 0xd5,
	0xe3, 0x3f, 0x50, 0xcd, 0x4d, 0x57, 0x85, 0xd4, 0x2b, 0x35, 0xac, 0xea, 0xe8, 0x1e, 0x74, 0x22,
	0xc1, 0xbf, 0x6c, 0xab, 0xa3, 0x21, 0xf0, 0x15, 0x5d, 0x99, 0x7b, 0xb5, 0x63, 0x1d, 0x8e, 0xbe,
	0x7b, 0xd0, 0xdd, 0x55, 0xc0, 0x18, 0xea, 0x39, 0x15, 0x37, 0xa6, 0xa6, 0x11, 0x9b, 0x18, 0x5f,
	0x80, 0xaf, 0xae, 0xb9, 0x79, 0x3b, 0xbd, 0xf0, 0xd9, 0x5f, 0xf8, 0x21, 0xc9, 0x34, 0x7a, 0x49,
	0xa5, 0xba, 0xdc, 0xb0, 0x62, 0xae, 0xa8, 0x62, 0xb1, 0x16, 0xc2, 0xff, 0x41, 0x53, 0x6d, 0xb4,
	0x2d, 0xb3, 0xde, 0x56, 0x5c, 0x7d, 0x8d, 0x9e, 0x02, 0xba, 0xdd, 0x80, 0x5b, 0x50, 0x7f, 0x2e,
	0xcf, 0x24, 0x3a, 0xc0, 0x00, 0xcd, 0x17, 0x05, 0x5d, 0xae, 0x19, 0xf2, 0x70, 0x07, 0x0e, 0x67,
	0x99, 0x34, 0x1f, 0xb5, 0x47, 0x57, 0xd0, 0xff, 0xe5, 0x6d, 0xe1, 0x43, 0xf0, 0x93, 0x69, 0x84,
	0x0e, 0x74, 0xf0, 0x76, 0x16, 0x21, 0x4f, 0x2b, 0x9d, 0xbf, 0x9e, 0x46, 0xa8, 0x86, 0x8f, 0xa0,
	0xfd, 0x9e, 0x2d, 0x4b, 0xa3, 0xc8, 0xd7, 0x89, 0x57, 0x49, 0x12, 0xa1, 0x3a, 0x46, 0xd0, 0x9d,
	0xf1, 0x9c, 0x66, 0x45, 0x95, 0x6b, 0x4c, 0x2e, 0xe1, 0x7e, 0xca, 0xf3, 0xfd, 0xc8, 0x91, 0x77,
	0xd5, 0xb2, 0xf1, 0xd7, 0xda, 0xc9, 0xbb, 0x30, 0xa6, 0x5b, 0x32, 0xd5, 0xb5, 0xce, 0x16, 0x39,
	0xab, 0xf2, 0xcb, 0xa6, 0x79, 0xce, 0x4f, 0x7e, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x1d, 0xbb,
	0xfd, 0x9d, 0x04, 0x00, 0x00,
}
