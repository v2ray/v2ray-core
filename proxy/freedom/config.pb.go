package freedom

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import protocol "v2ray.com/core/common/protocol"

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
	Config_AS_IS  Config_DomainStrategy = 0
	Config_USE_IP Config_DomainStrategy = 1
)

var Config_DomainStrategy_name = map[int32]string{
	0: "AS_IS",
	1: "USE_IP",
}
var Config_DomainStrategy_value = map[string]int32{
	"AS_IS":  0,
	"USE_IP": 1,
}

func (x Config_DomainStrategy) String() string {
	return proto.EnumName(Config_DomainStrategy_name, int32(x))
}
func (Config_DomainStrategy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_9714d01f3d402353, []int{1, 0}
}

type DestinationOverride struct {
	Server               *protocol.ServerEndpoint `protobuf:"bytes,1,opt,name=server" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *DestinationOverride) Reset()         { *m = DestinationOverride{} }
func (m *DestinationOverride) String() string { return proto.CompactTextString(m) }
func (*DestinationOverride) ProtoMessage()    {}
func (*DestinationOverride) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_9714d01f3d402353, []int{0}
}
func (m *DestinationOverride) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DestinationOverride.Unmarshal(m, b)
}
func (m *DestinationOverride) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DestinationOverride.Marshal(b, m, deterministic)
}
func (dst *DestinationOverride) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DestinationOverride.Merge(dst, src)
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
	DomainStrategy       Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,enum=v2ray.core.proxy.freedom.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	Timeout              uint32                `protobuf:"varint,2,opt,name=timeout" json:"timeout,omitempty"` // Deprecated: Do not use.
	DestinationOverride  *DestinationOverride  `protobuf:"bytes,3,opt,name=destination_override,json=destinationOverride" json:"destination_override,omitempty"`
	UserLevel            uint32                `protobuf:"varint,4,opt,name=user_level,json=userLevel" json:"user_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_9714d01f3d402353, []int{1}
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
	proto.RegisterType((*DestinationOverride)(nil), "v2ray.core.proxy.freedom.DestinationOverride")
	proto.RegisterType((*Config)(nil), "v2ray.core.proxy.freedom.Config")
	proto.RegisterEnum("v2ray.core.proxy.freedom.Config_DomainStrategy", Config_DomainStrategy_name, Config_DomainStrategy_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/freedom/config.proto", fileDescriptor_config_9714d01f3d402353)
}

var fileDescriptor_config_9714d01f3d402353 = []byte{
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x6f, 0x4b, 0x83, 0x50,
	0x14, 0xc6, 0xd3, 0xca, 0xb1, 0x13, 0xad, 0xe1, 0x7a, 0x21, 0xb1, 0x60, 0xec, 0x4d, 0x2b, 0xe8,
	0x1a, 0xf6, 0x09, 0xda, 0x9f, 0x60, 0x10, 0x34, 0x94, 0xa2, 0x7a, 0x63, 0xa6, 0x67, 0x43, 0x98,
	0xf7, 0xc8, 0xf5, 0x4e, 0xf2, 0x2b, 0xed, 0x53, 0x86, 0x57, 0x47, 0x2d, 0xb6, 0x77, 0xfa, 0xf8,
	0x7b, 0x9e, 0x73, 0x9e, 0x23, 0x5c, 0xe7, 0x8e, 0x08, 0x0a, 0x16, 0x52, 0x62, 0x87, 0x24, 0xd0,
	0x4e, 0x05, 0x7d, 0x17, 0xf6, 0x5c, 0x20, 0x46, 0x4a, 0xe2, 0xf3, 0x78, 0xc1, 0x52, 0x41, 0x92,
	0x4c, 0x6b, 0x83, 0x0a, 0x64, 0x0a, 0x63, 0x35, 0x76, 0x71, 0xf7, 0x2f, 0x24, 0xa4, 0x24, 0x21,
	0x6e, 0x2b, 0x5b, 0x48, 0x4b, 0x3b, 0x43, 0x91, 0xa3, 0xf0, 0xb3, 0x14, 0xc3, 0x2a, 0xab, 0xff,
	0x0e, 0x9d, 0x31, 0x66, 0x32, 0xe6, 0x81, 0x8c, 0x89, 0x3f, 0xe7, 0x28, 0x44, 0x1c, 0xa1, 0x39,
	0x04, 0xa3, 0x62, 0x2d, 0xad, 0xa7, 0x0d, 0x4e, 0x9c, 0x1b, 0xf6, 0x67, 0x66, 0x95, 0xca, 0x36,
	0xa9, 0xcc, 0x53, 0xe4, 0x84, 0x47, 0x29, 0xc5, 0x5c, 0xba, 0xb5, 0xb3, 0xbf, 0xd6, 0xc1, 0x18,
	0xa9, 0xbd, 0xcd, 0x37, 0x38, 0x8b, 0x28, 0x09, 0x62, 0xee, 0x67, 0x52, 0x04, 0x12, 0x17, 0x85,
	0xca, 0x6d, 0x39, 0x36, 0xdb, 0xd7, 0x85, 0x55, 0x56, 0x36, 0x56, 0x3e, 0xaf, 0xb6, 0xb9, 0xad,
	0x68, 0xeb, 0xdd, 0xec, 0x42, 0x43, 0xc6, 0x09, 0xd2, 0x4a, 0x5a, 0x7a, 0x4f, 0x1b, 0x9c, 0x0e,
	0x75, 0x4b, 0x73, 0x37, 0x92, 0xf9, 0x09, 0xe7, 0xd1, 0x6f, 0x3b, 0x9f, 0xea, 0x7a, 0xd6, 0xa1,
	0x2a, 0x75, 0xbb, 0x7f, 0xf8, 0x8e, 0x9b, 0xb8, 0x9d, 0x68, 0xc7, 0xa1, 0x2e, 0x01, 0x56, 0x19,
	0x0a, 0x7f, 0x89, 0x39, 0x2e, 0xad, 0xa3, 0x72, 0x05, 0xb7, 0x59, 0x2a, 0x4f, 0xa5, 0xd0, 0xbf,
	0x82, 0xd6, 0x76, 0x01, 0xb3, 0x09, 0xc7, 0x0f, 0x9e, 0x3f, 0xf5, 0xda, 0x07, 0x26, 0x80, 0xf1,
	0xe2, 0x4d, 0xfc, 0xe9, 0xac, 0xad, 0x0d, 0xc7, 0xd0, 0x0d, 0x29, 0xd9, 0xbb, 0xd0, 0x4c, 0xfb,
	0x68, 0xd4, 0x8f, 0x6b, 0xdd, 0x7a, 0x75, 0xdc, 0xa0, 0x60, 0xa3, 0x92, 0x9a, 0x29, 0xea, 0xb1,
	0xfa, 0xf4, 0x65, 0xa8, 0x7f, 0x72, 0xff, 0x13, 0x00, 0x00, 0xff, 0xff, 0xb7, 0x5e, 0xda, 0x4d,
	0x4d, 0x02, 0x00, 0x00,
}
