package router

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	net "v2ray.com/core/common/net"
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

// Type of domain value.
type Domain_Type int32

const (
	// The value is used as is.
	Domain_Plain Domain_Type = 0
	// The value is used as a regular expression.
	Domain_Regex Domain_Type = 1
	// The value is a root domain.
	Domain_Domain Domain_Type = 2
	// The value is a domain.
	Domain_Full Domain_Type = 3
)

var Domain_Type_name = map[int32]string{
	0: "Plain",
	1: "Regex",
	2: "Domain",
	3: "Full",
}

var Domain_Type_value = map[string]int32{
	"Plain":  0,
	"Regex":  1,
	"Domain": 2,
	"Full":   3,
}

func (x Domain_Type) String() string {
	return proto.EnumName(Domain_Type_name, int32(x))
}

func (Domain_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{0, 0}
}

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

func (Config_DomainStrategy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{7, 0}
}

// Domain for routing decision.
type Domain struct {
	// Domain matching type.
	Type Domain_Type `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.router.Domain_Type" json:"type,omitempty"`
	// Domain value.
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Domain) Reset()         { *m = Domain{} }
func (m *Domain) String() string { return proto.CompactTextString(m) }
func (*Domain) ProtoMessage()    {}
func (*Domain) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{0}
}

func (m *Domain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Domain.Unmarshal(m, b)
}
func (m *Domain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Domain.Marshal(b, m, deterministic)
}
func (m *Domain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Domain.Merge(m, src)
}
func (m *Domain) XXX_Size() int {
	return xxx_messageInfo_Domain.Size(m)
}
func (m *Domain) XXX_DiscardUnknown() {
	xxx_messageInfo_Domain.DiscardUnknown(m)
}

var xxx_messageInfo_Domain proto.InternalMessageInfo

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
	Prefix               uint32   `protobuf:"varint,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CIDR) Reset()         { *m = CIDR{} }
func (m *CIDR) String() string { return proto.CompactTextString(m) }
func (*CIDR) ProtoMessage()    {}
func (*CIDR) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{1}
}

func (m *CIDR) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CIDR.Unmarshal(m, b)
}
func (m *CIDR) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CIDR.Marshal(b, m, deterministic)
}
func (m *CIDR) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CIDR.Merge(m, src)
}
func (m *CIDR) XXX_Size() int {
	return xxx_messageInfo_CIDR.Size(m)
}
func (m *CIDR) XXX_DiscardUnknown() {
	xxx_messageInfo_CIDR.DiscardUnknown(m)
}

var xxx_messageInfo_CIDR proto.InternalMessageInfo

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
	CountryCode          string   `protobuf:"bytes,1,opt,name=country_code,json=countryCode,proto3" json:"country_code,omitempty"`
	Cidr                 []*CIDR  `protobuf:"bytes,2,rep,name=cidr,proto3" json:"cidr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeoIP) Reset()         { *m = GeoIP{} }
func (m *GeoIP) String() string { return proto.CompactTextString(m) }
func (*GeoIP) ProtoMessage()    {}
func (*GeoIP) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{2}
}

func (m *GeoIP) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoIP.Unmarshal(m, b)
}
func (m *GeoIP) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoIP.Marshal(b, m, deterministic)
}
func (m *GeoIP) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoIP.Merge(m, src)
}
func (m *GeoIP) XXX_Size() int {
	return xxx_messageInfo_GeoIP.Size(m)
}
func (m *GeoIP) XXX_DiscardUnknown() {
	xxx_messageInfo_GeoIP.DiscardUnknown(m)
}

var xxx_messageInfo_GeoIP proto.InternalMessageInfo

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
	Entry                []*GeoIP `protobuf:"bytes,1,rep,name=entry,proto3" json:"entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeoIPList) Reset()         { *m = GeoIPList{} }
func (m *GeoIPList) String() string { return proto.CompactTextString(m) }
func (*GeoIPList) ProtoMessage()    {}
func (*GeoIPList) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{3}
}

func (m *GeoIPList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoIPList.Unmarshal(m, b)
}
func (m *GeoIPList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoIPList.Marshal(b, m, deterministic)
}
func (m *GeoIPList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoIPList.Merge(m, src)
}
func (m *GeoIPList) XXX_Size() int {
	return xxx_messageInfo_GeoIPList.Size(m)
}
func (m *GeoIPList) XXX_DiscardUnknown() {
	xxx_messageInfo_GeoIPList.DiscardUnknown(m)
}

var xxx_messageInfo_GeoIPList proto.InternalMessageInfo

func (m *GeoIPList) GetEntry() []*GeoIP {
	if m != nil {
		return m.Entry
	}
	return nil
}

type GeoSite struct {
	CountryCode          string    `protobuf:"bytes,1,opt,name=country_code,json=countryCode,proto3" json:"country_code,omitempty"`
	Domain               []*Domain `protobuf:"bytes,2,rep,name=domain,proto3" json:"domain,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GeoSite) Reset()         { *m = GeoSite{} }
func (m *GeoSite) String() string { return proto.CompactTextString(m) }
func (*GeoSite) ProtoMessage()    {}
func (*GeoSite) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{4}
}

func (m *GeoSite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoSite.Unmarshal(m, b)
}
func (m *GeoSite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoSite.Marshal(b, m, deterministic)
}
func (m *GeoSite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoSite.Merge(m, src)
}
func (m *GeoSite) XXX_Size() int {
	return xxx_messageInfo_GeoSite.Size(m)
}
func (m *GeoSite) XXX_DiscardUnknown() {
	xxx_messageInfo_GeoSite.DiscardUnknown(m)
}

var xxx_messageInfo_GeoSite proto.InternalMessageInfo

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
	Entry                []*GeoSite `protobuf:"bytes,1,rep,name=entry,proto3" json:"entry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *GeoSiteList) Reset()         { *m = GeoSiteList{} }
func (m *GeoSiteList) String() string { return proto.CompactTextString(m) }
func (*GeoSiteList) ProtoMessage()    {}
func (*GeoSiteList) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{5}
}

func (m *GeoSiteList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoSiteList.Unmarshal(m, b)
}
func (m *GeoSiteList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoSiteList.Marshal(b, m, deterministic)
}
func (m *GeoSiteList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoSiteList.Merge(m, src)
}
func (m *GeoSiteList) XXX_Size() int {
	return xxx_messageInfo_GeoSiteList.Size(m)
}
func (m *GeoSiteList) XXX_DiscardUnknown() {
	xxx_messageInfo_GeoSiteList.DiscardUnknown(m)
}

var xxx_messageInfo_GeoSiteList proto.InternalMessageInfo

func (m *GeoSiteList) GetEntry() []*GeoSite {
	if m != nil {
		return m.Entry
	}
	return nil
}

type RoutingRule struct {
	// Tag of outbound that this rule is pointing to.
	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	// List of domains for target domain matching.
	Domain []*Domain `protobuf:"bytes,2,rep,name=domain,proto3" json:"domain,omitempty"`
	// List of CIDRs for target IP address matching.
	// The list must be sorted beforehand.
	Cidr []*CIDR `protobuf:"bytes,3,rep,name=cidr,proto3" json:"cidr,omitempty"` // Deprecated: Do not use.
	// List of GeoIPs for target IP address matching. If this entry exists, the cidr above will have no effect.
	// GeoIP fields with the same country code are supposed to contain exactly same content. They will be merged during runtime.
	// For customized GeoIPs, please leave country code empty.
	// The CIDR list in the GeoIP must be sorted beforehand.
	Geoip       []*GeoIP         `protobuf:"bytes,10,rep,name=geoip,proto3" json:"geoip,omitempty"`
	PortRange   *net.PortRange   `protobuf:"bytes,4,opt,name=port_range,json=portRange,proto3" json:"port_range,omitempty"`
	NetworkList *net.NetworkList `protobuf:"bytes,5,opt,name=network_list,json=networkList,proto3" json:"network_list,omitempty"`
	// List of CIDRs for source IP address matching.
	// The list must be sorted beforehand.
	SourceCidr []*CIDR `protobuf:"bytes,6,rep,name=source_cidr,json=sourceCidr,proto3" json:"source_cidr,omitempty"` // Deprecated: Do not use.
	// List of GeoIPs for source IP address matching. If this entry exists, the source_cidr above will have no effect.
	// The CIDR list in the GeoIP must be sorted beforehand.
	SourceGeoip          []*GeoIP `protobuf:"bytes,11,rep,name=source_geoip,json=sourceGeoip,proto3" json:"source_geoip,omitempty"`
	UserEmail            []string `protobuf:"bytes,7,rep,name=user_email,json=userEmail,proto3" json:"user_email,omitempty"`
	InboundTag           []string `protobuf:"bytes,8,rep,name=inbound_tag,json=inboundTag,proto3" json:"inbound_tag,omitempty"`
	Protocol             []string `protobuf:"bytes,9,rep,name=protocol,proto3" json:"protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoutingRule) Reset()         { *m = RoutingRule{} }
func (m *RoutingRule) String() string { return proto.CompactTextString(m) }
func (*RoutingRule) ProtoMessage()    {}
func (*RoutingRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{6}
}

func (m *RoutingRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoutingRule.Unmarshal(m, b)
}
func (m *RoutingRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoutingRule.Marshal(b, m, deterministic)
}
func (m *RoutingRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoutingRule.Merge(m, src)
}
func (m *RoutingRule) XXX_Size() int {
	return xxx_messageInfo_RoutingRule.Size(m)
}
func (m *RoutingRule) XXX_DiscardUnknown() {
	xxx_messageInfo_RoutingRule.DiscardUnknown(m)
}

var xxx_messageInfo_RoutingRule proto.InternalMessageInfo

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

// Deprecated: Do not use.
func (m *RoutingRule) GetCidr() []*CIDR {
	if m != nil {
		return m.Cidr
	}
	return nil
}

func (m *RoutingRule) GetGeoip() []*GeoIP {
	if m != nil {
		return m.Geoip
	}
	return nil
}

func (m *RoutingRule) GetPortRange() *net.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *RoutingRule) GetNetworkList() *net.NetworkList {
	if m != nil {
		return m.NetworkList
	}
	return nil
}

// Deprecated: Do not use.
func (m *RoutingRule) GetSourceCidr() []*CIDR {
	if m != nil {
		return m.SourceCidr
	}
	return nil
}

func (m *RoutingRule) GetSourceGeoip() []*GeoIP {
	if m != nil {
		return m.SourceGeoip
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

func (m *RoutingRule) GetProtocol() []string {
	if m != nil {
		return m.Protocol
	}
	return nil
}

type Config struct {
	DomainStrategy       Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,proto3,enum=v2ray.core.app.router.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	Rule                 []*RoutingRule        `protobuf:"bytes,2,rep,name=rule,proto3" json:"rule,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{7}
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
	return Config_AsIs
}

func (m *Config) GetRule() []*RoutingRule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func init() {
	proto.RegisterEnum("v2ray.core.app.router.Domain_Type", Domain_Type_name, Domain_Type_value)
	proto.RegisterEnum("v2ray.core.app.router.Config_DomainStrategy", Config_DomainStrategy_name, Config_DomainStrategy_value)
	proto.RegisterType((*Domain)(nil), "v2ray.core.app.router.Domain")
	proto.RegisterType((*CIDR)(nil), "v2ray.core.app.router.CIDR")
	proto.RegisterType((*GeoIP)(nil), "v2ray.core.app.router.GeoIP")
	proto.RegisterType((*GeoIPList)(nil), "v2ray.core.app.router.GeoIPList")
	proto.RegisterType((*GeoSite)(nil), "v2ray.core.app.router.GeoSite")
	proto.RegisterType((*GeoSiteList)(nil), "v2ray.core.app.router.GeoSiteList")
	proto.RegisterType((*RoutingRule)(nil), "v2ray.core.app.router.RoutingRule")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.router.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/router/config.proto", fileDescriptor_6b1608360690c5fc)
}

var fileDescriptor_6b1608360690c5fc = []byte{
	// 691 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0xfd, 0xec, 0xfc, 0xb4, 0x1e, 0xe7, 0x0b, 0xd6, 0x8a, 0x22, 0x53, 0x28, 0x04, 0x0b, 0x41,
	0x2e, 0x90, 0x23, 0xa5, 0xc0, 0x1d, 0x2a, 0x6d, 0x5a, 0xa2, 0x48, 0x50, 0xa2, 0x6d, 0xcb, 0x05,
	0x5c, 0x44, 0xae, 0xb3, 0x35, 0x16, 0xce, 0xee, 0x6a, 0xbd, 0x2e, 0xcd, 0x2b, 0xf0, 0x20, 0x5c,
	0xf0, 0x54, 0x3c, 0x0a, 0xda, 0x9f, 0x40, 0x8b, 0x1a, 0x88, 0xb8, 0xdb, 0x59, 0x9f, 0x33, 0x73,
	0xe6, 0x78, 0x76, 0xe0, 0xd1, 0x79, 0x5f, 0x24, 0xf3, 0x38, 0x65, 0xb3, 0x5e, 0xca, 0x04, 0xe9,
	0x25, 0x9c, 0xf7, 0x04, 0xab, 0x24, 0x11, 0xbd, 0x94, 0xd1, 0xb3, 0x3c, 0x8b, 0xb9, 0x60, 0x92,
	0xa1, 0x8d, 0x05, 0x4e, 0x90, 0x38, 0xe1, 0x3c, 0x36, 0x98, 0xcd, 0x87, 0xbf, 0xd1, 0x53, 0x36,
	0x9b, 0x31, 0xda, 0xa3, 0x44, 0xf6, 0x38, 0x13, 0xd2, 0x90, 0x37, 0x1f, 0x2f, 0x47, 0x51, 0x22,
	0x3f, 0x33, 0xf1, 0xc9, 0x00, 0xa3, 0x2f, 0x0e, 0x34, 0xf7, 0xd9, 0x2c, 0xc9, 0x29, 0x7a, 0x0e,
	0x75, 0x39, 0xe7, 0x24, 0x74, 0x3a, 0x4e, 0xb7, 0xdd, 0x8f, 0xe2, 0x6b, 0xeb, 0xc7, 0x06, 0x1c,
	0x1f, 0xcf, 0x39, 0xc1, 0x1a, 0x8f, 0x6e, 0x42, 0xe3, 0x3c, 0x29, 0x2a, 0x12, 0xba, 0x1d, 0xa7,
	0xeb, 0x61, 0x13, 0x44, 0x7d, 0xa8, 0x2b, 0x0c, 0xf2, 0xa0, 0x31, 0x2e, 0x92, 0x9c, 0x06, 0xff,
	0xa9, 0x23, 0x26, 0x19, 0xb9, 0x08, 0x1c, 0x04, 0x8b, 0xaa, 0x81, 0x8b, 0xd6, 0xa1, 0xfe, 0xaa,
	0x2a, 0x8a, 0xa0, 0x16, 0xc5, 0x50, 0x1f, 0x8c, 0xf6, 0x31, 0x6a, 0x83, 0x9b, 0x73, 0xad, 0xa3,
	0x85, 0xdd, 0x9c, 0xa3, 0x5b, 0xd0, 0xe4, 0x82, 0x9c, 0xe5, 0x17, 0xba, 0xc4, 0xff, 0xd8, 0x46,
	0xd1, 0x07, 0x68, 0x0c, 0x09, 0x1b, 0x8d, 0xd1, 0x03, 0x68, 0xa5, 0xac, 0xa2, 0x52, 0xcc, 0x27,
	0x29, 0x9b, 0x9a, 0x16, 0x3c, 0xec, 0xdb, 0xbb, 0x01, 0x9b, 0x12, 0xd4, 0x83, 0x7a, 0x9a, 0x4f,
	0x45, 0xe8, 0x76, 0x6a, 0x5d, 0xbf, 0x7f, 0x67, 0x49, 0x77, 0xaa, 0x3c, 0xd6, 0xc0, 0x68, 0x07,
	0x3c, 0x9d, 0xfc, 0x75, 0x5e, 0x4a, 0xd4, 0x87, 0x06, 0x51, 0xa9, 0x42, 0x47, 0xd3, 0xef, 0x2e,
	0xa1, 0x6b, 0x02, 0x36, 0xd0, 0x28, 0x85, 0xb5, 0x21, 0x61, 0x47, 0xb9, 0x24, 0xab, 0xe8, 0x7b,
	0x06, 0xcd, 0xa9, 0x76, 0xc4, 0x2a, 0xdc, 0xfa, 0xa3, 0xff, 0xd8, 0x82, 0xa3, 0x01, 0xf8, 0xb6,
	0x88, 0xd6, 0xf9, 0xf4, 0xaa, 0xce, 0x7b, 0xcb, 0x75, 0x2a, 0xca, 0x42, 0xe9, 0xd7, 0x3a, 0xf8,
	0x98, 0x55, 0x32, 0xa7, 0x19, 0xae, 0x0a, 0x82, 0x02, 0xa8, 0xc9, 0x24, 0xb3, 0x2a, 0xd5, 0xf1,
	0x1f, 0xd5, 0xa1, 0x6d, 0x6b, 0x7a, 0xed, 0xaf, 0xa6, 0xef, 0xb9, 0xa1, 0x63, 0x8c, 0x57, 0x5e,
	0x67, 0x84, 0xe5, 0x3c, 0x84, 0x55, 0xbc, 0xd6, 0x50, 0xb4, 0x03, 0xa0, 0xa6, 0x7f, 0x22, 0x12,
	0x9a, 0x91, 0xb0, 0xde, 0x71, 0xba, 0x7e, 0xbf, 0x73, 0x99, 0x68, 0x1e, 0x40, 0x4c, 0x89, 0x8c,
	0xc7, 0x4c, 0x48, 0xac, 0x70, 0xd8, 0xe3, 0x8b, 0x23, 0x3a, 0x80, 0x96, 0x7d, 0x18, 0x93, 0x22,
	0x2f, 0x65, 0xd8, 0xd0, 0x29, 0xa2, 0x25, 0x29, 0x0e, 0x0d, 0x54, 0x59, 0x8e, 0x7d, 0xfa, 0x2b,
	0x40, 0x2f, 0xc1, 0x2f, 0x59, 0x25, 0x52, 0x32, 0xd1, 0x7d, 0x37, 0x57, 0xeb, 0x1b, 0x0c, 0x67,
	0xa0, 0xba, 0xdf, 0x81, 0x96, 0xcd, 0x60, 0x4c, 0xf0, 0x57, 0x30, 0xc1, 0xd6, 0x1c, 0x6a, 0x2b,
	0xb6, 0x00, 0xaa, 0x92, 0x88, 0x09, 0x99, 0x25, 0x79, 0x11, 0xae, 0x75, 0x6a, 0x5d, 0x0f, 0x7b,
	0xea, 0xe6, 0x40, 0x5d, 0xa0, 0xfb, 0xe0, 0xe7, 0xf4, 0x94, 0x55, 0x74, 0x3a, 0x51, 0xff, 0x78,
	0x5d, 0x7f, 0x07, 0x7b, 0x75, 0x9c, 0x64, 0x68, 0x13, 0xd6, 0xf5, 0x6a, 0x48, 0x59, 0x11, 0x7a,
	0xfa, 0xeb, 0xcf, 0x38, 0xfa, 0xee, 0x40, 0x73, 0xa0, 0x97, 0x14, 0x3a, 0x81, 0x1b, 0xe6, 0x27,
	0x4f, 0x4a, 0x29, 0x12, 0x49, 0xb2, 0xb9, 0x5d, 0x1c, 0x4f, 0x96, 0x75, 0x6b, 0x96, 0x9b, 0x99,
	0x90, 0x23, 0xcb, 0xc1, 0xed, 0xe9, 0x95, 0x58, 0x2d, 0x21, 0x51, 0x15, 0xc4, 0x8e, 0xd9, 0xb2,
	0x25, 0x74, 0x69, 0x58, 0xb1, 0xc6, 0x47, 0x43, 0x68, 0x5f, 0xcd, 0xac, 0xd6, 0xca, 0x6e, 0x39,
	0x2a, 0xcd, 0xde, 0x39, 0x29, 0xc9, 0x88, 0x07, 0x0e, 0x0a, 0xa0, 0x35, 0xe2, 0xa3, 0xb3, 0x43,
	0x46, 0xdf, 0x24, 0x32, 0xfd, 0x18, 0xb8, 0xa8, 0x0d, 0x30, 0xe2, 0x6f, 0xe9, 0x3e, 0x99, 0x25,
	0x74, 0x1a, 0xd4, 0xf6, 0x5e, 0xc0, 0xed, 0x94, 0xcd, 0xae, 0xaf, 0x3b, 0x76, 0xde, 0x37, 0xcd,
	0xe9, 0x9b, 0xbb, 0xf1, 0xae, 0x8f, 0x93, 0x79, 0x3c, 0x50, 0x88, 0x5d, 0xce, 0xb5, 0x24, 0x22,
	0x4e, 0x9b, 0xda, 0xab, 0xed, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x42, 0x5b, 0xf4, 0x78, 0xe6,
	0x05, 0x00, 0x00,
}
