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
	return fileDescriptor_config_2b74f789a4ce2a14, []int{0, 1, 0}
}

type Config struct {
	// Nameservers used by this DNS. Only traditional UDP servers are support at the moment.
	// A special value 'localhost' as a domain address can be set to use DNS on local system.
	NameServers []*net.Endpoint `protobuf:"bytes,1,rep,name=NameServers,proto3" json:"NameServers,omitempty"`
	// Static hosts. Domain to IP.
	// Deprecated. Use static_hosts.
	Hosts map[string]*net.IPOrDomain `protobuf:"bytes,2,rep,name=Hosts,proto3" json:"Hosts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // Deprecated: Do not use.
	// Client IP for EDNS client subnet. Must be 4 bytes (IPv4) or 16 bytes (IPv6).
	ClientIp             []byte                `protobuf:"bytes,3,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
	StaticHosts          []*Config_HostMapping `protobuf:"bytes,4,rep,name=static_hosts,json=staticHosts,proto3" json:"static_hosts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_2b74f789a4ce2a14, []int{0}
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

func (m *Config) GetClientIp() []byte {
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
	return fileDescriptor_config_2b74f789a4ce2a14, []int{0, 1}
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
	proto.RegisterType((*Config_HostMapping)(nil), "v2ray.core.app.dns.Config.HostMapping")
	proto.RegisterEnum("v2ray.core.app.dns.Config_HostMapping_Type", Config_HostMapping_Type_name, Config_HostMapping_Type_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dns/config.proto", fileDescriptor_config_2b74f789a4ce2a14)
}

var fileDescriptor_config_2b74f789a4ce2a14 = []byte{
	// 416 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xdd, 0x6a, 0x14, 0x31,
	0x14, 0xc7, 0xcd, 0xcc, 0x76, 0xe9, 0x9e, 0x59, 0xcb, 0x92, 0x8b, 0x32, 0xac, 0x17, 0x5d, 0x2b,
	0xea, 0x40, 0x21, 0x03, 0xa3, 0xa0, 0x78, 0x53, 0xfa, 0x25, 0xee, 0x85, 0x5a, 0x52, 0xf1, 0x42,
	0x2f, 0x4a, 0x3a, 0x89, 0x35, 0xb8, 0x73, 0x12, 0x92, 0xec, 0xc2, 0x3c, 0x89, 0xef, 0xe0, 0xab,
	0xf9, 0x12, 0xb2, 0x89, 0xe2, 0xa2, 0x16, 0x7b, 0x37, 0x1f, 0xff, 0xdf, 0xf9, 0xfd, 0x73, 0x08,
	0x3c, 0x58, 0x35, 0x4e, 0xf4, 0xac, 0x35, 0x5d, 0xdd, 0x1a, 0xa7, 0x6a, 0x61, 0x6d, 0x2d, 0xd1,
	0xd7, 0xad, 0xc1, 0x4f, 0xfa, 0x9a, 0x59, 0x67, 0x82, 0xa1, 0xf4, 0x57, 0xc8, 0x29, 0x26, 0xac,
	0x65, 0x12, 0xfd, 0xf4, 0xf1, 0x1f, 0x60, 0x6b, 0xba, 0xce, 0x60, 0x8d, 0x2a, 0xd4, 0x42, 0x4a,
	0xa7, 0xbc, 0x4f, 0xf0, 0xf4, 0xe0, 0xe6, 0xa0, 0x54, 0x3e, 0x68, 0x14, 0x41, 0x1b, 0x4c, 0xe1,
	0xfd, 0xef, 0x39, 0x0c, 0x4f, 0xa2, 0x9a, 0x1e, 0x41, 0xf1, 0x46, 0x74, 0xea, 0x42, 0xb9, 0x95,
	0x72, 0xbe, 0x24, 0xb3, 0xbc, 0x2a, 0x9a, 0x3d, 0xb6, 0x51, 0x25, 0x4d, 0x62, 0xa8, 0x02, 0x3b,
	0x43, 0x69, 0x8d, 0xc6, 0xc0, 0x37, 0x19, 0x7a, 0x08, 0x5b, 0xaf, 0x8c, 0x0f, 0xbe, 0xcc, 0x22,
	0xfc, 0x90, 0xfd, 0x7d, 0x0e, 0x96, 0x6c, 0x2c, 0xe6, 0xce, 0x30, 0xb8, 0xfe, 0x38, 0x2b, 0x09,
	0x4f, 0x1c, 0xbd, 0x07, 0xa3, 0x76, 0xa1, 0x15, 0x86, 0x4b, 0x6d, 0xcb, 0x7c, 0x46, 0xaa, 0x31,
	0xdf, 0x4e, 0x1f, 0xe6, 0x96, 0xce, 0x61, 0xec, 0x83, 0x08, 0xba, 0xbd, 0xfc, 0x1c, 0x25, 0x83,
	0x28, 0x79, 0xf4, 0x1f, 0xc9, 0x6b, 0x61, 0xad, 0xc6, 0x6b, 0x5e, 0x24, 0x36, 0x7a, 0xa6, 0x1f,
	0x01, 0x7e, 0x17, 0xa0, 0x13, 0xc8, 0xbf, 0xa8, 0xbe, 0x24, 0x33, 0x52, 0x8d, 0xf8, 0xfa, 0x91,
	0x3e, 0x83, 0xad, 0x95, 0x58, 0x2c, 0x55, 0x99, 0xcd, 0x48, 0x55, 0x34, 0xf7, 0x6f, 0xd8, 0xc2,
	0xfc, 0xfc, 0xad, 0x3b, 0x35, 0x9d, 0xd0, 0xc8, 0x53, 0xfe, 0x45, 0xf6, 0x9c, 0x4c, 0xbf, 0x12,
	0x28, 0x36, 0xcc, 0xf4, 0x10, 0x06, 0xa1, 0xb7, 0x2a, 0xce, 0xdf, 0x69, 0x0e, 0x6e, 0xd7, 0x97,
	0xbd, 0xeb, 0xad, 0xe2, 0x11, 0xa4, 0xbb, 0x30, 0x94, 0xd1, 0x12, 0xeb, 0x8c, 0xf8, 0xcf, 0x37,
	0xba, 0x03, 0x59, 0x5c, 0x53, 0x5e, 0x8d, 0x79, 0xa6, 0xed, 0xfe, 0x1e, 0x0c, 0xd6, 0x14, 0xdd,
	0x86, 0xc1, 0xcb, 0xe5, 0x62, 0x31, 0xb9, 0x43, 0xef, 0xc2, 0xe8, 0x62, 0x79, 0x95, 0x2a, 0x4e,
	0xc8, 0xf1, 0x53, 0xd8, 0x6d, 0x4d, 0xf7, 0x8f, 0x02, 0xe7, 0xe4, 0x43, 0x2e, 0xd1, 0x7f, 0xcb,
	0xe8, 0xfb, 0x86, 0x8b, 0x9e, 0x9d, 0xac, 0xff, 0x1d, 0x59, 0xcb, 0x4e, 0xd1, 0x5f, 0x0d, 0xe3,
	0x55, 0x79, 0xf2, 0x23, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x30, 0x09, 0xcd, 0xbb, 0x02, 0x00, 0x00,
}
