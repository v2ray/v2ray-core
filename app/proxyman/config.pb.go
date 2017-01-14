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

type StreamReceiverConfig struct {
	PortRange          *v2ray_core_common_net1.PortRange           `protobuf:"bytes,1,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
	Listen             *v2ray_core_common_net.IPOrDomain           `protobuf:"bytes,2,opt,name=listen" json:"listen,omitempty"`
	AllocationStrategy *AllocationStrategy                         `protobuf:"bytes,3,opt,name=allocation_strategy,json=allocationStrategy" json:"allocation_strategy,omitempty"`
	StreamSettings     *v2ray_core_transport_internet.StreamConfig `protobuf:"bytes,4,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
}

func (m *StreamReceiverConfig) Reset()                    { *m = StreamReceiverConfig{} }
func (m *StreamReceiverConfig) String() string            { return proto.CompactTextString(m) }
func (*StreamReceiverConfig) ProtoMessage()               {}
func (*StreamReceiverConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *StreamReceiverConfig) GetPortRange() *v2ray_core_common_net1.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *StreamReceiverConfig) GetListen() *v2ray_core_common_net.IPOrDomain {
	if m != nil {
		return m.Listen
	}
	return nil
}

func (m *StreamReceiverConfig) GetAllocationStrategy() *AllocationStrategy {
	if m != nil {
		return m.AllocationStrategy
	}
	return nil
}

func (m *StreamReceiverConfig) GetStreamSettings() *v2ray_core_transport_internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

type DatagramReceiverConfig struct {
	PortRange          *v2ray_core_common_net1.PortRange `protobuf:"bytes,1,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
	Listen             *v2ray_core_common_net.IPOrDomain `protobuf:"bytes,2,opt,name=listen" json:"listen,omitempty"`
	AllocationStrategy *AllocationStrategy               `protobuf:"bytes,3,opt,name=allocation_strategy,json=allocationStrategy" json:"allocation_strategy,omitempty"`
}

func (m *DatagramReceiverConfig) Reset()                    { *m = DatagramReceiverConfig{} }
func (m *DatagramReceiverConfig) String() string            { return proto.CompactTextString(m) }
func (*DatagramReceiverConfig) ProtoMessage()               {}
func (*DatagramReceiverConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DatagramReceiverConfig) GetPortRange() *v2ray_core_common_net1.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *DatagramReceiverConfig) GetListen() *v2ray_core_common_net.IPOrDomain {
	if m != nil {
		return m.Listen
	}
	return nil
}

func (m *DatagramReceiverConfig) GetAllocationStrategy() *AllocationStrategy {
	if m != nil {
		return m.AllocationStrategy
	}
	return nil
}

type InboundHandlerConfig struct {
	Tag              string                                   `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	ReceiverSettings []*v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,rep,name=receiver_settings,json=receiverSettings" json:"receiver_settings,omitempty"`
	ProxySettings    *v2ray_core_common_serial.TypedMessage   `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
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

func (m *InboundHandlerConfig) GetReceiverSettings() []*v2ray_core_common_serial.TypedMessage {
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

type StreamSenderConfig struct {
	// Send traffic through the given IP. Only IP is allowed.
	Via            *v2ray_core_common_net.IPOrDomain           `protobuf:"bytes,1,opt,name=via" json:"via,omitempty"`
	StreamSettings *v2ray_core_transport_internet.StreamConfig `protobuf:"bytes,2,opt,name=stream_settings,json=streamSettings" json:"stream_settings,omitempty"`
	ProxySettings  *v2ray_core_transport_internet.ProxyConfig  `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
}

func (m *StreamSenderConfig) Reset()                    { *m = StreamSenderConfig{} }
func (m *StreamSenderConfig) String() string            { return proto.CompactTextString(m) }
func (*StreamSenderConfig) ProtoMessage()               {}
func (*StreamSenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *StreamSenderConfig) GetVia() *v2ray_core_common_net.IPOrDomain {
	if m != nil {
		return m.Via
	}
	return nil
}

func (m *StreamSenderConfig) GetStreamSettings() *v2ray_core_transport_internet.StreamConfig {
	if m != nil {
		return m.StreamSettings
	}
	return nil
}

func (m *StreamSenderConfig) GetProxySettings() *v2ray_core_transport_internet.ProxyConfig {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

type DatagramSenderConfig struct {
	// Send traffic through the given IP. Only IP is allowed.
	Via           *v2ray_core_common_net.IPOrDomain          `protobuf:"bytes,1,opt,name=via" json:"via,omitempty"`
	ProxySettings *v2ray_core_transport_internet.ProxyConfig `protobuf:"bytes,2,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
}

func (m *DatagramSenderConfig) Reset()                    { *m = DatagramSenderConfig{} }
func (m *DatagramSenderConfig) String() string            { return proto.CompactTextString(m) }
func (*DatagramSenderConfig) ProtoMessage()               {}
func (*DatagramSenderConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *DatagramSenderConfig) GetVia() *v2ray_core_common_net.IPOrDomain {
	if m != nil {
		return m.Via
	}
	return nil
}

func (m *DatagramSenderConfig) GetProxySettings() *v2ray_core_transport_internet.ProxyConfig {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

type OutboundHandlerConfig struct {
	Tag            string                                   `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	SenderSettings []*v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,rep,name=sender_settings,json=senderSettings" json:"sender_settings,omitempty"`
	ProxySettings  *v2ray_core_common_serial.TypedMessage   `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
}

func (m *OutboundHandlerConfig) Reset()                    { *m = OutboundHandlerConfig{} }
func (m *OutboundHandlerConfig) String() string            { return proto.CompactTextString(m) }
func (*OutboundHandlerConfig) ProtoMessage()               {}
func (*OutboundHandlerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *OutboundHandlerConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *OutboundHandlerConfig) GetSenderSettings() []*v2ray_core_common_serial.TypedMessage {
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

func init() {
	proto.RegisterType((*InboundConfig)(nil), "v2ray.core.app.proxyman.InboundConfig")
	proto.RegisterType((*AllocationStrategy)(nil), "v2ray.core.app.proxyman.AllocationStrategy")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyConcurrency)(nil), "v2ray.core.app.proxyman.AllocationStrategy.AllocationStrategyConcurrency")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyRefresh)(nil), "v2ray.core.app.proxyman.AllocationStrategy.AllocationStrategyRefresh")
	proto.RegisterType((*StreamReceiverConfig)(nil), "v2ray.core.app.proxyman.StreamReceiverConfig")
	proto.RegisterType((*DatagramReceiverConfig)(nil), "v2ray.core.app.proxyman.DatagramReceiverConfig")
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.app.proxyman.InboundHandlerConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterType((*StreamSenderConfig)(nil), "v2ray.core.app.proxyman.StreamSenderConfig")
	proto.RegisterType((*DatagramSenderConfig)(nil), "v2ray.core.app.proxyman.DatagramSenderConfig")
	proto.RegisterType((*OutboundHandlerConfig)(nil), "v2ray.core.app.proxyman.OutboundHandlerConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 683 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe4, 0x55, 0xdd, 0x6e, 0xd3, 0x30,
	0x18, 0x25, 0xe9, 0x18, 0xdb, 0x57, 0xd6, 0x15, 0x53, 0x58, 0x29, 0x42, 0x2a, 0x15, 0x82, 0x0a,
	0x50, 0x32, 0x3a, 0x71, 0xc1, 0x15, 0xda, 0x9f, 0xc4, 0x2e, 0xc6, 0x8a, 0x3b, 0x71, 0x81, 0x90,
	0x2a, 0x2f, 0xf1, 0x42, 0x44, 0x62, 0x47, 0xb6, 0x5b, 0x96, 0x97, 0xe1, 0x01, 0x78, 0x0a, 0xae,
	0x90, 0x90, 0x78, 0x9a, 0x3d, 0x01, 0x72, 0x9c, 0x74, 0xd5, 0xda, 0xb2, 0x8d, 0x8a, 0x2b, 0xee,
	0x5c, 0xf7, 0x3b, 0xc7, 0x3e, 0xe7, 0x7c, 0xfe, 0x02, 0xed, 0x61, 0x47, 0x90, 0xd4, 0xf1, 0x78,
	0xec, 0x7a, 0x5c, 0x50, 0x97, 0x24, 0x89, 0x9b, 0x08, 0x7e, 0x92, 0xc6, 0x84, 0xb9, 0x1e, 0x67,
	0xc7, 0x61, 0xe0, 0x24, 0x82, 0x2b, 0x8e, 0xd6, 0x8a, 0x4a, 0x41, 0x1d, 0x92, 0x24, 0x4e, 0x51,
	0xd5, 0x58, 0x3f, 0x47, 0xe1, 0xf1, 0x38, 0xe6, 0xcc, 0x95, 0x54, 0x84, 0x24, 0x72, 0x55, 0x9a,
	0x50, 0xbf, 0x1f, 0x53, 0x29, 0x49, 0x40, 0x0d, 0x55, 0xe3, 0xc9, 0x74, 0x04, 0xa3, 0xca, 0x25,
	0xbe, 0x2f, 0xa8, 0x94, 0x79, 0xe1, 0xa3, 0xd9, 0x85, 0x09, 0x17, 0x2a, 0xaf, 0x72, 0xce, 0x55,
	0x29, 0x41, 0x98, 0xd4, 0xff, 0xbb, 0x21, 0x53, 0x54, 0xe8, 0xea, 0x71, 0x25, 0xad, 0x55, 0x58,
	0xd9, 0x63, 0x47, 0x7c, 0xc0, 0xfc, 0xed, 0x6c, 0xbb, 0xf5, 0xbd, 0x04, 0x68, 0x33, 0x8a, 0xb8,
	0x47, 0x54, 0xc8, 0x59, 0x4f, 0x09, 0xa2, 0x68, 0x90, 0xa2, 0x1d, 0x58, 0xd0, 0xb7, 0xaf, 0x5b,
	0x4d, 0xab, 0x5d, 0xe9, 0xac, 0x3b, 0x33, 0x0c, 0x70, 0x26, 0xa1, 0xce, 0x61, 0x9a, 0x50, 0x9c,
	0xa1, 0xd1, 0x67, 0x28, 0x7b, 0x9c, 0x79, 0x03, 0x21, 0x28, 0xf3, 0xd2, 0xba, 0xdd, 0xb4, 0xda,
	0xe5, 0xce, 0xde, 0x55, 0xc8, 0x26, 0xb7, 0xb6, 0xcf, 0x08, 0xf1, 0x38, 0x3b, 0xea, 0xc3, 0x0d,
	0x41, 0x8f, 0x05, 0x95, 0x9f, 0xea, 0xa5, 0xec, 0xa0, 0xdd, 0xf9, 0x0e, 0xc2, 0x86, 0x0c, 0x17,
	0xac, 0x8d, 0x97, 0xf0, 0xe0, 0x8f, 0xd7, 0x41, 0x35, 0xb8, 0x3e, 0x24, 0xd1, 0xc0, 0xb8, 0xb6,
	0x82, 0xcd, 0x8f, 0xc6, 0x0b, 0xb8, 0x37, 0x93, 0x7c, 0x3a, 0xa4, 0xf5, 0x1c, 0x16, 0xb4, 0x8b,
	0x08, 0x60, 0x71, 0x33, 0xfa, 0x42, 0x52, 0x59, 0xbd, 0xa6, 0xd7, 0x98, 0x30, 0x9f, 0xc7, 0x55,
	0x0b, 0xdd, 0x84, 0xa5, 0xdd, 0x13, 0x1d, 0x2f, 0x89, 0xaa, 0x76, 0xeb, 0x87, 0x0d, 0xb5, 0x9e,
	0x12, 0x94, 0xc4, 0x98, 0x7a, 0x34, 0x1c, 0x52, 0x61, 0xb2, 0x45, 0xaf, 0x01, 0x74, 0x2b, 0xf4,
	0x05, 0x61, 0x81, 0x39, 0xa1, 0xdc, 0x69, 0x8e, 0x9b, 0x62, 0x7a, 0xca, 0x61, 0x54, 0x39, 0x5d,
	0x2e, 0x14, 0xd6, 0x75, 0x78, 0x39, 0x29, 0x96, 0xe8, 0x15, 0x2c, 0x46, 0xa1, 0x54, 0x94, 0xe5,
	0xd1, 0x3d, 0x9c, 0x01, 0xde, 0xeb, 0x1e, 0x88, 0x1d, 0x1e, 0x93, 0x90, 0xe1, 0x1c, 0x80, 0x3e,
	0xc2, 0x6d, 0x32, 0x52, 0xdd, 0x97, 0xb9, 0xec, 0x3c, 0x99, 0x67, 0x57, 0x48, 0x06, 0x23, 0x32,
	0xd9, 0x9e, 0x87, 0xb0, 0x2a, 0x33, 0xc5, 0x7d, 0x49, 0x95, 0x0a, 0x59, 0x20, 0xeb, 0x0b, 0x93,
	0xcc, 0xa3, 0xc7, 0xe0, 0x14, 0x8f, 0xc1, 0x31, 0x3e, 0x19, 0x7f, 0x70, 0xc5, 0x70, 0xf4, 0x72,
	0x8a, 0xd6, 0xa9, 0x05, 0x77, 0x77, 0x88, 0x22, 0x81, 0xf8, 0x7f, 0xac, 0x6c, 0xfd, 0xb2, 0xa0,
	0x96, 0x8f, 0x84, 0x37, 0x84, 0xf9, 0xd1, 0x48, 0x72, 0x15, 0x4a, 0x8a, 0x04, 0x99, 0xd6, 0x65,
	0xac, 0x97, 0xa8, 0x07, 0xb7, 0x44, 0x6e, 0xcb, 0x99, 0xef, 0x76, 0xb3, 0xd4, 0x2e, 0x77, 0x1e,
	0x4f, 0x91, 0x63, 0xa6, 0x60, 0x36, 0x0f, 0xfc, 0x7d, 0x33, 0x04, 0x71, 0xb5, 0x20, 0x28, 0x4c,
	0x47, 0xfb, 0x50, 0xc9, 0xae, 0x7c, 0xc6, 0x68, 0x84, 0x5d, 0x96, 0x71, 0x25, 0x43, 0x8f, 0x32,
	0xac, 0x42, 0xe5, 0x60, 0xa0, 0xc6, 0x27, 0xdc, 0xa9, 0x05, 0xa8, 0x97, 0x07, 0xcd, 0xfc, 0x91,
	0xbc, 0x0d, 0x28, 0x0d, 0x43, 0x92, 0x47, 0x79, 0x89, 0x34, 0x74, 0xf5, 0xb4, 0xbe, 0xb3, 0xe7,
	0xee, 0x3b, 0xf4, 0x6e, 0x86, 0x05, 0x4f, 0x2f, 0x20, 0xed, 0x6a, 0x50, 0xce, 0x79, 0xce, 0x86,
	0xaf, 0x16, 0xd4, 0x8a, 0x56, 0x9e, 0x5f, 0xf6, 0xe4, 0x05, 0xed, 0x79, 0x2f, 0xf8, 0xd3, 0x82,
	0x3b, 0x45, 0x50, 0x17, 0xf5, 0xdd, 0x01, 0xac, 0xca, 0x4c, 0xc3, 0xdf, 0x76, 0x5d, 0xc5, 0xc0,
	0xff, 0x51, 0xcf, 0x6d, 0xbd, 0x85, 0xfb, 0x1e, 0x8f, 0x67, 0x3d, 0xc4, 0xad, 0xb2, 0x11, 0xd6,
	0xd5, 0x1f, 0xe0, 0x0f, 0x4b, 0xc5, 0xf6, 0x37, 0x7b, 0xed, 0x7d, 0x07, 0x93, 0xd4, 0xd9, 0xd6,
	0x80, 0xcd, 0x24, 0x31, 0x6e, 0xc5, 0x84, 0x1d, 0x2d, 0x66, 0xdf, 0xea, 0x8d, 0xdf, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x18, 0xce, 0x1c, 0x24, 0xa1, 0x08, 0x00, 0x00,
}
