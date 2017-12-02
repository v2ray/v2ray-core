package router

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_net "v2ray.com/core/common/net"
import v2ray_core_common_net1 "v2ray.com/core/common/net"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Type of domain value.
type Domain_Type int32

const (
	// The value is used as is.
	Domain_Plain Domain_Type = 0
	// The value is used as a regular expression.
	Domain_Regex Domain_Type = 1
	// The value is a domain.
	Domain_Domain Domain_Type = 2
)

var Domain_Type_name = map[int32]string{
	0: "Plain",
	1: "Regex",
	2: "Domain",
}
var Domain_Type_value = map[string]int32{
	"Plain":  0,
	"Regex":  1,
	"Domain": 2,
}

func (x Domain_Type) String() string {
	return proto.EnumName(Domain_Type_name, int32(x))
}
func (Domain_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Config_DomainStrategy int32

const (
	// Use domain as is.
	Config_AsIs Config_DomainStrategy = 0
	// Always resolve IP for domains.
	Config_UseIp Config_DomainStrategy = 1
	// Resolve to IP if the domain doesn't match any rules.
	Config_IpIfNonMatch Config_DomainStrategy = 2
	// Resolve to IP if any rule requires IP matching.
	Config_IpOnDemand Config_DomainStrategy = 3
)

var Config_DomainStrategy_name = map[int32]string{
	0: "AsIs",
	1: "UseIp",
	2: "IpIfNonMatch",
	3: "IpOnDemand",
}
var Config_DomainStrategy_value = map[string]int32{
	"AsIs":         0,
	"UseIp":        1,
	"IpIfNonMatch": 2,
	"IpOnDemand":   3,
}

func (x Config_DomainStrategy) String() string {
	return proto.EnumName(Config_DomainStrategy_name, int32(x))
}
func (Config_DomainStrategy) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{7, 0} }

// Domain for routing decision.
type Domain struct {
	// Domain matching type.
	Type Domain_Type `protobuf:"varint,1,opt,name=type,enum=v2ray.core.app.router.Domain_Type" json:"type,omitempty"`
	// Domain value.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Domain) Reset()                    { *m = Domain{} }
func (m *Domain) String() string            { return proto.CompactTextString(m) }
func (*Domain) ProtoMessage()               {}
func (*Domain) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Domain) GetType() Domain_Type {
	if m != nil {
		return m.Type
	}
	return Domain_Plain
}

func (m *Domain) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// IP for routing decision, in CIDR form.
type CIDR struct {
	// IP address, should be either 4 or 16 bytes.
	Ip []byte `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	// Number of leading ones in the network mask.
	Prefix uint32 `protobuf:"varint,2,opt,name=prefix" json:"prefix,omitempty"`
}

func (m *CIDR) Reset()                    { *m = CIDR{} }
func (m *CIDR) String() string            { return proto.CompactTextString(m) }
func (*CIDR) ProtoMessage()               {}
func (*CIDR) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CIDR) GetIp() []byte {
	if m != nil {
		return m.Ip
	}
	return nil
}

func (m *CIDR) GetPrefix() uint32 {
	if m != nil {
		return m.Prefix
	}
	return 0
}

type GeoIP struct {
	CountryCode string  `protobuf:"bytes,1,opt,name=country_code,json=countryCode" json:"country_code,omitempty"`
	Cidr        []*CIDR `protobuf:"bytes,2,rep,name=cidr" json:"cidr,omitempty"`
}

func (m *GeoIP) Reset()                    { *m = GeoIP{} }
func (m *GeoIP) String() string            { return proto.CompactTextString(m) }
func (*GeoIP) ProtoMessage()               {}
func (*GeoIP) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GeoIP) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

func (m *GeoIP) GetCidr() []*CIDR {
	if m != nil {
		return m.Cidr
	}
	return nil
}

type GeoIPList struct {
	Entry []*GeoIP `protobuf:"bytes,1,rep,name=entry" json:"entry,omitempty"`
}

func (m *GeoIPList) Reset()                    { *m = GeoIPList{} }
func (m *GeoIPList) String() string            { return proto.CompactTextString(m) }
func (*GeoIPList) ProtoMessage()               {}
func (*GeoIPList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GeoIPList) GetEntry() []*GeoIP {
	if m != nil {
		return m.Entry
	}
	return nil
}

type GeoSite struct {
	CountryCode string    `protobuf:"bytes,1,opt,name=country_code,json=countryCode" json:"country_code,omitempty"`
	Domain      []*Domain `protobuf:"bytes,2,rep,name=domain" json:"domain,omitempty"`
}

func (m *GeoSite) Reset()                    { *m = GeoSite{} }
func (m *GeoSite) String() string            { return proto.CompactTextString(m) }
func (*GeoSite) ProtoMessage()               {}
func (*GeoSite) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GeoSite) GetCountryCode() string {
	if m != nil {
		return m.CountryCode
	}
	return ""
}

func (m *GeoSite) GetDomain() []*Domain {
	if m != nil {
		return m.Domain
	}
	return nil
}

type GeoSiteList struct {
	Entry []*GeoSite `protobuf:"bytes,1,rep,name=entry" json:"entry,omitempty"`
}

func (m *GeoSiteList) Reset()                    { *m = GeoSiteList{} }
func (m *GeoSiteList) String() string            { return proto.CompactTextString(m) }
func (*GeoSiteList) ProtoMessage()               {}
func (*GeoSiteList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GeoSiteList) GetEntry() []*GeoSite {
	if m != nil {
		return m.Entry
	}
	return nil
}

type RoutingRule struct {
	Tag         string                              `protobuf:"bytes,1,opt,name=tag" json:"tag,omitempty"`
	Domain      []*Domain                           `protobuf:"bytes,2,rep,name=domain" json:"domain,omitempty"`
	Cidr        []*CIDR                             `protobuf:"bytes,3,rep,name=cidr" json:"cidr,omitempty"`
	PortRange   *v2ray_core_common_net.PortRange    `protobuf:"bytes,4,opt,name=port_range,json=portRange" json:"port_range,omitempty"`
	NetworkList *v2ray_core_common_net1.NetworkList `protobuf:"bytes,5,opt,name=network_list,json=networkList" json:"network_list,omitempty"`
	SourceCidr  []*CIDR                             `protobuf:"bytes,6,rep,name=source_cidr,json=sourceCidr" json:"source_cidr,omitempty"`
	UserEmail   []string                            `protobuf:"bytes,7,rep,name=user_email,json=userEmail" json:"user_email,omitempty"`
	InboundTag  []string                            `protobuf:"bytes,8,rep,name=inbound_tag,json=inboundTag" json:"inbound_tag,omitempty"`
}

func (m *RoutingRule) Reset()                    { *m = RoutingRule{} }
func (m *RoutingRule) String() string            { return proto.CompactTextString(m) }
func (*RoutingRule) ProtoMessage()               {}
func (*RoutingRule) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *RoutingRule) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *RoutingRule) GetDomain() []*Domain {
	if m != nil {
		return m.Domain
	}
	return nil
}

func (m *RoutingRule) GetCidr() []*CIDR {
	if m != nil {
		return m.Cidr
	}
	return nil
}

func (m *RoutingRule) GetPortRange() *v2ray_core_common_net.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *RoutingRule) GetNetworkList() *v2ray_core_common_net1.NetworkList {
	if m != nil {
		return m.NetworkList
	}
	return nil
}

func (m *RoutingRule) GetSourceCidr() []*CIDR {
	if m != nil {
		return m.SourceCidr
	}
	return nil
}

func (m *RoutingRule) GetUserEmail() []string {
	if m != nil {
		return m.UserEmail
	}
	return nil
}

func (m *RoutingRule) GetInboundTag() []string {
	if m != nil {
		return m.InboundTag
	}
	return nil
}

type Config struct {
	DomainStrategy Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,enum=v2ray.core.app.router.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	Rule           []*RoutingRule        `protobuf:"bytes,2,rep,name=rule" json:"rule,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Config) GetDomainStrategy() Config_DomainStrategy {
	if m != nil {
		return m.DomainStrategy
	}
	return Config_AsIs
}

func (m *Config) GetRule() []*RoutingRule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func init() {
	proto.RegisterType((*Domain)(nil), "v2ray.core.app.router.Domain")
	proto.RegisterType((*CIDR)(nil), "v2ray.core.app.router.CIDR")
	proto.RegisterType((*GeoIP)(nil), "v2ray.core.app.router.GeoIP")
	proto.RegisterType((*GeoIPList)(nil), "v2ray.core.app.router.GeoIPList")
	proto.RegisterType((*GeoSite)(nil), "v2ray.core.app.router.GeoSite")
	proto.RegisterType((*GeoSiteList)(nil), "v2ray.core.app.router.GeoSiteList")
	proto.RegisterType((*RoutingRule)(nil), "v2ray.core.app.router.RoutingRule")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.router.Config")
	proto.RegisterEnum("v2ray.core.app.router.Domain_Type", Domain_Type_name, Domain_Type_value)
	proto.RegisterEnum("v2ray.core.app.router.Config_DomainStrategy", Config_DomainStrategy_name, Config_DomainStrategy_value)
}

func init() { proto.RegisterFile("v2ray.com/core/app/router/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 640 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xcd, 0x6e, 0xd4, 0x3a,
	0x14, 0xc7, 0x6f, 0xe6, 0xab, 0x9d, 0x93, 0xb9, 0x73, 0x23, 0xeb, 0x16, 0x0d, 0x85, 0xc2, 0x10,
	0x21, 0x98, 0x05, 0x4a, 0xa4, 0xe1, 0x63, 0x05, 0xaa, 0xca, 0xb4, 0xaa, 0x22, 0x41, 0x19, 0xb9,
	0x2d, 0x0b, 0x58, 0x44, 0x6e, 0xe2, 0x86, 0x88, 0x89, 0x6d, 0x39, 0x4e, 0xe9, 0xec, 0x78, 0x01,
	0x5e, 0x84, 0xa7, 0xe2, 0x51, 0x90, 0xed, 0x0c, 0xb4, 0xa8, 0x81, 0x8a, 0x9d, 0xed, 0xfc, 0xfe,
	0xe7, 0xfc, 0x73, 0x7c, 0x8e, 0xe1, 0xc1, 0xd9, 0x54, 0x92, 0x65, 0x90, 0xf0, 0x22, 0x4c, 0xb8,
	0xa4, 0x21, 0x11, 0x22, 0x94, 0xbc, 0x52, 0x54, 0x86, 0x09, 0x67, 0xa7, 0x79, 0x16, 0x08, 0xc9,
	0x15, 0x47, 0x1b, 0x2b, 0x4e, 0xd2, 0x80, 0x08, 0x11, 0x58, 0x66, 0xf3, 0xfe, 0x2f, 0xf2, 0x84,
	0x17, 0x05, 0x67, 0x21, 0xa3, 0x2a, 0x14, 0x5c, 0x2a, 0x2b, 0xde, 0x7c, 0xd8, 0x4c, 0x31, 0xaa,
	0x3e, 0x71, 0xf9, 0xd1, 0x82, 0xfe, 0x67, 0x07, 0x7a, 0xbb, 0xbc, 0x20, 0x39, 0x43, 0xcf, 0xa0,
	0xa3, 0x96, 0x82, 0x8e, 0x9c, 0xb1, 0x33, 0x19, 0x4e, 0xfd, 0xe0, 0xca, 0xfc, 0x81, 0x85, 0x83,
	0xa3, 0xa5, 0xa0, 0xd8, 0xf0, 0xe8, 0x7f, 0xe8, 0x9e, 0x91, 0x45, 0x45, 0x47, 0xad, 0xb1, 0x33,
	0xe9, 0x63, 0xbb, 0xf1, 0x27, 0xd0, 0xd1, 0x0c, 0xea, 0x43, 0x77, 0xbe, 0x20, 0x39, 0xf3, 0xfe,
	0xd1, 0x4b, 0x4c, 0x33, 0x7a, 0xee, 0x39, 0x08, 0x56, 0x59, 0xbd, 0x96, 0x1f, 0x40, 0x67, 0x16,
	0xed, 0x62, 0x34, 0x84, 0x56, 0x2e, 0x4c, 0xf6, 0x01, 0x6e, 0xe5, 0x02, 0xdd, 0x80, 0x9e, 0x90,
	0xf4, 0x34, 0x3f, 0x37, 0x81, 0xff, 0xc5, 0xf5, 0xce, 0x7f, 0x0f, 0xdd, 0x7d, 0xca, 0xa3, 0x39,
	0xba, 0x07, 0x83, 0x84, 0x57, 0x4c, 0xc9, 0x65, 0x9c, 0xf0, 0xd4, 0x1a, 0xef, 0x63, 0xb7, 0x3e,
	0x9b, 0xf1, 0x94, 0xa2, 0x10, 0x3a, 0x49, 0x9e, 0xca, 0x51, 0x6b, 0xdc, 0x9e, 0xb8, 0xd3, 0x5b,
	0x0d, 0xff, 0xa4, 0xd3, 0x63, 0x03, 0xfa, 0xdb, 0xd0, 0x37, 0xc1, 0x5f, 0xe5, 0xa5, 0x42, 0x53,
	0xe8, 0x52, 0x1d, 0x6a, 0xe4, 0x18, 0xf9, 0xed, 0x06, 0xb9, 0x11, 0x60, 0x8b, 0xfa, 0x09, 0xac,
	0xed, 0x53, 0x7e, 0x98, 0x2b, 0x7a, 0x1d, 0x7f, 0x4f, 0xa1, 0x97, 0x9a, 0x3a, 0xd4, 0x0e, 0xb7,
	0x7e, 0x5b, 0x75, 0x5c, 0xc3, 0xfe, 0x0c, 0xdc, 0x3a, 0x89, 0xf1, 0xf9, 0xe4, 0xb2, 0xcf, 0x3b,
	0xcd, 0x3e, 0xb5, 0x64, 0xe5, 0xf4, 0x4b, 0x1b, 0x5c, 0xcc, 0x2b, 0x95, 0xb3, 0x0c, 0x57, 0x0b,
	0x8a, 0x3c, 0x68, 0x2b, 0x92, 0xd5, 0x2e, 0xf5, 0xf2, 0x2f, 0xdd, 0xfd, 0x28, 0x7a, 0xfb, 0x9a,
	0x45, 0x47, 0xdb, 0x00, 0xba, 0x77, 0x63, 0x49, 0x58, 0x46, 0x47, 0x9d, 0xb1, 0x33, 0x71, 0xa7,
	0xe3, 0x8b, 0x32, 0xdb, 0xbe, 0x01, 0xa3, 0x2a, 0x98, 0x73, 0xa9, 0xb0, 0xe6, 0x70, 0x5f, 0xac,
	0x96, 0x68, 0x0f, 0x06, 0x75, 0x5b, 0xc7, 0x8b, 0xbc, 0x54, 0xa3, 0xae, 0x09, 0xe1, 0x37, 0x84,
	0x38, 0xb0, 0xa8, 0x2e, 0x1d, 0x76, 0xd9, 0xcf, 0x0d, 0x7a, 0x0e, 0x6e, 0xc9, 0x2b, 0x99, 0xd0,
	0xd8, 0xf8, 0xef, 0xfd, 0xd9, 0x3f, 0x58, 0x7e, 0xa6, 0xff, 0x62, 0x0b, 0xa0, 0x2a, 0xa9, 0x8c,
	0x69, 0x41, 0xf2, 0xc5, 0x68, 0x6d, 0xdc, 0x9e, 0xf4, 0x71, 0x5f, 0x9f, 0xec, 0xe9, 0x03, 0x74,
	0x17, 0xdc, 0x9c, 0x9d, 0xf0, 0x8a, 0xa5, 0xb1, 0x2e, 0xf3, 0xba, 0xf9, 0x0e, 0xf5, 0xd1, 0x11,
	0xc9, 0xfc, 0x6f, 0x0e, 0xf4, 0x66, 0xe6, 0x05, 0x40, 0xc7, 0xf0, 0x9f, 0xad, 0x65, 0x5c, 0x2a,
	0x49, 0x14, 0xcd, 0x96, 0xf5, 0x54, 0x3e, 0x6a, 0x32, 0x63, 0x5f, 0x0e, 0x7b, 0x11, 0x87, 0xb5,
	0x06, 0x0f, 0xd3, 0x4b, 0x7b, 0x3d, 0xe1, 0xb2, 0x5a, 0xd0, 0xfa, 0x36, 0x9b, 0x26, 0xfc, 0x42,
	0x4f, 0x60, 0xc3, 0xfb, 0xfb, 0x30, 0xbc, 0x1c, 0x19, 0xad, 0x43, 0x67, 0xa7, 0x8c, 0x4a, 0x3b,
	0xd4, 0xc7, 0x25, 0x8d, 0x84, 0xe7, 0x20, 0x0f, 0x06, 0x91, 0x88, 0x4e, 0x0f, 0x38, 0x7b, 0x4d,
	0x54, 0xf2, 0xc1, 0x6b, 0xa1, 0x21, 0x40, 0x24, 0xde, 0xb0, 0x5d, 0x5a, 0x10, 0x96, 0x7a, 0xed,
	0x97, 0x2f, 0xe0, 0x66, 0xc2, 0x8b, 0xab, 0xf3, 0xce, 0x9d, 0x77, 0x3d, 0xbb, 0xfa, 0xda, 0xda,
	0x78, 0x3b, 0xc5, 0x64, 0x19, 0xcc, 0x34, 0xb1, 0x23, 0x84, 0xb1, 0x44, 0xe5, 0x49, 0xcf, 0xbc,
	0x59, 0x8f, 0xbf, 0x07, 0x00, 0x00, 0xff, 0xff, 0x53, 0x7c, 0xa8, 0x94, 0x43, 0x05, 0x00, 0x00,
}
