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

type SessionFrame_FrameCommand int32

const (
	SessionFrame_SessionNew  SessionFrame_FrameCommand = 0
	SessionFrame_SessionKeep SessionFrame_FrameCommand = 1
	SessionFrame_SessionEnd  SessionFrame_FrameCommand = 2
)

var SessionFrame_FrameCommand_name = map[int32]string{
	0: "SessionNew",
	1: "SessionKeep",
	2: "SessionEnd",
}
var SessionFrame_FrameCommand_value = map[string]int32{
	"SessionNew":  0,
	"SessionKeep": 1,
	"SessionEnd":  2,
}

func (x SessionFrame_FrameCommand) String() string {
	return proto.EnumName(SessionFrame_FrameCommand_name, int32(x))
}
func (SessionFrame_FrameCommand) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{8, 0} }

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
	Enabled     bool   `protobuf:"varint,1,opt,name=enabled" json:"enabled,omitempty"`
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

type SessionFrame struct {
	Id      uint32                           `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Command SessionFrame_FrameCommand        `protobuf:"varint,2,opt,name=command,enum=v2ray.core.app.proxyman.SessionFrame_FrameCommand" json:"command,omitempty"`
	Target  *v2ray_core_common_net2.Endpoint `protobuf:"bytes,3,opt,name=target" json:"target,omitempty"`
	Payload []byte                           `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *SessionFrame) Reset()                    { *m = SessionFrame{} }
func (m *SessionFrame) String() string            { return proto.CompactTextString(m) }
func (*SessionFrame) ProtoMessage()               {}
func (*SessionFrame) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *SessionFrame) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SessionFrame) GetCommand() SessionFrame_FrameCommand {
	if m != nil {
		return m.Command
	}
	return SessionFrame_SessionNew
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
	proto.RegisterType((*SessionFrame)(nil), "v2ray.core.app.proxyman.SessionFrame")
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.SessionFrame_FrameCommand", SessionFrame_FrameCommand_name, SessionFrame_FrameCommand_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 899 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0x5d, 0x6f, 0x1b, 0x45,
	0x14, 0xed, 0xda, 0xa9, 0x9b, 0xdc, 0x24, 0x8e, 0x3b, 0x94, 0xd6, 0x18, 0x10, 0xc6, 0x42, 0x10,
	0x51, 0xb4, 0x2e, 0xae, 0x10, 0x42, 0x42, 0x2a, 0xa9, 0x13, 0x44, 0x04, 0x21, 0x66, 0x5c, 0xf1,
	0x50, 0x21, 0x59, 0x93, 0xdd, 0xe9, 0x32, 0x62, 0x77, 0x66, 0x34, 0x33, 0x4e, 0xb2, 0x6f, 0xfc,
	0x1e, 0x7e, 0x05, 0x8f, 0x3c, 0xf0, 0x8f, 0x78, 0x41, 0xf3, 0xb1, 0x8e, 0x53, 0x67, 0x5b, 0x42,
	0xd4, 0x97, 0x68, 0xc6, 0x39, 0xe7, 0xcc, 0xdc, 0x73, 0xee, 0x1d, 0x1b, 0x76, 0x4f, 0x47, 0x8a,
	0x94, 0x71, 0x22, 0x8a, 0x61, 0x22, 0x14, 0x1d, 0x12, 0x29, 0x87, 0x52, 0x89, 0xf3, 0xb2, 0x20,
	0x7c, 0x98, 0x08, 0xfe, 0x82, 0x65, 0xb1, 0x54, 0xc2, 0x08, 0xf4, 0xa0, 0x42, 0x2a, 0x1a, 0x13,
	0x29, 0xe3, 0x0a, 0xd5, 0x7b, 0xf4, 0x92, 0x44, 0x22, 0x8a, 0x42, 0xf0, 0xa1, 0xa6, 0x8a, 0x91,
	0x7c, 0x68, 0x4a, 0x49, 0xd3, 0x59, 0x41, 0xb5, 0x26, 0x19, 0xf5, 0x52, 0xbd, 0x4f, 0xae, 0x66,
	0x70, 0x6a, 0x86, 0x24, 0x4d, 0x15, 0xd5, 0x3a, 0x00, 0x1f, 0xd6, 0x03, 0x53, 0xaa, 0x0d, 0xe3,
	0xc4, 0x30, 0xc1, 0x03, 0xf8, 0xa3, 0x7a, 0xb0, 0x14, 0xca, 0x04, 0x54, 0xfc, 0x12, 0xca, 0x28,
	0xc2, 0xb5, 0xfd, 0xff, 0x90, 0x71, 0x43, 0x95, 0x45, 0x2f, 0x97, 0x3d, 0xd8, 0x81, 0xed, 0x43,
	0x7e, 0x22, 0xe6, 0x3c, 0x1d, 0xbb, 0x8f, 0x07, 0x7f, 0x36, 0x01, 0xed, 0xe5, 0xb9, 0x48, 0xdc,
	0xd9, 0x53, 0xa3, 0x88, 0xa1, 0x59, 0x89, 0xf6, 0x61, 0xcd, 0x96, 0xda, 0x8d, 0xfa, 0xd1, 0x6e,
	0x7b, 0xf4, 0x28, 0xae, 0x71, 0x2b, 0x5e, 0xa5, 0xc6, 0xcf, 0x4a, 0x49, 0xb1, 0x63, 0xa3, 0xdf,
	0x60, 0x33, 0x11, 0x3c, 0x99, 0x2b, 0x45, 0x79, 0x52, 0x76, 0x1b, 0xfd, 0x68, 0x77, 0x73, 0x74,
	0x78, 0x1d, 0xb1, 0xd5, 0x8f, 0xc6, 0x17, 0x82, 0x78, 0x59, 0x1d, 0xcd, 0xe0, 0x8e, 0xa2, 0x2f,
	0x14, 0xd5, 0xbf, 0x76, 0x9b, 0xee, 0xa0, 0x83, 0x9b, 0x1d, 0x84, 0xbd, 0x18, 0xae, 0x54, 0x7b,
	0x5f, 0xc0, 0xfb, 0xaf, 0xbc, 0x0e, 0xba, 0x07, 0xb7, 0x4f, 0x49, 0x3e, 0xf7, 0xae, 0x6d, 0x63,
	0xbf, 0xe9, 0x7d, 0x0e, 0xef, 0xd4, 0x8a, 0x5f, 0x4d, 0x19, 0x7c, 0x06, 0x6b, 0xd6, 0x45, 0x04,
	0xd0, 0xda, 0xcb, 0xcf, 0x48, 0xa9, 0x3b, 0xb7, 0xec, 0x1a, 0x13, 0x9e, 0x8a, 0xa2, 0x13, 0xa1,
	0x2d, 0x58, 0x3f, 0x38, 0xb7, 0xf1, 0x92, 0xbc, 0xd3, 0xb0, 0x11, 0xb6, 0x31, 0x4d, 0x28, 0x3b,
	0xa5, 0xca, 0xa7, 0x8a, 0x9e, 0x00, 0xd8, 0x26, 0x98, 0x29, 0xc2, 0x33, 0xaf, 0xbd, 0x39, 0xea,
	0x2f, 0xdb, 0xe1, 0xbb, 0x29, 0xe6, 0xd4, 0xc4, 0x13, 0xa1, 0x0c, 0xb6, 0x38, 0xbc, 0x21, 0xab,
	0x25, 0xfa, 0x0a, 0x5a, 0x39, 0xd3, 0x86, 0xf2, 0x10, 0xda, 0x87, 0x35, 0xe4, 0xc3, 0xc9, 0xb1,
	0xda, 0x17, 0x05, 0x61, 0x1c, 0x07, 0x02, 0xfa, 0x05, 0xde, 0x22, 0x8b, 0x7a, 0x67, 0x3a, 0x14,
	0x1c, 0x32, 0x79, 0x78, 0x8d, 0x4c, 0x30, 0x22, 0xab, 0x8d, 0xf9, 0x0c, 0x76, 0xb4, 0x51, 0x94,
	0x14, 0x33, 0x4d, 0x8d, 0x61, 0x3c, 0xd3, 0xdd, 0xb5, 0x55, 0xe5, 0xc5, 0x18, 0xc4, 0xd5, 0x18,
	0xc4, 0x53, 0xc7, 0xf2, 0xfe, 0xe0, 0xb6, 0xd7, 0x98, 0x06, 0x09, 0xf4, 0x0d, 0xbc, 0xa7, 0xbc,
	0x83, 0x33, 0xa1, 0x58, 0xc6, 0x38, 0xc9, 0x67, 0x4b, 0x23, 0xd9, 0xbd, 0xdd, 0x8f, 0x76, 0xd7,
	0x71, 0x2f, 0x60, 0x8e, 0x03, 0x64, 0xff, 0x02, 0x81, 0xbe, 0x86, 0xae, 0xbd, 0xed, 0xd9, 0x4c,
	0x12, 0xad, 0xad, 0x4e, 0x22, 0x38, 0xa7, 0x89, 0x63, 0xb7, 0x2c, 0xfb, 0x69, 0xa3, 0x1b, 0xe1,
	0xfb, 0x0e, 0x33, 0xf1, 0x90, 0xf1, 0x02, 0x31, 0xf8, 0x3b, 0x82, 0x7b, 0x61, 0x2e, 0xbf, 0x23,
	0x3c, 0xcd, 0x17, 0x41, 0x76, 0xa0, 0x69, 0x48, 0xe6, 0x12, 0xdc, 0xc0, 0x76, 0x89, 0xa6, 0x70,
	0x37, 0x5c, 0x43, 0x5d, 0x58, 0xe0, 0x43, 0xfa, 0xf8, 0x8a, 0x90, 0xfc, 0xbb, 0xe5, 0x86, 0x32,
	0x3d, 0xf2, 0xcf, 0x16, 0xee, 0x54, 0x02, 0x8b, 0xfa, 0x8f, 0xa0, 0xed, 0x82, 0xb8, 0x50, 0x6c,
	0x5e, 0x4b, 0x71, 0xdb, 0xb1, 0x2b, 0xb9, 0x41, 0x07, 0xda, 0xc7, 0x73, 0xb3, 0xfc, 0xcc, 0xfc,
	0xd5, 0x80, 0xad, 0x29, 0xe5, 0xe9, 0xa2, 0xb0, 0xc7, 0xd0, 0x3c, 0x65, 0x24, 0xb4, 0xe6, 0x7f,
	0xe8, 0x2e, 0x8b, 0xbe, 0x2a, 0xfc, 0xc6, 0xcd, 0xc3, 0xff, 0xa9, 0xa6, 0xf8, 0x4f, 0x5f, 0x23,
	0x3a, 0xb1, 0xa4, 0xa0, 0x79, 0xd9, 0x00, 0xf4, 0x1c, 0x50, 0x31, 0xcf, 0x0d, 0x93, 0x39, 0x3d,
	0x7f, 0x65, 0xa3, 0x5e, 0x1a, 0x81, 0xa3, 0x8a, 0xc2, 0x78, 0x16, 0x74, 0xef, 0x2e, 0x64, 0x16,
	0xe6, 0xfe, 0x13, 0xc1, 0xdb, 0x95, 0xbb, 0xaf, 0x6b, 0x96, 0x63, 0xd8, 0xd1, 0xce, 0xf5, 0xff,
	0xdb, 0x2a, 0x6d, 0x4f, 0x7f, 0x43, 0x8d, 0x82, 0xee, 0x43, 0x8b, 0x9e, 0x4b, 0xa6, 0xa8, 0xf3,
	0xa6, 0x89, 0xc3, 0x0e, 0x75, 0xe1, 0x8e, 0x15, 0xa1, 0xdc, 0xb8, 0xd1, 0xdb, 0xc0, 0xd5, 0x76,
	0x30, 0x01, 0xb4, 0x6a, 0x93, 0xc5, 0x53, 0x4e, 0x4e, 0x72, 0x9a, 0xba, 0xea, 0xd7, 0x71, 0xb5,
	0x45, 0xfd, 0xd5, 0xaf, 0xa0, 0xed, 0x4b, 0xdf, 0x1b, 0x83, 0xdf, 0x5d, 0x6b, 0x6a, 0xcd, 0x04,
	0xff, 0x56, 0x91, 0x82, 0xa2, 0x36, 0x34, 0x58, 0x1a, 0x1e, 0xe4, 0x06, 0x4b, 0xd1, 0x0f, 0xfe,
	0x32, 0x84, 0xa7, 0x8e, 0xde, 0x1e, 0x8d, 0x6a, 0x13, 0x5c, 0xd6, 0x89, 0xdd, 0xdf, 0xb1, 0x67,
	0xe2, 0x4a, 0x02, 0x7d, 0x09, 0x2d, 0x43, 0x54, 0x46, 0x4d, 0x70, 0xee, 0x83, 0x9a, 0xde, 0x3f,
	0xe0, 0xa9, 0x14, 0x8c, 0x1b, 0x1c, 0xe0, 0xb6, 0x46, 0x49, 0xca, 0x5c, 0x90, 0xd4, 0x99, 0xb5,
	0x85, 0xab, 0xed, 0xe0, 0x09, 0x6c, 0x2d, 0x9f, 0x85, 0xda, 0x00, 0xe1, 0x22, 0x3f, 0xd2, 0xb3,
	0xce, 0x2d, 0xb4, 0x03, 0x9b, 0x61, 0xff, 0x3d, 0xa5, 0xb2, 0x13, 0x2d, 0x01, 0x0e, 0x78, 0xda,
	0x69, 0x3c, 0x1d, 0xc3, 0xbb, 0x89, 0x28, 0xea, 0xaa, 0x9a, 0x44, 0xcf, 0xd7, 0xab, 0xf5, 0x1f,
	0x8d, 0x07, 0x3f, 0x8f, 0x30, 0x29, 0xe3, 0xb1, 0x45, 0xed, 0x49, 0xe9, 0xa7, 0xa0, 0x20, 0xfc,
	0xa4, 0xe5, 0x7e, 0x61, 0x3c, 0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0x70, 0x61, 0x1e, 0xe3, 0x84,
	0x09, 0x00, 0x00,
}
