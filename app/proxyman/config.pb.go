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

func init() {
	proto.RegisterType((*InboundConfig)(nil), "v2ray.core.app.proxyman.InboundConfig")
	proto.RegisterType((*AllocationStrategy)(nil), "v2ray.core.app.proxyman.AllocationStrategy")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyConcurrency)(nil), "v2ray.core.app.proxyman.AllocationStrategy.AllocationStrategyConcurrency")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyRefresh)(nil), "v2ray.core.app.proxyman.AllocationStrategy.AllocationStrategyRefresh")
	proto.RegisterType((*StreamReceiverConfig)(nil), "v2ray.core.app.proxyman.StreamReceiverConfig")
	proto.RegisterType((*DatagramReceiverConfig)(nil), "v2ray.core.app.proxyman.DatagramReceiverConfig")
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.app.proxyman.InboundHandlerConfig")
	proto.RegisterType((*OutboundConfig)(nil), "v2ray.core.app.proxyman.OutboundConfig")
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 601 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe4, 0x54, 0x5f, 0x6f, 0xd3, 0x3e,
	0x14, 0xfd, 0xb5, 0xdd, 0x6f, 0x6c, 0xb7, 0xac, 0x2b, 0x66, 0x62, 0xa5, 0x08, 0xa9, 0x54, 0x08,
	0x2a, 0x81, 0x9c, 0x11, 0xc4, 0x03, 0x4f, 0x68, 0xff, 0x24, 0xf6, 0x30, 0x56, 0xb9, 0x13, 0x0f,
	0x08, 0xa9, 0xba, 0x4b, 0xbc, 0x10, 0x91, 0xd8, 0x96, 0xed, 0x8e, 0xe5, 0x2b, 0xf1, 0x29, 0x78,
	0xe2, 0x89, 0x4f, 0xc3, 0x27, 0x40, 0x89, 0x93, 0xae, 0x5a, 0x5b, 0xc4, 0xc4, 0x23, 0x6f, 0x8e,
	0x73, 0xce, 0xb1, 0xcf, 0xb9, 0xd7, 0x17, 0x06, 0x17, 0xbe, 0xc6, 0x8c, 0x06, 0x32, 0xf5, 0x02,
	0xa9, 0xb9, 0x87, 0x4a, 0x79, 0x4a, 0xcb, 0xcb, 0x2c, 0x45, 0xe1, 0x05, 0x52, 0x9c, 0xc7, 0x11,
	0x55, 0x5a, 0x5a, 0x49, 0xb6, 0x2b, 0xa4, 0xe6, 0x14, 0x95, 0xa2, 0x15, 0xaa, 0xbb, 0x73, 0x4d,
	0x22, 0x90, 0x69, 0x2a, 0x85, 0x67, 0xb8, 0x8e, 0x31, 0xf1, 0x6c, 0xa6, 0x78, 0x38, 0x4e, 0xb9,
	0x31, 0x18, 0x71, 0x27, 0xd5, 0x7d, 0xba, 0x98, 0x21, 0xb8, 0xf5, 0x30, 0x0c, 0x35, 0x37, 0xa6,
	0x04, 0x3e, 0x5e, 0x0e, 0x54, 0x52, 0xdb, 0x12, 0x45, 0xaf, 0xa1, 0xac, 0x46, 0x61, 0xf2, 0xff,
	0x5e, 0x2c, 0x2c, 0xd7, 0x39, 0x7a, 0xd6, 0x49, 0x7f, 0x13, 0x36, 0x8e, 0xc4, 0x99, 0x9c, 0x88,
	0x70, 0xbf, 0xd8, 0xee, 0x7f, 0x6b, 0x00, 0xd9, 0x4d, 0x12, 0x19, 0xa0, 0x8d, 0xa5, 0x18, 0x59,
	0x8d, 0x96, 0x47, 0x19, 0x39, 0x80, 0x95, 0xfc, 0xf6, 0x9d, 0x5a, 0xaf, 0x36, 0x68, 0xf9, 0x3b,
	0x74, 0x49, 0x00, 0x74, 0x9e, 0x4a, 0x4f, 0x33, 0xc5, 0x59, 0xc1, 0x26, 0x9f, 0xa1, 0x19, 0x48,
	0x11, 0x4c, 0xb4, 0xe6, 0x22, 0xc8, 0x3a, 0xf5, 0x5e, 0x6d, 0xd0, 0xf4, 0x8f, 0x6e, 0x22, 0x36,
	0xbf, 0xb5, 0x7f, 0x25, 0xc8, 0x66, 0xd5, 0xc9, 0x18, 0x6e, 0x69, 0x7e, 0xae, 0xb9, 0xf9, 0xd4,
	0x69, 0x14, 0x07, 0x1d, 0xfe, 0xdd, 0x41, 0xcc, 0x89, 0xb1, 0x4a, 0xb5, 0xfb, 0x0a, 0x1e, 0xfe,
	0xf6, 0x3a, 0x64, 0x0b, 0xfe, 0xbf, 0xc0, 0x64, 0xe2, 0x52, 0xdb, 0x60, 0xee, 0xa3, 0xfb, 0x02,
	0xee, 0x2f, 0x15, 0x5f, 0x4c, 0xe9, 0x3f, 0x87, 0x95, 0x3c, 0x45, 0x02, 0xb0, 0xba, 0x9b, 0x7c,
	0xc1, 0xcc, 0xb4, 0xff, 0xcb, 0xd7, 0x0c, 0x45, 0x28, 0xd3, 0x76, 0x8d, 0xdc, 0x86, 0xb5, 0xc3,
	0xcb, 0xbc, 0xbc, 0x98, 0xb4, 0xeb, 0xfd, 0xef, 0x75, 0xd8, 0x1a, 0x59, 0xcd, 0x31, 0x65, 0x3c,
	0xe0, 0xf1, 0x05, 0xd7, 0xae, 0xb6, 0xe4, 0x0d, 0x40, 0xde, 0x0a, 0x63, 0x8d, 0x22, 0x72, 0x27,
	0x34, 0xfd, 0xde, 0x6c, 0x28, 0xae, 0xa7, 0xa8, 0xe0, 0x96, 0x0e, 0xa5, 0xb6, 0x2c, 0xc7, 0xb1,
	0x75, 0x55, 0x2d, 0xc9, 0x6b, 0x58, 0x4d, 0x62, 0x63, 0xb9, 0x28, 0x4b, 0xf7, 0x68, 0x09, 0xf9,
	0x68, 0x78, 0xa2, 0x0f, 0x64, 0x8a, 0xb1, 0x60, 0x25, 0x81, 0x7c, 0x84, 0xbb, 0x38, 0x75, 0x3d,
	0x36, 0xa5, 0xed, 0xb2, 0x32, 0xcf, 0x6e, 0x50, 0x19, 0x46, 0x70, 0xbe, 0x3d, 0x4f, 0x61, 0xd3,
	0x14, 0x8e, 0xc7, 0x86, 0x5b, 0x1b, 0x8b, 0xc8, 0x74, 0x56, 0xe6, 0x95, 0xa7, 0x8f, 0x81, 0x56,
	0x8f, 0x81, 0xba, 0x9c, 0x5c, 0x3e, 0xac, 0xe5, 0x34, 0x46, 0xa5, 0x44, 0xff, 0x67, 0x0d, 0xee,
	0x1d, 0xa0, 0xc5, 0x48, 0xff, 0x3b, 0x51, 0xf6, 0x7f, 0xd4, 0x60, 0xab, 0x1c, 0x09, 0x6f, 0x51,
	0x84, 0xc9, 0xd4, 0x72, 0x1b, 0x1a, 0x16, 0xa3, 0xc2, 0xeb, 0x3a, 0xcb, 0x97, 0x64, 0x04, 0x77,
	0x74, 0x19, 0xcb, 0x55, 0xee, 0xf5, 0x5e, 0x63, 0xd0, 0xf4, 0x9f, 0x2c, 0xb0, 0xe3, 0xa6, 0x60,
	0x31, 0x0f, 0xc2, 0x63, 0x37, 0x04, 0x59, 0xbb, 0x12, 0xa8, 0x42, 0x27, 0xc7, 0xd0, 0x2a, 0xae,
	0x7c, 0xa5, 0xe8, 0x8c, 0xfd, 0xa9, 0xe2, 0x46, 0xc1, 0x9e, 0xd6, 0xb0, 0x0d, 0xad, 0x93, 0x89,
	0x9d, 0x99, 0x70, 0x7b, 0xef, 0xe0, 0x41, 0x20, 0xd3, 0x65, 0x31, 0xed, 0x35, 0x1d, 0x6c, 0x98,
	0x8f, 0xc7, 0x0f, 0x6b, 0xd5, 0xf6, 0xd7, 0xfa, 0xf6, 0x7b, 0x9f, 0x61, 0x46, 0xf7, 0x73, 0xc2,
	0xae, 0x52, 0x74, 0x58, 0xfe, 0x39, 0x5b, 0x2d, 0x26, 0xe9, 0xcb, 0x5f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x94, 0x72, 0xcb, 0xa7, 0x3f, 0x06, 0x00, 0x00,
}
