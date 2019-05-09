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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

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
	// This option is for UDP only.
	ReceiveOriginalDestAddress bool     `protobuf:"varint,4,opt,name=receive_original_dest_address,json=receiveOriginalDestAddress,proto3" json:"receive_original_dest_address,omitempty"`
	BindAddress                []byte   `protobuf:"bytes,5,opt,name=bind_address,json=bindAddress,proto3" json:"bind_address,omitempty"`
	BindPort                   uint32   `protobuf:"varint,6,opt,name=bind_port,json=bindPort,proto3" json:"bind_port,omitempty"`
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

func (m *SocketConfig) GetBindAddress() []byte {
	if m != nil {
		return m.BindAddress
	}
	return nil
}

func (m *SocketConfig) GetBindPort() uint32 {
	if m != nil {
		return m.BindPort
	}
	return 0
}

func init() {
	proto.RegisterEnum("v2ray.core.transport.internet.TransportProtocol", TransportProtocol_name, TransportProtocol_value)
	proto.RegisterEnum("v2ray.core.transport.internet.SocketConfig_TCPFastOpenState", SocketConfig_TCPFastOpenState_name, SocketConfig_TCPFastOpenState_value)
	proto.RegisterEnum("v2ray.core.transport.internet.SocketConfig_TProxyMode", SocketConfig_TProxyMode_name, SocketConfig_TProxyMode_value)
	proto.RegisterType((*TransportConfig)(nil), "v2ray.core.transport.internet.TransportConfig")
	proto.RegisterType((*StreamConfig)(nil), "v2ray.core.transport.internet.StreamConfig")
	proto.RegisterType((*ProxyConfig)(nil), "v2ray.core.transport.internet.ProxyConfig")
	proto.RegisterType((*SocketConfig)(nil), "v2ray.core.transport.internet.SocketConfig")
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/config.proto", fileDescriptor_91dbc815c3d97a05)
}

var fileDescriptor_91dbc815c3d97a05 = []byte{
	// 636 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xdd, 0x6a, 0xdb, 0x4c,
	0x10, 0x8d, 0x2c, 0xc7, 0x91, 0xc7, 0x8e, 0xa3, 0xec, 0x95, 0xc9, 0x47, 0xf8, 0x12, 0x17, 0x8a,
	0x69, 0x41, 0x0a, 0x2a, 0xed, 0x55, 0x6f, 0x12, 0xbb, 0xa5, 0xa1, 0x4d, 0x2c, 0x64, 0xb5, 0x85,
	0x40, 0x11, 0x6b, 0x69, 0x6c, 0x44, 0x2c, 0xad, 0xd9, 0xdd, 0x86, 0xfa, 0x95, 0x0a, 0xbd, 0xeb,
	0x43, 0xf4, 0xb1, 0xca, 0xae, 0x7e, 0x6a, 0xd2, 0x92, 0x26, 0xf4, 0x6e, 0x34, 0x73, 0xe6, 0xcc,
	0x9c, 0x39, 0x5a, 0x70, 0x6e, 0x3c, 0x4e, 0xd7, 0x4e, 0xcc, 0x32, 0x37, 0x66, 0x1c, 0x5d, 0xc9,
	0x69, 0x2e, 0x56, 0x8c, 0x4b, 0x37, 0xcd, 0x25, 0xf2, 0x1c, 0xa5, 0x1b, 0xb3, 0x7c, 0x9e, 0x2e,
	0x9c, 0x15, 0x67, 0x92, 0x91, 0xc3, 0x0a, 0xcf, 0xd1, 0xa9, 0xb1, 0x4e, 0x85, 0x3d, 0x38, 0xb9,
	0x45, 0x17, 0xb3, 0x2c, 0x63, 0xb9, 0x2b, 0x90, 0xa7, 0x74, 0xe9, 0xca, 0xf5, 0x0a, 0x93, 0x28,
	0x43, 0x21, 0xe8, 0x02, 0x0b, 0xc2, 0xc1, 0x0f, 0x03, 0xf6, 0xc2, 0x8a, 0x68, 0xa4, 0x47, 0x91,
	0x77, 0x60, 0xe9, 0x62, 0xcc, 0x96, 0x7d, 0xe3, 0xc8, 0x18, 0xf6, 0xbc, 0x13, 0xe7, 0xce, 0xb9,
	0x4e, 0xcd, 0xe0, 0x97, 0x7d, 0x41, 0xcd, 0x40, 0x1e, 0xc1, 0x6e, 0x15, 0x47, 0x39, 0xcd, 0xb0,
	0x6f, 0x1e, 0x19, 0xc3, 0x76, 0xd0, 0xad, 0x92, 0x97, 0x34, 0x43, 0x72, 0x06, 0x96, 0x40, 0x29,
	0xd3, 0x7c, 0x21, 0xfa, 0x8d, 0x23, 0x63, 0xd8, 0xf1, 0x1e, 0x6f, 0x8e, 0x2c, 0x74, 0x38, 0x85,
	0x0e, 0x27, 0x54, 0x3a, 0x2e, 0x0a, 0x19, 0x41, 0xdd, 0x37, 0xf8, 0x6e, 0x42, 0x77, 0x2a, 0x39,
	0xd2, 0xac, 0xd4, 0xe1, 0xff, 0xbb, 0x8e, 0xb3, 0x46, 0xdf, 0xb8, 0x4b, 0xcb, 0xf6, 0x1f, 0xb4,
	0x7c, 0x02, 0x52, 0x53, 0x47, 0x1b, 0xaa, 0xcc, 0x61, 0xc7, 0x73, 0xee, 0xbb, 0x40, 0x21, 0x21,
	0xd8, 0xaf, 0x31, 0xd3, 0x92, 0x48, 0xed, 0x20, 0x30, 0xfe, 0xcc, 0x53, 0xb9, 0x8e, 0x94, 0xa3,
	0xd5, 0x3d, 0xab, 0xa4, 0xba, 0x0e, 0x99, 0xc2, 0x7e, 0x0d, 0xaa, 0x57, 0x68, 0xea, 0x15, 0xee,
	0x7b, 0x58, 0xbb, 0x22, 0xa8, 0x27, 0x87, 0xb0, 0x27, 0x58, 0x7c, 0x8d, 0x1b, 0xaa, 0x5a, 0xda,
	0xab, 0xa7, 0x7f, 0x51, 0x35, 0xd5, 0x5d, 0xa5, 0xa4, 0x5e, 0xc1, 0x51, 0xb1, 0x0e, 0xfe, 0x87,
	0x8e, 0xcf, 0xd9, 0x97, 0x75, 0x69, 0x9a, 0x0d, 0xa6, 0xa4, 0x0b, 0xed, 0x57, 0x3b, 0x50, 0xe1,
	0xe0, 0x9b, 0xf2, 0x75, 0x83, 0x81, 0x10, 0x68, 0x66, 0x94, 0x5f, 0x6b, 0xcc, 0x76, 0xa0, 0x63,
	0x72, 0x09, 0xa6, 0x9c, 0x33, 0xfd, 0xef, 0xf4, 0xbc, 0x97, 0x0f, 0xd8, 0xc7, 0x09, 0x47, 0xfe,
	0x6b, 0x2a, 0xe4, 0x64, 0x85, 0xf9, 0x54, 0x52, 0x89, 0x81, 0x22, 0x22, 0x97, 0xd0, 0x92, 0x2b,
	0xb5, 0x96, 0x3e, 0x6f, 0xcf, 0x7b, 0xf1, 0x20, 0x4a, 0x2d, 0xe8, 0x82, 0x25, 0x18, 0x94, 0x2c,
	0xe4, 0x14, 0x0e, 0x39, 0xc6, 0x98, 0xde, 0x60, 0xc4, 0x78, 0xba, 0x48, 0x73, 0xba, 0x8c, 0x12,
	0x14, 0x32, 0xa2, 0x49, 0xc2, 0x51, 0x28, 0x73, 0x8c, 0xa1, 0x15, 0x1c, 0x94, 0xa0, 0x49, 0x89,
	0x19, 0xa3, 0x90, 0xa7, 0x05, 0x82, 0x1c, 0x43, 0x77, 0x96, 0xe6, 0x49, 0xdd, 0xa1, 0xfe, 0xbd,
	0x6e, 0xd0, 0x51, 0xb9, 0x0a, 0xf2, 0x1f, 0xb4, 0x35, 0x44, 0xed, 0xa6, 0xbd, 0xd9, 0x0d, 0x2c,
	0x95, 0xf0, 0x19, 0x97, 0x83, 0xe7, 0x60, 0xdf, 0xd6, 0x4a, 0x2c, 0x68, 0x9e, 0x8a, 0x73, 0x61,
	0x6f, 0x11, 0x80, 0xd6, 0xab, 0x9c, 0xce, 0x96, 0x68, 0x1b, 0xa4, 0x03, 0x3b, 0xe3, 0x54, 0xe8,
	0x8f, 0xc6, 0xc0, 0x05, 0xf8, 0xa5, 0x87, 0xec, 0x80, 0x39, 0x99, 0xcf, 0x0b, 0x7c, 0x91, 0xb6,
	0x0d, 0xd2, 0x05, 0x2b, 0xc0, 0x24, 0xe5, 0x18, 0x4b, 0xbb, 0xf1, 0xe4, 0x0a, 0xf6, 0x7f, 0x7b,
	0x47, 0xaa, 0x2f, 0x1c, 0xf9, 0xf6, 0x96, 0x0a, 0xde, 0x8f, 0x7d, 0xdb, 0x50, 0xa3, 0x2f, 0xde,
	0x8e, 0x7c, 0xbb, 0x41, 0x76, 0xa1, 0xfd, 0x11, 0x67, 0xc5, 0x05, 0x6d, 0x53, 0x15, 0xde, 0x84,
	0xa1, 0x6f, 0x37, 0x89, 0x0d, 0xdd, 0x31, 0xcb, 0x68, 0x9a, 0x97, 0xb5, 0xed, 0xb3, 0x09, 0x1c,
	0xc7, 0x2c, 0xbb, 0xdb, 0x0b, 0xdf, 0xb8, 0xb2, 0xaa, 0xf8, 0x6b, 0xe3, 0xf0, 0x83, 0x17, 0xd0,
	0xb5, 0x33, 0x52, 0xd8, 0x7a, 0x2d, 0xe7, 0xbc, 0xac, 0xcf, 0x5a, 0xfa, 0xe9, 0x3e, 0xfb, 0x19,
	0x00, 0x00, 0xff, 0xff, 0xce, 0xe9, 0xc8, 0x20, 0x89, 0x05, 0x00, 0x00,
}
