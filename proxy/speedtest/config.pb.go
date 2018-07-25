package speedtest

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Config struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_a46c59b79195c0b1, []int{0}
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

func init() {
	proto.RegisterType((*Config)(nil), "v2ray.core.proxy.speedtest.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/speedtest/config.proto", fileDescriptor_config_a46c59b79195c0b1)
}

var fileDescriptor_config_a46c59b79195c0b1 = []byte{
	// 131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2e, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x2f, 0x28, 0xca, 0xaf, 0xa8,
	0xd4, 0x2f, 0x2e, 0x48, 0x4d, 0x4d, 0x29, 0x49, 0x2d, 0x2e, 0xd1, 0x4f, 0xce, 0xcf, 0x4b, 0xcb,
	0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x82, 0x29, 0x2e, 0x4a, 0xd5, 0x03, 0x2b,
	0xd4, 0x83, 0x2b, 0x54, 0xe2, 0xe0, 0x62, 0x73, 0x06, 0xab, 0x75, 0xf2, 0xe2, 0x92, 0x4b, 0xce,
	0xcf, 0xd5, 0xc3, 0xad, 0x36, 0x80, 0x31, 0x8a, 0x13, 0xce, 0x59, 0xc5, 0x24, 0x15, 0x66, 0x14,
	0x94, 0x58, 0xa9, 0xe7, 0x0c, 0x52, 0x19, 0x00, 0x56, 0x19, 0x0c, 0x92, 0x0c, 0x49, 0x2d, 0x2e,
	0x49, 0x62, 0x03, 0x5b, 0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x07, 0x81, 0xb2, 0x36, 0xa7,
	0x00, 0x00, 0x00,
}
