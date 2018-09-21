package internet

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	serial "v2ray.com/core/common/serial"
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

type SocketConfig_TProxyMode int32

const (
	// TProxy is off.
	SocketConfig_Off SocketConfig_TProxyMode = 0
	// TProxy mode.
	SocketConfig_TProxy SocketConfig_TProxyMode = 1
	// Redirect mode.
	SocketConfig_Redirect SocketConfig_TProxyMode = 2
)

var SocketConfig_TProxyMode_name = map[int32]string{
	0: "Off",
	1: "TProxy",
	2: "Redirect",
}

var SocketConfig_TProxyMode_value = map[string]int32{
	"Off":      0,
	"TProxy":   1,
	"Redirect": 2,
}

func (x SocketConfig_TProxyMode) String() string {
	return proto.EnumName(SocketConfig_TProxyMode_name, int32(x))
}

func (SocketConfig_TProxyMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_91dbc815c3d97a05, []int{3, 1}
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
	Tproxy SocketConfig_TProxyMode `protobuf:"varint,3,opt,name=tproxy,proto3,enum=v2ray.core.transport.internet.SocketConfig_TProxyMode" json:"tproxy,omitempty"`
	// ReceiveOriginalDestAddress is for enabling IP_RECVORIGDSTADDR socket option.
	ReceiveOriginalDestAddress bool     `protobuf:"varint,4,opt,name=receive_original_dest_address,json=receiveOriginalDestAddress,proto3" json:"receive_original_dest_address,omitempty"`
	XXX_NoUnkeyedLiteral       struct{} `json:"-"`
	XXX_unrecognized           []byte   `json:"-"`
	XXX_sizecache              int32    `json:"-"`
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

func (m *SocketConfig) GetTproxy() SocketConfig_TProxyMode {
	if m != nil {
		return m.Tproxy
	}
	return SocketConfig_Off
}

func (m *SocketConfig) GetReceiveOriginalDestAddress() bool {
	if m != nil {
		return m.ReceiveOriginalDestAddress
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
	proto.RegisterEnum("v2ray.core.transport.internet.SocketConfig_TProxyMode", SocketConfig_TProxyMode_name, SocketConfig_TProxyMode_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_91dbc815c3d97a05)
}

var fileDescriptor_91dbc815c3d97a05 = []byte{
	// 607 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0xad, 0xed, 0x34, 0x4d, 0x27, 0x69, 0xea, 0xee, 0x55, 0x54, 0xa9, 0xfa, 0xfa, 0x05, 0x09,
	0x45, 0x20, 0xad, 0x2b, 0x23, 0xb8, 0xe2, 0xa6, 0x4d, 0x40, 0x54, 0xd0, 0xc6, 0x72, 0x0c, 0x48,
	0x95, 0x90, 0xb5, 0x75, 0x26, 0x91, 0xd5, 0xd8, 0x1b, 0xed, 0x2e, 0x15, 0x79, 0x25, 0xae, 0x79,
	0x08, 0x5e, 0x86, 0x77, 0x40, 0xbb, 0xfe, 0x21, 0x2a, 0x28, 0xb4, 0xe2, 0x6e, 0x3c, 0x73, 0xe6,
	0xcc, 0x39, 0x33, 0x5e, 0xa0, 0xb7, 0xbe, 0x60, 0x2b, 0x9a, 0xf0, 0xcc, 0x4b, 0xb8, 0x40, 0x4f,
	0x09, 0x96, 0xcb, 0x25, 0x17, 0xca, 0x4b, 0x73, 0x85, 0x22, 0x47, 0xe5, 0x25, 0x3c, 0x9f, 0xa5,
	0x73, 0xba, 0x14, 0x5c, 0x71, 0x72, 0x54, 0xe1, 0x05, 0xd2, 0x1a, 0x4b, 0x2b, 0xec, 0xe1, 0xc9,
	0x1d, 0xba, 0x84, 0x67, 0x19, 0xcf, 0x3d, 0x89, 0x22, 0x65, 0x0b, 0x4f, 0xad, 0x96, 0x38, 0x8d,
	0x33, 0x94, 0x92, 0xcd, 0xb1, 0x20, 0xec, 0x7f, 0xb7, 0x60, 0x3f, 0xaa, 0x88, 0x86, 0x66, 0x14,
	0x79, 0x07, 0x2d, 0x53, 0x4c, 0xf8, 0xa2, 0x67, 0x1d, 0x5b, 0x83, 0xae, 0x7f, 0x42, 0x37, 0xce,
	0xa5, 0x35, 0x43, 0x50, 0xf6, 0x85, 0x35, 0x03, 0x79, 0x04, 0x7b, 0x55, 0x1c, 0xe7, 0x2c, 0xc3,
	0x9e, 0x73, 0x6c, 0x0d, 0x76, 0xc3, 0x4e, 0x95, 0xbc, 0x64, 0x19, 0x92, 0x33, 0x68, 0x49, 0x54,
	0x2a, 0xcd, 0xe7, 0xb2, 0x67, 0x1f, 0x5b, 0x83, 0xb6, 0xff, 0x78, 0x7d, 0x64, 0xe1, 0x83, 0x16,
	0x3e, 0x68, 0xa4, 0x7d, 0x5c, 0x14, 0x36, 0xc2, 0xba, 0xaf, 0xff, 0xcd, 0x81, 0xce, 0x44, 0x09,
	0x64, 0x59, 0xe9, 0x23, 0xf8, 0x77, 0x1f, 0x67, 0x76, 0xcf, 0xda, 0xe4, 0x65, 0xfb, 0x0f, 0x5e,
	0x3e, 0x01, 0xa9, 0xa9, 0xe3, 0x35, 0x57, 0xce, 0xa0, 0xed, 0xd3, 0xfb, 0x0a, 0x28, 0x2c, 0x84,
	0x07, 0x35, 0x66, 0x52, 0x12, 0x69, 0x0d, 0x12, 0x93, 0xcf, 0x22, 0x55, 0xab, 0x58, 0x5f, 0xb4,
	0xda, 0x67, 0x95, 0xd4, 0xdb, 0x21, 0x13, 0x38, 0xa8, 0x41, 0xb5, 0x84, 0x86, 0x91, 0x70, 0xdf,
	0xc5, 0xba, 0x15, 0x41, 0x3d, 0x39, 0x82, 0x7d, 0xc9, 0x93, 0x1b, 0x5c, 0x73, 0xd5, 0x34, 0xb7,
	0x7a, 0xfa, 0x17, 0x57, 0x13, 0xd3, 0x55, 0x5a, 0xea, 0x16, 0x1c, 0x15, 0x6b, 0xff, 0x3f, 0x68,
	0x07, 0x82, 0x7f, 0x59, 0x95, 0x47, 0x73, 0xc1, 0x51, 0x6c, 0x6e, 0xee, 0xb5, 0x1b, 0xea, 0xb0,
	0xff, 0xc3, 0x86, 0xce, 0x3a, 0x03, 0x21, 0xd0, 0xc8, 0x98, 0xb8, 0x31, 0x98, 0xed, 0xd0, 0xc4,
	0xe4, 0x12, 0x1c, 0x35, 0xe3, 0xe6, 0xdf, 0xe9, 0xfa, 0x2f, 0x1f, 0xa0, 0x87, 0x46, 0xc3, 0xe0,
	0x35, 0x93, 0x6a, 0xbc, 0xc4, 0x7c, 0xa2, 0x98, 0xc2, 0x50, 0x13, 0x91, 0x4b, 0x68, 0xaa, 0xa5,
	0x96, 0x65, 0xd6, 0xdb, 0xf5, 0x5f, 0x3c, 0x88, 0xd2, 0x18, 0xba, 0xe0, 0x53, 0x0c, 0x4b, 0x16,
	0x72, 0x0a, 0x47, 0x02, 0x13, 0x4c, 0x6f, 0x31, 0xe6, 0x22, 0x9d, 0xa7, 0x39, 0x5b, 0xc4, 0x53,
	0x94, 0x2a, 0x66, 0xd3, 0xa9, 0x40, 0xa9, 0x8f, 0x63, 0x0d, 0x5a, 0xe1, 0x61, 0x09, 0x1a, 0x97,
	0x98, 0x11, 0x4a, 0x75, 0x5a, 0x20, 0xfa, 0xcf, 0xc1, 0xbd, 0xab, 0x95, 0xb4, 0xa0, 0x71, 0x2a,
	0xcf, 0xa5, 0xbb, 0x45, 0x00, 0x9a, 0xaf, 0x72, 0x76, 0xbd, 0x40, 0xd7, 0x22, 0x6d, 0xd8, 0x19,
	0xa5, 0xd2, 0x7c, 0xd8, 0x7d, 0x0f, 0xe0, 0x97, 0x1e, 0xb2, 0x03, 0xce, 0x78, 0x36, 0x2b, 0xf0,
	0x45, 0xda, 0xb5, 0x48, 0x07, 0x5a, 0x21, 0x4e, 0x53, 0x81, 0x89, 0x72, 0xed, 0x27, 0x57, 0x70,
	0xf0, 0xdb, 0x3b, 0xd0, 0x7d, 0xd1, 0x30, 0x70, 0xb7, 0x74, 0xf0, 0x7e, 0x14, 0xb8, 0x96, 0x1e,
	0x7d, 0xf1, 0x76, 0x18, 0xb8, 0x36, 0xd9, 0x83, 0xdd, 0x8f, 0x78, 0x5d, 0x6c, 0xc0, 0x75, 0x74,
	0xe1, 0x4d, 0x14, 0x05, 0x6e, 0x83, 0xb8, 0xd0, 0x19, 0xf1, 0x8c, 0xa5, 0x79, 0x59, 0xdb, 0x3e,
	0x1b, 0xc3, 0xff, 0x09, 0xcf, 0x36, 0xef, 0x32, 0xb0, 0xae, 0x5a, 0x55, 0xfc, 0xd5, 0x3e, 0xfa,
	0xe0, 0x87, 0x6c, 0x45, 0x87, 0x1a, 0x5b, 0xcb, 0xa2, 0xe7, 0x65, 0xfd, 0xba, 0x69, 0x9e, 0xde,
	0xb3, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf2, 0x99, 0xd5, 0x37, 0x49, 0x05, 0x00, 0x00,
}
