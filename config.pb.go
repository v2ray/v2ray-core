package core

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import serial "v2ray.com/core/common/serial"
import transport "v2ray.com/core/transport"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Config is the master config of V2Ray. V2Ray takes this config as input and functions accordingly.
type Config struct {
	// Inbound handler configurations. Must have at least one item.
	Inbound []*InboundHandlerConfig `protobuf:"bytes,1,rep,name=inbound" json:"inbound,omitempty"`
	// Outbound handler configurations. Must have at least one item. The first item is used as default for routing.
	Outbound []*OutboundHandlerConfig `protobuf:"bytes,2,rep,name=outbound" json:"outbound,omitempty"`
	// App is for configurations of all features in V2Ray. A feature must implement the Feature interface, and its config type must be registered through common.RegisterConfig.
	App []*serial.TypedMessage `protobuf:"bytes,4,rep,name=app" json:"app,omitempty"`
	// Transport settings.
	Transport *transport.Config `protobuf:"bytes,5,opt,name=transport" json:"transport,omitempty"`
	// Configuration for extensions. The config may not work if corresponding extension is not loaded into V2Ray.
	// V2Ray will ignore such config during initialization.
	Extension            []*serial.TypedMessage `protobuf:"bytes,6,rep,name=extension" json:"extension,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_aee233fb2b85ca90, []int{0}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (dst *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(dst, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetInbound() []*InboundHandlerConfig {
	if m != nil {
		return m.Inbound
	}
	return nil
}

func (m *Config) GetOutbound() []*OutboundHandlerConfig {
	if m != nil {
		return m.Outbound
	}
	return nil
}

func (m *Config) GetApp() []*serial.TypedMessage {
	if m != nil {
		return m.App
	}
	return nil
}

func (m *Config) GetTransport() *transport.Config {
	if m != nil {
		return m.Transport
	}
	return nil
}

func (m *Config) GetExtension() []*serial.TypedMessage {
	if m != nil {
		return m.Extension
	}
	return nil
}

// InboundHandlerConfig is the configuration for inbound handler.
type InboundHandlerConfig struct {
	// Tag of the inbound handler. The tag must be unique among all inbound handlers
	Tag string `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	// Settings for how this inbound proxy is handled.
	ReceiverSettings *serial.TypedMessage `protobuf:"bytes,2,opt,name=receiver_settings,json=receiverSettings" json:"receiver_settings,omitempty"`
	// Settings for inbound proxy. Must be one of the inbound proxies.
	ProxySettings        *serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InboundHandlerConfig) Reset()         { *m = InboundHandlerConfig{} }
func (m *InboundHandlerConfig) String() string { return proto.CompactTextString(m) }
func (*InboundHandlerConfig) ProtoMessage()    {}
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_aee233fb2b85ca90, []int{1}
}
func (m *InboundHandlerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InboundHandlerConfig.Unmarshal(m, b)
}
func (m *InboundHandlerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InboundHandlerConfig.Marshal(b, m, deterministic)
}
func (dst *InboundHandlerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InboundHandlerConfig.Merge(dst, src)
}
func (m *InboundHandlerConfig) XXX_Size() int {
	return xxx_messageInfo_InboundHandlerConfig.Size(m)
}
func (m *InboundHandlerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_InboundHandlerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_InboundHandlerConfig proto.InternalMessageInfo

func (m *InboundHandlerConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *InboundHandlerConfig) GetReceiverSettings() *serial.TypedMessage {
	if m != nil {
		return m.ReceiverSettings
	}
	return nil
}

func (m *InboundHandlerConfig) GetProxySettings() *serial.TypedMessage {
	if m != nil {
		return m.ProxySettings
	}
	return nil
}

// OutboundHandlerConfig is the configuration for outbound handler.
type OutboundHandlerConfig struct {
	// Tag of this outbound handler.
	Tag string `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	// Settings for how to dial connection for this outbound handler.
	SenderSettings *serial.TypedMessage `protobuf:"bytes,2,opt,name=sender_settings,json=senderSettings" json:"sender_settings,omitempty"`
	// Settings for this outbound proxy. Must be one of the outbound proxies.
	ProxySettings *serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	// If not zero, this outbound will be expired in seconds. Not used for now.
	Expire int64 `protobuf:"varint,4,opt,name=expire" json:"expire,omitempty"`
	// Comment of this outbound handler. Not used for now.
	Comment              string   `protobuf:"bytes,5,opt,name=comment" json:"comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OutboundHandlerConfig) Reset()         { *m = OutboundHandlerConfig{} }
func (m *OutboundHandlerConfig) String() string { return proto.CompactTextString(m) }
func (*OutboundHandlerConfig) ProtoMessage()    {}
func (*OutboundHandlerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_aee233fb2b85ca90, []int{2}
}
func (m *OutboundHandlerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutboundHandlerConfig.Unmarshal(m, b)
}
func (m *OutboundHandlerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutboundHandlerConfig.Marshal(b, m, deterministic)
}
func (dst *OutboundHandlerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutboundHandlerConfig.Merge(dst, src)
}
func (m *OutboundHandlerConfig) XXX_Size() int {
	return xxx_messageInfo_OutboundHandlerConfig.Size(m)
}
func (m *OutboundHandlerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_OutboundHandlerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_OutboundHandlerConfig proto.InternalMessageInfo

func (m *OutboundHandlerConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *OutboundHandlerConfig) GetSenderSettings() *serial.TypedMessage {
	if m != nil {
		return m.SenderSettings
	}
	return nil
}

func (m *OutboundHandlerConfig) GetProxySettings() *serial.TypedMessage {
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
	proto.RegisterType((*Config)(nil), "v2ray.core.Config")
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.InboundHandlerConfig")
	proto.RegisterType((*OutboundHandlerConfig)(nil), "v2ray.core.OutboundHandlerConfig")
}

func init() { proto.RegisterFile("v2ray.com/core/config.proto", fileDescriptor_config_aee233fb2b85ca90) }

var fileDescriptor_config_aee233fb2b85ca90 = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x41, 0x6f, 0x94, 0x40,
	0x1c, 0xc5, 0x33, 0x0b, 0x6e, 0xcb, 0xbf, 0x71, 0x5d, 0x27, 0x6a, 0x26, 0xd5, 0x03, 0x6e, 0x52,
	0xc3, 0x69, 0x30, 0x78, 0x31, 0x4d, 0xbc, 0x58, 0x0f, 0x6a, 0xd2, 0xd4, 0x50, 0xe3, 0xc1, 0x4b,
	0x33, 0x85, 0xbf, 0x84, 0xa4, 0xcc, 0x90, 0x99, 0x69, 0xb3, 0x7c, 0x25, 0xbf, 0x87, 0x37, 0xbf,
	0x91, 0x17, 0x03, 0x03, 0x0b, 0xd5, 0x3d, 0xb8, 0x26, 0x3d, 0xc1, 0xf0, 0xe6, 0xf7, 0xde, 0xff,
	0x01, 0x03, 0x4f, 0x6f, 0x12, 0x2d, 0x1a, 0x9e, 0xa9, 0x2a, 0xce, 0x94, 0xc6, 0x38, 0x53, 0xf2,
	0x5b, 0x59, 0xf0, 0x5a, 0x2b, 0xab, 0x28, 0x0c, 0xa2, 0xc6, 0xc3, 0x97, 0x7f, 0x6d, 0xac, 0x2a,
	0x25, 0x63, 0x83, 0xba, 0x14, 0x57, 0xb1, 0x6d, 0x6a, 0xcc, 0x2f, 0x2a, 0x34, 0x46, 0x14, 0xe8,
	0xe8, 0xc3, 0xa3, 0x3f, 0x08, 0xab, 0x85, 0x34, 0xb5, 0xd2, 0xf6, 0x56, 0xc8, 0xea, 0xc7, 0x0c,
	0xe6, 0x27, 0xdd, 0x03, 0x7a, 0x0c, 0x7b, 0xa5, 0xbc, 0x54, 0xd7, 0x32, 0x67, 0x24, 0xf4, 0xa2,
	0x83, 0x24, 0xe4, 0xe3, 0x04, 0xfc, 0x83, 0x93, 0xde, 0x0b, 0x99, 0x5f, 0xa1, 0x76, 0x48, 0x3a,
	0x00, 0xf4, 0x0d, 0xec, 0xab, 0x6b, 0xeb, 0xe0, 0x59, 0x07, 0x3f, 0x9f, 0xc2, 0x67, 0xbd, 0x76,
	0x9b, 0xde, 0x20, 0xf4, 0x35, 0x78, 0xa2, 0xae, 0x99, 0xdf, 0x91, 0x2f, 0xa6, 0xa4, 0x2b, 0xca,
	0x5d, 0x51, 0xfe, 0xb9, 0x2d, 0x7a, 0xea, 0x7a, 0xa6, 0x2d, 0x42, 0x8f, 0x21, 0xd8, 0x34, 0x63,
	0xf7, 0x42, 0x12, 0x1d, 0x24, 0xcf, 0xa6, 0xfc, 0x46, 0xe4, 0x7d, 0xe8, 0xb8, 0x9d, 0xbe, 0x83,
	0x00, 0xd7, 0x16, 0xa5, 0x29, 0x95, 0x64, 0xf3, 0x9d, 0xb2, 0x47, 0xf0, 0xa3, 0xbf, 0xef, 0x2d,
	0xfd, 0xd5, 0x4f, 0x02, 0x8f, 0xb6, 0xbd, 0x22, 0xba, 0x04, 0xcf, 0x8a, 0x82, 0x91, 0x90, 0x44,
	0x41, 0xda, 0xde, 0xd2, 0x73, 0x78, 0xa8, 0x31, 0xc3, 0xf2, 0x06, 0xf5, 0x85, 0x41, 0x6b, 0x4b,
	0x59, 0x18, 0x36, 0xeb, 0x46, 0xff, 0xd7, 0xf8, 0xe5, 0x60, 0x70, 0xde, 0xf3, 0xf4, 0x14, 0x16,
	0xb5, 0x56, 0xeb, 0x66, 0x74, 0xf4, 0x76, 0x72, 0xbc, 0xdf, 0xd1, 0x83, 0xdd, 0xea, 0x17, 0x81,
	0xc7, 0x5b, 0x3f, 0xda, 0x96, 0x3e, 0x67, 0xf0, 0xc0, 0xa0, 0xcc, 0xff, 0xbf, 0xcd, 0xc2, 0xe1,
	0x77, 0xd4, 0x85, 0x3e, 0x81, 0x39, 0xae, 0xeb, 0x52, 0x23, 0xf3, 0x43, 0x12, 0x79, 0x69, 0xbf,
	0xa2, 0x0c, 0xf6, 0x5a, 0x13, 0x94, 0xee, 0xc7, 0x09, 0xd2, 0x61, 0xf9, 0xf6, 0x08, 0x16, 0x99,
	0xaa, 0x26, 0x69, 0x9f, 0xc8, 0x57, 0xbf, 0xbd, 0x7e, 0x9f, 0xc1, 0x97, 0x24, 0x15, 0x0d, 0x3f,
	0x51, 0x1a, 0x2f, 0xe7, 0xdd, 0x11, 0x7a, 0xf5, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xf0, 0x15, 0xee,
	0x31, 0xc6, 0x03, 0x00, 0x00,
}
