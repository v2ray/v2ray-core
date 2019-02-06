package dns

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Config struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_c49bb2d51e576d57, []int{0}
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

func init() {
	proto.RegisterType((*Config)(nil), "v2ray.core.proxy.dns.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/dns/config.proto", fileDescriptor_c49bb2d51e576d57)
}

var fileDescriptor_c49bb2d51e576d57 = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2d, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x2f, 0x28, 0xca, 0xaf, 0xa8,
	0xd4, 0x4f, 0xc9, 0x2b, 0xd6, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x12, 0x81, 0x29, 0x2b, 0x4a, 0xd5, 0x03, 0x2b, 0xd1, 0x4b, 0xc9, 0x2b, 0x56, 0xe2,
	0xe0, 0x62, 0x73, 0x06, 0xab, 0x72, 0xb2, 0xe0, 0x92, 0x48, 0xce, 0xcf, 0xd5, 0xc3, 0xa6, 0x2a,
	0x80, 0x31, 0x8a, 0x39, 0x25, 0xaf, 0x78, 0x15, 0x93, 0x48, 0x98, 0x51, 0x50, 0x62, 0xa5, 0x9e,
	0x33, 0x48, 0x36, 0x00, 0x2c, 0xeb, 0x92, 0x57, 0x9c, 0xc4, 0x06, 0xb6, 0xc0, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0xee, 0x22, 0xde, 0xc9, 0x89, 0x00, 0x00, 0x00,
}
