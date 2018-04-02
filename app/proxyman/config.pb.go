package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_net "v2ray.com/core/common/net"
import v2ray_core_common_net1 "v2ray.com/core/common/net"
import v2ray_core_transport_internet "v2ray.com/core/transport/internet"

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

type OutboundConfig struct {
}

func (m *OutboundConfig) Reset()                    { *m = OutboundConfig{} }
func (m *OutboundConfig) String() string            { return proto.CompactTextString(m) }
func (*OutboundConfig) ProtoMessage()               {}
func (*OutboundConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
func (*SenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

type MultiplexingConfig struct {
	// Whether or not Mux is enabled.
	Enabled bool `protobuf:"varint,1,opt,name=enabled" json:"enabled,omitempty"`
	// Max number of concurrent connections that one Mux connection can handle.
	Concurrency uint32 `protobuf:"varint,2,opt,name=concurrency" json:"concurrency,omitempty"`
}

func (m *MultiplexingConfig) Reset()                    { *m = MultiplexingConfig{} }
func (m *MultiplexingConfig) String() string            { return proto.CompactTextString(m) }
func (*MultiplexingConfig) ProtoMessage()               {}
func (*MultiplexingConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterType((*SenderConfig)(nil), "v2ray.core.app.proxyman.SenderConfig")
	proto.RegisterType((*MultiplexingConfig)(nil), "v2ray.core.app.proxyman.MultiplexingConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.KnownProtocols", KnownProtocols_name, KnownProtocols_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 691 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0x5d, 0x6f, 0xd3, 0x3a,
	0x1c, 0xc6, 0x97, 0xb6, 0x6b, 0x7b, 0xfe, 0x5b, 0xb3, 0x1c, 0x9f, 0x23, 0x2d, 0xa7, 0x07, 0xa4,
	0x52, 0x90, 0x56, 0x0d, 0x94, 0x40, 0x27, 0x2e, 0xb8, 0x82, 0xd1, 0x4d, 0x62, 0xbc, 0xa8, 0xc1,
	0xad, 0xb8, 0x98, 0x90, 0x22, 0x2f, 0xf1, 0x8a, 0x45, 0x62, 0x47, 0x8e, 0xdb, 0x2d, 0x5f, 0x89,
	0x4f, 0xc1, 0x25, 0x9f, 0x81, 0x4f, 0x83, 0xf2, 0xd6, 0x75, 0xeb, 0x3a, 0x98, 0x76, 0xe7, 0xa6,
	0xcf, 0xf3, 0xb3, 0xfd, 0x3c, 0xb6, 0xa1, 0x37, 0xeb, 0x4b, 0x92, 0x58, 0x9e, 0x08, 0x6d, 0x4f,
	0x48, 0x6a, 0x93, 0x28, 0xb2, 0x23, 0x29, 0xce, 0x93, 0x90, 0x70, 0xdb, 0x13, 0xfc, 0x94, 0x4d,
	0xac, 0x48, 0x0a, 0x25, 0xd0, 0x76, 0xa9, 0x94, 0xd4, 0x22, 0x51, 0x64, 0x95, 0xaa, 0xf6, 0xce,
	0x15, 0x84, 0x27, 0xc2, 0x50, 0x70, 0x9b, 0x53, 0x65, 0x13, 0xdf, 0x97, 0x34, 0x8e, 0x73, 0x42,
	0xfb, 0xd1, 0x6a, 0x61, 0x24, 0xa4, 0x2a, 0x54, 0xd6, 0x15, 0x95, 0x92, 0x84, 0xc7, 0xe9, 0xff,
	0x36, 0xe3, 0x8a, 0xca, 0x54, 0xbd, 0xb8, 0xae, 0xee, 0x16, 0xb4, 0x8e, 0xf8, 0x89, 0x98, 0x72,
	0x7f, 0x90, 0x7d, 0xee, 0x7e, 0xaf, 0x02, 0xda, 0x0f, 0x02, 0xe1, 0x11, 0xc5, 0x04, 0x1f, 0x29,
	0x49, 0x14, 0x9d, 0x24, 0xe8, 0x00, 0x6a, 0x2a, 0x89, 0xa8, 0xa9, 0x75, 0xb4, 0x9e, 0xde, 0x7f,
	0x6a, 0xad, 0xd8, 0x8e, 0xb5, 0x6c, 0xb5, 0xc6, 0x49, 0x44, 0x71, 0xe6, 0x46, 0x5f, 0x61, 0xc3,
	0x13, 0xdc, 0x9b, 0x4a, 0x49, 0xb9, 0x97, 0x98, 0x95, 0x8e, 0xd6, 0xdb, 0xe8, 0x1f, 0xdd, 0x06,
	0xb6, 0xfc, 0x69, 0x70, 0x01, 0xc4, 0x8b, 0x74, 0xe4, 0x42, 0x43, 0xd2, 0x53, 0x49, 0xe3, 0x2f,
	0x66, 0x35, 0x9b, 0xe8, 0xf0, 0x6e, 0x13, 0xe1, 0x1c, 0x86, 0x4b, 0x6a, 0xfb, 0x39, 0xdc, 0xbf,
	0x71, 0x39, 0xe8, 0x5f, 0x58, 0x9f, 0x91, 0x60, 0x9a, 0xa7, 0xd6, 0xc2, 0xf9, 0x8f, 0xf6, 0x33,
	0xf8, 0x6f, 0x25, 0xfc, 0x7a, 0x4b, 0xf7, 0x09, 0xd4, 0xd2, 0x14, 0x11, 0x40, 0x7d, 0x3f, 0x38,
	0x23, 0x49, 0x6c, 0xac, 0xa5, 0x63, 0x4c, 0xb8, 0x2f, 0x42, 0x43, 0x43, 0x9b, 0xd0, 0x3c, 0x3c,
	0x4f, 0xeb, 0x25, 0x81, 0x51, 0xe9, 0xfe, 0xac, 0x82, 0x8e, 0xa9, 0x47, 0xd9, 0x8c, 0xca, 0xbc,
	0x55, 0xf4, 0x12, 0x20, 0x3d, 0x04, 0xae, 0x24, 0x7c, 0x92, 0xb3, 0x37, 0xfa, 0x9d, 0xc5, 0x38,
	0xf2, 0xd3, 0x64, 0x71, 0xaa, 0x2c, 0x47, 0x48, 0x85, 0x53, 0x1d, 0xfe, 0x2b, 0x2a, 0x87, 0xe8,
	0x05, 0xd4, 0x03, 0x16, 0x2b, 0xca, 0x8b, 0xd2, 0x1e, 0xac, 0x30, 0x1f, 0x39, 0x43, 0x79, 0x20,
	0x42, 0xc2, 0x38, 0x2e, 0x0c, 0xe8, 0x33, 0xfc, 0x43, 0xe6, 0xfb, 0x75, 0xe3, 0x62, 0xc3, 0x45,
	0x27, 0x8f, 0x6f, 0xd1, 0x09, 0x46, 0x64, 0xf9, 0x60, 0x8e, 0x61, 0x2b, 0x56, 0x92, 0x92, 0xd0,
	0x8d, 0xa9, 0x52, 0x8c, 0x4f, 0x62, 0xb3, 0xb6, 0x4c, 0x9e, 0x5f, 0x03, 0xab, 0xbc, 0x06, 0xd6,
	0x28, 0x73, 0xe5, 0xf9, 0x60, 0x3d, 0x67, 0x8c, 0x0a, 0x04, 0x7a, 0x05, 0xf7, 0x64, 0x9e, 0xa0,
	0x2b, 0x24, 0x9b, 0x30, 0x4e, 0x02, 0xd7, 0xa7, 0xb1, 0x62, 0x3c, 0x9b, 0xdd, 0x5c, 0xef, 0x68,
	0xbd, 0x26, 0x6e, 0x17, 0x9a, 0x61, 0x21, 0x39, 0xb8, 0x50, 0x20, 0x07, 0xb6, 0xfc, 0x2c, 0x07,
	0x57, 0xcc, 0xa8, 0x94, 0xcc, 0xa7, 0x66, 0xa3, 0x53, 0xed, 0xe9, 0xfd, 0x9d, 0x95, 0x3b, 0x7e,
	0xc7, 0xc5, 0x19, 0x77, 0xd2, 0x6b, 0xe9, 0x89, 0x20, 0xc6, 0x7a, 0xee, 0x1f, 0x16, 0xf6, 0xb7,
	0xb5, 0x66, 0xdd, 0x68, 0x74, 0x0d, 0xd0, 0x87, 0x53, 0xb5, 0x78, 0x63, 0x7f, 0x54, 0x60, 0x73,
	0x44, 0xb9, 0x3f, 0x2f, 0x7b, 0x0f, 0xaa, 0x33, 0x46, 0x8a, 0x96, 0xff, 0xa0, 0xa8, 0x54, 0x7d,
	0x5d, 0x8e, 0x95, 0xbb, 0xe7, 0xf8, 0x11, 0xf4, 0x6c, 0x7b, 0x17, 0xd0, 0xbc, 0xf6, 0xdd, 0xdf,
	0x40, 0x9d, 0xd4, 0x54, 0x30, 0x5b, 0x19, 0x61, 0x8e, 0x3c, 0x06, 0x14, 0x4e, 0x03, 0xc5, 0xa2,
	0x80, 0x9e, 0xdf, 0xd8, 0xf9, 0xa5, 0x6c, 0x3f, 0x94, 0x16, 0xc6, 0x27, 0x05, 0xf7, 0xef, 0x39,
	0xa6, 0x64, 0x77, 0x1d, 0x40, 0xcb, 0x42, 0x64, 0x42, 0x83, 0x72, 0x72, 0x12, 0x50, 0x3f, 0xcb,
	0xb4, 0x89, 0xcb, 0x9f, 0xa8, 0xb3, 0xfc, 0x9e, 0xb5, 0x2e, 0x3d, 0x42, 0xbb, 0x0f, 0x41, 0xbf,
	0x5c, 0x2b, 0x6a, 0x42, 0xed, 0xcd, 0x78, 0xec, 0x18, 0x6b, 0xa8, 0x01, 0xd5, 0xf1, 0xfb, 0x91,
	0xa1, 0xbd, 0x1e, 0xc0, 0xff, 0x9e, 0x08, 0x57, 0xad, 0xdd, 0xd1, 0x8e, 0x9b, 0xe5, 0xf8, 0x5b,
	0x65, 0xfb, 0x53, 0x1f, 0x93, 0xc4, 0x1a, 0xa4, 0xaa, 0xfd, 0x28, 0xca, 0x93, 0x0a, 0x09, 0x3f,
	0xa9, 0x67, 0x0f, 0xfa, 0xde, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x77, 0x7f, 0xdc, 0x8e, 0x94,
	0x06, 0x00, 0x00,
}
