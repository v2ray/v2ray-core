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
func (SessionFrame_FrameCommand) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{9, 0} }

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
	DispatchSettings *DispatchConfig                        `protobuf:"bytes,4,opt,name=dispatch_settings,json=dispatchSettings" json:"dispatch_settings,omitempty"`
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

func (m *InboundHandlerConfig) GetDispatchSettings() *DispatchConfig {
	if m != nil {
		return m.DispatchSettings
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
	Tag              string                                 `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	SenderSettings   *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=sender_settings,json=senderSettings" json:"sender_settings,omitempty"`
	ProxySettings    *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	Expire           int64                                  `protobuf:"varint,4,opt,name=expire" json:"expire,omitempty"`
	Comment          string                                 `protobuf:"bytes,5,opt,name=comment" json:"comment,omitempty"`
	DispatchSettings *DispatchConfig                        `protobuf:"bytes,6,opt,name=dispatch_settings,json=dispatchSettings" json:"dispatch_settings,omitempty"`
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

func (m *OutboundHandlerConfig) GetDispatchSettings() *DispatchConfig {
	if m != nil {
		return m.DispatchSettings
	}
	return nil
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
	Command SessionFrame_FrameCommand        `protobuf:"varint,2,opt,name=command,enum=v2ray.core.app.proxyman.SessionFrame_FrameCommand" json:"command,omitempty"`
	Target  *v2ray_core_common_net2.Endpoint `protobuf:"bytes,3,opt,name=target" json:"target,omitempty"`
	Payload []byte                           `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
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
	proto.RegisterType((*DispatchConfig)(nil), "v2ray.core.app.proxyman.DispatchConfig")
	proto.RegisterType((*SessionFrame)(nil), "v2ray.core.app.proxyman.SessionFrame")
	proto.RegisterEnum("v2ray.core.app.proxyman.AllocationStrategy_Type", AllocationStrategy_Type_name, AllocationStrategy_Type_value)
	proto.RegisterEnum("v2ray.core.app.proxyman.SessionFrame_FrameCommand", SessionFrame_FrameCommand_name, SessionFrame_FrameCommand_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/proxyman/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 929 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x56, 0xd1, 0x6e, 0x23, 0x35,
	0x14, 0xdd, 0x4c, 0xba, 0x69, 0x7b, 0x93, 0x4e, 0x67, 0xcd, 0xb2, 0x1b, 0x02, 0x88, 0x12, 0x21,
	0xb6, 0x62, 0xd1, 0x64, 0xc9, 0x0a, 0x01, 0x4f, 0x4b, 0x49, 0x8b, 0xa8, 0xa0, 0xdb, 0xe0, 0x54,
	0x3c, 0x20, 0xa4, 0xe0, 0xce, 0x78, 0xb3, 0x16, 0x33, 0xb6, 0x65, 0x3b, 0x6d, 0xe6, 0x8d, 0x6f,
	0xe1, 0x91, 0x6f, 0xe0, 0x81, 0x5f, 0xe0, 0x3f, 0xf8, 0x08, 0x34, 0x1e, 0x4f, 0x92, 0x6d, 0x32,
	0x94, 0x52, 0xed, 0x4b, 0x65, 0x8f, 0xcf, 0x3d, 0xb6, 0xcf, 0x39, 0xd7, 0x0d, 0xec, 0x5f, 0xf4,
	0x15, 0xc9, 0xc2, 0x48, 0xa4, 0xbd, 0x48, 0x28, 0xda, 0x23, 0x52, 0xf6, 0xa4, 0x12, 0xb3, 0x2c,
	0x25, 0xbc, 0x17, 0x09, 0xfe, 0x82, 0x4d, 0x42, 0xa9, 0x84, 0x11, 0xe8, 0x61, 0x89, 0x54, 0x34,
	0x24, 0x52, 0x86, 0x25, 0xaa, 0xf3, 0xe4, 0x0a, 0x45, 0x24, 0xd2, 0x54, 0xf0, 0x9e, 0xa6, 0x8a,
	0x91, 0xa4, 0x67, 0x32, 0x49, 0xe3, 0x71, 0x4a, 0xb5, 0x26, 0x13, 0x5a, 0x50, 0x75, 0x1e, 0xad,
	0xaf, 0xe0, 0xd4, 0xf4, 0x48, 0x1c, 0x2b, 0xaa, 0xb5, 0x03, 0x3e, 0xae, 0x06, 0xc6, 0x54, 0x1b,
	0xc6, 0x89, 0x61, 0x82, 0x3b, 0xf0, 0x07, 0xd5, 0x60, 0x29, 0x94, 0x71, 0xa8, 0xf0, 0x0a, 0xca,
	0x28, 0xc2, 0x75, 0xbe, 0xde, 0x63, 0xdc, 0x50, 0x95, 0xa3, 0x97, 0xaf, 0xdd, 0xdd, 0x85, 0x9d,
	0x63, 0x7e, 0x2e, 0xa6, 0x3c, 0x1e, 0xd8, 0xcf, 0xdd, 0x3f, 0xeb, 0x80, 0x0e, 0x92, 0x44, 0x44,
	0x76, 0xef, 0x91, 0x51, 0xc4, 0xd0, 0x49, 0x86, 0x0e, 0x61, 0x23, 0xbf, 0x6a, 0xbb, 0xb6, 0x57,
	0xdb, 0xf7, 0xfb, 0x4f, 0xc2, 0x0a, 0xb5, 0xc2, 0xd5, 0xd2, 0xf0, 0x2c, 0x93, 0x14, 0xdb, 0x6a,
	0xf4, 0x0b, 0x34, 0x23, 0xc1, 0xa3, 0xa9, 0x52, 0x94, 0x47, 0x59, 0xdb, 0xdb, 0xab, 0xed, 0x37,
	0xfb, 0xc7, 0x37, 0x21, 0x5b, 0xfd, 0x34, 0x58, 0x10, 0xe2, 0x65, 0x76, 0x34, 0x86, 0x4d, 0x45,
	0x5f, 0x28, 0xaa, 0x5f, 0xb6, 0xeb, 0x76, 0xa3, 0xa3, 0xdb, 0x6d, 0x84, 0x0b, 0x32, 0x5c, 0xb2,
	0x76, 0x3e, 0x85, 0x77, 0xff, 0xf5, 0x38, 0xe8, 0x3e, 0xdc, 0xbd, 0x20, 0xc9, 0xb4, 0x50, 0x6d,
	0x07, 0x17, 0x93, 0xce, 0x27, 0xf0, 0x56, 0x25, 0xf9, 0xfa, 0x92, 0xee, 0xc7, 0xb0, 0x91, 0xab,
	0x88, 0x00, 0x1a, 0x07, 0xc9, 0x25, 0xc9, 0x74, 0x70, 0x27, 0x1f, 0x63, 0xc2, 0x63, 0x91, 0x06,
	0x35, 0xd4, 0x82, 0xad, 0xa3, 0x59, 0x6e, 0x2f, 0x49, 0x02, 0xaf, 0xfb, 0x47, 0x1d, 0x7c, 0x4c,
	0x23, 0xca, 0x2e, 0xa8, 0x2a, 0x5c, 0x45, 0xcf, 0x00, 0xf2, 0x10, 0x8c, 0x15, 0xe1, 0x93, 0x82,
	0xbb, 0xd9, 0xdf, 0x5b, 0x96, 0xa3, 0x48, 0x53, 0xc8, 0xa9, 0x09, 0x87, 0x42, 0x19, 0x9c, 0xe3,
	0xf0, 0xb6, 0x2c, 0x87, 0xe8, 0x0b, 0x68, 0x24, 0x4c, 0x1b, 0xca, 0x9d, 0x69, 0xef, 0x57, 0x14,
	0x1f, 0x0f, 0x4f, 0xd5, 0xa1, 0x48, 0x09, 0xe3, 0xd8, 0x15, 0xa0, 0x9f, 0xe0, 0x0d, 0x32, 0xbf,
	0xef, 0x58, 0xbb, 0x0b, 0x3b, 0x4f, 0x1e, 0xdf, 0xc0, 0x13, 0x8c, 0xc8, 0x6a, 0x30, 0xcf, 0x60,
	0x57, 0x1b, 0x45, 0x49, 0x3a, 0xd6, 0xd4, 0x18, 0xc6, 0x27, 0xba, 0xbd, 0xb1, 0xca, 0x3c, 0x6f,
	0x83, 0xb0, 0x6c, 0x83, 0x70, 0x64, 0xab, 0x0a, 0x7d, 0xb0, 0x5f, 0x70, 0x8c, 0x1c, 0x05, 0xfa,
	0x12, 0xde, 0x51, 0x85, 0x82, 0x63, 0xa1, 0xd8, 0x84, 0x71, 0x92, 0x8c, 0x97, 0x5a, 0xb2, 0x7d,
	0x77, 0xaf, 0xb6, 0xbf, 0x85, 0x3b, 0x0e, 0x73, 0xea, 0x20, 0x87, 0x0b, 0x04, 0xfa, 0x1c, 0xda,
	0xf9, 0x69, 0x2f, 0xc7, 0x92, 0x68, 0x9d, 0xf3, 0x44, 0x82, 0x73, 0x1a, 0xd9, 0xea, 0x86, 0xad,
	0x7e, 0x60, 0xd7, 0x87, 0xc5, 0xf2, 0x60, 0xbe, 0xda, 0xfd, 0xcd, 0x83, 0xfb, 0xae, 0x27, 0xbf,
	0x21, 0x3c, 0x4e, 0xe6, 0x26, 0x06, 0x50, 0x37, 0x64, 0x62, 0xdd, 0xdb, 0xc6, 0xf9, 0x10, 0x8d,
	0xe0, 0x9e, 0x3b, 0x82, 0x5a, 0x5c, 0xbf, 0x30, 0xe8, 0xc3, 0x35, 0x06, 0x15, 0x6f, 0x96, 0x6d,
	0xc8, 0xf8, 0xa4, 0x78, 0xb2, 0x70, 0x50, 0x12, 0xcc, 0xef, 0x7e, 0x02, 0xbe, 0x35, 0x61, 0xc1,
	0x58, 0xbf, 0x11, 0xe3, 0x8e, 0xad, 0x9e, 0xd3, 0x9d, 0xc1, 0xbd, 0x98, 0x69, 0x49, 0x4c, 0xf4,
	0xf2, 0xaa, 0x45, 0x8f, 0x2a, 0xcd, 0x3f, 0x74, 0x15, 0xce, 0x9e, 0xa0, 0x64, 0x28, 0x59, 0xbb,
	0x01, 0xf8, 0xa7, 0x53, 0xb3, 0xfc, 0x70, 0xfd, 0x5d, 0x83, 0xd6, 0x88, 0xf2, 0x78, 0x2e, 0xd7,
	0x53, 0xa8, 0x5f, 0x30, 0xe2, 0xc2, 0xfe, 0x1f, 0xf2, 0x9a, 0xa3, 0xd7, 0xc5, 0xc9, 0xbb, 0x7d,
	0x9c, 0xbe, 0xaf, 0x90, 0xf4, 0xa3, 0x6b, 0x48, 0x87, 0x79, 0x91, 0xe3, 0x7c, 0x55, 0xd6, 0xee,
	0x5f, 0x1e, 0xbc, 0x59, 0x2a, 0x70, 0x5d, 0x4c, 0x4e, 0x61, 0x57, 0x5b, 0x65, 0xfe, 0x6f, 0x48,
	0xfc, 0xa2, 0xfc, 0x75, 0x45, 0xe4, 0x01, 0x34, 0xe8, 0x4c, 0x32, 0x45, 0x6d, 0x2e, 0xea, 0xd8,
	0xcd, 0x50, 0x1b, 0x36, 0x73, 0x12, 0xca, 0x8d, 0x6d, 0xb8, 0x6d, 0x5c, 0x4e, 0xd7, 0x87, 0xaa,
	0x71, 0xdb, 0x50, 0x85, 0x80, 0x4e, 0xa6, 0x89, 0x61, 0x32, 0xa1, 0x33, 0xc6, 0x27, 0x4e, 0xcf,
	0x36, 0x6c, 0x52, 0x4e, 0xce, 0x13, 0x1a, 0x5b, 0x4d, 0xb7, 0x70, 0x39, 0xed, 0xfe, 0x0c, 0xfe,
	0xab, 0x9c, 0xe8, 0x39, 0xb4, 0xd2, 0xe9, 0x6c, 0x71, 0xa4, 0xda, 0x35, 0x8f, 0xdc, 0xea, 0x76,
	0xb8, 0x99, 0x4e, 0x67, 0xf3, 0x13, 0xfd, 0xea, 0xe5, 0xa1, 0xd6, 0x9a, 0x09, 0xfe, 0xb5, 0x22,
	0x29, 0x45, 0x3e, 0x78, 0x2c, 0x76, 0xff, 0x1c, 0x3c, 0x16, 0xa3, 0xef, 0x0a, 0x89, 0x08, 0x8f,
	0xad, 0xa5, 0x7e, 0xbf, 0x5f, 0xb9, 0xd7, 0x32, 0x4f, 0x68, 0xff, 0x0e, 0x8a, 0x4a, 0x5c, 0x52,
	0xa0, 0xcf, 0xa0, 0x61, 0x88, 0x9a, 0x50, 0xe3, 0xfc, 0x7c, 0xaf, 0xa2, 0x6b, 0x8e, 0x78, 0x2c,
	0x05, 0xe3, 0x06, 0x3b, 0x78, 0xae, 0x91, 0x24, 0x59, 0x22, 0x48, 0x6c, 0x2d, 0x6c, 0xe1, 0x72,
	0xda, 0x7d, 0x06, 0xad, 0xe5, 0xbd, 0x90, 0x0f, 0xe0, 0x0e, 0xf2, 0x9c, 0x5e, 0x06, 0x77, 0xd0,
	0x2e, 0x34, 0xdd, 0xfc, 0x5b, 0x4a, 0x65, 0x50, 0x5b, 0x02, 0x1c, 0xf1, 0x38, 0xf0, 0xbe, 0x1a,
	0xc0, 0xdb, 0x91, 0x48, 0xab, 0x6e, 0x35, 0xac, 0xfd, 0xb8, 0x55, 0x8e, 0x7f, 0xf7, 0x1e, 0xfe,
	0xd0, 0xc7, 0x24, 0x0b, 0x07, 0x39, 0xea, 0x40, 0xca, 0xa2, 0x7f, 0x52, 0xc2, 0xcf, 0x1b, 0xf6,
	0xd7, 0xce, 0xd3, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x78, 0xee, 0xf7, 0xe3, 0x10, 0x0a, 0x00,
	0x00,
}
