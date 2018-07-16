package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import net "v2ray.com/core/common/net"
import serial "v2ray.com/core/common/serial"
import internet "v2ray.com/core/transport/internet"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type KnownProtocols int32

const (
	KnownProtocols_HTTP KnownProtocols = 0
	KnownProtocols_TLS  KnownProtocols = 1
)

var KnownProtocols_name = map[int32]string{
	0: "HTTP",
	1: "TLS",
}
var KnownProtocols_value = map[string]int32{
	"HTTP": 0,
	"TLS":  1,
}

func (x KnownProtocols) String() string {
	return proto.EnumName(KnownProtocols_name, int32(x))
}
func (KnownProtocols) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{0}
}

type AllocationStrategy_Type int32

const (
	// Always allocate all connection handlers.
	AllocationStrategy_Always AllocationStrategy_Type = 0
	// Randomly allocate specific range of handlers.
	AllocationStrategy_Random AllocationStrategy_Type = 1
	// External. Not supported yet.
	AllocationStrategy_External AllocationStrategy_Type = 2
)

var AllocationStrategy_Type_name = map[int32]string{
	0: "Always",
	1: "Random",
	2: "External",
}
var AllocationStrategy_Type_value = map[string]int32{
	"Always":   0,
	"Random":   1,
	"External": 2,
}

func (x AllocationStrategy_Type) String() string {
	return proto.EnumName(AllocationStrategy_Type_name, int32(x))
}
func (AllocationStrategy_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{1, 0}
}

type InboundConfig struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InboundConfig) Reset()         { *m = InboundConfig{} }
func (m *InboundConfig) String() string { return proto.CompactTextString(m) }
func (*InboundConfig) ProtoMessage()    {}
func (*InboundConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{0}
}
func (m *InboundConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InboundConfig.Unmarshal(m, b)
}
func (m *InboundConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InboundConfig.Marshal(b, m, deterministic)
}
func (dst *InboundConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InboundConfig.Merge(dst, src)
}
func (m *InboundConfig) XXX_Size() int {
	return xxx_messageInfo_InboundConfig.Size(m)
}
func (m *InboundConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_InboundConfig.DiscardUnknown(m)
}

var xxx_messageInfo_InboundConfig proto.InternalMessageInfo

type AllocationStrategy struct {
	Type AllocationStrategy_Type `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.proxyman.AllocationStrategy_Type" json:"type,omitempty"`
	// Number of handlers (ports) running in parallel.
	// Default value is 3 if unset.
	Concurrency *AllocationStrategy_AllocationStrategyConcurrency `protobuf:"bytes,2,opt,name=concurrency,proto3" json:"concurrency,omitempty"`
	// Number of minutes before a handler is regenerated.
	// Default value is 5 if unset.
	Refresh              *AllocationStrategy_AllocationStrategyRefresh `protobuf:"bytes,3,opt,name=refresh,proto3" json:"refresh,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *AllocationStrategy) Reset()         { *m = AllocationStrategy{} }
func (m *AllocationStrategy) String() string { return proto.CompactTextString(m) }
func (*AllocationStrategy) ProtoMessage()    {}
func (*AllocationStrategy) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{1}
}
func (m *AllocationStrategy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AllocationStrategy.Unmarshal(m, b)
}
func (m *AllocationStrategy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AllocationStrategy.Marshal(b, m, deterministic)
}
func (dst *AllocationStrategy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllocationStrategy.Merge(dst, src)
}
func (m *AllocationStrategy) XXX_Size() int {
	return xxx_messageInfo_AllocationStrategy.Size(m)
}
func (m *AllocationStrategy) XXX_DiscardUnknown() {
	xxx_messageInfo_AllocationStrategy.DiscardUnknown(m)
}

var xxx_messageInfo_AllocationStrategy proto.InternalMessageInfo

func (m *AllocationStrategy) GetType() AllocationStrategy_Type {
	if m != nil {
		return m.Type
	}
	return AllocationStrategy_Always
}

func (m *AllocationStrategy) GetConcurrency() *AllocationStrategy_AllocationStrategyConcurrency {
	if m != nil {
		return m.Concurrency
	}
	return nil
}

func (m *AllocationStrategy) GetRefresh() *AllocationStrategy_AllocationStrategyRefresh {
	if m != nil {
		return m.Refresh
	}
	return nil
}

type AllocationStrategy_AllocationStrategyConcurrency struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AllocationStrategy_AllocationStrategyConcurrency) Reset() {
	*m = AllocationStrategy_AllocationStrategyConcurrency{}
}
func (m *AllocationStrategy_AllocationStrategyConcurrency) String() string {
	return proto.CompactTextString(m)
}
func (*AllocationStrategy_AllocationStrategyConcurrency) ProtoMessage() {}
func (*AllocationStrategy_AllocationStrategyConcurrency) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{1, 0}
}
func (m *AllocationStrategy_AllocationStrategyConcurrency) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AllocationStrategy_AllocationStrategyConcurrency.Unmarshal(m, b)
}
func (m *AllocationStrategy_AllocationStrategyConcurrency) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AllocationStrategy_AllocationStrategyConcurrency.Marshal(b, m, deterministic)
}
func (dst *AllocationStrategy_AllocationStrategyConcurrency) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllocationStrategy_AllocationStrategyConcurrency.Merge(dst, src)
}
func (m *AllocationStrategy_AllocationStrategyConcurrency) XXX_Size() int {
	return xxx_messageInfo_AllocationStrategy_AllocationStrategyConcurrency.Size(m)
}
func (m *AllocationStrategy_AllocationStrategyConcurrency) XXX_DiscardUnknown() {
	xxx_messageInfo_AllocationStrategy_AllocationStrategyConcurrency.DiscardUnknown(m)
}

var xxx_messageInfo_AllocationStrategy_AllocationStrategyConcurrency proto.InternalMessageInfo

func (m *AllocationStrategy_AllocationStrategyConcurrency) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type AllocationStrategy_AllocationStrategyRefresh struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AllocationStrategy_AllocationStrategyRefresh) Reset() {
	*m = AllocationStrategy_AllocationStrategyRefresh{}
}
func (m *AllocationStrategy_AllocationStrategyRefresh) String() string {
	return proto.CompactTextString(m)
}
func (*AllocationStrategy_AllocationStrategyRefresh) ProtoMessage() {}
func (*AllocationStrategy_AllocationStrategyRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{1, 1}
}
func (m *AllocationStrategy_AllocationStrategyRefresh) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AllocationStrategy_AllocationStrategyRefresh.Unmarshal(m, b)
}
func (m *AllocationStrategy_AllocationStrategyRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AllocationStrategy_AllocationStrategyRefresh.Marshal(b, m, deterministic)
}
func (dst *AllocationStrategy_AllocationStrategyRefresh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllocationStrategy_AllocationStrategyRefresh.Merge(dst, src)
}
func (m *AllocationStrategy_AllocationStrategyRefresh) XXX_Size() int {
	return xxx_messageInfo_AllocationStrategy_AllocationStrategyRefresh.Size(m)
}
func (m *AllocationStrategy_AllocationStrategyRefresh) XXX_DiscardUnknown() {
	xxx_messageInfo_AllocationStrategy_AllocationStrategyRefresh.DiscardUnknown(m)
}

var xxx_messageInfo_AllocationStrategy_AllocationStrategyRefresh proto.InternalMessageInfo

func (m *AllocationStrategy_AllocationStrategyRefresh) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type SniffingConfig struct {
	// Whether or not to enable content sniffing on an inbound connection.
	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// Override target domain if sniff'ed protocol is in the given list.
	// Supported values are "http", "tls".
	DomainOverride       []string `protobuf:"bytes,2,rep,name=domain_override,json=domainOverride,proto3" json:"domain_override,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SniffingConfig) Reset()         { *m = SniffingConfig{} }
func (m *SniffingConfig) String() string { return proto.CompactTextString(m) }
func (*SniffingConfig) ProtoMessage()    {}
func (*SniffingConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{2}
}
func (m *SniffingConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SniffingConfig.Unmarshal(m, b)
}
func (m *SniffingConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SniffingConfig.Marshal(b, m, deterministic)
}
func (dst *SniffingConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SniffingConfig.Merge(dst, src)
}
func (m *SniffingConfig) XXX_Size() int {
	return xxx_messageInfo_SniffingConfig.Size(m)
}
func (m *SniffingConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SniffingConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SniffingConfig proto.InternalMessageInfo

func (m *SniffingConfig) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *SniffingConfig) GetDomainOverride() []string {
	if m != nil {
		return m.DomainOverride
	}
	return nil
}

type ReceiverConfig struct {
	// PortRange specifies the ports which the Receiver should listen on.
	PortRange *net.PortRange `protobuf:"bytes,1,opt,name=port_range,json=portRange,proto3" json:"port_range,omitempty"`
	// Listen specifies the IP address that the Receiver should listen on.
	Listen                     *net.IPOrDomain        `protobuf:"bytes,2,opt,name=listen,proto3" json:"listen,omitempty"`
	AllocationStrategy         *AllocationStrategy    `protobuf:"bytes,3,opt,name=allocation_strategy,json=allocationStrategy,proto3" json:"allocation_strategy,omitempty"`
	StreamSettings             *internet.StreamConfig `protobuf:"bytes,4,opt,name=stream_settings,json=streamSettings,proto3" json:"stream_settings,omitempty"`
	ReceiveOriginalDestination bool                   `protobuf:"varint,5,opt,name=receive_original_destination,json=receiveOriginalDestination,proto3" json:"receive_original_destination,omitempty"`
	// Override domains for the given protocol.
	// Deprecated. Use sniffing_settings.
	DomainOverride       []KnownProtocols `protobuf:"varint,7,rep,packed,name=domain_override,json=domainOverride,proto3,enum=v2ray.core.app.proxyman.KnownProtocols" json:"domain_override,omitempty"` // Deprecated: Do not use.
	SniffingSettings     *SniffingConfig  `protobuf:"bytes,8,opt,name=sniffing_settings,json=sniffingSettings,proto3" json:"sniffing_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ReceiverConfig) Reset()         { *m = ReceiverConfig{} }
func (m *ReceiverConfig) String() string { return proto.CompactTextString(m) }
func (*ReceiverConfig) ProtoMessage()    {}
func (*ReceiverConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{3}
}
func (m *ReceiverConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiverConfig.Unmarshal(m, b)
}
func (m *ReceiverConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiverConfig.Marshal(b, m, deterministic)
}
func (dst *ReceiverConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiverConfig.Merge(dst, src)
}
func (m *ReceiverConfig) XXX_Size() int {
	return xxx_messageInfo_ReceiverConfig.Size(m)
}
func (m *ReceiverConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiverConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiverConfig proto.InternalMessageInfo

func (m *ReceiverConfig) GetPortRange() *net.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *ReceiverConfig) GetListen() *net.IPOrDomain {
	if m != nil {
		return m.Listen
	}
	return nil
}

func (m *ReceiverConfig) GetAllocationStrategy() *AllocationStrategy {
	if m != nil {
		return m.AllocationStrategy
	}
	return nil
}

func (m *ReceiverConfig) GetStreamSettings() *internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

func (m *ReceiverConfig) GetReceiveOriginalDestination() bool {
	if m != nil {
		return m.ReceiveOriginalDestination
	}
	return false
}

// Deprecated: Do not use.
func (m *ReceiverConfig) GetDomainOverride() []KnownProtocols {
	if m != nil {
		return m.DomainOverride
	}
	return nil
}

func (m *ReceiverConfig) GetSniffingSettings() *SniffingConfig {
	if m != nil {
		return m.SniffingSettings
	}
	return nil
}

type InboundHandlerConfig struct {
	Tag                  string               `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	ReceiverSettings     *serial.TypedMessage `protobuf:"bytes,2,opt,name=receiver_settings,json=receiverSettings,proto3" json:"receiver_settings,omitempty"`
	ProxySettings        *serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings,proto3" json:"proxy_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InboundHandlerConfig) Reset()         { *m = InboundHandlerConfig{} }
func (m *InboundHandlerConfig) String() string { return proto.CompactTextString(m) }
func (*InboundHandlerConfig) ProtoMessage()    {}
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{4}
}
func (m *InboundHandlerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InboundHandlerConfig.Unmarshal(m, b)
}
func (m *InboundHandlerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InboundHandlerConfig.Marshal(b, m, deterministic)
}
func (dst *InboundHandlerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InboundHandlerConfig.Merge(dst, src)
}
func (m *InboundHandlerConfig) XXX_Size() int {
	return xxx_messageInfo_InboundHandlerConfig.Size(m)
}
func (m *InboundHandlerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_InboundHandlerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_InboundHandlerConfig proto.InternalMessageInfo

func (m *InboundHandlerConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *InboundHandlerConfig) GetReceiverSettings() *serial.TypedMessage {
	if m != nil {
		return m.ReceiverSettings
	}
	return nil
}

func (m *InboundHandlerConfig) GetProxySettings() *serial.TypedMessage {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

type OutboundConfig struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OutboundConfig) Reset()         { *m = OutboundConfig{} }
func (m *OutboundConfig) String() string { return proto.CompactTextString(m) }
func (*OutboundConfig) ProtoMessage()    {}
func (*OutboundConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{5}
}
func (m *OutboundConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutboundConfig.Unmarshal(m, b)
}
func (m *OutboundConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutboundConfig.Marshal(b, m, deterministic)
}
func (dst *OutboundConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutboundConfig.Merge(dst, src)
}
func (m *OutboundConfig) XXX_Size() int {
	return xxx_messageInfo_OutboundConfig.Size(m)
}
func (m *OutboundConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_OutboundConfig.DiscardUnknown(m)
}

var xxx_messageInfo_OutboundConfig proto.InternalMessageInfo

type SenderConfig struct {
	// Send traffic through the given IP. Only IP is allowed.
	Via                  *net.IPOrDomain        `protobuf:"bytes,1,opt,name=via,proto3" json:"via,omitempty"`
	StreamSettings       *internet.StreamConfig `protobuf:"bytes,2,opt,name=stream_settings,json=streamSettings,proto3" json:"stream_settings,omitempty"`
	ProxySettings        *internet.ProxyConfig  `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings,proto3" json:"proxy_settings,omitempty"`
	MultiplexSettings    *MultiplexingConfig    `protobuf:"bytes,4,opt,name=multiplex_settings,json=multiplexSettings,proto3" json:"multiplex_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *SenderConfig) Reset()         { *m = SenderConfig{} }
func (m *SenderConfig) String() string { return proto.CompactTextString(m) }
func (*SenderConfig) ProtoMessage()    {}
func (*SenderConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{6}
}
func (m *SenderConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SenderConfig.Unmarshal(m, b)
}
func (m *SenderConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SenderConfig.Marshal(b, m, deterministic)
}
func (dst *SenderConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SenderConfig.Merge(dst, src)
}
func (m *SenderConfig) XXX_Size() int {
	return xxx_messageInfo_SenderConfig.Size(m)
}
func (m *SenderConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SenderConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SenderConfig proto.InternalMessageInfo

func (m *SenderConfig) GetVia() *net.IPOrDomain {
	if m != nil {
		return m.Via
	}
	return nil
}

func (m *SenderConfig) GetStreamSettings() *internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

func (m *SenderConfig) GetProxySettings() *internet.ProxyConfig {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

func (m *SenderConfig) GetMultiplexSettings() *MultiplexingConfig {
	if m != nil {
		return m.MultiplexSettings
	}
	return nil
}

type MultiplexingConfig struct {
	// Whether or not Mux is enabled.
	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// Max number of concurrent connections that one Mux connection can handle.
	Concurrency          uint32   `protobuf:"varint,2,opt,name=concurrency,proto3" json:"concurrency,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MultiplexingConfig) Reset()         { *m = MultiplexingConfig{} }
func (m *MultiplexingConfig) String() string { return proto.CompactTextString(m) }
func (*MultiplexingConfig) ProtoMessage()    {}
func (*MultiplexingConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_bab73170eded8a95, []int{7}
}
func (m *MultiplexingConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MultiplexingConfig.Unmarshal(m, b)
}
func (m *MultiplexingConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MultiplexingConfig.Marshal(b, m, deterministic)
}
func (dst *MultiplexingConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiplexingConfig.Merge(dst, src)
}
func (m *MultiplexingConfig) XXX_Size() int {
	return xxx_messageInfo_MultiplexingConfig.Size(m)
}
func (m *MultiplexingConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiplexingConfig.DiscardUnknown(m)
}

var xxx_messageInfo_MultiplexingConfig proto.InternalMessageInfo

func (m *MultiplexingConfig) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *MultiplexingConfig) GetConcurrency() uint32 {
	if m != nil {
		return m.Concurrency
	}
	return 0
}

func init() {
	proto.RegisterType((*InboundConfig)(nil), "v2ray.core.app.proxyman.InboundConfig")
	proto.RegisterType((*AllocationStrategy)(nil), "v2ray.core.app.proxyman.AllocationStrategy")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyConcurrency)(nil), "v2ray.core.app.proxyman.AllocationStrategy.AllocationStrategyConcurrency")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyRefresh)(nil), "v2ray.core.app.proxyman.AllocationStrategy.AllocationStrategyRefresh")
	proto.RegisterType((*SniffingConfig)(nil), "v2ray.core.app.proxyman.SniffingConfig")
	proto.RegisterType((*ReceiverConfig)(nil), "v2ray.core.app.proxyman.ReceiverConfig")
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.app.proxyman.InboundHandlerConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterType((*SenderConfig)(nil), "v2ray.core.app.proxyman.SenderConfig")
	proto.RegisterType((*MultiplexingConfig)(nil), "v2ray.core.app.proxyman.MultiplexingConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.KnownProtocols", KnownProtocols_name, KnownProtocols_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor_config_bab73170eded8a95)
}

var fileDescriptor_config_bab73170eded8a95 = []byte{
	// 819 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0x41, 0x6f, 0xe3, 0x44,
	0x14, 0xc7, 0x37, 0x71, 0xb6, 0x49, 0x5f, 0xb7, 0xae, 0x3b, 0xac, 0xb4, 0x26, 0x80, 0x14, 0x02,
	0xa2, 0xd1, 0x82, 0xec, 0x25, 0x2b, 0x0e, 0x9c, 0xa0, 0xdb, 0xae, 0xb4, 0x05, 0xaa, 0x86, 0x71,
	0xc4, 0x61, 0x85, 0x64, 0x4d, 0xed, 0xa9, 0x19, 0x61, 0xcf, 0x58, 0x33, 0x93, 0x6c, 0xfd, 0x95,
	0x38, 0xf3, 0x01, 0x38, 0x72, 0xe0, 0x43, 0x21, 0x7b, 0xec, 0xa4, 0x69, 0xea, 0x42, 0xb5, 0xb7,
	0x89, 0xf3, 0x7f, 0x3f, 0xbf, 0xf7, 0x9f, 0xff, 0x8c, 0x61, 0xb2, 0x9c, 0x4a, 0x52, 0x78, 0x91,
	0xc8, 0xfc, 0x48, 0x48, 0xea, 0x93, 0x3c, 0xf7, 0x73, 0x29, 0xae, 0x8b, 0x8c, 0x70, 0x3f, 0x12,
	0xfc, 0x8a, 0x25, 0x5e, 0x2e, 0x85, 0x16, 0xe8, 0x59, 0xa3, 0x94, 0xd4, 0x23, 0x79, 0xee, 0x35,
	0xaa, 0xe1, 0xd1, 0x2d, 0x44, 0x24, 0xb2, 0x4c, 0x70, 0x9f, 0x53, 0xed, 0x93, 0x38, 0x96, 0x54,
	0x29, 0x43, 0x18, 0x7e, 0xde, 0x2e, 0xcc, 0x85, 0xd4, 0xb5, 0xca, 0xbb, 0xa5, 0xd2, 0x92, 0x70,
	0x55, 0xfe, 0xef, 0x33, 0xae, 0xa9, 0x2c, 0xd5, 0x37, 0xfb, 0x1a, 0xbe, 0xb8, 0x9b, 0xaa, 0xa8,
	0x64, 0x24, 0xf5, 0x75, 0x91, 0xd3, 0x38, 0xcc, 0xa8, 0x52, 0x24, 0xa1, 0xa6, 0x62, 0x7c, 0x00,
	0xfb, 0x67, 0xfc, 0x52, 0x2c, 0x78, 0x7c, 0x52, 0x81, 0xc6, 0x7f, 0x59, 0x80, 0x8e, 0xd3, 0x54,
	0x44, 0x44, 0x33, 0xc1, 0x03, 0x2d, 0x89, 0xa6, 0x49, 0x81, 0x4e, 0xa1, 0x57, 0x96, 0xbb, 0x9d,
	0x51, 0x67, 0x62, 0x4f, 0x5f, 0x78, 0x2d, 0x06, 0x78, 0xdb, 0xa5, 0xde, 0xbc, 0xc8, 0x29, 0xae,
	0xaa, 0xd1, 0xef, 0xb0, 0x17, 0x09, 0x1e, 0x2d, 0xa4, 0xa4, 0x3c, 0x2a, 0xdc, 0xee, 0xa8, 0x33,
	0xd9, 0x9b, 0x9e, 0x3d, 0x04, 0xb6, 0xfd, 0xe8, 0x64, 0x0d, 0xc4, 0x37, 0xe9, 0x28, 0x84, 0xbe,
	0xa4, 0x57, 0x92, 0xaa, 0xdf, 0x5c, 0xab, 0x7a, 0xd1, 0xeb, 0xf7, 0x7b, 0x11, 0x36, 0x30, 0xdc,
	0x50, 0x87, 0xdf, 0xc0, 0x27, 0xf7, 0xb6, 0x83, 0x9e, 0xc2, 0xe3, 0x25, 0x49, 0x17, 0xc6, 0xb5,
	0x7d, 0x6c, 0x7e, 0x0c, 0xbf, 0x86, 0x0f, 0x5b, 0xe1, 0x77, 0x97, 0x8c, 0xbf, 0x82, 0x5e, 0xe9,
	0x22, 0x02, 0xd8, 0x39, 0x4e, 0xdf, 0x91, 0x42, 0x39, 0x8f, 0xca, 0x35, 0x26, 0x3c, 0x16, 0x99,
	0xd3, 0x41, 0x4f, 0x60, 0xf0, 0xfa, 0xba, 0x0c, 0x04, 0x49, 0x9d, 0xee, 0x38, 0x00, 0x3b, 0xe0,
	0xec, 0xea, 0x8a, 0xf1, 0xc4, 0x6c, 0x2a, 0x72, 0xa1, 0x4f, 0x39, 0xb9, 0x4c, 0x69, 0x5c, 0x71,
	0x07, 0xb8, 0xf9, 0x89, 0x8e, 0xe0, 0x20, 0x16, 0x19, 0x61, 0x3c, 0x14, 0x4b, 0x2a, 0x25, 0x8b,
	0xa9, 0xdb, 0x1d, 0x59, 0x93, 0x5d, 0x6c, 0x9b, 0xc7, 0x17, 0xf5, 0xd3, 0xf1, 0x9f, 0x3d, 0xb0,
	0x31, 0x8d, 0x28, 0x5b, 0x52, 0x59, 0x53, 0xbf, 0x03, 0x28, 0xb3, 0x18, 0x4a, 0xc2, 0x13, 0xd3,
	0xf0, 0xde, 0x74, 0x74, 0xd3, 0x63, 0x13, 0x3f, 0x8f, 0x53, 0xed, 0xcd, 0x84, 0xd4, 0xb8, 0xd4,
	0xe1, 0xdd, 0xbc, 0x59, 0xa2, 0x6f, 0x61, 0x27, 0x65, 0x4a, 0x53, 0x5e, 0x27, 0xe1, 0xd3, 0x96,
	0xe2, 0xb3, 0xd9, 0x85, 0x3c, 0xad, 0xda, 0xc1, 0x75, 0x01, 0xfa, 0x15, 0x3e, 0x20, 0x2b, 0x13,
	0x43, 0x55, 0xbb, 0x58, 0x6f, 0xf4, 0x97, 0x0f, 0xd8, 0x68, 0x8c, 0xc8, 0x76, 0xda, 0xe7, 0x70,
	0xa0, 0xb4, 0xa4, 0x24, 0x0b, 0x15, 0xd5, 0x9a, 0xf1, 0x44, 0xb9, 0xbd, 0x6d, 0xf2, 0xea, 0x34,
	0x7a, 0xcd, 0x69, 0xf4, 0x82, 0xaa, 0xca, 0xf8, 0x83, 0x6d, 0xc3, 0x08, 0x6a, 0x04, 0xfa, 0x1e,
	0x3e, 0x96, 0xc6, 0xc1, 0x50, 0x48, 0x96, 0x30, 0x4e, 0xd2, 0x30, 0xa6, 0x4a, 0x33, 0x5e, 0xbd,
	0xdd, 0x7d, 0x5c, 0x6d, 0xcd, 0xb0, 0xd6, 0x5c, 0xd4, 0x92, 0xd3, 0xb5, 0xa2, 0xec, 0xeb, 0xf6,
	0x6e, 0xf5, 0x47, 0xd6, 0xc4, 0x9e, 0x1e, 0xb5, 0x4e, 0xfc, 0x23, 0x17, 0xef, 0xf8, 0xac, 0x3c,
	0xeb, 0x91, 0x48, 0xd5, 0xab, 0xae, 0xdb, 0xb9, 0xbd, 0xb5, 0x68, 0x0e, 0x87, 0xaa, 0xce, 0xcb,
	0x7a, 0xde, 0x41, 0x35, 0x6f, 0x3b, 0x77, 0x33, 0x61, 0xd8, 0x69, 0x08, 0xcd, 0xb4, 0x3f, 0xf4,
	0x06, 0x3b, 0x4e, 0x7f, 0xfc, 0x4f, 0x07, 0x9e, 0xd6, 0x17, 0xcc, 0x1b, 0xc2, 0xe3, 0x74, 0x15,
	0x1e, 0x07, 0x2c, 0x4d, 0x92, 0x2a, 0x35, 0xbb, 0xb8, 0x5c, 0xa2, 0x00, 0x0e, 0xeb, 0xd1, 0xe5,
	0xba, 0x0d, 0x13, 0x8c, 0x2f, 0xee, 0x08, 0x86, 0xb9, 0xd4, 0xaa, 0xdb, 0x25, 0x3e, 0x37, 0x77,
	0x1a, 0x76, 0x1a, 0xc0, 0xca, 0xf3, 0x73, 0xb0, 0xab, 0x96, 0xd7, 0x44, 0xeb, 0x41, 0xc4, 0xfd,
	0xaa, 0xba, 0xc1, 0x8d, 0x1d, 0xb0, 0x2f, 0x16, 0xfa, 0xe6, 0x7d, 0xf9, 0x77, 0x17, 0x9e, 0x04,
	0x94, 0xc7, 0xab, 0xc1, 0x5e, 0x82, 0xb5, 0x64, 0xa4, 0x3e, 0x0e, 0xff, 0x23, 0xd1, 0xa5, 0xfa,
	0xae, 0xc0, 0x75, 0xdf, 0x3f, 0x70, 0x3f, 0xb7, 0x0c, 0xff, 0xfc, 0x3f, 0xa0, 0xb3, 0xb2, 0xa8,
	0x66, 0x6e, 0x1a, 0x80, 0xde, 0x02, 0xca, 0x16, 0xa9, 0x66, 0x79, 0x4a, 0xaf, 0xef, 0x3d, 0x1c,
	0x1b, 0x61, 0x39, 0x6f, 0x4a, 0xd6, 0x81, 0x39, 0x5c, 0x61, 0x56, 0xe6, 0xce, 0x00, 0x6d, 0x0b,
	0xef, 0xb9, 0xbb, 0x46, 0xdb, 0x5f, 0x93, 0xfd, 0x8d, 0x4f, 0xc0, 0xf3, 0xcf, 0xc0, 0xde, 0xcc,
	0x3f, 0x1a, 0x40, 0xef, 0xcd, 0x7c, 0x3e, 0x73, 0x1e, 0xa1, 0x3e, 0x58, 0xf3, 0x9f, 0x02, 0xa7,
	0xf3, 0xea, 0x04, 0x3e, 0x8a, 0x44, 0xd6, 0xd6, 0xfb, 0xac, 0xf3, 0x76, 0xd0, 0xac, 0xff, 0xe8,
	0x3e, 0xfb, 0x65, 0x8a, 0x49, 0xe1, 0x9d, 0x94, 0xaa, 0xe3, 0x3c, 0x37, 0x4e, 0x65, 0x84, 0x5f,
	0xee, 0x54, 0x9f, 0xd3, 0x97, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x88, 0x8f, 0xff, 0x44,
	0x08, 0x00, 0x00,
}
