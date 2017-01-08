package receiver

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
func (AllocationStrategy_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type AllocationStrategy struct {
	Type AllocationStrategy_Type `protobuf:"varint,1,opt,name=type,enum=v2ray.core.app.receiver.AllocationStrategy_Type" json:"type,omitempty"`
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
func (*AllocationStrategy) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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
	return fileDescriptor0, []int{0, 0}
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
	return fileDescriptor0, []int{0, 1}
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
func (*StreamReceiverConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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
func (*DatagramReceiverConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

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

type PerProxyConfig struct {
	Tag      string                                   `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	Settings []*v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,rep,name=settings" json:"settings,omitempty"`
}

func (m *PerProxyConfig) Reset()                    { *m = PerProxyConfig{} }
func (m *PerProxyConfig) String() string            { return proto.CompactTextString(m) }
func (*PerProxyConfig) ProtoMessage()               {}
func (*PerProxyConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PerProxyConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *PerProxyConfig) GetSettings() []*v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.Settings
	}
	return nil
}

type Config struct {
	Settings []*PerProxyConfig `protobuf:"bytes,1,rep,name=settings" json:"settings,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Config) GetSettings() []*PerProxyConfig {
	if m != nil {
		return m.Settings
	}
	return nil
}

func init() {
	proto.RegisterType((*AllocationStrategy)(nil), "v2ray.core.app.receiver.AllocationStrategy")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyConcurrency)(nil), "v2ray.core.app.receiver.AllocationStrategy.AllocationStrategyConcurrency")
	proto.RegisterType((*AllocationStrategy_AllocationStrategyRefresh)(nil), "v2ray.core.app.receiver.AllocationStrategy.AllocationStrategyRefresh")
	proto.RegisterType((*StreamReceiverConfig)(nil), "v2ray.core.app.receiver.StreamReceiverConfig")
	proto.RegisterType((*DatagramReceiverConfig)(nil), "v2ray.core.app.receiver.DatagramReceiverConfig")
	proto.RegisterType((*PerProxyConfig)(nil), "v2ray.core.app.receiver.PerProxyConfig")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.receiver.Config")
	proto.RegisterEnum("v2ray.core.app.receiver.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/receiver/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 572 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe4, 0x94, 0xd1, 0x6e, 0xd3, 0x3e,
	0x14, 0xc6, 0xff, 0x69, 0xf7, 0x2f, 0xdb, 0x29, 0x8c, 0xca, 0x4c, 0xac, 0x14, 0x21, 0x95, 0x08,
	0xb1, 0x4a, 0x20, 0x67, 0x04, 0x71, 0xc1, 0x15, 0xda, 0xda, 0x5d, 0xec, 0x62, 0x50, 0xb9, 0x13,
	0x17, 0x08, 0xa9, 0x32, 0xa9, 0x1b, 0x22, 0x12, 0x3b, 0x3a, 0xf6, 0xca, 0xf2, 0x4a, 0x3c, 0x05,
	0x57, 0x3c, 0x10, 0x4f, 0x80, 0x12, 0x27, 0x5d, 0x59, 0x1b, 0xa4, 0x89, 0x4b, 0xee, 0x5c, 0xf7,
	0xfb, 0x7e, 0xf6, 0x77, 0x8e, 0x73, 0x60, 0xb0, 0xf0, 0x91, 0x67, 0x34, 0x50, 0x89, 0x17, 0x28,
	0x14, 0x1e, 0x4f, 0x53, 0x0f, 0x45, 0x20, 0xa2, 0x85, 0x40, 0x2f, 0x50, 0x72, 0x1e, 0x85, 0x34,
	0x45, 0x65, 0x14, 0xd9, 0xaf, 0x94, 0x28, 0x28, 0x4f, 0x53, 0x5a, 0xa9, 0x7a, 0x87, 0xd7, 0x10,
	0x81, 0x4a, 0x12, 0x25, 0x3d, 0x2d, 0x30, 0xe2, 0xb1, 0x67, 0xb2, 0x54, 0xcc, 0xa6, 0x89, 0xd0,
	0x9a, 0x87, 0xc2, 0xa2, 0x7a, 0x07, 0x9b, 0x1d, 0x52, 0x18, 0x8f, 0xcf, 0x66, 0x28, 0xb4, 0x2e,
	0x85, 0x4f, 0xea, 0x85, 0xa9, 0x42, 0x53, 0xaa, 0xe8, 0x35, 0x95, 0x41, 0x2e, 0x75, 0xfe, 0xbf,
	0x17, 0x49, 0x23, 0x30, 0x57, 0xaf, 0x26, 0x71, 0xbf, 0x37, 0x81, 0x1c, 0xc5, 0xb1, 0x0a, 0xb8,
	0x89, 0x94, 0x9c, 0x18, 0xe4, 0x46, 0x84, 0x19, 0x19, 0xc1, 0x56, 0x7e, 0xd9, 0xae, 0xd3, 0x77,
	0x06, 0xbb, 0xfe, 0x21, 0xad, 0xc9, 0x4b, 0xd7, 0xad, 0xf4, 0x3c, 0x4b, 0x05, 0x2b, 0xdc, 0xe4,
	0x0b, 0xb4, 0x03, 0x25, 0x83, 0x0b, 0x44, 0x21, 0x83, 0xac, 0xdb, 0xe8, 0x3b, 0x83, 0xb6, 0x7f,
	0x7a, 0x13, 0xd8, 0xfa, 0xd6, 0xf0, 0x0a, 0xc8, 0x56, 0xe9, 0x64, 0x0a, 0xb7, 0x50, 0xcc, 0x51,
	0xe8, 0xcf, 0xdd, 0x66, 0x71, 0xd0, 0xc9, 0xdf, 0x1d, 0xc4, 0x2c, 0x8c, 0x55, 0xd4, 0xde, 0x2b,
	0x78, 0xf4, 0xc7, 0xeb, 0x90, 0x3d, 0xf8, 0x7f, 0xc1, 0xe3, 0x0b, 0x5b, 0xb5, 0x3b, 0xcc, 0xfe,
	0xe8, 0xbd, 0x80, 0x07, 0xb5, 0xf0, 0xcd, 0x16, 0xf7, 0x39, 0x6c, 0xe5, 0x55, 0x24, 0x00, 0xad,
	0xa3, 0xf8, 0x2b, 0xcf, 0x74, 0xe7, 0xbf, 0x7c, 0xcd, 0xb8, 0x9c, 0xa9, 0xa4, 0xe3, 0x90, 0xdb,
	0xb0, 0x7d, 0x72, 0x99, 0x77, 0x93, 0xc7, 0x9d, 0x86, 0xfb, 0xa3, 0x01, 0x7b, 0x13, 0x83, 0x82,
	0x27, 0xac, 0x0c, 0x38, 0x2c, 0x3a, 0x4c, 0xde, 0x00, 0xe4, 0x9d, 0x9f, 0x22, 0x97, 0xa1, 0x3d,
	0xa1, 0xed, 0xf7, 0x57, 0x8b, 0x62, 0x9f, 0x10, 0x95, 0xc2, 0xd0, 0xb1, 0x42, 0xc3, 0x72, 0x1d,
	0xdb, 0x49, 0xab, 0x25, 0x79, 0x0d, 0xad, 0x38, 0xd2, 0x46, 0xc8, 0xb2, 0x75, 0x8f, 0x6b, 0xcc,
	0xa7, 0xe3, 0x77, 0x38, 0x52, 0x09, 0x8f, 0x24, 0x2b, 0x0d, 0xe4, 0x23, 0xdc, 0xe3, 0xcb, 0xd4,
	0x53, 0x5d, 0xc6, 0x2e, 0x3b, 0xf3, 0xec, 0x06, 0x9d, 0x61, 0x84, 0xaf, 0x3f, 0xcf, 0x73, 0xb8,
	0xab, 0x8b, 0xc4, 0x53, 0x2d, 0x8c, 0x89, 0x64, 0xa8, 0xbb, 0x5b, 0xeb, 0xe4, 0xe5, 0xdb, 0xa7,
	0xd5, 0xdb, 0xa7, 0xb6, 0x4e, 0xb6, 0x3e, 0x6c, 0xd7, 0x32, 0x26, 0x25, 0xc2, 0xfd, 0xe9, 0xc0,
	0xfd, 0x11, 0x37, 0x3c, 0xc4, 0x7f, 0xa7, 0x94, 0xee, 0x1c, 0x76, 0xc7, 0x02, 0xc7, 0xa8, 0x2e,
	0xb3, 0x32, 0x6b, 0x07, 0x9a, 0x86, 0x87, 0x45, 0xc8, 0x1d, 0x96, 0x2f, 0xc9, 0x31, 0x6c, 0x2f,
	0xeb, 0xdc, 0xe8, 0x37, 0x07, 0x6d, 0xff, 0xe9, 0x86, 0xeb, 0xdb, 0x21, 0x57, 0x7c, 0xff, 0xb3,
	0x33, 0x3b, 0xe3, 0xd8, 0xd2, 0xe7, 0x9e, 0x41, 0xab, 0xe4, 0x0f, 0x57, 0x68, 0x4e, 0x41, 0x3b,
	0xa8, 0x0d, 0xf1, 0xfb, 0xd5, 0xae, 0x70, 0xc7, 0x6f, 0xe1, 0x61, 0xa0, 0x92, 0x3a, 0xdf, 0x71,
	0xdb, 0x1a, 0xc6, 0xf9, 0x8c, 0xfb, 0xb0, 0x5d, 0x6d, 0x7f, 0x6b, 0xec, 0xbf, 0xf7, 0x19, 0xcf,
	0xe8, 0x30, 0x37, 0x1c, 0xa5, 0x29, 0xad, 0xda, 0xfc, 0xa9, 0x55, 0x8c, 0xc3, 0x97, 0xbf, 0x02,
	0x00, 0x00, 0xff, 0xff, 0x3a, 0x23, 0xe2, 0x82, 0x04, 0x06, 0x00, 0x00,
}
