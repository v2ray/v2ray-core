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

type ReceiverConfig struct {
	PortRange                  *v2ray_core_common_net1.PortRange           `protobuf:"bytes,1,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
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
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 737 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x55, 0xd1, 0x6e, 0xdb, 0x36,
	0x14, 0x9d, 0xec, 0xc4, 0x49, 0xae, 0x13, 0xc7, 0xe3, 0xb2, 0xc4, 0xf3, 0x36, 0xc0, 0x33, 0x86,
	0xcd, 0xd8, 0x06, 0x29, 0x73, 0x30, 0x60, 0x7b, 0xda, 0x12, 0x27, 0xc0, 0xf2, 0x90, 0xd9, 0xa3,
	0x83, 0x3e, 0x14, 0x05, 0x04, 0x46, 0x62, 0x5c, 0xa1, 0x12, 0x29, 0x90, 0xb4, 0x13, 0xfd, 0x52,
	0xbf, 0xa1, 0x0f, 0xfd, 0x80, 0x7e, 0x4a, 0xff, 0xa0, 0x2f, 0x05, 0x49, 0xc9, 0x36, 0x62, 0xbb,
	0x69, 0x1a, 0xf4, 0x4d, 0x24, 0xcf, 0x39, 0xe4, 0x3d, 0x87, 0x97, 0x82, 0xce, 0xa4, 0x2b, 0x48,
	0xe6, 0x06, 0x3c, 0xf1, 0x02, 0x2e, 0xa8, 0x47, 0xd2, 0xd4, 0x4b, 0x05, 0xbf, 0xcd, 0x12, 0xc2,
	0xbc, 0x80, 0xb3, 0xeb, 0x68, 0xe4, 0xa6, 0x82, 0x2b, 0x8e, 0x0e, 0x0a, 0xa4, 0xa0, 0x2e, 0x49,
	0x53, 0xb7, 0x40, 0x35, 0x0f, 0xef, 0x48, 0x04, 0x3c, 0x49, 0x38, 0xf3, 0x24, 0x15, 0x11, 0x89,
	0x3d, 0x95, 0xa5, 0x34, 0xf4, 0x13, 0x2a, 0x25, 0x19, 0x51, 0x2b, 0xd5, 0xfc, 0x79, 0x39, 0x83,
	0x51, 0xe5, 0x91, 0x30, 0x14, 0x54, 0xca, 0x1c, 0xf8, 0xe3, 0x6a, 0x60, 0xca, 0x85, 0xca, 0x51,
	0xee, 0x1d, 0x94, 0x12, 0x84, 0x49, 0xbd, 0xee, 0x45, 0x4c, 0x51, 0xa1, 0xd1, 0xf3, 0x95, 0xb4,
	0x77, 0x61, 0xe7, 0x9c, 0x5d, 0xf1, 0x31, 0x0b, 0x7b, 0x66, 0xba, 0xfd, 0xba, 0x0c, 0xe8, 0x38,
	0x8e, 0x79, 0x40, 0x54, 0xc4, 0xd9, 0x50, 0x09, 0xa2, 0xe8, 0x28, 0x43, 0xa7, 0xb0, 0xa6, 0x4f,
	0xdf, 0x70, 0x5a, 0x4e, 0xa7, 0xd6, 0x3d, 0x74, 0x57, 0x18, 0xe0, 0x2e, 0x52, 0xdd, 0xcb, 0x2c,
	0xa5, 0xd8, 0xb0, 0xd1, 0x0b, 0xa8, 0x06, 0x9c, 0x05, 0x63, 0x21, 0x28, 0x0b, 0xb2, 0x46, 0xa9,
	0xe5, 0x74, 0xaa, 0xdd, 0xf3, 0x87, 0x88, 0x2d, 0x4e, 0xf5, 0x66, 0x82, 0x78, 0x5e, 0x1d, 0xf9,
	0xb0, 0x21, 0xe8, 0xb5, 0xa0, 0xf2, 0x79, 0xa3, 0x6c, 0x36, 0x3a, 0x7b, 0xdc, 0x46, 0xd8, 0x8a,
	0xe1, 0x42, 0xb5, 0xf9, 0x07, 0x7c, 0xff, 0xc1, 0xe3, 0xa0, 0x3d, 0x58, 0x9f, 0x90, 0x78, 0x6c,
	0x5d, 0xdb, 0xc1, 0x76, 0xd0, 0xfc, 0x1d, 0xbe, 0x59, 0x29, 0xbe, 0x9c, 0xd2, 0xfe, 0x0d, 0xd6,
	0xb4, 0x8b, 0x08, 0xa0, 0x72, 0x1c, 0xdf, 0x90, 0x4c, 0xd6, 0xbf, 0xd0, 0xdf, 0x98, 0xb0, 0x90,
	0x27, 0x75, 0x07, 0x6d, 0xc3, 0xe6, 0xd9, 0xad, 0x8e, 0x97, 0xc4, 0xf5, 0x52, 0xfb, 0x55, 0x19,
	0x6a, 0x98, 0x06, 0x34, 0x9a, 0x50, 0x61, 0x53, 0x45, 0x7f, 0x03, 0xe8, 0x4b, 0xe0, 0x0b, 0xc2,
	0x46, 0x56, 0xbb, 0xda, 0x6d, 0xcd, 0xdb, 0x61, 0x6f, 0x93, 0xcb, 0xa8, 0x72, 0x07, 0x5c, 0x28,
	0xac, 0x71, 0x78, 0x2b, 0x2d, 0x3e, 0xd1, 0x5f, 0x50, 0x89, 0x23, 0xa9, 0x28, 0xcb, 0x43, 0xfb,
	0x61, 0x05, 0xf9, 0x7c, 0xd0, 0x17, 0xa7, 0x3c, 0x21, 0x11, 0xc3, 0x39, 0x01, 0x3d, 0x83, 0xaf,
	0xc8, 0xb4, 0x5e, 0x5f, 0xe6, 0x05, 0xe7, 0x99, 0xfc, 0xfa, 0x80, 0x4c, 0x30, 0x22, 0x8b, 0x17,
	0xf3, 0x12, 0x76, 0xa5, 0x12, 0x94, 0x24, 0xbe, 0xa4, 0x4a, 0x45, 0x6c, 0x24, 0x1b, 0x6b, 0x8b,
	0xca, 0xd3, 0x36, 0x70, 0x8b, 0x36, 0x70, 0x87, 0x86, 0x65, 0xfd, 0xc1, 0x35, 0xab, 0x31, 0xcc,
	0x25, 0xd0, 0x3f, 0xf0, 0x9d, 0xb0, 0x0e, 0xfa, 0x5c, 0x44, 0xa3, 0x88, 0x91, 0xd8, 0x0f, 0xa9,
	0x54, 0x11, 0x33, 0xbb, 0x37, 0xd6, 0x5b, 0x4e, 0x67, 0x13, 0x37, 0x73, 0x4c, 0x3f, 0x87, 0x9c,
	0xce, 0x10, 0xe8, 0x4f, 0x68, 0xe8, 0xd3, 0xde, 0xf8, 0x29, 0x91, 0x52, 0xeb, 0x04, 0x9c, 0x31,
	0x1a, 0x18, 0x76, 0xc5, 0xb0, 0xf7, 0xcd, 0xfa, 0xc0, 0x2e, 0xf7, 0xa6, 0xab, 0xed, 0x37, 0x0e,
	0xec, 0xe5, 0x3d, 0xf9, 0x2f, 0x61, 0x61, 0x3c, 0x0d, 0xb1, 0x0e, 0x65, 0x45, 0x46, 0x26, 0xbd,
	0x2d, 0xac, 0x3f, 0xd1, 0x10, 0xbe, 0xcc, 0x8f, 0x20, 0x66, 0xe5, 0xdb, 0x80, 0x7e, 0x5a, 0x12,
	0x90, 0x7d, 0x86, 0x4c, 0x43, 0x86, 0x17, 0xf6, 0x15, 0xc2, 0xf5, 0x42, 0x60, 0x5a, 0xfb, 0x05,
	0xd4, 0x4c, 0x08, 0x33, 0xc5, 0xf2, 0x83, 0x14, 0x77, 0x0c, 0xbb, 0x90, 0x6b, 0xd7, 0xa1, 0xd6,
	0x1f, 0xab, 0xf9, 0x27, 0xe6, 0xad, 0x03, 0xdb, 0x43, 0xca, 0xc2, 0x69, 0x61, 0x47, 0x50, 0x9e,
	0x44, 0x24, 0xbf, 0x96, 0x1f, 0x71, 0xb3, 0x34, 0x7a, 0x59, 0xf0, 0xa5, 0xc7, 0x07, 0xff, 0xff,
	0x8a, 0xe2, 0x7f, 0xb9, 0x47, 0x74, 0xa0, 0x49, 0xb9, 0xe6, 0x1d, 0x03, 0xde, 0x39, 0xf0, 0x75,
	0xe1, 0xc0, 0x7d, 0x81, 0xf6, 0x61, 0x57, 0x1a, 0x67, 0x3e, 0x35, 0xce, 0x9a, 0xa5, 0x7f, 0xa6,
	0x30, 0xd1, 0x3e, 0x54, 0xe8, 0x6d, 0x1a, 0x09, 0x6a, 0x9a, 0xac, 0x8c, 0xf3, 0x11, 0x6a, 0xc0,
	0x86, 0x16, 0xa1, 0x4c, 0x99, 0xd6, 0xd8, 0xc2, 0xc5, 0xf0, 0xe4, 0x3f, 0xf8, 0x36, 0xe0, 0xc9,
	0xaa, 0x2e, 0x3f, 0xa9, 0x5a, 0x2b, 0x06, 0xfa, 0x67, 0xf4, 0x74, 0xb3, 0x98, 0x7e, 0x59, 0x3a,
	0x78, 0xd2, 0xc5, 0x24, 0x73, 0x7b, 0x9a, 0x70, 0x9c, 0xa6, 0xd6, 0xdf, 0x84, 0xb0, 0xab, 0x8a,
	0xf9, 0x6f, 0x1d, 0xbd, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x26, 0xa3, 0xe8, 0x2c, 0xad, 0x07, 0x00,
	0x00,
}
