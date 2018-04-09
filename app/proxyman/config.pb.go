package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_net "v2ray.com/core/common/net"
import v2ray_core_common_net1 "v2ray.com/core/common/net"
import v2ray_core_transport_internet "v2ray.com/core/transport/internet"
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
	// PortRange specifies the ports which the Receiver should listen on.
	PortRange *v2ray_core_common_net1.PortRange `protobuf:"bytes,1,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
	// Listen specifies the IP address that the Receiver should listen on.
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

type InboundHandlerConfig struct {
	Tag              string                                 `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	ReceiverSettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=receiver_settings,json=receiverSettings" json:"receiver_settings,omitempty"`
	ProxySettings    *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
}

func (m *InboundHandlerConfig) Reset()                    { *m = InboundHandlerConfig{} }
func (m *InboundHandlerConfig) String() string            { return proto.CompactTextString(m) }
func (*InboundHandlerConfig) ProtoMessage()               {}
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
func (*OutboundConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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
func (*SenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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
func (*MultiplexingConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.app.proxyman.InboundHandlerConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterType((*SenderConfig)(nil), "v2ray.core.app.proxyman.SenderConfig")
	proto.RegisterType((*MultiplexingConfig)(nil), "v2ray.core.app.proxyman.MultiplexingConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.KnownProtocols", KnownProtocols_name, KnownProtocols_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 772 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0x5f, 0x6f, 0xeb, 0x34,
	0x18, 0xc6, 0x4f, 0x9a, 0x9e, 0xb6, 0xe7, 0xdd, 0x69, 0x96, 0x63, 0x26, 0xad, 0x14, 0x90, 0x4a,
	0x41, 0xac, 0x1a, 0x28, 0x19, 0x9d, 0xb8, 0xe0, 0x0a, 0x46, 0x37, 0x69, 0x03, 0xa6, 0x06, 0xb7,
	0xe2, 0x62, 0x42, 0x8a, 0xbc, 0xc4, 0x0b, 0x11, 0x89, 0x1d, 0x39, 0x6e, 0xb7, 0x7c, 0x25, 0x3e,
	0x05, 0x97, 0x5c, 0xf0, 0x09, 0xf8, 0x34, 0x28, 0x71, 0xd2, 0x3f, 0x6b, 0x3b, 0xce, 0xb4, 0x3b,
	0x67, 0x7b, 0x9e, 0x9f, 0xed, 0xe7, 0x7d, 0xfd, 0x16, 0x06, 0xf3, 0xa1, 0x20, 0x99, 0xe5, 0xf1,
	0xd8, 0xf6, 0xb8, 0xa0, 0x36, 0x49, 0x12, 0x3b, 0x11, 0xfc, 0x21, 0x8b, 0x09, 0xb3, 0x3d, 0xce,
	0xee, 0xc2, 0xc0, 0x4a, 0x04, 0x97, 0x1c, 0x1d, 0x56, 0x4a, 0x41, 0x2d, 0x92, 0x24, 0x56, 0xa5,
	0xea, 0x1e, 0x3d, 0x42, 0x78, 0x3c, 0x8e, 0x39, 0xb3, 0x19, 0x95, 0x36, 0xf1, 0x7d, 0x41, 0xd3,
	0x54, 0x11, 0xba, 0x9f, 0xef, 0x16, 0x26, 0x5c, 0xc8, 0x52, 0x65, 0x3d, 0x52, 0x49, 0x41, 0x58,
	0x9a, 0xff, 0xdf, 0x0e, 0x99, 0xa4, 0x22, 0x57, 0xaf, 0x9e, 0xab, 0x7b, 0xb2, 0x9d, 0x9a, 0x52,
	0x11, 0x92, 0xc8, 0x96, 0x59, 0x42, 0x7d, 0x37, 0xa6, 0x69, 0x4a, 0x02, 0xaa, 0x1c, 0xfd, 0x7d,
	0x68, 0x5f, 0xb1, 0x5b, 0x3e, 0x63, 0xfe, 0xa8, 0x00, 0xf5, 0xff, 0xd2, 0x01, 0x9d, 0x45, 0x11,
	0xf7, 0x88, 0x0c, 0x39, 0x9b, 0x48, 0x41, 0x24, 0x0d, 0x32, 0x74, 0x0e, 0xf5, 0xdc, 0xde, 0xd1,
	0x7a, 0xda, 0xc0, 0x18, 0x9e, 0x58, 0x3b, 0x02, 0xb0, 0x36, 0xad, 0xd6, 0x34, 0x4b, 0x28, 0x2e,
	0xdc, 0xe8, 0x0f, 0xd8, 0xf3, 0x38, 0xf3, 0x66, 0x42, 0x50, 0xe6, 0x65, 0x9d, 0x5a, 0x4f, 0x1b,
	0xec, 0x0d, 0xaf, 0x9e, 0x03, 0xdb, 0xfc, 0xd3, 0x68, 0x09, 0xc4, 0xab, 0x74, 0xe4, 0x42, 0x53,
	0xd0, 0x3b, 0x41, 0xd3, 0xdf, 0x3b, 0x7a, 0xb1, 0xd1, 0xc5, 0xcb, 0x36, 0xc2, 0x0a, 0x86, 0x2b,
	0x6a, 0xf7, 0x1b, 0xf8, 0xe4, 0xc9, 0xe3, 0xa0, 0x03, 0x78, 0x3d, 0x27, 0xd1, 0x4c, 0xa5, 0xd6,
	0xc6, 0xea, 0xa3, 0xfb, 0x35, 0x7c, 0xb8, 0x13, 0xbe, 0xdd, 0xd2, 0xff, 0x0a, 0xea, 0x79, 0x8a,
	0x08, 0xa0, 0x71, 0x16, 0xdd, 0x93, 0x2c, 0x35, 0x5f, 0xe5, 0x6b, 0x4c, 0x98, 0xcf, 0x63, 0x53,
	0x43, 0x6f, 0xa1, 0x75, 0xf1, 0x90, 0x37, 0x04, 0x89, 0xcc, 0x5a, 0xff, 0x5f, 0x1d, 0x0c, 0x4c,
	0x3d, 0x1a, 0xce, 0xa9, 0x50, 0x55, 0x45, 0xdf, 0x01, 0xe4, 0x6d, 0xe3, 0x0a, 0xc2, 0x02, 0xc5,
	0xde, 0x1b, 0xf6, 0x56, 0xe3, 0x50, 0x9d, 0x62, 0x31, 0x2a, 0x2d, 0x87, 0x0b, 0x89, 0x73, 0x1d,
	0x7e, 0x93, 0x54, 0x4b, 0xf4, 0x2d, 0x34, 0xa2, 0x30, 0x95, 0x94, 0x95, 0x45, 0xfb, 0x74, 0x87,
	0xf9, 0xca, 0x19, 0x8b, 0x73, 0x1e, 0x93, 0x90, 0xe1, 0xd2, 0x80, 0x7e, 0x83, 0x0f, 0xc8, 0xe2,
	0xbe, 0x6e, 0x5a, 0x5e, 0xb8, 0xac, 0xc9, 0x97, 0xcf, 0xa8, 0x09, 0x46, 0x64, 0xb3, 0x31, 0xa7,
	0xb0, 0x9f, 0x4a, 0x41, 0x49, 0xec, 0xa6, 0x54, 0xca, 0x90, 0x05, 0x69, 0xa7, 0xbe, 0x49, 0x5e,
	0x3c, 0x1c, 0xab, 0x7a, 0x38, 0xd6, 0xa4, 0x70, 0xa9, 0x7c, 0xb0, 0xa1, 0x18, 0x93, 0x12, 0x81,
	0xbe, 0x87, 0x8f, 0x85, 0x4a, 0xd0, 0xe5, 0x22, 0x0c, 0x42, 0x46, 0x22, 0xd7, 0xa7, 0xa9, 0x0c,
	0x59, 0xb1, 0x7b, 0xe7, 0x75, 0x4f, 0x1b, 0xb4, 0x70, 0xb7, 0xd4, 0x8c, 0x4b, 0xc9, 0xf9, 0x52,
	0x81, 0x1c, 0xd8, 0xf7, 0x8b, 0x1c, 0x5c, 0x3e, 0xa7, 0x42, 0x84, 0x3e, 0xed, 0x34, 0x7b, 0xfa,
	0xc0, 0x18, 0x1e, 0xed, 0xbc, 0xf1, 0x4f, 0x8c, 0xdf, 0x33, 0x27, 0x7f, 0x96, 0x1e, 0x8f, 0x52,
	0x6c, 0x28, 0xff, 0xb8, 0xb4, 0xff, 0x58, 0x6f, 0x35, 0xcc, 0x66, 0xff, 0x1f, 0x0d, 0x0e, 0xca,
	0x17, 0x7b, 0x49, 0x98, 0x1f, 0x2d, 0x4a, 0x6c, 0x82, 0x2e, 0x49, 0x50, 0xd4, 0xf6, 0x0d, 0xce,
	0x97, 0x68, 0x02, 0xef, 0xca, 0x03, 0x8a, 0x65, 0x38, 0xaa, 0x7c, 0x5f, 0x6c, 0x29, 0x9f, 0x9a,
	0x12, 0xc5, 0x73, 0xf5, 0xaf, 0xd5, 0x90, 0xc0, 0x66, 0x05, 0x58, 0x24, 0x73, 0x0d, 0x46, 0x71,
	0xe0, 0x25, 0x51, 0x7f, 0x16, 0xb1, 0x5d, 0xb8, 0x2b, 0x5c, 0xdf, 0x04, 0x63, 0x3c, 0x93, 0xab,
	0x03, 0xe8, 0xef, 0x1a, 0xbc, 0x9d, 0x50, 0xe6, 0x2f, 0x2e, 0x76, 0x0a, 0xfa, 0x3c, 0x24, 0x65,
	0xd3, 0xbe, 0x47, 0xdf, 0xe5, 0xea, 0x6d, 0x6d, 0x51, 0x7b, 0x79, 0x5b, 0xfc, 0xb2, 0xe3, 0xf2,
	0xc7, 0xff, 0x03, 0x75, 0x72, 0x53, 0xc9, 0x5c, 0x0f, 0x00, 0xdd, 0x00, 0x8a, 0x67, 0x91, 0x0c,
	0x93, 0x88, 0x3e, 0x3c, 0xd9, 0xc2, 0x6b, 0xad, 0x72, 0x5d, 0x59, 0x42, 0x16, 0x94, 0xdc, 0x77,
	0x0b, 0xcc, 0x22, 0x5c, 0x07, 0xd0, 0xa6, 0x10, 0x75, 0xa0, 0x49, 0x19, 0xb9, 0x8d, 0xa8, 0x5f,
	0x64, 0xda, 0xc2, 0xd5, 0x27, 0xea, 0x6d, 0x8e, 0xe7, 0xf6, 0xda, 0x4c, 0x3d, 0xfe, 0x0c, 0x8c,
	0xf5, 0x2e, 0x45, 0x2d, 0xa8, 0x5f, 0x4e, 0xa7, 0x8e, 0xf9, 0x0a, 0x35, 0x41, 0x9f, 0xfe, 0x3c,
	0x31, 0xb5, 0x1f, 0x46, 0xf0, 0x91, 0xc7, 0xe3, 0x5d, 0x67, 0x77, 0xb4, 0x9b, 0x56, 0xb5, 0xfe,
	0xb3, 0x76, 0xf8, 0xeb, 0x10, 0x93, 0xcc, 0x1a, 0xe5, 0xaa, 0xb3, 0x24, 0x51, 0x49, 0xc5, 0x84,
	0xdd, 0x36, 0x8a, 0xdf, 0xa7, 0xd3, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x6a, 0xe5, 0x6f, 0xed,
	0x95, 0x07, 0x00, 0x00,
}
