package dns

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import net "v2ray.com/core/common/net"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Config_HostMapping_Type int32

const (
	Config_HostMapping_Full      Config_HostMapping_Type = 0
	Config_HostMapping_SubDomain Config_HostMapping_Type = 1
)

var Config_HostMapping_Type_name = map[int32]string{
	0: "Full",
	1: "SubDomain",
}
var Config_HostMapping_Type_value = map[string]int32{
	"Full":      0,
	"SubDomain": 1,
}

func (x Config_HostMapping_Type) String() string {
	return proto.EnumName(Config_HostMapping_Type_name, int32(x))
}
func (Config_HostMapping_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_299ae69458dda3a0, []int{0, 2, 0}
}

type Config struct {
	// Nameservers used by this DNS. Only traditional UDP servers are support at the moment.
	// A special value 'localhost' as a domain address can be set to use DNS on local system.
	NameServers []*net.Endpoint `protobuf:"bytes,1,rep,name=NameServers,proto3" json:"NameServers,omitempty"`
	// Static hosts. Domain to IP.
	// Deprecated. Use static_hosts.
	Hosts map[string]*net.IPOrDomain `protobuf:"bytes,2,rep,name=Hosts,proto3" json:"Hosts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // Deprecated: Do not use.
	// Client IP for EDNS client subnet.
	ClientIp             *Config_ClientIP      `protobuf:"bytes,3,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
	StaticHosts          []*Config_HostMapping `protobuf:"bytes,4,rep,name=static_hosts,json=staticHosts,proto3" json:"static_hosts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_299ae69458dda3a0, []int{0}
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

func (m *Config) GetNameServers() []*net.Endpoint {
	if m != nil {
		return m.NameServers
	}
	return nil
}

// Deprecated: Do not use.
func (m *Config) GetHosts() map[string]*net.IPOrDomain {
	if m != nil {
		return m.Hosts
	}
	return nil
}

func (m *Config) GetClientIp() *Config_ClientIP {
	if m != nil {
		return m.ClientIp
	}
	return nil
}

func (m *Config) GetStaticHosts() []*Config_HostMapping {
	if m != nil {
		return m.StaticHosts
	}
	return nil
}

type Config_ClientIP struct {
	// IPv4 address of the client. Must be 4 bytes.
	V4 []byte `protobuf:"bytes,1,opt,name=v4,proto3" json:"v4,omitempty"`
	// IPv6 address of the client. Must be 4 bytes.
	V6                   []byte   `protobuf:"bytes,2,opt,name=v6,proto3" json:"v6,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config_ClientIP) Reset()         { *m = Config_ClientIP{} }
func (m *Config_ClientIP) String() string { return proto.CompactTextString(m) }
func (*Config_ClientIP) ProtoMessage()    {}
func (*Config_ClientIP) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_299ae69458dda3a0, []int{0, 1}
}
func (m *Config_ClientIP) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config_ClientIP.Unmarshal(m, b)
}
func (m *Config_ClientIP) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config_ClientIP.Marshal(b, m, deterministic)
}
func (dst *Config_ClientIP) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config_ClientIP.Merge(dst, src)
}
func (m *Config_ClientIP) XXX_Size() int {
	return xxx_messageInfo_Config_ClientIP.Size(m)
}
func (m *Config_ClientIP) XXX_DiscardUnknown() {
	xxx_messageInfo_Config_ClientIP.DiscardUnknown(m)
}

var xxx_messageInfo_Config_ClientIP proto.InternalMessageInfo

func (m *Config_ClientIP) GetV4() []byte {
	if m != nil {
		return m.V4
	}
	return nil
}

func (m *Config_ClientIP) GetV6() []byte {
	if m != nil {
		return m.V6
	}
	return nil
}

type Config_HostMapping struct {
	Type                 Config_HostMapping_Type `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.dns.Config_HostMapping_Type" json:"type,omitempty"`
	Domain               string                  `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Ip                   [][]byte                `protobuf:"bytes,3,rep,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *Config_HostMapping) Reset()         { *m = Config_HostMapping{} }
func (m *Config_HostMapping) String() string { return proto.CompactTextString(m) }
func (*Config_HostMapping) ProtoMessage()    {}
func (*Config_HostMapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_299ae69458dda3a0, []int{0, 2}
}
func (m *Config_HostMapping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config_HostMapping.Unmarshal(m, b)
}
func (m *Config_HostMapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config_HostMapping.Marshal(b, m, deterministic)
}
func (dst *Config_HostMapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config_HostMapping.Merge(dst, src)
}
func (m *Config_HostMapping) XXX_Size() int {
	return xxx_messageInfo_Config_HostMapping.Size(m)
}
func (m *Config_HostMapping) XXX_DiscardUnknown() {
	xxx_messageInfo_Config_HostMapping.DiscardUnknown(m)
}

var xxx_messageInfo_Config_HostMapping proto.InternalMessageInfo

func (m *Config_HostMapping) GetType() Config_HostMapping_Type {
	if m != nil {
		return m.Type
	}
	return Config_HostMapping_Full
}

func (m *Config_HostMapping) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *Config_HostMapping) GetIp() [][]byte {
	if m != nil {
		return m.Ip
	}
	return nil
}

func init() {
	proto.RegisterType((*Config)(nil), "v2ray.core.app.dns.Config")
	proto.RegisterMapType((map[string]*net.IPOrDomain)(nil), "v2ray.core.app.dns.Config.HostsEntry")
	proto.RegisterType((*Config_ClientIP)(nil), "v2ray.core.app.dns.Config.ClientIP")
	proto.RegisterType((*Config_HostMapping)(nil), "v2ray.core.app.dns.Config.HostMapping")
	proto.RegisterEnum("v2ray.core.app.dns.Config_HostMapping_Type", Config_HostMapping_Type_name, Config_HostMapping_Type_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dns/config.proto", fileDescriptor_config_299ae69458dda3a0)
}

var fileDescriptor_config_299ae69458dda3a0 = []byte{
	// 444 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x86, 0x71, 0x92, 0x55, 0xed, 0xc9, 0x98, 0x2a, 0x5f, 0x4c, 0x51, 0x6e, 0x56, 0x36, 0x01,
	0x15, 0x93, 0x1c, 0x29, 0x54, 0x03, 0x71, 0x33, 0xb6, 0x6e, 0x88, 0x5e, 0x00, 0x95, 0x87, 0xb8,
	0x80, 0x8b, 0xc9, 0x4b, 0xcc, 0xb0, 0x68, 0x6c, 0xcb, 0x76, 0x23, 0xe5, 0x49, 0x78, 0x07, 0x9e,
	0x81, 0x87, 0x43, 0xb1, 0x99, 0xa8, 0x80, 0x01, 0x77, 0x3d, 0xea, 0xff, 0x9d, 0xff, 0x3b, 0x56,
	0xe0, 0xa0, 0x2d, 0x0d, 0xeb, 0x48, 0xa5, 0x9a, 0xa2, 0x52, 0x86, 0x17, 0x4c, 0xeb, 0xa2, 0x96,
	0xb6, 0xa8, 0x94, 0xfc, 0x28, 0xae, 0x89, 0x36, 0xca, 0x29, 0x8c, 0x6f, 0x42, 0x86, 0x13, 0xa6,
	0x35, 0xa9, 0xa5, 0xcd, 0x1f, 0xfe, 0x02, 0x56, 0xaa, 0x69, 0x94, 0x2c, 0x24, 0x77, 0x05, 0xab,
	0x6b, 0xc3, 0xad, 0x0d, 0x70, 0x7e, 0x78, 0x7b, 0xb0, 0xe6, 0xd6, 0x09, 0xc9, 0x9c, 0x50, 0x32,
	0x84, 0xf7, 0xbf, 0x25, 0x30, 0x98, 0xfb, 0x6a, 0x7c, 0x02, 0xe9, 0x6b, 0xd6, 0xf0, 0x0b, 0x6e,
	0x5a, 0x6e, 0x6c, 0x86, 0x26, 0xf1, 0x34, 0x2d, 0xf7, 0xc8, 0x86, 0x4a, 0xd8, 0x44, 0x24, 0x77,
	0xe4, 0x5c, 0xd6, 0x5a, 0x09, 0xe9, 0xe8, 0x26, 0x83, 0x8f, 0x61, 0xeb, 0xa5, 0xb2, 0xce, 0x66,
	0x91, 0x87, 0xef, 0x93, 0xdf, 0xef, 0x20, 0xa1, 0x8d, 0xf8, 0xdc, 0xb9, 0x74, 0xa6, 0x3b, 0x8d,
	0x32, 0x44, 0x03, 0x87, 0x9f, 0xc3, 0xa8, 0x5a, 0x09, 0x2e, 0xdd, 0xa5, 0xd0, 0x59, 0x3c, 0x41,
	0xd3, 0xb4, 0x3c, 0xf8, 0xcb, 0x92, 0xb9, 0xcf, 0x2e, 0x96, 0x74, 0x18, 0xa8, 0x85, 0xc6, 0x0b,
	0xd8, 0xb6, 0x8e, 0x39, 0x51, 0x5d, 0x7e, 0xf2, 0x26, 0x89, 0x37, 0x79, 0xf0, 0x0f, 0x93, 0x57,
	0x4c, 0x6b, 0x21, 0xaf, 0x69, 0x1a, 0x58, 0x2f, 0x93, 0x7f, 0x00, 0xf8, 0x69, 0x89, 0xc7, 0x10,
	0x7f, 0xe6, 0x5d, 0x86, 0x26, 0x68, 0x3a, 0xa2, 0xfd, 0x4f, 0xfc, 0x04, 0xb6, 0x5a, 0xb6, 0x5a,
	0xf3, 0x2c, 0xf2, 0xa2, 0xf7, 0x6e, 0x79, 0xaa, 0xc5, 0xf2, 0x8d, 0x39, 0x53, 0x0d, 0x13, 0x92,
	0x86, 0xfc, 0xb3, 0xe8, 0x29, 0xca, 0x1f, 0xc1, 0xf0, 0xc6, 0x1e, 0xef, 0x40, 0xd4, 0xce, 0xfc,
	0xe6, 0x6d, 0x1a, 0xb5, 0x33, 0x3f, 0x1f, 0xf9, 0xad, 0xfd, 0x7c, 0x94, 0x7f, 0x41, 0x90, 0x6e,
	0x58, 0xe2, 0x63, 0x48, 0x5c, 0xa7, 0xb9, 0x27, 0x76, 0xca, 0xc3, 0xff, 0xbb, 0x8d, 0xbc, 0xed,
	0x34, 0xa7, 0x1e, 0xc4, 0xbb, 0x30, 0xa8, 0xbd, 0x91, 0x2f, 0x19, 0xd1, 0x1f, 0x53, 0x5f, 0xec,
	0xdf, 0x3d, 0xee, 0x8b, 0x85, 0xde, 0xdf, 0x83, 0xa4, 0xa7, 0xf0, 0x10, 0x92, 0x17, 0xeb, 0xd5,
	0x6a, 0x7c, 0x07, 0xdf, 0x85, 0xd1, 0xc5, 0xfa, 0x2a, 0x9c, 0x33, 0x46, 0xa7, 0x33, 0xd8, 0xad,
	0x54, 0xf3, 0x07, 0x81, 0x25, 0x7a, 0x1f, 0xd7, 0xd2, 0x7e, 0x8d, 0xf0, 0xbb, 0x92, 0xb2, 0x8e,
	0xcc, 0xfb, 0xff, 0x4e, 0xb4, 0x26, 0x67, 0xd2, 0x5e, 0x0d, 0xfc, 0xb7, 0xf7, 0xf8, 0x7b, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xee, 0x90, 0x3f, 0xc9, 0x0c, 0x03, 0x00, 0x00,
}
