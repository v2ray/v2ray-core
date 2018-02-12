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

type MultiplexingConfig struct {
	// Whether or not Mux is enabled.
	Enabled bool `protobuf:"varint,1,opt,name=enabled" json:"enabled,omitempty"`
	// Max number of concurrent connections that one Mux connection can handle.
	Concurrency uint32 `protobuf:"varint,2,opt,name=concurrency" json:"concurrency,omitempty"`
}

func (m *MultiplexingConfig) Reset()                    { *m = MultiplexingConfig{} }
func (m *MultiplexingConfig) String() string            { return proto.CompactTextString(m) }
func (*MultiplexingConfig) ProtoMessage()               {}
func (*MultiplexingConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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
	proto.RegisterType((*MultiplexingConfig)(nil), "v2ray.core.app.proxyman.MultiplexingConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.KnownProtocols", KnownProtocols_name, KnownProtocols_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 847 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xdd, 0x6e, 0x1b, 0x45,
	0x18, 0xed, 0x7a, 0x1d, 0xdb, 0xfd, 0xd2, 0x38, 0xdb, 0xa1, 0x52, 0x8d, 0x01, 0xc9, 0x18, 0x44,
	0xad, 0x82, 0xd6, 0xc5, 0x15, 0x48, 0x88, 0x0b, 0x08, 0x49, 0xa5, 0x26, 0x10, 0xc5, 0x8c, 0x0d,
	0x17, 0x15, 0x92, 0x35, 0xd9, 0x9d, 0x2c, 0xa3, 0xec, 0xce, 0xac, 0x66, 0xc6, 0x6e, 0xf6, 0x95,
	0x78, 0x06, 0x2e, 0xb8, 0xe4, 0x82, 0x27, 0xe0, 0x69, 0xd0, 0xee, 0xcc, 0x6e, 0x1c, 0xff, 0xb4,
	0x44, 0x55, 0xee, 0x66, 0xed, 0x73, 0xce, 0xcc, 0x39, 0xdf, 0xf7, 0xcd, 0xc0, 0x60, 0x31, 0x92,
	0x24, 0xf3, 0x03, 0x91, 0x0c, 0x03, 0x21, 0xe9, 0x90, 0xa4, 0xe9, 0x30, 0x95, 0xe2, 0x2a, 0x4b,
	0x08, 0x1f, 0x06, 0x82, 0x5f, 0xb0, 0xc8, 0x4f, 0xa5, 0xd0, 0x02, 0x3d, 0x2e, 0x91, 0x92, 0xfa,
	0x24, 0x4d, 0xfd, 0x12, 0xd5, 0x7d, 0xb2, 0x22, 0x11, 0x88, 0x24, 0x11, 0x7c, 0xc8, 0xa9, 0x1e,
	0x92, 0x30, 0x94, 0x54, 0x29, 0xa3, 0xd0, 0xfd, 0x74, 0x3b, 0x30, 0x15, 0x52, 0x5b, 0x94, 0xbf,
	0x82, 0xd2, 0x92, 0x70, 0x95, 0xff, 0x3f, 0x64, 0x5c, 0x53, 0x99, 0xa3, 0x97, 0xcf, 0xd5, 0xfd,
	0xf6, 0xed, 0xf8, 0x50, 0x24, 0x84, 0x71, 0x25, 0x82, 0xcb, 0x55, 0xf2, 0xb3, 0xcd, 0x47, 0x52,
	0x54, 0x32, 0x12, 0x0f, 0x75, 0x96, 0xd2, 0x70, 0x96, 0x50, 0xa5, 0x48, 0x44, 0x0d, 0xa3, 0xbf,
	0x0f, 0x7b, 0xc7, 0xfc, 0x5c, 0xcc, 0x79, 0x78, 0x58, 0x08, 0xf5, 0xff, 0x72, 0x01, 0x1d, 0xc4,
	0xb1, 0x08, 0x88, 0x66, 0x82, 0x4f, 0xb4, 0x24, 0x9a, 0x46, 0x19, 0x3a, 0x82, 0x7a, 0x4e, 0xef,
	0x38, 0x3d, 0x67, 0xd0, 0x1e, 0x3d, 0xf3, 0xb7, 0xa4, 0xe7, 0xaf, 0x53, 0xfd, 0x69, 0x96, 0x52,
	0x5c, 0xb0, 0xd1, 0x25, 0xec, 0x06, 0x82, 0x07, 0x73, 0x29, 0x29, 0x0f, 0xb2, 0x4e, 0xad, 0xe7,
	0x0c, 0x76, 0x47, 0xc7, 0xb7, 0x11, 0x5b, 0xff, 0xe9, 0xf0, 0x5a, 0x10, 0x2f, 0xab, 0xa3, 0x19,
	0x34, 0x25, 0xbd, 0x90, 0x54, 0xfd, 0xde, 0x71, 0x8b, 0x8d, 0x5e, 0xbc, 0xdb, 0x46, 0xd8, 0x88,
	0xe1, 0x52, 0xb5, 0xfb, 0x15, 0x7c, 0xf4, 0xc6, 0xe3, 0xa0, 0x47, 0xb0, 0xb3, 0x20, 0xf1, 0xdc,
	0xa4, 0xb6, 0x87, 0xcd, 0x47, 0xf7, 0x4b, 0x78, 0x7f, 0xab, 0xf8, 0x66, 0x4a, 0xff, 0x0b, 0xa8,
	0xe7, 0x29, 0x22, 0x80, 0xc6, 0x41, 0xfc, 0x9a, 0x64, 0xca, 0xbb, 0x97, 0xaf, 0x31, 0xe1, 0xa1,
	0x48, 0x3c, 0x07, 0x3d, 0x80, 0xd6, 0x8b, 0xab, 0xbc, 0x3b, 0x48, 0xec, 0xd5, 0xfa, 0xff, 0xba,
	0xd0, 0xc6, 0x34, 0xa0, 0x6c, 0x41, 0xa5, 0xa9, 0x2a, 0xfa, 0x0e, 0x20, 0xef, 0xa1, 0x99, 0x24,
	0x3c, 0x32, 0xda, 0xbb, 0xa3, 0xde, 0x72, 0x1c, 0xa6, 0x53, 0x7c, 0x4e, 0xb5, 0x3f, 0x16, 0x52,
	0xe3, 0x1c, 0x87, 0xef, 0xa7, 0xe5, 0x12, 0x7d, 0x03, 0x8d, 0x98, 0x29, 0x4d, 0xb9, 0x2d, 0xda,
	0xc7, 0x5b, 0xc8, 0xc7, 0xe3, 0x33, 0x79, 0x54, 0xf4, 0x27, 0xb6, 0x04, 0xf4, 0x1b, 0xbc, 0x47,
	0x2a, 0xbf, 0x33, 0x65, 0x0d, 0xdb, 0x9a, 0x7c, 0x7e, 0x8b, 0x9a, 0x60, 0x44, 0xd6, 0x1b, 0x73,
	0x0a, 0xfb, 0x4a, 0x4b, 0x4a, 0x92, 0x99, 0xa2, 0x5a, 0x33, 0x1e, 0xa9, 0x4e, 0x7d, 0x5d, 0xb9,
	0x9a, 0x22, 0xbf, 0x9c, 0x22, 0x7f, 0x52, 0xb0, 0x4c, 0x3e, 0xb8, 0x6d, 0x34, 0x26, 0x56, 0x02,
	0x7d, 0x0f, 0x1f, 0x4a, 0x93, 0xe0, 0x4c, 0x48, 0x16, 0x31, 0x4e, 0xe2, 0x59, 0x48, 0x95, 0x66,
	0xbc, 0xd8, 0xbd, 0xb3, 0xd3, 0x73, 0x06, 0x2d, 0xdc, 0xb5, 0x98, 0x33, 0x0b, 0x39, 0xba, 0x46,
	0xa0, 0x31, 0xec, 0x9b, 0x39, 0x9d, 0x89, 0x05, 0x95, 0x92, 0x85, 0xb4, 0xd3, 0xec, 0xb9, 0x83,
	0xf6, 0xe8, 0xc9, 0x56, 0xc7, 0x3f, 0x72, 0xf1, 0x9a, 0x8f, 0xf3, 0xb1, 0x0c, 0x44, 0xac, 0x70,
	0xdb, 0xf0, 0xcf, 0x2c, 0xfd, 0xa4, 0xde, 0x6a, 0x78, 0xcd, 0xfe, 0x9f, 0x35, 0x40, 0xbf, 0x70,
	0x76, 0xb5, 0x52, 0xe0, 0x0b, 0x40, 0x06, 0x3e, 0x11, 0xc1, 0x65, 0x69, 0xc3, 0xd6, 0xea, 0xeb,
	0xe5, 0x1d, 0x2b, 0xff, 0xcb, 0xb7, 0x88, 0x7f, 0x54, 0x71, 0xa9, 0x2e, 0xd9, 0x78, 0x83, 0xe2,
	0x1d, 0xc5, 0x7d, 0x17, 0x61, 0x39, 0x5e, 0xed, 0xa4, 0xde, 0x72, 0xbd, 0xfa, 0x49, 0xbd, 0xb5,
	0xe3, 0x35, 0x6c, 0x7c, 0xff, 0x38, 0xf0, 0xc8, 0x5e, 0x78, 0x2f, 0x09, 0x0f, 0xe3, 0x2a, 0x40,
	0x0f, 0x5c, 0x4d, 0xa2, 0x62, 0x34, 0xee, 0xe3, 0x7c, 0x89, 0x26, 0xf0, 0xd0, 0xd6, 0x57, 0x5e,
	0x9b, 0x35, 0x89, 0x7e, 0xb6, 0xa1, 0xfb, 0xcd, 0x25, 0x5b, 0xdc, 0x76, 0xe1, 0xa9, 0xb9, 0x63,
	0xb1, 0x57, 0x0a, 0x54, 0x4e, 0x4f, 0xa1, 0x5d, 0x58, 0xb8, 0x56, 0x74, 0x6f, 0xa5, 0xb8, 0x57,
	0xb0, 0x4b, 0xb9, 0xbe, 0x07, 0xed, 0xb3, 0xb9, 0x5e, 0xbe, 0xbf, 0xff, 0xae, 0xc1, 0x83, 0x09,
	0xe5, 0x61, 0x65, 0xec, 0x39, 0xb8, 0x0b, 0x46, 0xec, 0xcc, 0xff, 0x8f, 0xb1, 0xcd, 0xd1, 0x9b,
	0xca, 0x5c, 0x7b, 0xf7, 0x32, 0xff, 0xbc, 0xc5, 0xfc, 0xd3, 0xb7, 0x88, 0x8e, 0x73, 0x92, 0xd5,
	0xbc, 0x19, 0x00, 0x7a, 0x05, 0x28, 0x99, 0xc7, 0x9a, 0xa5, 0x31, 0xbd, 0x7a, 0x63, 0x4b, 0xde,
	0x68, 0x9e, 0xd3, 0x92, 0xc2, 0x78, 0x64, 0x75, 0x1f, 0x56, 0x32, 0x55, 0xb8, 0x63, 0x40, 0xeb,
	0x40, 0xd4, 0x81, 0x26, 0xe5, 0xe4, 0x3c, 0xa6, 0x61, 0x91, 0x69, 0x0b, 0x97, 0x9f, 0xa8, 0xb7,
	0xfe, 0xba, 0xed, 0xdd, 0x78, 0x92, 0x9e, 0x7e, 0x02, 0xed, 0x9b, 0x7d, 0x8b, 0x5a, 0x50, 0x7f,
	0x39, 0x9d, 0x8e, 0xbd, 0x7b, 0xa8, 0x09, 0xee, 0xf4, 0xa7, 0x89, 0xe7, 0xfc, 0x70, 0x08, 0x1f,
	0x04, 0x22, 0xd9, 0x76, 0xf6, 0xb1, 0xf3, 0xaa, 0x55, 0xae, 0xff, 0xa8, 0x3d, 0xfe, 0x75, 0x84,
	0x49, 0xe6, 0x1f, 0xe6, 0xa8, 0x83, 0x34, 0x35, 0x49, 0x25, 0x84, 0x9f, 0x37, 0x8a, 0xe7, 0xfd,
	0xf9, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc0, 0x60, 0xa6, 0x99, 0x11, 0x09, 0x00, 0x00,
}
