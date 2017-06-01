package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_serial "v2ray.com/core/common/serial"
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

type OutboundHandlerConfig struct {
	Tag            string                                 `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	SenderSettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=sender_settings,json=senderSettings" json:"sender_settings,omitempty"`
	ProxySettings  *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	Expire         int64                                  `protobuf:"varint,4,opt,name=expire" json:"expire,omitempty"`
	Comment        string                                 `protobuf:"bytes,5,opt,name=comment" json:"comment,omitempty"`
}

func (m *OutboundHandlerConfig) Reset()                    { *m = OutboundHandlerConfig{} }
func (m *OutboundHandlerConfig) String() string            { return proto.CompactTextString(m) }
func (*OutboundHandlerConfig) ProtoMessage()               {}
func (*OutboundHandlerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *OutboundHandlerConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *OutboundHandlerConfig) GetSenderSettings() *v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.SenderSettings
	}
	return nil
}

func (m *OutboundHandlerConfig) GetProxySettings() *v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

func (m *OutboundHandlerConfig) GetExpire() int64 {
	if m != nil {
		return m.Expire
	}
	return 0
}

func (m *OutboundHandlerConfig) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
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
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.app.proxyman.InboundHandlerConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterType((*SenderConfig)(nil), "v2ray.core.app.proxyman.SenderConfig")
	proto.RegisterType((*OutboundHandlerConfig)(nil), "v2ray.core.app.proxyman.OutboundHandlerConfig")
	proto.RegisterType((*MultiplexingConfig)(nil), "v2ray.core.app.proxyman.MultiplexingConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.KnownProtocols", KnownProtocols_name, KnownProtocols_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 822 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xd1, 0x8e, 0xdb, 0x44,
	0x14, 0xad, 0xe3, 0x34, 0xc9, 0xde, 0xed, 0x7a, 0xdd, 0xa1, 0xd0, 0x10, 0x40, 0x0a, 0x01, 0xd1,
	0xa8, 0x20, 0xa7, 0xa4, 0xe2, 0x81, 0x27, 0x58, 0x76, 0x2b, 0x75, 0x81, 0x55, 0xcc, 0x24, 0xe2,
	0xa1, 0x42, 0xb2, 0x66, 0xed, 0xa9, 0x19, 0x61, 0xcf, 0x58, 0x33, 0x93, 0x74, 0xfd, 0x4b, 0x7c,
	0x05, 0x8f, 0x3c, 0xf0, 0x05, 0xfc, 0x0a, 0x2f, 0xc8, 0x9e, 0x71, 0x76, 0xb7, 0x49, 0x5a, 0x96,
	0xaa, 0x6f, 0x33, 0xc9, 0x39, 0xc7, 0x73, 0xcf, 0x3d, 0x77, 0x06, 0xc6, 0xab, 0xa9, 0x24, 0x65,
	0x10, 0x8b, 0x7c, 0x12, 0x0b, 0x49, 0x27, 0xa4, 0x28, 0x26, 0x85, 0x14, 0x17, 0x65, 0x4e, 0xf8,
	0x24, 0x16, 0xfc, 0x39, 0x4b, 0x83, 0x42, 0x0a, 0x2d, 0xd0, 0xfd, 0x06, 0x29, 0x69, 0x40, 0x8a,
	0x22, 0x68, 0x50, 0x83, 0x47, 0x2f, 0x49, 0xc4, 0x22, 0xcf, 0x05, 0x9f, 0x28, 0x2a, 0x19, 0xc9,
	0x26, 0xba, 0x2c, 0x68, 0x12, 0xe5, 0x54, 0x29, 0x92, 0x52, 0x23, 0x35, 0x78, 0xb0, 0x9d, 0xc1,
	0xa9, 0x9e, 0x90, 0x24, 0x91, 0x54, 0x29, 0x0b, 0xfc, 0x74, 0x37, 0xb0, 0x10, 0x52, 0x5b, 0x54,
	0xf0, 0x12, 0x4a, 0x4b, 0xc2, 0x55, 0xf5, 0xff, 0x84, 0x71, 0x4d, 0x65, 0x85, 0xbe, 0x5a, 0xc9,
	0xe8, 0x10, 0x0e, 0x4e, 0xf9, 0xb9, 0x58, 0xf2, 0xe4, 0xb8, 0xfe, 0x79, 0xf4, 0x87, 0x0b, 0xe8,
	0x28, 0xcb, 0x44, 0x4c, 0x34, 0x13, 0x7c, 0xae, 0x25, 0xd1, 0x34, 0x2d, 0xd1, 0x09, 0xb4, 0xab,
	0xd3, 0xf7, 0x9d, 0xa1, 0x33, 0xf6, 0xa6, 0x8f, 0x82, 0x1d, 0x06, 0x04, 0x9b, 0xd4, 0x60, 0x51,
	0x16, 0x14, 0xd7, 0x6c, 0xf4, 0x1b, 0xec, 0xc7, 0x82, 0xc7, 0x4b, 0x29, 0x29, 0x8f, 0xcb, 0x7e,
	0x6b, 0xe8, 0x8c, 0xf7, 0xa7, 0xa7, 0x37, 0x11, 0xdb, 0xfc, 0xe9, 0xf8, 0x52, 0x10, 0x5f, 0x55,
	0x47, 0x11, 0x74, 0x25, 0x7d, 0x2e, 0xa9, 0xfa, 0xb5, 0xef, 0xd6, 0x1f, 0x7a, 0xf2, 0x66, 0x1f,
	0xc2, 0x46, 0x0c, 0x37, 0xaa, 0x83, 0xaf, 0xe0, 0xa3, 0x57, 0x1e, 0x07, 0xdd, 0x83, 0xdb, 0x2b,
	0x92, 0x2d, 0x8d, 0x6b, 0x07, 0xd8, 0x6c, 0x06, 0x5f, 0xc2, 0xfb, 0x3b, 0xc5, 0xb7, 0x53, 0x46,
	0x5f, 0x40, 0xbb, 0x72, 0x11, 0x01, 0x74, 0x8e, 0xb2, 0x17, 0xa4, 0x54, 0xfe, 0xad, 0x6a, 0x8d,
	0x09, 0x4f, 0x44, 0xee, 0x3b, 0xe8, 0x0e, 0xf4, 0x9e, 0x5c, 0x54, 0xed, 0x25, 0x99, 0xdf, 0x1a,
	0xfd, 0xed, 0x82, 0x87, 0x69, 0x4c, 0xd9, 0x8a, 0x4a, 0xd3, 0x55, 0xf4, 0x0d, 0x40, 0x15, 0x82,
	0x48, 0x12, 0x9e, 0x1a, 0xed, 0xfd, 0xe9, 0xf0, 0xaa, 0x1d, 0x26, 0x4d, 0x01, 0xa7, 0x3a, 0x08,
	0x85, 0xd4, 0xb8, 0xc2, 0xe1, 0xbd, 0xa2, 0x59, 0xa2, 0xaf, 0xa1, 0x93, 0x31, 0xa5, 0x29, 0xb7,
	0x4d, 0xfb, 0x78, 0x07, 0xf9, 0x34, 0x9c, 0xc9, 0x13, 0x91, 0x13, 0xc6, 0xb1, 0x25, 0xa0, 0x5f,
	0xe0, 0x1d, 0xb2, 0xae, 0x37, 0x52, 0xb6, 0x60, 0xdb, 0x93, 0xcf, 0x6f, 0xd0, 0x13, 0x8c, 0xc8,
	0x66, 0x30, 0x17, 0x70, 0xa8, 0xb4, 0xa4, 0x24, 0x8f, 0x14, 0xd5, 0x9a, 0xf1, 0x54, 0xf5, 0xdb,
	0x9b, 0xca, 0xeb, 0x31, 0x08, 0x9a, 0x31, 0x08, 0xe6, 0x35, 0xcb, 0xf8, 0x83, 0x3d, 0xa3, 0x31,
	0xb7, 0x12, 0xe8, 0x5b, 0xf8, 0x50, 0x1a, 0x07, 0x23, 0x21, 0x59, 0xca, 0x38, 0xc9, 0xa2, 0x84,
	0x2a, 0xcd, 0x78, 0xfd, 0xf5, 0xfe, 0xed, 0xa1, 0x33, 0xee, 0xe1, 0x81, 0xc5, 0xcc, 0x2c, 0xe4,
	0xe4, 0x12, 0x81, 0x42, 0x38, 0x4c, 0x6a, 0x1f, 0x22, 0xb1, 0xa2, 0x52, 0xb2, 0x84, 0xf6, 0xbb,
	0x43, 0x77, 0xec, 0x4d, 0x1f, 0xec, 0xac, 0xf8, 0x07, 0x2e, 0x5e, 0xf0, 0xb0, 0x1a, 0xcb, 0x58,
	0x64, 0x0a, 0x7b, 0x86, 0x3f, 0xb3, 0xf4, 0xef, 0xdb, 0xbd, 0x8e, 0xdf, 0x1d, 0xfd, 0xe5, 0xc0,
	0x3d, 0x3b, 0xb1, 0x4f, 0x09, 0x4f, 0xb2, 0x75, 0x8b, 0x7d, 0x70, 0x35, 0x49, 0xeb, 0xde, 0xee,
	0xe1, 0x6a, 0x89, 0xe6, 0x70, 0xd7, 0x1e, 0x50, 0x5e, 0x9a, 0x63, 0xda, 0xf7, 0xd9, 0x96, 0xf6,
	0x99, 0x4b, 0xaa, 0x1e, 0xd7, 0xe4, 0xcc, 0xdc, 0x51, 0xd8, 0x6f, 0x04, 0xd6, 0xce, 0x9c, 0x81,
	0x57, 0x1f, 0xf8, 0x52, 0xd1, 0xbd, 0x91, 0xe2, 0x41, 0xcd, 0x6e, 0xe4, 0x46, 0x3e, 0x78, 0xb3,
	0xa5, 0xbe, 0x7a, 0x01, 0xfd, 0xd9, 0x82, 0x3b, 0x73, 0xca, 0x93, 0x75, 0x61, 0x8f, 0xc1, 0x5d,
	0x31, 0x62, 0x43, 0xfb, 0x1f, 0x72, 0x57, 0xa1, 0xb7, 0xc5, 0xa2, 0xf5, 0xe6, 0xb1, 0xf8, 0x69,
	0x47, 0xf1, 0x0f, 0x5f, 0x23, 0x1a, 0x56, 0x24, 0xab, 0x79, 0xdd, 0x00, 0xf4, 0x0c, 0x50, 0xbe,
	0xcc, 0x34, 0x2b, 0x32, 0x7a, 0xf1, 0xca, 0x08, 0x5f, 0x8b, 0xca, 0x59, 0x43, 0x61, 0x3c, 0xb5,
	0xba, 0x77, 0xd7, 0x32, 0x6b, 0x73, 0xff, 0x71, 0xe0, 0xdd, 0xc6, 0xdd, 0xd7, 0x85, 0x65, 0x06,
	0x87, 0xaa, 0x76, 0xfd, 0xff, 0x46, 0xc5, 0x33, 0xf4, 0xb7, 0x14, 0x14, 0xf4, 0x1e, 0x74, 0xe8,
	0x45, 0xc1, 0x24, 0xad, 0xbd, 0x71, 0xb1, 0xdd, 0xa1, 0x3e, 0x74, 0x2b, 0x11, 0xca, 0x75, 0x3d,
	0x94, 0x7b, 0xb8, 0xd9, 0x8e, 0x42, 0x40, 0x9b, 0x36, 0x55, 0x78, 0xca, 0xc9, 0x79, 0x46, 0x93,
	0xba, 0xfa, 0x1e, 0x6e, 0xb6, 0x68, 0xb8, 0xf9, 0x38, 0x1d, 0x5c, 0x7b, 0x51, 0x1e, 0x7e, 0x02,
	0xde, 0xf5, 0x19, 0x45, 0x3d, 0x68, 0x3f, 0x5d, 0x2c, 0x42, 0xff, 0x16, 0xea, 0x82, 0xbb, 0xf8,
	0x71, 0xee, 0x3b, 0xdf, 0x1d, 0xc3, 0x07, 0xb1, 0xc8, 0x77, 0x75, 0x2e, 0x74, 0x9e, 0xf5, 0x9a,
	0xf5, 0xef, 0xad, 0xfb, 0x3f, 0x4f, 0x31, 0x29, 0x83, 0xe3, 0x0a, 0x75, 0x54, 0x14, 0x26, 0x27,
	0x39, 0xe1, 0xe7, 0x9d, 0xfa, 0x75, 0x7e, 0xfc, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x25, 0xdf,
	0x6a, 0xb2, 0x93, 0x08, 0x00, 0x00,
}
