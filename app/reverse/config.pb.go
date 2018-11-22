package reverse

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Control_State int32

const (
	Control_ACTIVE Control_State = 0
	Control_DRAIN  Control_State = 1
)

var Control_State_name = map[int32]string{
	0: "ACTIVE",
	1: "DRAIN",
}

var Control_State_value = map[string]int32{
	"ACTIVE": 0,
	"DRAIN":  1,
}

func (x Control_State) String() string {
	return proto.EnumName(Control_State_name, int32(x))
}

func (Control_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_829a0eeb60380cbc, []int{0, 0}
}

type Control struct {
	State                Control_State `protobuf:"varint,1,opt,name=state,proto3,enum=v2ray.core.app.reverse.Control_State" json:"state,omitempty"`
	Random               []byte        `protobuf:"bytes,99,opt,name=random,proto3" json:"random,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Control) Reset()         { *m = Control{} }
func (m *Control) String() string { return proto.CompactTextString(m) }
func (*Control) ProtoMessage()    {}
func (*Control) Descriptor() ([]byte, []int) {
	return fileDescriptor_829a0eeb60380cbc, []int{0}
}

func (m *Control) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Control.Unmarshal(m, b)
}
func (m *Control) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Control.Marshal(b, m, deterministic)
}
func (m *Control) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Control.Merge(m, src)
}
func (m *Control) XXX_Size() int {
	return xxx_messageInfo_Control.Size(m)
}
func (m *Control) XXX_DiscardUnknown() {
	xxx_messageInfo_Control.DiscardUnknown(m)
}

var xxx_messageInfo_Control proto.InternalMessageInfo

func (m *Control) GetState() Control_State {
	if m != nil {
		return m.State
	}
	return Control_ACTIVE
}

func (m *Control) GetRandom() []byte {
	if m != nil {
		return m.Random
	}
	return nil
}

type BridgeConfig struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Domain               string   `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BridgeConfig) Reset()         { *m = BridgeConfig{} }
func (m *BridgeConfig) String() string { return proto.CompactTextString(m) }
func (*BridgeConfig) ProtoMessage()    {}
func (*BridgeConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_829a0eeb60380cbc, []int{1}
}

func (m *BridgeConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BridgeConfig.Unmarshal(m, b)
}
func (m *BridgeConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BridgeConfig.Marshal(b, m, deterministic)
}
func (m *BridgeConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgeConfig.Merge(m, src)
}
func (m *BridgeConfig) XXX_Size() int {
	return xxx_messageInfo_BridgeConfig.Size(m)
}
func (m *BridgeConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgeConfig.DiscardUnknown(m)
}

var xxx_messageInfo_BridgeConfig proto.InternalMessageInfo

func (m *BridgeConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *BridgeConfig) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

type PortalConfig struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Domain               string   `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PortalConfig) Reset()         { *m = PortalConfig{} }
func (m *PortalConfig) String() string { return proto.CompactTextString(m) }
func (*PortalConfig) ProtoMessage()    {}
func (*PortalConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_829a0eeb60380cbc, []int{2}
}

func (m *PortalConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PortalConfig.Unmarshal(m, b)
}
func (m *PortalConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PortalConfig.Marshal(b, m, deterministic)
}
func (m *PortalConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PortalConfig.Merge(m, src)
}
func (m *PortalConfig) XXX_Size() int {
	return xxx_messageInfo_PortalConfig.Size(m)
}
func (m *PortalConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_PortalConfig.DiscardUnknown(m)
}

var xxx_messageInfo_PortalConfig proto.InternalMessageInfo

func (m *PortalConfig) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *PortalConfig) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

type Config struct {
	BridgeConfig         []*BridgeConfig `protobuf:"bytes,1,rep,name=bridge_config,json=bridgeConfig,proto3" json:"bridge_config,omitempty"`
	PortalConfig         []*PortalConfig `protobuf:"bytes,2,rep,name=portal_config,json=portalConfig,proto3" json:"portal_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_829a0eeb60380cbc, []int{3}
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

func (m *Config) GetBridgeConfig() []*BridgeConfig {
	if m != nil {
		return m.BridgeConfig
	}
	return nil
}

func (m *Config) GetPortalConfig() []*PortalConfig {
	if m != nil {
		return m.PortalConfig
	}
	return nil
}

func init() {
	proto.RegisterEnum("v2ray.core.app.reverse.Control_State", Control_State_name, Control_State_value)
	proto.RegisterType((*Control)(nil), "v2ray.core.app.reverse.Control")
	proto.RegisterType((*BridgeConfig)(nil), "v2ray.core.app.reverse.BridgeConfig")
	proto.RegisterType((*PortalConfig)(nil), "v2ray.core.app.reverse.PortalConfig")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.reverse.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/reverse/config.proto", fileDescriptor_829a0eeb60380cbc)
}

var fileDescriptor_829a0eeb60380cbc = []byte{
	// 310 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xc1, 0x4b, 0xfb, 0x30,
	0x14, 0xc7, 0x7f, 0xd9, 0x58, 0xc7, 0xde, 0xaf, 0x93, 0x91, 0xc3, 0xe8, 0x41, 0x64, 0x0c, 0xc5,
	0x9d, 0x52, 0xa8, 0x17, 0xc1, 0xd3, 0xd6, 0x79, 0xe8, 0x45, 0x46, 0x94, 0x1d, 0xbc, 0x48, 0xd6,
	0xc5, 0x52, 0x58, 0xfb, 0x42, 0x16, 0x86, 0xbd, 0xf8, 0xa7, 0xf8, 0x07, 0xf8, 0x57, 0x4a, 0xd3,
	0x14, 0x7a, 0x50, 0xc1, 0xdb, 0x7b, 0xc9, 0xe7, 0xf3, 0xf2, 0x4d, 0x02, 0xd7, 0xa7, 0x48, 0x8b,
	0x8a, 0xa5, 0x58, 0x84, 0x29, 0x6a, 0x19, 0x0a, 0xa5, 0x42, 0x2d, 0x4f, 0x52, 0x1f, 0x65, 0x98,
	0x62, 0xf9, 0x9a, 0x67, 0x4c, 0x69, 0x34, 0x48, 0xa7, 0x2d, 0xa8, 0x25, 0x13, 0x4a, 0x31, 0x07,
	0xcd, 0xdf, 0x61, 0x18, 0x63, 0x69, 0x34, 0x1e, 0xe8, 0x1d, 0x0c, 0x8e, 0x46, 0x18, 0x19, 0x90,
	0x19, 0x59, 0x9c, 0x45, 0x57, 0xec, 0x7b, 0x85, 0x39, 0x9e, 0x3d, 0xd6, 0x30, 0x6f, 0x1c, 0x3a,
	0x05, 0x4f, 0x8b, 0x72, 0x8f, 0x45, 0x90, 0xce, 0xc8, 0xc2, 0xe7, 0xae, 0x9b, 0x5f, 0xc0, 0xc0,
	0x72, 0x14, 0xc0, 0x5b, 0xc6, 0x4f, 0xc9, 0xf6, 0x7e, 0xf2, 0x8f, 0x8e, 0x60, 0xb0, 0xe6, 0xcb,
	0xe4, 0x61, 0x42, 0xe6, 0xb7, 0xe0, 0xaf, 0x74, 0xbe, 0xcf, 0x64, 0x6c, 0xd3, 0xd2, 0x09, 0xf4,
	0x8d, 0xc8, 0x6c, 0x84, 0x11, 0xaf, 0xcb, 0x7a, 0xf2, 0x1e, 0x0b, 0x91, 0x97, 0x41, 0xcf, 0x2e,
	0xba, 0xae, 0x36, 0x37, 0xa8, 0x8d, 0x38, 0xfc, 0xd9, 0xfc, 0x20, 0xe0, 0x39, 0x29, 0x81, 0xf1,
	0xce, 0x1e, 0xff, 0xd2, 0xbc, 0x56, 0x40, 0x66, 0xfd, 0xc5, 0xff, 0xe8, 0xf2, 0xa7, 0xbb, 0x77,
	0xb3, 0x72, 0x7f, 0xd7, 0x4d, 0x9e, 0xc0, 0x58, 0xd9, 0x3c, 0xed, 0xa8, 0xde, 0xef, 0xa3, 0xba,
	0xe1, 0xb9, 0xaf, 0x3a, 0xdd, 0x6a, 0x0d, 0xe7, 0x29, 0x16, 0x5d, 0x51, 0x69, 0x7c, 0xab, 0x5a,
	0x75, 0x43, 0x9e, 0x87, 0xae, 0xfc, 0xec, 0x05, 0xdb, 0x88, 0x8b, 0x8a, 0xc5, 0x35, 0xb5, 0xb1,
	0x14, 0x6f, 0xb6, 0x76, 0x9e, 0xfd, 0xf9, 0x9b, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb7, 0x83,
	0x30, 0xdb, 0x24, 0x02, 0x00, 0x00,
}
