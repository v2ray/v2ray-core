package router

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
func (Domain_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_227cf1ddacaf1282, []int{0, 0}
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{7, 0}
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{0}
}
func (m *Domain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Domain.Unmarshal(m, b)
}
func (m *Domain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Domain.Marshal(b, m, deterministic)
}
func (dst *Domain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Domain.Merge(dst, src)
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{1}
}
func (m *CIDR) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CIDR.Unmarshal(m, b)
}
func (m *CIDR) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CIDR.Marshal(b, m, deterministic)
}
func (dst *CIDR) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CIDR.Merge(dst, src)
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{2}
}
func (m *GeoIP) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoIP.Unmarshal(m, b)
}
func (m *GeoIP) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoIP.Marshal(b, m, deterministic)
}
func (dst *GeoIP) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoIP.Merge(dst, src)
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{3}
}
func (m *GeoIPList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoIPList.Unmarshal(m, b)
}
func (m *GeoIPList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoIPList.Marshal(b, m, deterministic)
}
func (dst *GeoIPList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoIPList.Merge(dst, src)
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{4}
}
func (m *GeoSite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoSite.Unmarshal(m, b)
}
func (m *GeoSite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoSite.Marshal(b, m, deterministic)
}
func (dst *GeoSite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoSite.Merge(dst, src)
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{5}
}
func (m *GeoSiteList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoSiteList.Unmarshal(m, b)
}
func (m *GeoSiteList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoSiteList.Marshal(b, m, deterministic)
}
func (dst *GeoSiteList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoSiteList.Merge(dst, src)
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
	Tag                  string           `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Domain               []*Domain        `protobuf:"bytes,2,rep,name=domain,proto3" json:"domain,omitempty"`
	Cidr                 []*CIDR          `protobuf:"bytes,3,rep,name=cidr,proto3" json:"cidr,omitempty"`
	PortRange            *net.PortRange   `protobuf:"bytes,4,opt,name=port_range,json=portRange,proto3" json:"port_range,omitempty"`
	NetworkList          *net.NetworkList `protobuf:"bytes,5,opt,name=network_list,json=networkList,proto3" json:"network_list,omitempty"`
	SourceCidr           []*CIDR          `protobuf:"bytes,6,rep,name=source_cidr,json=sourceCidr,proto3" json:"source_cidr,omitempty"`
	UserEmail            []string         `protobuf:"bytes,7,rep,name=user_email,json=userEmail,proto3" json:"user_email,omitempty"`
	InboundTag           []string         `protobuf:"bytes,8,rep,name=inbound_tag,json=inboundTag,proto3" json:"inbound_tag,omitempty"`
	Protocol             []string         `protobuf:"bytes,9,rep,name=protocol,proto3" json:"protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *RoutingRule) Reset()         { *m = RoutingRule{} }
func (m *RoutingRule) String() string { return proto.CompactTextString(m) }
func (*RoutingRule) ProtoMessage()    {}
func (*RoutingRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_227cf1ddacaf1282, []int{6}
}
func (m *RoutingRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoutingRule.Unmarshal(m, b)
}
func (m *RoutingRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoutingRule.Marshal(b, m, deterministic)
}
func (dst *RoutingRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoutingRule.Merge(dst, src)
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

func (m *RoutingRule) GetCidr() []*CIDR {
	if m != nil {
		return m.Cidr
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
	return fileDescriptor_config_227cf1ddacaf1282, []int{7}
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

func init() {
	proto.RegisterFile("v2ray.com/core/app/router/config.proto", fileDescriptor_config_227cf1ddacaf1282)
}

var fileDescriptor_config_227cf1ddacaf1282 = []byte{
	// 651 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x5d, 0x6f, 0xd3, 0x4a,
	0x10, 0xbd, 0xce, 0x57, 0xeb, 0x71, 0x6e, 0xae, 0xb5, 0xba, 0xbd, 0xf2, 0x2d, 0x14, 0x82, 0x85,
	0x20, 0x0f, 0xc8, 0x91, 0xc2, 0xc7, 0x13, 0xa8, 0x2a, 0x69, 0x55, 0x45, 0x82, 0x12, 0x6d, 0x5b,
	0x1e, 0xe0, 0x21, 0xda, 0xda, 0x5b, 0x63, 0x61, 0xef, 0xae, 0xd6, 0xeb, 0xd2, 0xbc, 0xf1, 0x77,
	0xe0, 0x57, 0xf1, 0x53, 0xd0, 0xee, 0x3a, 0xa5, 0x45, 0x35, 0x54, 0xbc, 0xcd, 0x8c, 0xcf, 0x99,
	0x39, 0x3e, 0x9e, 0x31, 0x3c, 0x38, 0x9b, 0x48, 0xb2, 0x8c, 0x62, 0x5e, 0x8c, 0x63, 0x2e, 0xe9,
	0x98, 0x08, 0x31, 0x96, 0xbc, 0x52, 0x54, 0x8e, 0x63, 0xce, 0x4e, 0xb3, 0x34, 0x12, 0x92, 0x2b,
	0x8e, 0x36, 0x56, 0x38, 0x49, 0x23, 0x22, 0x44, 0x64, 0x31, 0x9b, 0xf7, 0x7f, 0xa2, 0xc7, 0xbc,
	0x28, 0x38, 0x1b, 0x33, 0xaa, 0xc6, 0x82, 0x4b, 0x65, 0xc9, 0x9b, 0x0f, 0x9b, 0x51, 0x8c, 0xaa,
	0x4f, 0x5c, 0x7e, 0xb4, 0xc0, 0xf0, 0xb3, 0x03, 0xbd, 0x5d, 0x5e, 0x90, 0x8c, 0xa1, 0x67, 0xd0,
	0x51, 0x4b, 0x41, 0x03, 0x67, 0xe8, 0x8c, 0x06, 0x93, 0x30, 0xba, 0x76, 0x7e, 0x64, 0xc1, 0xd1,
	0xd1, 0x52, 0x50, 0x6c, 0xf0, 0xe8, 0x5f, 0xe8, 0x9e, 0x91, 0xbc, 0xa2, 0x41, 0x6b, 0xe8, 0x8c,
	0x5c, 0x6c, 0x93, 0x70, 0x04, 0x1d, 0x8d, 0x41, 0x2e, 0x74, 0xe7, 0x39, 0xc9, 0x98, 0xff, 0x97,
	0x0e, 0x31, 0x4d, 0xe9, 0xb9, 0xef, 0x20, 0x58, 0x4d, 0xf5, 0x5b, 0x61, 0x04, 0x9d, 0xe9, 0x6c,
	0x17, 0xa3, 0x01, 0xb4, 0x32, 0x61, 0xa6, 0xf7, 0x71, 0x2b, 0x13, 0xe8, 0x3f, 0xe8, 0x09, 0x49,
	0x4f, 0xb3, 0x73, 0xd3, 0xf8, 0x6f, 0x5c, 0x67, 0xe1, 0x7b, 0xe8, 0xee, 0x53, 0x3e, 0x9b, 0xa3,
	0x7b, 0xd0, 0x8f, 0x79, 0xc5, 0x94, 0x5c, 0x2e, 0x62, 0x9e, 0x58, 0xe1, 0x2e, 0xf6, 0xea, 0xda,
	0x94, 0x27, 0x14, 0x8d, 0xa1, 0x13, 0x67, 0x89, 0x0c, 0x5a, 0xc3, 0xf6, 0xc8, 0x9b, 0xdc, 0x6a,
	0x78, 0x27, 0x3d, 0x1e, 0x1b, 0x60, 0xb8, 0x0d, 0xae, 0x69, 0xfe, 0x2a, 0x2b, 0x15, 0x9a, 0x40,
	0x97, 0xea, 0x56, 0x81, 0x63, 0xe8, 0xb7, 0x1b, 0xe8, 0x86, 0x80, 0x2d, 0x34, 0x8c, 0x61, 0x6d,
	0x9f, 0xf2, 0xc3, 0x4c, 0xd1, 0x9b, 0xe8, 0x7b, 0x0a, 0xbd, 0xc4, 0xf8, 0x50, 0x2b, 0xdc, 0xfa,
	0xa5, 0xeb, 0xb8, 0x06, 0x87, 0x53, 0xf0, 0xea, 0x21, 0x46, 0xe7, 0x93, 0xab, 0x3a, 0xef, 0x34,
	0xeb, 0xd4, 0x94, 0x95, 0xd2, 0x2f, 0x6d, 0xf0, 0x30, 0xaf, 0x54, 0xc6, 0x52, 0x5c, 0xe5, 0x14,
	0xf9, 0xd0, 0x56, 0x24, 0xad, 0x55, 0xea, 0xf0, 0x0f, 0xd5, 0x5d, 0x98, 0xde, 0xbe, 0xa1, 0xe9,
	0x68, 0x1b, 0x40, 0xef, 0xee, 0x42, 0x12, 0x96, 0xd2, 0xa0, 0x33, 0x74, 0x46, 0xde, 0x64, 0x78,
	0x99, 0x66, 0xd7, 0x37, 0x62, 0x54, 0x45, 0x73, 0x2e, 0x15, 0xd6, 0x38, 0xec, 0x8a, 0x55, 0x88,
	0xf6, 0xa0, 0x5f, 0xaf, 0xf5, 0x22, 0xcf, 0x4a, 0x15, 0x74, 0x4d, 0x8b, 0xb0, 0xa1, 0xc5, 0x81,
	0x85, 0x6a, 0xeb, 0xb0, 0xc7, 0x7e, 0x24, 0xe8, 0x39, 0x78, 0x25, 0xaf, 0x64, 0x4c, 0x17, 0x46,
	0x7f, 0xef, 0xf7, 0xfa, 0xc1, 0xe2, 0xa7, 0xfa, 0x2d, 0xb6, 0x00, 0xaa, 0x92, 0xca, 0x05, 0x2d,
	0x48, 0x96, 0x07, 0x6b, 0xc3, 0xf6, 0xc8, 0xc5, 0xae, 0xae, 0xec, 0xe9, 0x02, 0xba, 0x0b, 0x5e,
	0xc6, 0x4e, 0x78, 0xc5, 0x92, 0x85, 0xb6, 0x79, 0xdd, 0x3c, 0x87, 0xba, 0x74, 0x44, 0x52, 0xb4,
	0x09, 0xeb, 0xe6, 0x26, 0x63, 0x9e, 0x07, 0xae, 0x79, 0x7a, 0x91, 0x87, 0xdf, 0x1c, 0xe8, 0x4d,
	0xcd, 0xdf, 0x01, 0x1d, 0xc3, 0x3f, 0xd6, 0xe7, 0x45, 0xa9, 0x24, 0x51, 0x34, 0x5d, 0xd6, 0x17,
	0xfb, 0xa8, 0x49, 0xa8, 0xfd, 0xab, 0xd8, 0x8f, 0x74, 0x58, 0x73, 0xf0, 0x20, 0xb9, 0x92, 0xeb,
	0xeb, 0x97, 0x55, 0x4e, 0xeb, 0x2f, 0xdd, 0x74, 0xfd, 0x97, 0xf6, 0x05, 0x1b, 0x7c, 0xb8, 0x0f,
	0x83, 0xab, 0x9d, 0xd1, 0x3a, 0x74, 0x76, 0xca, 0x59, 0x69, 0x0f, 0xfe, 0xb8, 0xa4, 0x33, 0xe1,
	0x3b, 0xc8, 0x87, 0xfe, 0x4c, 0xcc, 0x4e, 0x0f, 0x38, 0x7b, 0x4d, 0x54, 0xfc, 0xc1, 0x6f, 0xa1,
	0x01, 0xc0, 0x4c, 0xbc, 0x61, 0xbb, 0xb4, 0x20, 0x2c, 0xf1, 0xdb, 0x2f, 0x5f, 0xc0, 0xff, 0x31,
	0x2f, 0xae, 0x9f, 0x3b, 0x77, 0xde, 0xf5, 0x6c, 0xf4, 0xb5, 0xb5, 0xf1, 0x76, 0x82, 0xc9, 0x32,
	0x9a, 0x6a, 0xc4, 0x8e, 0x10, 0x46, 0x12, 0x95, 0x27, 0x3d, 0xe3, 0xd5, 0xe3, 0xef, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xa5, 0x5f, 0xde, 0x29, 0x5f, 0x05, 0x00, 0x00,
}
