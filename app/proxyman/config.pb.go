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
	DomainSockSettings *v2ray_core_internet_domainsocket.DomainSocketSettings `protobuf:"bytes,1,opt,name=domainSockSettings" json:"domainSockSettings,omitempty"`
	StreamSettings     *v2ray_core_transport_internet.StreamConfig            `protobuf:"bytes,2,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	ProxySettings      *v2ray_core_transport_internet.ProxyConfig             `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	MultiplexSettings  *MultiplexingConfig                                    `protobuf:"bytes,4,opt,name=multiplex_settings,json=multiplexSettings" json:"multiplex_settings,omitempty"`
}

func (m *UnixSenderConfig) Reset()                    { *m = UnixSenderConfig{} }
func (m *UnixSenderConfig) String() string            { return proto.CompactTextString(m) }
func (*UnixSenderConfig) ProtoMessage()               {}
func (*UnixSenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UnixSenderConfig) GetDomainSockSettings() *v2ray_core_internet_domainsocket.DomainSocketSettings {
	if m != nil {
		return m.DomainSockSettings
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
	// 866 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe4, 0x56, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0xee, 0x7a, 0x1d, 0xdb, 0x3d, 0x69, 0x9c, 0xed, 0x50, 0xa9, 0xc6, 0x80, 0x64, 0x0c, 0xa2,
	0x56, 0x41, 0xeb, 0xe2, 0x0a, 0x24, 0xc4, 0x05, 0x84, 0xa4, 0x52, 0x13, 0x88, 0x62, 0xc6, 0x86,
	0x8b, 0x0a, 0xc9, 0x9a, 0xec, 0x4e, 0xcc, 0x28, 0xbb, 0x33, 0xab, 0x99, 0xb1, 0x9b, 0x7d, 0x25,
	0x9e, 0x81, 0x0b, 0x2e, 0xb9, 0xe0, 0x09, 0x78, 0x19, 0xd0, 0xee, 0xcc, 0x6e, 0xfc, 0xdb, 0x12,
	0x55, 0xe1, 0xa6, 0x77, 0xb3, 0xf6, 0xf7, 0x7d, 0x73, 0xce, 0x77, 0x7e, 0x34, 0xd0, 0x9b, 0x0f,
	0x24, 0x49, 0xfd, 0x40, 0xc4, 0xfd, 0x40, 0x48, 0xda, 0x27, 0x49, 0xd2, 0x4f, 0xa4, 0xb8, 0x4a,
	0x63, 0xc2, 0xfb, 0x81, 0xe0, 0x17, 0x6c, 0xea, 0x27, 0x52, 0x68, 0x81, 0x1e, 0x16, 0x48, 0x49,
	0x7d, 0x92, 0x24, 0x7e, 0x81, 0x6a, 0x3f, 0x5a, 0x91, 0x08, 0x44, 0x1c, 0x0b, 0xde, 0xe7, 0x54,
	0xf7, 0x49, 0x18, 0x4a, 0xaa, 0x94, 0x51, 0x68, 0x7f, 0xbc, 0x1d, 0x98, 0x08, 0xa9, 0x2d, 0xca,
	0x5f, 0x41, 0x69, 0x49, 0xb8, 0xca, 0xfe, 0xef, 0x33, 0xae, 0xa9, 0xcc, 0xd0, 0x8b, 0x71, 0xb5,
	0xbf, 0x7e, 0x3d, 0x3e, 0x14, 0x31, 0x61, 0x5c, 0x89, 0xe0, 0x72, 0x95, 0xfc, 0x64, 0x73, 0x48,
	0x8a, 0x4a, 0x46, 0xa2, 0xbe, 0x4e, 0x13, 0x1a, 0x4e, 0x62, 0xaa, 0x14, 0x99, 0x52, 0xc3, 0xe8,
	0xee, 0xc3, 0xde, 0x31, 0x3f, 0x17, 0x33, 0x1e, 0x1e, 0xe6, 0x42, 0xdd, 0x3f, 0x5c, 0x40, 0x07,
	0x51, 0x24, 0x02, 0xa2, 0x99, 0xe0, 0x23, 0x2d, 0x89, 0xa6, 0xd3, 0x14, 0x1d, 0x41, 0x35, 0xa3,
	0xb7, 0x9c, 0x8e, 0xd3, 0x6b, 0x0e, 0x9e, 0xf8, 0x5b, 0xdc, 0xf3, 0xd7, 0xa9, 0xfe, 0x38, 0x4d,
	0x28, 0xce, 0xd9, 0xe8, 0x12, 0x76, 0x03, 0xc1, 0x83, 0x99, 0x94, 0x94, 0x07, 0x69, 0xab, 0xd2,
	0x71, 0x7a, 0xbb, 0x83, 0xe3, 0x9b, 0x88, 0xad, 0xff, 0x74, 0x78, 0x2d, 0x88, 0x17, 0xd5, 0xd1,
	0x04, 0xea, 0x92, 0x5e, 0x48, 0xaa, 0x7e, 0x6d, 0xb9, 0xf9, 0x45, 0xcf, 0xde, 0xec, 0x22, 0x6c,
	0xc4, 0x70, 0xa1, 0xda, 0xfe, 0x02, 0x3e, 0x78, 0x65, 0x38, 0xe8, 0x01, 0xec, 0xcc, 0x49, 0x34,
	0x33, 0xae, 0xed, 0x61, 0xf3, 0xd1, 0xfe, 0x1c, 0xde, 0xdd, 0x2a, 0xbe, 0x99, 0xd2, 0xfd, 0x0c,
	0xaa, 0x99, 0x8b, 0x08, 0xa0, 0x76, 0x10, 0xbd, 0x24, 0xa9, 0xf2, 0xee, 0x64, 0x67, 0x4c, 0x78,
	0x28, 0x62, 0xcf, 0x41, 0xf7, 0xa0, 0xf1, 0xec, 0x2a, 0xeb, 0x0e, 0x12, 0x79, 0x95, 0xee, 0xdf,
	0x2e, 0x34, 0x31, 0x0d, 0x28, 0x9b, 0x53, 0x69, 0xaa, 0x8a, 0xbe, 0x01, 0xc8, 0x7a, 0x68, 0x22,
	0x09, 0x9f, 0x1a, 0xed, 0xdd, 0x41, 0x67, 0xd1, 0x0e, 0xd3, 0x29, 0x3e, 0xa7, 0xda, 0x1f, 0x0a,
	0xa9, 0x71, 0x86, 0xc3, 0x77, 0x93, 0xe2, 0x88, 0xbe, 0x82, 0x5a, 0xc4, 0x94, 0xa6, 0xdc, 0x16,
	0xed, 0xc3, 0x2d, 0xe4, 0xe3, 0xe1, 0x99, 0x3c, 0xca, 0xfb, 0x13, 0x5b, 0x02, 0xfa, 0x05, 0xde,
	0x21, 0x65, 0xbe, 0x13, 0x65, 0x13, 0xb6, 0x35, 0xf9, 0xf4, 0x06, 0x35, 0xc1, 0x88, 0xac, 0x37,
	0xe6, 0x18, 0xf6, 0x95, 0x96, 0x94, 0xc4, 0x13, 0x45, 0xb5, 0x66, 0x7c, 0xaa, 0x5a, 0xd5, 0x75,
	0xe5, 0x72, 0x8a, 0xfc, 0x62, 0x8a, 0xfc, 0x51, 0xce, 0x32, 0xfe, 0xe0, 0xa6, 0xd1, 0x18, 0x59,
	0x09, 0xf4, 0x2d, 0xbc, 0x2f, 0x8d, 0x83, 0x13, 0x21, 0xd9, 0x94, 0x71, 0x12, 0x4d, 0x42, 0xaa,
	0x34, 0xe3, 0xf9, 0xed, 0xad, 0x9d, 0x8e, 0xd3, 0x6b, 0xe0, 0xb6, 0xc5, 0x9c, 0x59, 0xc8, 0xd1,
	0x35, 0x02, 0x0d, 0x61, 0xdf, 0xcc, 0xe9, 0x44, 0xcc, 0xa9, 0x94, 0x2c, 0xa4, 0xad, 0x7a, 0xc7,
	0xed, 0x35, 0x07, 0x8f, 0xb6, 0x66, 0xfc, 0x3d, 0x17, 0x2f, 0xf9, 0x30, 0x1b, 0xcb, 0x40, 0x44,
	0x0a, 0x37, 0x0d, 0xff, 0xcc, 0xd2, 0x4f, 0xaa, 0x8d, 0x9a, 0x57, 0xef, 0xfe, 0x5e, 0x01, 0xf4,
	0x13, 0x67, 0x57, 0x2b, 0x05, 0xbe, 0x00, 0x64, 0xe0, 0x23, 0x11, 0x5c, 0x16, 0x69, 0xd8, 0x5a,
	0x7d, 0xb9, 0x78, 0x63, 0x99, 0xff, 0xe2, 0x16, 0xf1, 0x8f, 0x4a, 0x2e, 0xd5, 0x05, 0x1b, 0x6f,
	0x50, 0xbc, 0x25, 0xbb, 0x6f, 0xc3, 0x2c, 0xc7, 0xab, 0x9c, 0x54, 0x1b, 0xae, 0x57, 0x3d, 0xa9,
	0x36, 0x76, 0xbc, 0x9a, 0xb5, 0xef, 0x2f, 0x07, 0x1e, 0xd8, 0x85, 0xf7, 0x9c, 0xf0, 0x30, 0x2a,
	0x0d, 0xf4, 0xc0, 0xd5, 0x64, 0x9a, 0x8f, 0xc6, 0x5d, 0x9c, 0x1d, 0xd1, 0x08, 0xee, 0xdb, 0xfa,
	0xca, 0xeb, 0x64, 0x8d, 0xa3, 0x9f, 0x6c, 0xe8, 0x7e, 0xb3, 0x64, 0xf3, 0x6d, 0x17, 0x9e, 0x9a,
	0x1d, 0x8b, 0xbd, 0x42, 0xa0, 0xcc, 0xf4, 0x14, 0x9a, 0x79, 0x0a, 0xd7, 0x8a, 0xee, 0x8d, 0x14,
	0xf7, 0x72, 0x76, 0x21, 0xd7, 0xf5, 0xa0, 0x79, 0x36, 0xd3, 0x8b, 0xfb, 0xfb, 0xcf, 0x0a, 0xdc,
	0x1b, 0x51, 0x1e, 0x96, 0x89, 0x3d, 0x05, 0x77, 0xce, 0x88, 0x9d, 0xf9, 0xff, 0x30, 0xb6, 0x19,
	0x7a, 0x53, 0x99, 0x2b, 0x6f, 0x5e, 0xe6, 0x1f, 0xb7, 0x24, 0xff, 0xf8, 0x35, 0xa2, 0xc3, 0x8c,
	0x64, 0x35, 0x97, 0x0d, 0x40, 0x2f, 0x00, 0xc5, 0xb3, 0x48, 0xb3, 0x24, 0xa2, 0x57, 0xaf, 0x6c,
	0xc9, 0xa5, 0xe6, 0x39, 0x2d, 0x28, 0x8c, 0x4f, 0xad, 0xee, 0xfd, 0x52, 0xa6, 0x34, 0xf7, 0x9f,
	0x0a, 0x78, 0xd9, 0xa8, 0x2d, 0xd9, 0xb9, 0x79, 0xd0, 0x9c, 0xff, 0x63, 0xd0, 0xde, 0xbe, 0x0a,
	0x0c, 0x01, 0xad, 0x03, 0x51, 0x0b, 0xea, 0x94, 0x93, 0xf3, 0x88, 0x86, 0xb9, 0xef, 0x0d, 0x5c,
	0x7c, 0xa2, 0xce, 0xfa, 0xfb, 0x62, 0x6f, 0xe9, 0x51, 0xf0, 0xf8, 0x23, 0x68, 0x2e, 0x6f, 0x0e,
	0xd4, 0x80, 0xea, 0xf3, 0xf1, 0x78, 0xe8, 0xdd, 0x41, 0x75, 0x70, 0xc7, 0x3f, 0x8c, 0x3c, 0xe7,
	0xbb, 0x43, 0x78, 0x2f, 0x10, 0xf1, 0xb6, 0xd8, 0x87, 0xce, 0x8b, 0x46, 0x71, 0xfe, 0xad, 0xf2,
	0xf0, 0xe7, 0x01, 0x26, 0xa9, 0x7f, 0x98, 0xa1, 0x0e, 0x92, 0xc4, 0x38, 0x15, 0x13, 0x7e, 0x5e,
	0xcb, 0x1f, 0x58, 0x4f, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x97, 0xb5, 0xcf, 0xbd, 0x93, 0x0a,
	0x00, 0x00,
}
