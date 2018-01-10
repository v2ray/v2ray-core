package core

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_serial "v2ray.com/core/common/serial"
import v2ray_core_transport "v2ray.com/core/transport"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Configuration serialization format.
type ConfigFormat int32

const (
	ConfigFormat_Protobuf ConfigFormat = 0
	ConfigFormat_JSON     ConfigFormat = 1
)

var ConfigFormat_name = map[int32]string{
	0: "Protobuf",
	1: "JSON",
}
var ConfigFormat_value = map[string]int32{
	"Protobuf": 0,
	"JSON":     1,
}

func (x ConfigFormat) String() string {
	return proto.EnumName(ConfigFormat_name, int32(x))
}
func (ConfigFormat) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Master config of V2Ray. V2Ray takes this config as input and functions accordingly.
type Config struct {
	// Inbound handler configurations. Must have at least one item.
	Inbound []*InboundHandlerConfig `protobuf:"bytes,1,rep,name=inbound" json:"inbound,omitempty"`
	// Outbound handler configurations. Must have at least one item. The first item is used as default for routing.
	Outbound []*OutboundHandlerConfig `protobuf:"bytes,2,rep,name=outbound" json:"outbound,omitempty"`
	// App configuration. Must be one in the app directory.
	App []*v2ray_core_common_serial.TypedMessage `protobuf:"bytes,4,rep,name=app" json:"app,omitempty"`
	// Transport settings.
	Transport *v2ray_core_transport.Config `protobuf:"bytes,5,opt,name=transport" json:"transport,omitempty"`
	// Configuration for extensions. The config may not work if corresponding extension is not loaded into V2Ray.
	// V2Ray will ignore such config during initialization.
	Extension []*v2ray_core_common_serial.TypedMessage `protobuf:"bytes,6,rep,name=extension" json:"extension,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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

func (m *Config) GetApp() []*v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.App
	}
	return nil
}

func (m *Config) GetTransport() *v2ray_core_transport.Config {
	if m != nil {
		return m.Transport
	}
	return nil
}

func (m *Config) GetExtension() []*v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.Extension
	}
	return nil
}

type InboundHandlerConfig struct {
	// Tag of the inbound handler.
	Tag string `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	// Settings for how this inbound proxy is handled. Must be ReceiverConfig above.
	ReceiverSettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=receiver_settings,json=receiverSettings" json:"receiver_settings,omitempty"`
	// Settings for inbound proxy. Must be one of the inbound proxies.
	ProxySettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
}

func (m *InboundHandlerConfig) Reset()                    { *m = InboundHandlerConfig{} }
func (m *InboundHandlerConfig) String() string            { return proto.CompactTextString(m) }
func (*InboundHandlerConfig) ProtoMessage()               {}
func (*InboundHandlerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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

type OutboundHandlerConfig struct {
	// Tag of this outbound handler.
	Tag string `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	// Settings for how to dial connection for this outbound handler. Must be SenderConfig above.
	SenderSettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=sender_settings,json=senderSettings" json:"sender_settings,omitempty"`
	// Settings for this outbound proxy. Must be one of the outbound proxies.
	ProxySettings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,3,opt,name=proxy_settings,json=proxySettings" json:"proxy_settings,omitempty"`
	// If not zero, this outbound will be expired in seconds. Not used for now.
	Expire int64 `protobuf:"varint,4,opt,name=expire" json:"expire,omitempty"`
	// Comment of this outbound handler. Not used for now.
	Comment string `protobuf:"bytes,5,opt,name=comment" json:"comment,omitempty"`
}

func (m *OutboundHandlerConfig) Reset()                    { *m = OutboundHandlerConfig{} }
func (m *OutboundHandlerConfig) String() string            { return proto.CompactTextString(m) }
func (*OutboundHandlerConfig) ProtoMessage()               {}
func (*OutboundHandlerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

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
	proto.RegisterType((*Config)(nil), "v2ray.core.Config")
	proto.RegisterType((*InboundHandlerConfig)(nil), "v2ray.core.InboundHandlerConfig")
	proto.RegisterType((*OutboundHandlerConfig)(nil), "v2ray.core.OutboundHandlerConfig")
	proto.RegisterEnum("v2ray.core.ConfigFormat", ConfigFormat_name, ConfigFormat_value)
}

func init() { proto.RegisterFile("v2ray.com/core/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 436 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x71, 0x13, 0xba, 0xf6, 0x6d, 0x94, 0x60, 0x01, 0xb2, 0x06, 0x87, 0x50, 0x69, 0x53,
	0xc5, 0xc1, 0x45, 0xe5, 0x82, 0x26, 0x71, 0x61, 0x08, 0xc1, 0xa4, 0xd1, 0x29, 0x45, 0x1c, 0xb8,
	0x4c, 0x6e, 0xfa, 0x56, 0x45, 0x5a, 0xec, 0xc8, 0x76, 0xa7, 0xe6, 0x2b, 0xf1, 0x3d, 0xb8, 0xf1,
	0x8d, 0xb8, 0xa0, 0xc4, 0x49, 0x93, 0x41, 0x0f, 0x14, 0x69, 0xa7, 0xc4, 0x79, 0xfe, 0xfd, 0xdf,
	0xfb, 0x25, 0x31, 0x3c, 0xbb, 0x99, 0x68, 0x91, 0xf3, 0x58, 0xa5, 0xe3, 0x58, 0x69, 0x1c, 0xc7,
	0x4a, 0x5e, 0x25, 0x4b, 0x9e, 0x69, 0x65, 0x15, 0x85, 0xba, 0xa8, 0xf1, 0xf0, 0xd5, 0x5f, 0x1b,
	0xd3, 0x54, 0xc9, 0xb1, 0x41, 0x9d, 0x88, 0xeb, 0xb1, 0xcd, 0x33, 0x5c, 0x5c, 0xa6, 0x68, 0x8c,
	0x58, 0xa2, 0xa3, 0x0f, 0x8f, 0xfe, 0x20, 0xac, 0x16, 0xd2, 0x64, 0x4a, 0xdb, 0x5b, 0x4d, 0x86,
	0x3f, 0x3a, 0xd0, 0x3d, 0x2d, 0x1f, 0xd0, 0x13, 0xd8, 0x4b, 0xe4, 0x5c, 0xad, 0xe4, 0x82, 0x91,
	0xd0, 0x1b, 0xed, 0x4f, 0x42, 0xde, 0x4c, 0xc0, 0x3f, 0xb9, 0xd2, 0x47, 0x21, 0x17, 0xd7, 0xa8,
	0x1d, 0x12, 0xd5, 0x00, 0x7d, 0x0b, 0x3d, 0xb5, 0xb2, 0x0e, 0xee, 0x94, 0xf0, 0x8b, 0x36, 0x3c,
	0xad, 0x6a, 0xb7, 0xe9, 0x0d, 0x42, 0xdf, 0x80, 0x27, 0xb2, 0x8c, 0xf9, 0x25, 0x79, 0xdc, 0x26,
	0x9d, 0x28, 0x77, 0xa2, 0xfc, 0x4b, 0x21, 0x7a, 0xee, 0x3c, 0xa3, 0x02, 0xa1, 0x27, 0xd0, 0xdf,
	0x98, 0xb1, 0xfb, 0x21, 0x19, 0xed, 0x4f, 0x9e, 0xb7, 0xf9, 0x4d, 0x91, 0x57, 0x4d, 0x9b, 0xed,
	0xf4, 0x3d, 0xf4, 0x71, 0x6d, 0x51, 0x9a, 0x44, 0x49, 0xd6, 0xdd, 0xa9, 0x77, 0x03, 0x9e, 0xf9,
	0x3d, 0x2f, 0xf0, 0x87, 0x3f, 0x09, 0x3c, 0xde, 0xf6, 0x8a, 0x68, 0x00, 0x9e, 0x15, 0x4b, 0x46,
	0x42, 0x32, 0xea, 0x47, 0xc5, 0x2d, 0x9d, 0xc1, 0x23, 0x8d, 0x31, 0x26, 0x37, 0xa8, 0x2f, 0x0d,
	0x5a, 0x9b, 0xc8, 0xa5, 0x61, 0x9d, 0x72, 0xf4, 0x7f, 0x6d, 0x1f, 0xd4, 0x01, 0xb3, 0x8a, 0xa7,
	0xe7, 0x30, 0xc8, 0xb4, 0x5a, 0xe7, 0x4d, 0xa2, 0xb7, 0x53, 0xe2, 0x83, 0x92, 0xae, 0xe3, 0x86,
	0xbf, 0x08, 0x3c, 0xd9, 0xfa, 0xd1, 0xb6, 0xf8, 0x4c, 0xe1, 0xa1, 0x41, 0xb9, 0xf8, 0x7f, 0x9b,
	0x81, 0xc3, 0xef, 0xc8, 0x85, 0x3e, 0x85, 0x2e, 0xae, 0xb3, 0x44, 0x23, 0xf3, 0x43, 0x32, 0xf2,
	0xa2, 0x6a, 0x45, 0x19, 0xec, 0x15, 0x21, 0x28, 0xdd, 0x8f, 0xd3, 0x8f, 0xea, 0xe5, 0xcb, 0x63,
	0x38, 0x70, 0xb6, 0x1f, 0x94, 0x4e, 0x85, 0xa5, 0x07, 0xd0, 0xbb, 0x28, 0x4e, 0xcb, 0x7c, 0x75,
	0x15, 0xdc, 0xa3, 0x3d, 0xf0, 0xcf, 0x66, 0xd3, 0xcf, 0x01, 0x79, 0x77, 0x04, 0x83, 0x58, 0xa5,
	0xad, 0xa9, 0x2e, 0xc8, 0x37, 0xbf, 0xb8, 0x7e, 0xef, 0xc0, 0xd7, 0x49, 0x24, 0x72, 0x7e, 0xaa,
	0x34, 0xce, 0xbb, 0xe5, 0x51, 0x7b, 0xfd, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x01, 0x08, 0x59, 0x1b,
	0xee, 0x03, 0x00, 0x00,
}
