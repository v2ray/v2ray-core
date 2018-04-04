package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_net "v2ray.com/core/common/net"
import v2ray_core_common_net1 "v2ray.com/core/common/net"
import v2ray_core_transport_internet "v2ray.com/core/transport/internet"
import v2ray_core_internet_domainsocket "v2ray.com/core/transport/internet/domainsocket"
import v2ray_core_common_serial "v2ray.com/core/common/serial"

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
func (KnownProtocols) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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
func (AllocationStrategy_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type InboundConfig struct {
}

func (m *InboundConfig) Reset()                    { *m = InboundConfig{} }
func (m *InboundConfig) String() string            { return proto.CompactTextString(m) }
func (*InboundConfig) ProtoMessage()               {}
func (*InboundConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type AllocationStrategy struct {
	Type AllocationStrategy_Type `protobuf:"varint,1,opt,name=type,enum=v2ray.core.app.proxyman.AllocationStrategy_Type" json:"type,omitempty"`
	// Number of handlers (ports) running in parallel.
	// Default value is 3 if unset.
	Concurrency *AllocationStrategy_AllocationStrategyConcurrency `protobuf:"bytes,2,opt,name=concurrency" json:"concurrency,omitempty"`
	// Number of minutes before a handler is regenerated.
	// Default value is 5 if unset.
	Refresh *AllocationStrategy_AllocationStrategyRefresh `protobuf:"bytes,3,opt,name=refresh" json:"refresh,omitempty"`
}

func (m *AllocationStrategy) Reset()                    { *m = AllocationStrategy{} }
func (m *AllocationStrategy) String() string            { return proto.CompactTextString(m) }
func (*AllocationStrategy) ProtoMessage()               {}
func (*AllocationStrategy) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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
	Value uint32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *AllocationStrategy_AllocationStrategyConcurrency) Reset() {
	*m = AllocationStrategy_AllocationStrategyConcurrency{}
}
func (m *AllocationStrategy_AllocationStrategyConcurrency) String() string {
	return proto.CompactTextString(m)
}
func (*AllocationStrategy_AllocationStrategyConcurrency) ProtoMessage() {}
func (*AllocationStrategy_AllocationStrategyConcurrency) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1, 0}
}

func (m *AllocationStrategy_AllocationStrategyConcurrency) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type AllocationStrategy_AllocationStrategyRefresh struct {
	Value uint32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *AllocationStrategy_AllocationStrategyRefresh) Reset() {
	*m = AllocationStrategy_AllocationStrategyRefresh{}
}
func (m *AllocationStrategy_AllocationStrategyRefresh) String() string {
	return proto.CompactTextString(m)
}
func (*AllocationStrategy_AllocationStrategyRefresh) ProtoMessage() {}
func (*AllocationStrategy_AllocationStrategyRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1, 1}
}

func (m *AllocationStrategy_AllocationStrategyRefresh) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type ReceiverConfig struct {
	PortRange                  *v2ray_core_common_net1.PortRange           `protobuf:"bytes,1,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
	Listen                     *v2ray_core_common_net.IPOrDomain           `protobuf:"bytes,2,opt,name=listen" json:"listen,omitempty"`
	AllocationStrategy         *AllocationStrategy                         `protobuf:"bytes,3,opt,name=allocation_strategy,json=allocationStrategy" json:"allocation_strategy,omitempty"`
	StreamSettings             *v2ray_core_transport_internet.StreamConfig `protobuf:"bytes,4,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	ReceiveOriginalDestination bool                                        `protobuf:"varint,5,opt,name=receive_original_destination,json=receiveOriginalDestination" json:"receive_original_destination,omitempty"`
	DomainOverride             []KnownProtocols                            `protobuf:"varint,7,rep,packed,name=domain_override,json=domainOverride,enum=v2ray.core.app.proxyman.KnownProtocols" json:"domain_override,omitempty"`
}

func (m *ReceiverConfig) Reset()                    { *m = ReceiverConfig{} }
func (m *ReceiverConfig) String() string            { return proto.CompactTextString(m) }
func (*ReceiverConfig) ProtoMessage()               {}
func (*ReceiverConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ReceiverConfig) GetPortRange() *v2ray_core_common_net1.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *ReceiverConfig) GetListen() *v2ray_core_common_net.IPOrDomain {
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

func (m *ReceiverConfig) GetStreamSettings() *v2ray_core_transport_internet.StreamConfig {
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

func (m *ReceiverConfig) GetDomainOverride() []KnownProtocols {
	if m != nil {
		return m.DomainOverride
	}
	return nil
}

type UnixReceiverConfig struct {
	DomainSockSettings *v2ray_core_internet_domainsocket.DomainSocketSettings `protobuf:"bytes,2,opt,name=domainSockSettings" json:"domainSockSettings,omitempty"`
	StreamSettings     *v2ray_core_transport_internet.StreamConfig            `protobuf:"bytes,4,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	DomainOverride     []KnownProtocols                                       `protobuf:"varint,7,rep,packed,name=domain_override,json=domainOverride,enum=v2ray.core.app.proxyman.KnownProtocols" json:"domain_override,omitempty"`
}

func (m *UnixReceiverConfig) Reset()                    { *m = UnixReceiverConfig{} }
func (m *UnixReceiverConfig) String() string            { return proto.CompactTextString(m) }
func (*UnixReceiverConfig) ProtoMessage()               {}
func (*UnixReceiverConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *UnixReceiverConfig) GetDomainSockSettings() *v2ray_core_internet_domainsocket.DomainSocketSettings {
	if m != nil {
		return m.DomainSockSettings
	}
	return nil
}

func (m *UnixReceiverConfig) GetStreamSettings() *v2ray_core_transport_internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

func (m *UnixReceiverConfig) GetDomainOverride() []KnownProtocols {
	if m != nil {
		return m.DomainOverride
	}
	return nil
}

type InboundHandlerConfig struct {
	Tag              string                                 `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	ReceiverSettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=receiver_settings,json=receiverSettings" json:"receiver_settings,omitempty"`
	ProxySettings    *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
}

func (m *InboundHandlerConfig) Reset()                    { *m = InboundHandlerConfig{} }
func (m *InboundHandlerConfig) String() string            { return proto.CompactTextString(m) }
func (*InboundHandlerConfig) ProtoMessage()               {}
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *InboundHandlerConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *InboundHandlerConfig) GetReceiverSettings() *v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.ReceiverSettings
	}
	return nil
}

func (m *InboundHandlerConfig) GetProxySettings() *v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

type OutboundConfig struct {
}

func (m *OutboundConfig) Reset()                    { *m = OutboundConfig{} }
func (m *OutboundConfig) String() string            { return proto.CompactTextString(m) }
func (*OutboundConfig) ProtoMessage()               {}
func (*OutboundConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type SenderConfig struct {
	// Send traffic through the given IP. Only IP is allowed.
	Via               *v2ray_core_common_net.IPOrDomain           `protobuf:"bytes,1,opt,name=via" json:"via,omitempty"`
	StreamSettings    *v2ray_core_transport_internet.StreamConfig `protobuf:"bytes,2,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	ProxySettings     *v2ray_core_transport_internet.ProxyConfig  `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	MultiplexSettings *MultiplexingConfig                         `protobuf:"bytes,4,opt,name=multiplex_settings,json=multiplexSettings" json:"multiplex_settings,omitempty"`
}

func (m *SenderConfig) Reset()                    { *m = SenderConfig{} }
func (m *SenderConfig) String() string            { return proto.CompactTextString(m) }
func (*SenderConfig) ProtoMessage()               {}
func (*SenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *SenderConfig) GetVia() *v2ray_core_common_net.IPOrDomain {
	if m != nil {
		return m.Via
	}
	return nil
}

func (m *SenderConfig) GetStreamSettings() *v2ray_core_transport_internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

func (m *SenderConfig) GetProxySettings() *v2ray_core_transport_internet.ProxyConfig {
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

type UnixSenderConfig struct {
	Via               *v2ray_core_internet_domainsocket.DomainSocketSettings `protobuf:"bytes,1,opt,name=via" json:"via,omitempty"`
	StreamSettings    *v2ray_core_transport_internet.StreamConfig            `protobuf:"bytes,2,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	ProxySettings     *v2ray_core_transport_internet.ProxyConfig             `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	MultiplexSettings *MultiplexingConfig                                    `protobuf:"bytes,4,opt,name=multiplex_settings,json=multiplexSettings" json:"multiplex_settings,omitempty"`
}

func (m *UnixSenderConfig) Reset()                    { *m = UnixSenderConfig{} }
func (m *UnixSenderConfig) String() string            { return proto.CompactTextString(m) }
func (*UnixSenderConfig) ProtoMessage()               {}
func (*UnixSenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UnixSenderConfig) GetVia() *v2ray_core_internet_domainsocket.DomainSocketSettings {
	if m != nil {
		return m.Via
	}
	return nil
}

func (m *UnixSenderConfig) GetStreamSettings() *v2ray_core_transport_internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

func (m *UnixSenderConfig) GetProxySettings() *v2ray_core_transport_internet.ProxyConfig {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

func (m *UnixSenderConfig) GetMultiplexSettings() *MultiplexingConfig {
	if m != nil {
		return m.MultiplexSettings
	}
	return nil
}

type MultiplexingConfig struct {
	// Whether or not Mux is enabled.
	Enabled bool `protobuf:"varint,1,opt,name=enabled" json:"enabled,omitempty"`
	// Max number of concurrent connections that one Mux connection can handle.
	Concurrency uint32 `protobuf:"varint,2,opt,name=concurrency" json:"concurrency,omitempty"`
}

func (m *MultiplexingConfig) Reset()                    { *m = MultiplexingConfig{} }
func (m *MultiplexingConfig) String() string            { return proto.CompactTextString(m) }
func (*MultiplexingConfig) ProtoMessage()               {}
func (*MultiplexingConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

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
	proto.RegisterType((*ReceiverConfig)(nil), "v2ray.core.app.proxyman.ReceiverConfig")
	proto.RegisterType((*UnixReceiverConfig)(nil), "v2ray.core.app.proxyman.UnixReceiverConfig")
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.app.proxyman.InboundHandlerConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterType((*SenderConfig)(nil), "v2ray.core.app.proxyman.SenderConfig")
	proto.RegisterType((*UnixSenderConfig)(nil), "v2ray.core.app.proxyman.UnixSenderConfig")
	proto.RegisterType((*MultiplexingConfig)(nil), "v2ray.core.app.proxyman.MultiplexingConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.KnownProtocols", KnownProtocols_name, KnownProtocols_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 868 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x56, 0x4d, 0x6f, 0x1b, 0x45,
	0x18, 0xee, 0x7a, 0x1d, 0xdb, 0x7d, 0xd3, 0x38, 0xdb, 0xa1, 0x52, 0x8d, 0x01, 0xc9, 0x18, 0x44,
	0xad, 0x82, 0xd6, 0xc5, 0x15, 0x48, 0x88, 0x03, 0x84, 0xa4, 0x52, 0x12, 0x88, 0x62, 0xc6, 0x86,
	0x43, 0x85, 0x64, 0x4d, 0x76, 0x27, 0x66, 0x94, 0xdd, 0x99, 0xd5, 0xcc, 0xd8, 0xcd, 0xfe, 0x25,
	0x7e, 0x03, 0x07, 0x8e, 0x1c, 0xf8, 0x05, 0xdc, 0xf8, 0x27, 0x68, 0x77, 0x66, 0x37, 0xfe, 0x6c,
	0x09, 0x55, 0x6f, 0xbd, 0xcd, 0xda, 0xcf, 0xf3, 0xcc, 0xfb, 0x3e, 0xef, 0xc7, 0x2e, 0xf4, 0xe6,
	0x03, 0x49, 0x52, 0x3f, 0x10, 0x71, 0x3f, 0x10, 0x92, 0xf6, 0x49, 0x92, 0xf4, 0x13, 0x29, 0xae,
	0xd3, 0x98, 0xf0, 0x7e, 0x20, 0xf8, 0x25, 0x9b, 0xfa, 0x89, 0x14, 0x5a, 0xa0, 0x87, 0x05, 0x52,
	0x52, 0x9f, 0x24, 0x89, 0x5f, 0xa0, 0xda, 0x8f, 0x56, 0x24, 0x02, 0x11, 0xc7, 0x82, 0xf7, 0x39,
	0xd5, 0x7d, 0x12, 0x86, 0x92, 0x2a, 0x65, 0x14, 0xda, 0x1f, 0x6f, 0x07, 0x26, 0x42, 0x6a, 0x8b,
	0xf2, 0x57, 0x50, 0x5a, 0x12, 0xae, 0xb2, 0xff, 0xfb, 0x8c, 0x6b, 0x2a, 0x33, 0xf4, 0x62, 0x5c,
	0xed, 0xaf, 0x5f, 0x8d, 0x0f, 0x45, 0x4c, 0x18, 0x57, 0x22, 0xb8, 0x5a, 0x25, 0x3f, 0xd9, 0x1c,
	0x92, 0xa2, 0x92, 0x91, 0xa8, 0xaf, 0xd3, 0x84, 0x86, 0x93, 0x98, 0x2a, 0x45, 0xa6, 0xd4, 0x30,
	0xba, 0xfb, 0xb0, 0x77, 0xc2, 0x2f, 0xc4, 0x8c, 0x87, 0x87, 0xb9, 0x50, 0xf7, 0x0f, 0x17, 0xd0,
	0x41, 0x14, 0x89, 0x80, 0x68, 0x26, 0xf8, 0x48, 0x4b, 0xa2, 0xe9, 0x34, 0x45, 0x47, 0x50, 0xcd,
	0xe8, 0x2d, 0xa7, 0xe3, 0xf4, 0x9a, 0x83, 0x27, 0xfe, 0x16, 0xf7, 0xfc, 0x75, 0xaa, 0x3f, 0x4e,
	0x13, 0x8a, 0x73, 0x36, 0xba, 0x82, 0xdd, 0x40, 0xf0, 0x60, 0x26, 0x25, 0xe5, 0x41, 0xda, 0xaa,
	0x74, 0x9c, 0xde, 0xee, 0xe0, 0xe4, 0x36, 0x62, 0xeb, 0x3f, 0x1d, 0xde, 0x08, 0xe2, 0x45, 0x75,
	0x34, 0x81, 0xba, 0xa4, 0x97, 0x92, 0xaa, 0x5f, 0x5b, 0x6e, 0x7e, 0xd1, 0xb3, 0xd7, 0xbb, 0x08,
	0x1b, 0x31, 0x5c, 0xa8, 0xb6, 0xbf, 0x80, 0x0f, 0x5e, 0x1a, 0x0e, 0x7a, 0x00, 0x3b, 0x73, 0x12,
	0xcd, 0x8c, 0x6b, 0x7b, 0xd8, 0x3c, 0xb4, 0x3f, 0x87, 0x77, 0xb7, 0x8a, 0x6f, 0xa6, 0x74, 0x3f,
	0x83, 0x6a, 0xe6, 0x22, 0x02, 0xa8, 0x1d, 0x44, 0x2f, 0x48, 0xaa, 0xbc, 0x3b, 0xd9, 0x19, 0x13,
	0x1e, 0x8a, 0xd8, 0x73, 0xd0, 0x3d, 0x68, 0x3c, 0xbb, 0xce, 0xba, 0x83, 0x44, 0x5e, 0xa5, 0xfb,
	0xb7, 0x0b, 0x4d, 0x4c, 0x03, 0xca, 0xe6, 0x54, 0x9a, 0xaa, 0xa2, 0x6f, 0x00, 0xb2, 0x1e, 0x9a,
	0x48, 0xc2, 0xa7, 0x46, 0x7b, 0x77, 0xd0, 0x59, 0xb4, 0xc3, 0x74, 0x8a, 0xcf, 0xa9, 0xf6, 0x87,
	0x42, 0x6a, 0x9c, 0xe1, 0xf0, 0xdd, 0xa4, 0x38, 0xa2, 0xaf, 0xa0, 0x16, 0x31, 0xa5, 0x29, 0xb7,
	0x45, 0xfb, 0x70, 0x0b, 0xf9, 0x64, 0x78, 0x2e, 0x8f, 0xf2, 0xfe, 0xc4, 0x96, 0x80, 0x7e, 0x81,
	0x77, 0x48, 0x99, 0xef, 0x44, 0xd9, 0x84, 0x6d, 0x4d, 0x3e, 0xbd, 0x45, 0x4d, 0x30, 0x22, 0xeb,
	0x8d, 0x39, 0x86, 0x7d, 0xa5, 0x25, 0x25, 0xf1, 0x44, 0x51, 0xad, 0x19, 0x9f, 0xaa, 0x56, 0x75,
	0x5d, 0xb9, 0x9c, 0x22, 0xbf, 0x98, 0x22, 0x7f, 0x94, 0xb3, 0x8c, 0x3f, 0xb8, 0x69, 0x34, 0x46,
	0x56, 0x02, 0x7d, 0x0b, 0xef, 0x4b, 0xe3, 0xe0, 0x44, 0x48, 0x36, 0x65, 0x9c, 0x44, 0x93, 0x90,
	0x2a, 0xcd, 0x78, 0x7e, 0x7b, 0x6b, 0xa7, 0xe3, 0xf4, 0x1a, 0xb8, 0x6d, 0x31, 0xe7, 0x16, 0x72,
	0x74, 0x83, 0x40, 0x43, 0xd8, 0x37, 0x73, 0x3a, 0x11, 0x73, 0x2a, 0x25, 0x0b, 0x69, 0xab, 0xde,
	0x71, 0x7b, 0xcd, 0xc1, 0xa3, 0xad, 0x19, 0x7f, 0xcf, 0xc5, 0x0b, 0x3e, 0xcc, 0xc6, 0x32, 0x10,
	0x91, 0xc2, 0x4d, 0xc3, 0x3f, 0xb7, 0xf4, 0xd3, 0x6a, 0xa3, 0xe6, 0xd5, 0xbb, 0xbf, 0x57, 0x00,
	0xfd, 0xc4, 0xd9, 0xf5, 0x4a, 0x81, 0x2f, 0x01, 0x19, 0xf8, 0x48, 0x04, 0x57, 0x45, 0x1a, 0xb6,
	0x56, 0x5f, 0x2e, 0xde, 0x58, 0xe6, 0xbf, 0xb8, 0x45, 0xfc, 0xa3, 0x92, 0x4b, 0x75, 0xc1, 0xc6,
	0x1b, 0x14, 0xdf, 0x90, 0xdd, 0x6f, 0xc2, 0x2c, 0xc7, 0xab, 0x9c, 0x56, 0x1b, 0xae, 0x57, 0x3d,
	0xad, 0x36, 0x76, 0xbc, 0x9a, 0xb5, 0xef, 0x2f, 0x07, 0x1e, 0xd8, 0x85, 0x77, 0x4c, 0x78, 0x18,
	0x95, 0x06, 0x7a, 0xe0, 0x6a, 0x32, 0xcd, 0x47, 0xe3, 0x2e, 0xce, 0x8e, 0x68, 0x04, 0xf7, 0x6d,
	0x7d, 0xe5, 0x4d, 0xb2, 0xc6, 0xd1, 0x4f, 0x36, 0x74, 0xbf, 0x59, 0xb2, 0xf9, 0xb6, 0x0b, 0xcf,
	0xcc, 0x8e, 0xc5, 0x5e, 0x21, 0x50, 0x66, 0x7a, 0x06, 0xcd, 0x3c, 0x85, 0x1b, 0x45, 0xf7, 0x56,
	0x8a, 0x7b, 0x39, 0xbb, 0x90, 0xeb, 0x7a, 0xd0, 0x3c, 0x9f, 0xe9, 0xc5, 0xfd, 0xfd, 0x67, 0x05,
	0xee, 0x8d, 0x28, 0x0f, 0xcb, 0xc4, 0x9e, 0x82, 0x3b, 0x67, 0xc4, 0xce, 0xfc, 0x7f, 0x18, 0xdb,
	0x0c, 0xbd, 0xa9, 0xcc, 0x95, 0xd7, 0x2f, 0xf3, 0x8f, 0x5b, 0x92, 0x7f, 0xfc, 0x0a, 0xd1, 0x61,
	0x46, 0xb2, 0x9a, 0xcb, 0x06, 0xa0, 0xe7, 0x80, 0xe2, 0x59, 0xa4, 0x59, 0x12, 0xd1, 0xeb, 0x97,
	0xb6, 0xe4, 0x52, 0xf3, 0x9c, 0x15, 0x14, 0xc6, 0xa7, 0x56, 0xf7, 0x7e, 0x29, 0x53, 0x9a, 0xfb,
	0x4f, 0x05, 0xbc, 0x6c, 0xd4, 0x96, 0xec, 0x3c, 0x5e, 0xb4, 0xf3, 0xff, 0x4e, 0xd6, 0x5b, 0x8f,
	0x0b, 0x8f, 0x87, 0x80, 0xd6, 0x81, 0xa8, 0x05, 0x75, 0xca, 0xc9, 0x45, 0x44, 0xc3, 0xdc, 0xe8,
	0x06, 0x2e, 0x1e, 0x51, 0x67, 0xfd, 0x0b, 0x62, 0x6f, 0xe9, 0xb5, 0xff, 0xf8, 0x23, 0x68, 0x2e,
	0xef, 0x06, 0xd4, 0x80, 0xea, 0xf1, 0x78, 0x3c, 0xf4, 0xee, 0xa0, 0x3a, 0xb8, 0xe3, 0x1f, 0x46,
	0x9e, 0xf3, 0xdd, 0x21, 0xbc, 0x17, 0x88, 0x78, 0x5b, 0xec, 0x43, 0xe7, 0x79, 0xa3, 0x38, 0xff,
	0x56, 0x79, 0xf8, 0xf3, 0x00, 0x93, 0xd4, 0x3f, 0xcc, 0x50, 0x07, 0x49, 0x62, 0x9c, 0x8a, 0x09,
	0xbf, 0xa8, 0xe5, 0x9f, 0x50, 0x4f, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x04, 0x59, 0x29, 0xf6,
	0x75, 0x0a, 0x00, 0x00,
}
