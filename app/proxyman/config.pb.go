package proxyman

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_serial "v2ray.com/core/common/serial"
import v2ray_core_common_net "v2ray.com/core/common/net"
import v2ray_core_common_net2 "v2ray.com/core/common/net"
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
	AllowPassiveConnection     bool                                        `protobuf:"varint,6,opt,name=allow_passive_connection,json=allowPassiveConnection" json:"allow_passive_connection,omitempty"`
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

func (m *ReceiverConfig) GetAllowPassiveConnection() bool {
	if m != nil {
		return m.AllowPassiveConnection
	}
	return false
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
	Via            *v2ray_core_common_net.IPOrDomain           `protobuf:"bytes,1,opt,name=via" json:"via,omitempty"`
	StreamSettings *v2ray_core_transport_internet.StreamConfig `protobuf:"bytes,2,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	ProxySettings  *v2ray_core_transport_internet.ProxyConfig  `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
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
	Enabled bool `protobuf:"varint,1,opt,name=enabled" json:"enabled,omitempty"`
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

type DispatchConfig struct {
	MuxSettings *MultiplexingConfig `protobuf:"bytes,1,opt,name=mux_settings,json=muxSettings" json:"mux_settings,omitempty"`
}

func (m *DispatchConfig) Reset()                    { *m = DispatchConfig{} }
func (m *DispatchConfig) String() string            { return proto.CompactTextString(m) }
func (*DispatchConfig) ProtoMessage()               {}
func (*DispatchConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *DispatchConfig) GetMuxSettings() *MultiplexingConfig {
	if m != nil {
		return m.MuxSettings
	}
	return nil
}

type SessionFrame struct {
	Id      uint32                           `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Target  *v2ray_core_common_net2.Endpoint `protobuf:"bytes,2,opt,name=target" json:"target,omitempty"`
	Payload []byte                           `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *SessionFrame) Reset()                    { *m = SessionFrame{} }
func (m *SessionFrame) String() string            { return proto.CompactTextString(m) }
func (*SessionFrame) ProtoMessage()               {}
func (*SessionFrame) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *SessionFrame) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SessionFrame) GetTarget() *v2ray_core_common_net2.Endpoint {
	if m != nil {
		return m.Target
	}
	return nil
}

func (m *SessionFrame) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
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
	proto.RegisterType((*DispatchConfig)(nil), "v2ray.core.app.proxyman.DispatchConfig")
	proto.RegisterType((*SessionFrame)(nil), "v2ray.core.app.proxyman.SessionFrame")
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 846 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x55, 0xdd, 0x6e, 0x23, 0x35,
	0x18, 0x65, 0x92, 0x6e, 0xda, 0x7e, 0x6d, 0xd3, 0x60, 0x96, 0xdd, 0x10, 0x40, 0x94, 0x08, 0x41,
	0xc5, 0xa2, 0xc9, 0x92, 0x15, 0x02, 0xae, 0xa0, 0xa4, 0x45, 0xf4, 0xa2, 0x34, 0x38, 0x2b, 0x2e,
	0x10, 0x52, 0x70, 0x67, 0xbc, 0xb3, 0x16, 0x33, 0xb6, 0xb1, 0x9d, 0x6e, 0xe6, 0x95, 0x78, 0x06,
	0x2e, 0x78, 0x00, 0x1e, 0x85, 0x37, 0xe0, 0x06, 0xf9, 0x67, 0x92, 0xa8, 0xc9, 0x50, 0xca, 0x8a,
	0xbb, 0x71, 0x7c, 0xce, 0xb1, 0xcf, 0x77, 0xbe, 0xcf, 0x81, 0xe3, 0xeb, 0xa1, 0x22, 0x65, 0x9c,
	0x88, 0x62, 0x90, 0x08, 0x45, 0x07, 0x44, 0xca, 0x81, 0x54, 0x62, 0x5e, 0x16, 0x84, 0x0f, 0x12,
	0xc1, 0x9f, 0xb1, 0x2c, 0x96, 0x4a, 0x18, 0x81, 0x1e, 0x56, 0x48, 0x45, 0x63, 0x22, 0x65, 0x5c,
	0xa1, 0x7a, 0x8f, 0x6f, 0x48, 0x24, 0xa2, 0x28, 0x04, 0x1f, 0x68, 0xaa, 0x18, 0xc9, 0x07, 0xa6,
	0x94, 0x34, 0x9d, 0x16, 0x54, 0x6b, 0x92, 0x51, 0x2f, 0xd5, 0xfb, 0x60, 0x33, 0x83, 0x53, 0x33,
	0x20, 0x69, 0xaa, 0xa8, 0xd6, 0x01, 0xf8, 0xa8, 0x1e, 0x98, 0x52, 0x6d, 0x18, 0x27, 0x86, 0x09,
	0x1e, 0xc0, 0xef, 0xd5, 0x83, 0xa5, 0x50, 0x26, 0xa0, 0xe2, 0x1b, 0x28, 0xa3, 0x08, 0xd7, 0x76,
	0x7f, 0xc0, 0xb8, 0xa1, 0xca, 0xa2, 0x57, 0x6d, 0xf7, 0x0f, 0xe1, 0xe0, 0x9c, 0x5f, 0x89, 0x19,
	0x4f, 0x47, 0xee, 0xe7, 0xfe, 0xef, 0x4d, 0x40, 0x27, 0x79, 0x2e, 0x12, 0x77, 0xf6, 0xc4, 0x28,
	0x62, 0x68, 0x56, 0xa2, 0x53, 0xd8, 0xb2, 0x56, 0xbb, 0xd1, 0x51, 0x74, 0xdc, 0x1e, 0x3e, 0x8e,
	0x6b, 0xaa, 0x15, 0xaf, 0x53, 0xe3, 0xa7, 0xa5, 0xa4, 0xd8, 0xb1, 0xd1, 0xcf, 0xb0, 0x97, 0x08,
	0x9e, 0xcc, 0x94, 0xa2, 0x3c, 0x29, 0xbb, 0x8d, 0xa3, 0xe8, 0x78, 0x6f, 0x78, 0x7e, 0x17, 0xb1,
	0xf5, 0x9f, 0x46, 0x4b, 0x41, 0xbc, 0xaa, 0x8e, 0xa6, 0xb0, 0xad, 0xe8, 0x33, 0x45, 0xf5, 0xf3,
	0x6e, 0xd3, 0x1d, 0x74, 0xf6, 0x72, 0x07, 0x61, 0x2f, 0x86, 0x2b, 0xd5, 0xde, 0x27, 0xf0, 0xf6,
	0x3f, 0x5e, 0x07, 0xdd, 0x87, 0x7b, 0xd7, 0x24, 0x9f, 0xf9, 0xaa, 0x1d, 0x60, 0xbf, 0xe8, 0x7d,
	0x0c, 0x6f, 0xd4, 0x8a, 0x6f, 0xa6, 0xf4, 0x3f, 0x82, 0x2d, 0x5b, 0x45, 0x04, 0xd0, 0x3a, 0xc9,
	0x5f, 0x90, 0x52, 0x77, 0x5e, 0xb1, 0xdf, 0x98, 0xf0, 0x54, 0x14, 0x9d, 0x08, 0xed, 0xc3, 0xce,
	0xd9, 0xdc, 0xc6, 0x4b, 0xf2, 0x4e, 0xa3, 0xff, 0x5b, 0x13, 0xda, 0x98, 0x26, 0x94, 0x5d, 0x53,
	0xe5, 0x53, 0x45, 0x5f, 0x00, 0xd8, 0x26, 0x98, 0x2a, 0xc2, 0x33, 0xaf, 0xbd, 0x37, 0x3c, 0x5a,
	0x2d, 0x87, 0xef, 0xa6, 0x98, 0x53, 0x13, 0x8f, 0x85, 0x32, 0xd8, 0xe2, 0xf0, 0xae, 0xac, 0x3e,
	0xd1, 0xe7, 0xd0, 0xca, 0x99, 0x36, 0x94, 0x87, 0xd0, 0xde, 0xad, 0x21, 0x9f, 0x8f, 0x2f, 0xd5,
	0xa9, 0x28, 0x08, 0xe3, 0x38, 0x10, 0xd0, 0x8f, 0xf0, 0x1a, 0x59, 0xf8, 0x9d, 0xea, 0x60, 0x38,
	0x64, 0xf2, 0xe8, 0x0e, 0x99, 0x60, 0x44, 0xd6, 0x1b, 0xf3, 0x29, 0x1c, 0x6a, 0xa3, 0x28, 0x29,
	0xa6, 0x9a, 0x1a, 0xc3, 0x78, 0xa6, 0xbb, 0x5b, 0xeb, 0xca, 0x8b, 0x31, 0x88, 0xab, 0x31, 0x88,
	0x27, 0x8e, 0xe5, 0xeb, 0x83, 0xdb, 0x5e, 0x63, 0x12, 0x24, 0xd0, 0x97, 0xf0, 0x96, 0xf2, 0x15,
	0x9c, 0x0a, 0xc5, 0x32, 0xc6, 0x49, 0x3e, 0x5d, 0x19, 0xc9, 0xee, 0xbd, 0xa3, 0xe8, 0x78, 0x07,
	0xf7, 0x02, 0xe6, 0x32, 0x40, 0x4e, 0x97, 0x08, 0xf4, 0x19, 0x74, 0xed, 0x6d, 0x5f, 0x4c, 0x25,
	0xd1, 0xda, 0xea, 0x24, 0x82, 0x73, 0x9a, 0x38, 0x76, 0xcb, 0xb1, 0x1f, 0xb8, 0xfd, 0xb1, 0xdf,
	0x1e, 0x2d, 0x76, 0xfb, 0x7f, 0x44, 0x70, 0x3f, 0xcc, 0xe4, 0x37, 0x84, 0xa7, 0xf9, 0x22, 0xc4,
	0x0e, 0x34, 0x0d, 0xc9, 0x5c, 0x7a, 0xbb, 0xd8, 0x7e, 0xa2, 0x09, 0xbc, 0x1a, 0xae, 0xa0, 0x96,
	0xf6, 0x7d, 0x40, 0xef, 0x6f, 0x08, 0xc8, 0xbf, 0x59, 0x6e, 0x20, 0xd3, 0x0b, 0xff, 0x64, 0xe1,
	0x4e, 0x25, 0xb0, 0xf0, 0x7e, 0x01, 0x6d, 0x17, 0xc2, 0x52, 0xb1, 0x79, 0x27, 0xc5, 0x03, 0xc7,
	0xae, 0xe4, 0xfa, 0x1d, 0x68, 0x5f, 0xce, 0xcc, 0xea, 0x13, 0xf3, 0x67, 0x04, 0xfb, 0x13, 0xca,
	0xd3, 0x85, 0xb1, 0x27, 0xd0, 0xbc, 0x66, 0x24, 0xb4, 0xe5, 0xbf, 0xe8, 0x2c, 0x8b, 0xde, 0x14,
	0x7c, 0xe3, 0xe5, 0x83, 0xff, 0xae, 0xc6, 0xfc, 0x87, 0xb7, 0x88, 0x8e, 0x2d, 0x29, 0x68, 0xde,
	0x28, 0xc0, 0x5f, 0x11, 0xbc, 0x5e, 0x55, 0xe0, 0xb6, 0x40, 0x2f, 0xe1, 0x50, 0xbb, 0xca, 0xfc,
	0xd7, 0x38, 0xdb, 0x9e, 0xfe, 0x3f, 0x85, 0x89, 0x1e, 0x40, 0x8b, 0xce, 0x25, 0x53, 0xd4, 0x0d,
	0x59, 0x13, 0x87, 0x15, 0xea, 0xc2, 0xb6, 0x15, 0xa1, 0xdc, 0xb8, 0xd1, 0xd8, 0xc5, 0xd5, 0xb2,
	0x1f, 0x03, 0xba, 0x98, 0xe5, 0x86, 0xc9, 0x9c, 0xce, 0x19, 0xcf, 0x82, 0xf3, 0x2e, 0x6c, 0x53,
	0x4e, 0xae, 0x72, 0x9a, 0x3a, 0xf7, 0x3b, 0xb8, 0x5a, 0xf6, 0x7f, 0x82, 0xf6, 0x29, 0xd3, 0x92,
	0x98, 0xe4, 0x79, 0xc0, 0x7e, 0x0b, 0xfb, 0xc5, 0x6c, 0xbe, 0x34, 0x10, 0xdd, 0xf2, 0x70, 0xac,
	0x1f, 0x87, 0xf7, 0x8a, 0xd9, 0x7c, 0x91, 0xc7, 0x2f, 0xb6, 0xfb, 0xb4, 0x66, 0x82, 0x7f, 0xad,
	0x48, 0x41, 0x51, 0x1b, 0x1a, 0x2c, 0x0d, 0xef, 0x6d, 0x83, 0xa5, 0xe8, 0x53, 0x68, 0x19, 0xa2,
	0x32, 0x6a, 0x42, 0xe9, 0xdf, 0xa9, 0x69, 0xc8, 0x33, 0x9e, 0x4a, 0xc1, 0xb8, 0xc1, 0x01, 0x6e,
	0x4d, 0x49, 0x52, 0xe6, 0x82, 0xa4, 0xae, 0xc8, 0xfb, 0xb8, 0x5a, 0x7e, 0x35, 0x82, 0x37, 0x13,
	0x51, 0xd4, 0xdd, 0x78, 0x1c, 0xfd, 0xb0, 0x53, 0x7d, 0xff, 0xda, 0x78, 0xf8, 0xfd, 0x10, 0x93,
	0x32, 0x1e, 0x59, 0xd4, 0x89, 0x94, 0xbe, 0xb3, 0x0a, 0xc2, 0xaf, 0x5a, 0xee, 0x1f, 0xfb, 0xc9,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xcf, 0x2c, 0xca, 0xd4, 0x08, 0x00, 0x00,
}
