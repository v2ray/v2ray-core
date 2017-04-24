package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_serial "v2ray.com/core/common/serial"
import v2ray_core_common_net "v2ray.com/core/common/net"
import _ "v2ray.com/core/common/net"
import v2ray_core_common_net3 "v2ray.com/core/common/net"
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
	PortRange                  *v2ray_core_common_net3.PortRange           `protobuf:"bytes,1,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
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

func (m *ReceiverConfig) GetPortRange() *v2ray_core_common_net3.PortRange {
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
	// 829 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xd1, 0x8e, 0xdb, 0x44,
	0x14, 0xad, 0xe3, 0x34, 0xc9, 0xde, 0xed, 0x7a, 0xdd, 0xa1, 0xd0, 0x10, 0x40, 0x0a, 0x01, 0xd1,
	0xa8, 0x45, 0x4e, 0x49, 0xc5, 0x03, 0x4f, 0xb0, 0xec, 0x56, 0xea, 0x02, 0xab, 0x84, 0x49, 0xc4,
	0x43, 0x85, 0x64, 0xcd, 0xda, 0x53, 0x33, 0xc2, 0x9e, 0xb1, 0x66, 0x26, 0xe9, 0xfa, 0x97, 0xf8,
	0x0a, 0x1e, 0x79, 0xe0, 0x0b, 0xf8, 0x15, 0x5e, 0x90, 0x3d, 0xe3, 0x6c, 0xb6, 0x59, 0xb7, 0x2c,
	0x15, 0x6f, 0x33, 0xc9, 0x39, 0xc7, 0x73, 0xcf, 0x3d, 0x77, 0x06, 0xc6, 0xeb, 0xa9, 0x24, 0x45,
	0x10, 0x89, 0x6c, 0x12, 0x09, 0x49, 0x27, 0x24, 0xcf, 0x27, 0xb9, 0x14, 0x17, 0x45, 0x46, 0xf8,
	0x24, 0x12, 0xfc, 0x05, 0x4b, 0x82, 0x5c, 0x0a, 0x2d, 0xd0, 0xfd, 0x1a, 0x29, 0x69, 0x40, 0xf2,
	0x3c, 0xa8, 0x51, 0x83, 0xc7, 0xaf, 0x48, 0x44, 0x22, 0xcb, 0x04, 0x9f, 0x28, 0x2a, 0x19, 0x49,
	0x27, 0xba, 0xc8, 0x69, 0x1c, 0x66, 0x54, 0x29, 0x92, 0x50, 0x23, 0x35, 0x78, 0x70, 0x3d, 0x83,
	0x53, 0x3d, 0x21, 0x71, 0x2c, 0xa9, 0x52, 0x16, 0xf8, 0xa8, 0x19, 0x18, 0x53, 0xa5, 0x19, 0x27,
	0x9a, 0x09, 0x6e, 0xc1, 0x9f, 0x36, 0x83, 0x73, 0x21, 0xb5, 0x45, 0x05, 0xaf, 0xa0, 0xb4, 0x24,
	0x5c, 0x95, 0xff, 0x4f, 0x18, 0xd7, 0x54, 0x96, 0xe8, 0xed, 0xb2, 0x47, 0x87, 0x70, 0x70, 0xca,
	0xcf, 0xc5, 0x8a, 0xc7, 0xc7, 0xd5, 0xcf, 0xa3, 0xdf, 0x5d, 0x40, 0x47, 0x69, 0x2a, 0xa2, 0xea,
	0xdb, 0x0b, 0x2d, 0x89, 0xa6, 0x49, 0x81, 0x4e, 0xa0, 0x5d, 0x96, 0xda, 0x77, 0x86, 0xce, 0xd8,
	0x9b, 0x3e, 0x0e, 0x1a, 0xdc, 0x0a, 0x76, 0xa9, 0xc1, 0xb2, 0xc8, 0x29, 0xae, 0xd8, 0xe8, 0x57,
	0xd8, 0x8f, 0x04, 0x8f, 0x56, 0x52, 0x52, 0x1e, 0x15, 0xfd, 0xd6, 0xd0, 0x19, 0xef, 0x4f, 0x4f,
	0x6f, 0x22, 0xb6, 0xfb, 0xd3, 0xf1, 0xa5, 0x20, 0xde, 0x56, 0x47, 0x21, 0x74, 0x25, 0x7d, 0x21,
	0xa9, 0xfa, 0xa5, 0xef, 0x56, 0x1f, 0x7a, 0xfa, 0x76, 0x1f, 0xc2, 0x46, 0x0c, 0xd7, 0xaa, 0x83,
	0x2f, 0xe1, 0xa3, 0xd7, 0x1e, 0x07, 0xdd, 0x83, 0xdb, 0x6b, 0x92, 0xae, 0x8c, 0x6b, 0x07, 0xd8,
	0x6c, 0x06, 0x5f, 0xc0, 0xfb, 0x8d, 0xe2, 0xd7, 0x53, 0x46, 0x9f, 0x43, 0xbb, 0x74, 0x11, 0x01,
	0x74, 0x8e, 0xd2, 0x97, 0xa4, 0x50, 0xfe, 0xad, 0x72, 0x8d, 0x09, 0x8f, 0x45, 0xe6, 0x3b, 0xe8,
	0x0e, 0xf4, 0x9e, 0x5e, 0x94, 0xed, 0x25, 0xa9, 0xdf, 0x1a, 0xfd, 0xe5, 0x82, 0x87, 0x69, 0x44,
	0xd9, 0x9a, 0x4a, 0xd3, 0x55, 0xf4, 0x35, 0x40, 0x19, 0x82, 0x50, 0x12, 0x9e, 0x18, 0xed, 0xfd,
	0xe9, 0x70, 0xdb, 0x0e, 0x93, 0xa6, 0x80, 0x53, 0x1d, 0xcc, 0x85, 0xd4, 0xb8, 0xc4, 0xe1, 0xbd,
	0xbc, 0x5e, 0xa2, 0xaf, 0xa0, 0x93, 0x32, 0xa5, 0x29, 0xb7, 0x4d, 0xfb, 0xb8, 0x81, 0x7c, 0x3a,
	0x9f, 0xc9, 0x13, 0x91, 0x11, 0xc6, 0xb1, 0x25, 0xa0, 0x9f, 0xe1, 0x1d, 0xb2, 0xa9, 0x37, 0x54,
	0xb6, 0x60, 0xdb, 0x93, 0x47, 0x37, 0xe8, 0x09, 0x46, 0x64, 0x37, 0x98, 0x4b, 0x38, 0x54, 0x5a,
	0x52, 0x92, 0x85, 0x8a, 0x6a, 0xcd, 0x78, 0xa2, 0xfa, 0xed, 0x5d, 0xe5, 0xcd, 0x18, 0x04, 0xf5,
	0x18, 0x04, 0x8b, 0x8a, 0x65, 0xfc, 0xc1, 0x9e, 0xd1, 0x58, 0x58, 0x09, 0xf4, 0x0d, 0x7c, 0x28,
	0x8d, 0x83, 0xa1, 0x90, 0x2c, 0x61, 0x9c, 0xa4, 0xe1, 0xd6, 0x48, 0xf6, 0x6f, 0x0f, 0x9d, 0x71,
	0x0f, 0x0f, 0x2c, 0x66, 0x66, 0x21, 0x27, 0x97, 0x08, 0x34, 0x87, 0xc3, 0xb8, 0xf2, 0x21, 0x14,
	0x6b, 0x2a, 0x25, 0x8b, 0x69, 0xbf, 0x3b, 0x74, 0xc7, 0xde, 0xf4, 0x41, 0x63, 0xc5, 0xdf, 0x73,
	0xf1, 0x92, 0xcf, 0xcb, 0xb1, 0x8c, 0x44, 0xaa, 0xb0, 0x67, 0xf8, 0x33, 0x4b, 0xff, 0xae, 0xdd,
	0xeb, 0xf8, 0xdd, 0xd1, 0x9f, 0x0e, 0xdc, 0xb3, 0x13, 0xfb, 0x8c, 0xf0, 0x38, 0xdd, 0xb4, 0xd8,
	0x07, 0x57, 0x93, 0xa4, 0xea, 0xed, 0x1e, 0x2e, 0x97, 0x68, 0x01, 0x77, 0xed, 0x01, 0xe5, 0xa5,
	0x39, 0xa6, 0x7d, 0x9f, 0x5d, 0xd3, 0x3e, 0x73, 0xa3, 0x55, 0xe3, 0x1a, 0x9f, 0x99, 0x0b, 0x0d,
	0xfb, 0xb5, 0xc0, 0xc6, 0x99, 0x33, 0xf0, 0xaa, 0x03, 0x5f, 0x2a, 0xba, 0x37, 0x52, 0x3c, 0xa8,
	0xd8, 0xb5, 0xdc, 0xc8, 0x07, 0x6f, 0xb6, 0xd2, 0xdb, 0x17, 0xd0, 0x1f, 0x2d, 0xb8, 0xb3, 0xa0,
	0x3c, 0xde, 0x14, 0xf6, 0x04, 0xdc, 0x35, 0x23, 0x36, 0xb4, 0xff, 0x22, 0x77, 0x25, 0xfa, 0xba,
	0x58, 0xb4, 0xde, 0x3e, 0x16, 0x3f, 0x36, 0x14, 0xff, 0xf0, 0x0d, 0xa2, 0xf3, 0x92, 0x64, 0x35,
	0xaf, 0x1a, 0x80, 0x9e, 0x03, 0xca, 0x56, 0xa9, 0x66, 0x79, 0x4a, 0x2f, 0x5e, 0x1b, 0xe1, 0x2b,
	0x51, 0x39, 0xab, 0x29, 0x8c, 0x27, 0x56, 0xf7, 0xee, 0x46, 0x66, 0x63, 0xee, 0xdf, 0x0e, 0xbc,
	0x5b, 0xbb, 0xfb, 0xa6, 0xb0, 0xcc, 0xe0, 0x50, 0x55, 0xae, 0xff, 0xd7, 0xa8, 0x78, 0x86, 0xfe,
	0x3f, 0x05, 0x05, 0xbd, 0x07, 0x1d, 0x7a, 0x91, 0x33, 0x49, 0x2b, 0x6f, 0x5c, 0x6c, 0x77, 0xa8,
	0x0f, 0xdd, 0x52, 0x84, 0x72, 0x5d, 0x0d, 0xe5, 0x1e, 0xae, 0xb7, 0xa3, 0x39, 0xa0, 0x5d, 0x9b,
	0x4a, 0x3c, 0xe5, 0xe4, 0x3c, 0xa5, 0x71, 0x55, 0x7d, 0x0f, 0xd7, 0x5b, 0x34, 0xdc, 0x7d, 0x9c,
	0x0e, 0xae, 0xbc, 0x28, 0x0f, 0x3f, 0x01, 0xef, 0xea, 0x8c, 0xa2, 0x1e, 0xb4, 0x9f, 0x2d, 0x97,
	0x73, 0xff, 0x16, 0xea, 0x82, 0xbb, 0xfc, 0x61, 0xe1, 0x3b, 0xdf, 0x1e, 0xc3, 0x07, 0x91, 0xc8,
	0x9a, 0x3a, 0x37, 0x77, 0x9e, 0xf7, 0xea, 0xf5, 0x6f, 0xad, 0xfb, 0x3f, 0x4d, 0x31, 0x29, 0x82,
	0xe3, 0x12, 0x75, 0x94, 0xe7, 0x26, 0x27, 0x19, 0xe1, 0xe7, 0x9d, 0xea, 0x75, 0x7e, 0xf2, 0x4f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x17, 0x26, 0x33, 0xf4, 0xc0, 0x08, 0x00, 0x00,
}
