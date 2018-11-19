package freedom

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	protocol "v2ray.com/core/common/protocol"
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

type Config_DomainStrategy int32

const (
	Config_AS_IS   Config_DomainStrategy = 0
	Config_USE_IP  Config_DomainStrategy = 1
	Config_USE_IP4 Config_DomainStrategy = 2
	Config_USE_IP6 Config_DomainStrategy = 3
)

var Config_DomainStrategy_name = map[int32]string{
	0: "AS_IS",
	1: "USE_IP",
	2: "USE_IP4",
	3: "USE_IP6",
}

var Config_DomainStrategy_value = map[string]int32{
	"AS_IS":   0,
	"USE_IP":  1,
	"USE_IP4": 2,
	"USE_IP6": 3,
}

func (x Config_DomainStrategy) String() string {
	return proto.EnumName(Config_DomainStrategy_name, int32(x))
}

func (Config_DomainStrategy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_66807b6fe2cca4da, []int{1, 0}
}

type DestinationOverride struct {
	Server               *protocol.ServerEndpoint `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *DestinationOverride) Reset()         { *m = DestinationOverride{} }
func (m *DestinationOverride) String() string { return proto.CompactTextString(m) }
func (*DestinationOverride) ProtoMessage()    {}
func (*DestinationOverride) Descriptor() ([]byte, []int) {
	return fileDescriptor_66807b6fe2cca4da, []int{0}
}

func (m *DestinationOverride) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DestinationOverride.Unmarshal(m, b)
}
func (m *DestinationOverride) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DestinationOverride.Marshal(b, m, deterministic)
}
func (m *DestinationOverride) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DestinationOverride.Merge(m, src)
}
func (m *DestinationOverride) XXX_Size() int {
	return xxx_messageInfo_DestinationOverride.Size(m)
}
func (m *DestinationOverride) XXX_DiscardUnknown() {
	xxx_messageInfo_DestinationOverride.DiscardUnknown(m)
}

var xxx_messageInfo_DestinationOverride proto.InternalMessageInfo

func (m *DestinationOverride) GetServer() *protocol.ServerEndpoint {
	if m != nil {
		return m.Server
	}
	return nil
}

type Config struct {
	DomainStrategy       Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,proto3,enum=v2ray.core.proxy.freedom.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	Timeout              uint32                `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"` // Deprecated: Do not use.
	DestinationOverride  *DestinationOverride  `protobuf:"bytes,3,opt,name=destination_override,json=destinationOverride,proto3" json:"destination_override,omitempty"`
	UserLevel            uint32                `protobuf:"varint,4,opt,name=user_level,json=userLevel,proto3" json:"user_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_66807b6fe2cca4da, []int{1}
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

func (m *Config) GetDomainStrategy() Config_DomainStrategy {
	if m != nil {
		return m.DomainStrategy
	}
	return Config_AS_IS
}

// Deprecated: Do not use.
func (m *Config) GetTimeout() uint32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *Config) GetDestinationOverride() *DestinationOverride {
	if m != nil {
		return m.DestinationOverride
	}
	return nil
}

func (m *Config) GetUserLevel() uint32 {
	if m != nil {
		return m.UserLevel
	}
	return 0
}

func init() {
	proto.RegisterEnum("v2ray.core.proxy.freedom.Config_DomainStrategy", Config_DomainStrategy_name, Config_DomainStrategy_value)
	proto.RegisterType((*DestinationOverride)(nil), "v2ray.core.proxy.freedom.DestinationOverride")
	proto.RegisterType((*Config)(nil), "v2ray.core.proxy.freedom.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/freedom/config.proto", fileDescriptor_66807b6fe2cca4da)
}

var fileDescriptor_66807b6fe2cca4da = []byte{
	// 357 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0x5d, 0x4b, 0xeb, 0x40,
	0x10, 0xbd, 0x49, 0xef, 0x4d, 0xe9, 0x94, 0xdb, 0x1b, 0xb6, 0xf7, 0x21, 0x48, 0x85, 0xd2, 0xa7,
	0x2a, 0xb8, 0x91, 0x28, 0xbe, 0xf7, 0x4b, 0x28, 0x08, 0x96, 0x04, 0x45, 0x7d, 0x89, 0x31, 0x99,
	0x96, 0x40, 0xb3, 0x1b, 0x36, 0xdb, 0x60, 0xfe, 0x92, 0xbf, 0xc3, 0x1f, 0x26, 0xd9, 0xa4, 0xd4,
	0x4a, 0xfb, 0x36, 0x73, 0xf6, 0x9c, 0x33, 0x73, 0x66, 0xe1, 0x2c, 0x77, 0x44, 0x50, 0xd0, 0x90,
	0x27, 0x76, 0xc8, 0x05, 0xda, 0xa9, 0xe0, 0xef, 0x85, 0xbd, 0x14, 0x88, 0x91, 0x82, 0xd8, 0x32,
	0x5e, 0xd1, 0x54, 0x70, 0xc9, 0x89, 0xb5, 0xa5, 0x0a, 0xa4, 0x8a, 0x46, 0x6b, 0xda, 0xc9, 0xe5,
	0x0f, 0x93, 0x90, 0x27, 0x09, 0x67, 0xb6, 0x92, 0x85, 0x7c, 0x6d, 0x67, 0x28, 0x72, 0x14, 0x7e,
	0x96, 0x62, 0x58, 0x79, 0x0d, 0x9e, 0xa1, 0x3b, 0xc5, 0x4c, 0xc6, 0x2c, 0x90, 0x31, 0x67, 0xf7,
	0x39, 0x0a, 0x11, 0x47, 0x48, 0xc6, 0x60, 0x54, 0x5c, 0x4b, 0xeb, 0x6b, 0xc3, 0xb6, 0x73, 0x4e,
	0xbf, 0xcd, 0xac, 0x5c, 0xe9, 0xd6, 0x95, 0x7a, 0x8a, 0x39, 0x63, 0x51, 0xca, 0x63, 0x26, 0xdd,
	0x5a, 0x39, 0xf8, 0xd4, 0xc1, 0x98, 0xa8, 0xbd, 0xc9, 0x13, 0xfc, 0x8b, 0x78, 0x12, 0xc4, 0xcc,
	0xcf, 0xa4, 0x08, 0x24, 0xae, 0x0a, 0xe5, 0xdb, 0x71, 0x6c, 0x7a, 0x2c, 0x0b, 0xad, 0xa4, 0x74,
	0xaa, 0x74, 0x5e, 0x2d, 0x73, 0x3b, 0xd1, 0x5e, 0x4f, 0x7a, 0xd0, 0x94, 0x71, 0x82, 0x7c, 0x23,
	0x2d, 0xbd, 0xaf, 0x0d, 0xff, 0x8e, 0x75, 0x4b, 0x73, 0xb7, 0x10, 0x79, 0x85, 0xff, 0xd1, 0x2e,
	0x9d, 0xcf, 0xeb, 0x78, 0x56, 0x43, 0x85, 0xba, 0x38, 0x3e, 0xfc, 0xc0, 0x4d, 0xdc, 0x6e, 0x74,
	0xe0, 0x50, 0xa7, 0x00, 0x9b, 0x0c, 0x85, 0xbf, 0xc6, 0x1c, 0xd7, 0xd6, 0xef, 0x72, 0x05, 0xb7,
	0x55, 0x22, 0x77, 0x25, 0x30, 0x18, 0x41, 0x67, 0x3f, 0x00, 0x69, 0xc1, 0x9f, 0x91, 0xe7, 0xcf,
	0x3d, 0xf3, 0x17, 0x01, 0x30, 0x1e, 0xbc, 0x99, 0x3f, 0x5f, 0x98, 0x1a, 0x69, 0x43, 0xb3, 0xaa,
	0xaf, 0x4d, 0x7d, 0xd7, 0xdc, 0x98, 0x8d, 0xf1, 0x14, 0x7a, 0x21, 0x4f, 0x8e, 0xae, 0xba, 0xd0,
	0x5e, 0x9a, 0x75, 0xf9, 0xa1, 0x5b, 0x8f, 0x8e, 0x1b, 0x14, 0x74, 0x52, 0xb2, 0x16, 0x8a, 0x75,
	0x5b, 0x3d, 0xbd, 0x19, 0xea, 0xb7, 0xae, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x72, 0x98, 0x5b,
	0x40, 0x67, 0x02, 0x00, 0x00,
}
