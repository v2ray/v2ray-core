package core

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	serial "v2ray.com/core/common/serial"
	transport "v2ray.com/core/transport"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Config is the master config of V2Ray. V2Ray takes this config as input and functions accordingly.
type Config struct {
	// Inbound handler configurations. Must have at least one item.
	Inbound []*InboundHandlerConfig `protobuf:"bytes,1,rep,name=inbound,proto3" json:"inbound,omitempty"`
	// Outbound handler configurations. Must have at least one item. The first item is used as default for routing.
	Outbound []*OutboundHandlerConfig `protobuf:"bytes,2,rep,name=outbound,proto3" json:"outbound,omitempty"`
	// App is for configurations of all features in V2Ray. A feature must implement the Feature interface, and its config type must be registered through common.RegisterConfig.
	App []*serial.TypedMessage `protobuf:"bytes,4,rep,name=app,proto3" json:"app,omitempty"`
	// Transport settings.
	// Deprecated. Each inbound and outbound should choose their own transport config.
	// Date to remove: 2020-01-13
	Transport *transport.Config `protobuf:"bytes,5,opt,name=transport,proto3" json:"transport,omitempty"` // Deprecated: Do not use.
	// Configuration for extensions. The config may not work if corresponding extension is not loaded into V2Ray.
	// V2Ray will ignore such config during initialization.
	Extension            []*serial.TypedMessage `protobuf:"bytes,6,rep,name=extension,proto3" json:"extension,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_13704400b1045c6b, []int{0}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
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

// Deprecated: Do not use.
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
	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	// Settings for how this inbound proxy is handled.
	ReceiverSettings *serial.TypedMessage `protobuf:"bytes,2,opt,name=receiver_settings,json=receiverSettings,proto3" json:"receiver_settings,omitempty"`
	// Settings for inbound proxy. Must be one of the inbound proxies.
	ProxySettings        *serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings,proto3" json:"proxy_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InboundHandlerConfig) Reset()         { *m = InboundHandlerConfig{} }
func (m *InboundHandlerConfig) String() string { return proto.CompactTextString(m) }
func (*InboundHandlerConfig) ProtoMessage()    {}
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_13704400b1045c6b, []int{1}
}

func (m *InboundHandlerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InboundHandlerConfig.Unmarshal(m, b)
}
func (m *InboundHandlerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InboundHandlerConfig.Marshal(b, m, deterministic)
}
func (m *InboundHandlerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InboundHandlerConfig.Merge(m, src)
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
	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	// Settings for how to dial connection for this outbound handler.
	SenderSettings *serial.TypedMessage `protobuf:"bytes,2,opt,name=sender_settings,json=senderSettings,proto3" json:"sender_settings,omitempty"`
	// Settings for this outbound proxy. Must be one of the outbound proxies.
	ProxySettings *serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings,proto3" json:"proxy_settings,omitempty"`
	// If not zero, this outbound will be expired in seconds. Not used for now.
	Expire int64 `protobuf:"varint,4,opt,name=expire,proto3" json:"expire,omitempty"`
	// Comment of this outbound handler. Not used for now.
	Comment              string   `protobuf:"bytes,5,opt,name=comment,proto3" json:"comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OutboundHandlerConfig) Reset()         { *m = OutboundHandlerConfig{} }
func (m *OutboundHandlerConfig) String() string { return proto.CompactTextString(m) }
func (*OutboundHandlerConfig) ProtoMessage()    {}
func (*OutboundHandlerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_13704400b1045c6b, []int{2}
}

func (m *OutboundHandlerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OutboundHandlerConfig.Unmarshal(m, b)
}
func (m *OutboundHandlerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OutboundHandlerConfig.Marshal(b, m, deterministic)
}
func (m *OutboundHandlerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutboundHandlerConfig.Merge(m, src)
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

func init() { proto.RegisterFile("v2ray.com/core/config.proto", fileDescriptor_13704400b1045c6b) }

var fileDescriptor_13704400b1045c6b = []byte{
	// 409 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0xd3, 0xcf, 0x6b, 0xd4, 0x40,
	0x14, 0x07, 0x70, 0x26, 0x89, 0xdb, 0xe6, 0x15, 0xd7, 0x75, 0x50, 0x19, 0xaa, 0x87, 0xb8, 0x50,
	0xc9, 0x69, 0x22, 0xf1, 0x22, 0x82, 0x1e, 0x5a, 0x0f, 0x2a, 0x94, 0x4a, 0x2a, 0x1e, 0xbc, 0x94,
	0x69, 0xf2, 0x0c, 0x81, 0x66, 0x26, 0xcc, 0x4c, 0xcb, 0xe6, 0x5f, 0xf2, 0x3f, 0x11, 0xfc, 0x8f,
	0xbc, 0x48, 0x32, 0x9b, 0x4d, 0xaa, 0x7b, 0x70, 0x05, 0x4f, 0xc9, 0xe4, 0xe5, 0xf3, 0xde, 0x7c,
	0xf3, 0x03, 0x1e, 0xdf, 0xa4, 0x5a, 0xb4, 0x3c, 0x57, 0x75, 0x92, 0x2b, 0x8d, 0x49, 0xae, 0xe4,
	0xd7, 0xaa, 0xe4, 0x8d, 0x56, 0x56, 0x51, 0x18, 0x8a, 0x1a, 0x0f, 0x9f, 0xff, 0x71, 0x63, 0x5d,
	0x2b, 0x99, 0x18, 0xd4, 0x95, 0xb8, 0x4a, 0x6c, 0xdb, 0x60, 0x71, 0x51, 0xa3, 0x31, 0xa2, 0x44,
	0xa7, 0x0f, 0x8f, 0x7e, 0x13, 0x56, 0x0b, 0x69, 0x1a, 0xa5, 0xed, 0xad, 0x21, 0xcb, 0xef, 0x1e,
	0xcc, 0x4e, 0xfa, 0x0b, 0xf4, 0x15, 0xec, 0x55, 0xf2, 0x52, 0x5d, 0xcb, 0x82, 0x91, 0xc8, 0x8f,
	0x0f, 0xd2, 0x88, 0x8f, 0x3b, 0xe0, 0xef, 0x5d, 0xe9, 0x9d, 0x90, 0xc5, 0x15, 0x6a, 0x47, 0xb2,
	0x01, 0xd0, 0xd7, 0xb0, 0xaf, 0xae, 0xad, 0xc3, 0x5e, 0x8f, 0x9f, 0x4e, 0xf1, 0xd9, 0xba, 0x76,
	0x5b, 0x6f, 0x08, 0x7d, 0x09, 0xbe, 0x68, 0x1a, 0x16, 0xf4, 0xf2, 0xd9, 0x54, 0xba, 0xa0, 0xdc,
	0x05, 0xe5, 0x9f, 0xba, 0xa0, 0xa7, 0x2e, 0x67, 0xd6, 0x11, 0xfa, 0x06, 0xc2, 0x4d, 0x32, 0x76,
	0x27, 0x22, 0xf1, 0x41, 0xfa, 0x64, 0xea, 0x37, 0x45, 0xee, 0x86, 0x1e, 0x7b, 0x8c, 0x64, 0x23,
	0xa1, 0x6f, 0x21, 0xc4, 0x95, 0x45, 0x69, 0x2a, 0x25, 0xd9, 0x6c, 0xa7, 0xf9, 0x23, 0xfc, 0x10,
	0xec, 0xfb, 0x8b, 0x60, 0xf9, 0x83, 0xc0, 0x83, 0x6d, 0x8f, 0x89, 0x2e, 0xc0, 0xb7, 0xa2, 0x64,
	0x24, 0x22, 0x71, 0x98, 0x75, 0xa7, 0xf4, 0x1c, 0xee, 0x6b, 0xcc, 0xb1, 0xba, 0x41, 0x7d, 0x61,
	0xd0, 0xda, 0x4a, 0x96, 0x86, 0x79, 0xfd, 0xf6, 0xff, 0x76, 0xfc, 0x62, 0x68, 0x70, 0xbe, 0xf6,
	0xf4, 0x14, 0xe6, 0x8d, 0x56, 0xab, 0x76, 0xec, 0xe8, 0xef, 0xd4, 0xf1, 0x6e, 0xaf, 0x87, 0x76,
	0xcb, 0x9f, 0x04, 0x1e, 0x6e, 0x7d, 0x71, 0x5b, 0xf2, 0x9c, 0xc1, 0x3d, 0x83, 0xb2, 0xf8, 0xf7,
	0x34, 0x73, 0xc7, 0xff, 0x53, 0x16, 0xfa, 0x08, 0x66, 0xb8, 0x6a, 0x2a, 0x8d, 0x2c, 0x88, 0x48,
	0xec, 0x67, 0xeb, 0x15, 0x65, 0xb0, 0xd7, 0x35, 0x41, 0xe9, 0x3e, 0x9e, 0x30, 0x1b, 0x96, 0xc7,
	0x47, 0x30, 0xcf, 0x55, 0x3d, 0x99, 0xf6, 0x91, 0x7c, 0x09, 0xba, 0xe3, 0x37, 0x0f, 0x3e, 0xa7,
	0x99, 0x68, 0xf9, 0x89, 0xd2, 0x78, 0x39, 0xeb, 0x7f, 0xa3, 0x17, 0xbf, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xde, 0x43, 0xe7, 0x49, 0xca, 0x03, 0x00, 0x00,
}
