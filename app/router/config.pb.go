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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

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
	return fileDescriptor_6b1608360690c5fc, []int{8, 0}
}

// Domain for routing decision.
type Domain struct {
	// Domain matching type.
	Type Domain_Type `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.router.Domain_Type" json:"type,omitempty"`
	// Domain value.
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// Attributes of this domain. May be used for filtering.
	Attribute            []*Domain_Attribute `protobuf:"bytes,3,rep,name=attribute,proto3" json:"attribute,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
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

func (m *Domain) GetAttribute() []*Domain_Attribute {
	if m != nil {
		return m.Attribute
	}
	return nil
}

type Domain_Attribute struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Types that are valid to be assigned to TypedValue:
	//	*Domain_Attribute_BoolValue
	//	*Domain_Attribute_IntValue
	TypedValue           isDomain_Attribute_TypedValue `protobuf_oneof:"typed_value"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *Domain_Attribute) Reset()         { *m = Domain_Attribute{} }
func (m *Domain_Attribute) String() string { return proto.CompactTextString(m) }
func (*Domain_Attribute) ProtoMessage()    {}
func (*Domain_Attribute) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{0, 0}
}

func (m *Domain_Attribute) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Domain_Attribute.Unmarshal(m, b)
}
func (m *Domain_Attribute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Domain_Attribute.Marshal(b, m, deterministic)
}
func (m *Domain_Attribute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Domain_Attribute.Merge(m, src)
}
func (m *Domain_Attribute) XXX_Size() int {
	return xxx_messageInfo_Domain_Attribute.Size(m)
}
func (m *Domain_Attribute) XXX_DiscardUnknown() {
	xxx_messageInfo_Domain_Attribute.DiscardUnknown(m)
}

var xxx_messageInfo_Domain_Attribute proto.InternalMessageInfo

func (m *Domain_Attribute) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type isDomain_Attribute_TypedValue interface {
	isDomain_Attribute_TypedValue()
}

type Domain_Attribute_BoolValue struct {
	BoolValue bool `protobuf:"varint,2,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type Domain_Attribute_IntValue struct {
	IntValue int64 `protobuf:"varint,3,opt,name=int_value,json=intValue,proto3,oneof"`
}

func (*Domain_Attribute_BoolValue) isDomain_Attribute_TypedValue() {}

func (*Domain_Attribute_IntValue) isDomain_Attribute_TypedValue() {}

func (m *Domain_Attribute) GetTypedValue() isDomain_Attribute_TypedValue {
	if m != nil {
		return m.TypedValue
	}
	return nil
}

func (m *Domain_Attribute) GetBoolValue() bool {
	if x, ok := m.GetTypedValue().(*Domain_Attribute_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (m *Domain_Attribute) GetIntValue() int64 {
	if x, ok := m.GetTypedValue().(*Domain_Attribute_IntValue); ok {
		return x.IntValue
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Domain_Attribute) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Domain_Attribute_BoolValue)(nil),
		(*Domain_Attribute_IntValue)(nil),
	}
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
	// Types that are valid to be assigned to TargetTag:
	//	*RoutingRule_Tag
	//	*RoutingRule_BalancingTag
	TargetTag isRoutingRule_TargetTag `protobuf_oneof:"target_tag"`
	// List of domains for target domain matching.
	Domain []*Domain `protobuf:"bytes,2,rep,name=domain,proto3" json:"domain,omitempty"`
	// List of CIDRs for target IP address matching.
	// Deprecated. Use geoip below.
	Cidr []*CIDR `protobuf:"bytes,3,rep,name=cidr,proto3" json:"cidr,omitempty"` // Deprecated: Do not use.
	// List of GeoIPs for target IP address matching. If this entry exists, the cidr above will have no effect.
	// GeoIP fields with the same country code are supposed to contain exactly same content. They will be merged during runtime.
	// For customized GeoIPs, please leave country code empty.
	Geoip []*GeoIP `protobuf:"bytes,10,rep,name=geoip,proto3" json:"geoip,omitempty"`
	// A range of port [from, to]. If the destination port is in this range, this rule takes effect.
	// Deprecated. Use port_list.
	PortRange *net.PortRange `protobuf:"bytes,4,opt,name=port_range,json=portRange,proto3" json:"port_range,omitempty"` // Deprecated: Do not use.
	// List of ports.
	PortList *net.PortList `protobuf:"bytes,14,opt,name=port_list,json=portList,proto3" json:"port_list,omitempty"`
	// List of networks. Deprecated. Use networks.
	NetworkList *net.NetworkList `protobuf:"bytes,5,opt,name=network_list,json=networkList,proto3" json:"network_list,omitempty"` // Deprecated: Do not use.
	// List of networks for matching.
	Networks []net.Network `protobuf:"varint,13,rep,packed,name=networks,proto3,enum=v2ray.core.common.net.Network" json:"networks,omitempty"`
	// List of CIDRs for source IP address matching.
	SourceCidr []*CIDR `protobuf:"bytes,6,rep,name=source_cidr,json=sourceCidr,proto3" json:"source_cidr,omitempty"` // Deprecated: Do not use.
	// List of GeoIPs for source IP address matching. If this entry exists, the source_cidr above will have no effect.
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

type isRoutingRule_TargetTag interface {
	isRoutingRule_TargetTag()
}

type RoutingRule_Tag struct {
	Tag string `protobuf:"bytes,1,opt,name=tag,proto3,oneof"`
}

type RoutingRule_BalancingTag struct {
	BalancingTag string `protobuf:"bytes,12,opt,name=balancing_tag,json=balancingTag,proto3,oneof"`
}

func (*RoutingRule_Tag) isRoutingRule_TargetTag() {}

func (*RoutingRule_BalancingTag) isRoutingRule_TargetTag() {}

func (m *RoutingRule) GetTargetTag() isRoutingRule_TargetTag {
	if m != nil {
		return m.TargetTag
	}
	return nil
}

func (m *RoutingRule) GetTag() string {
	if x, ok := m.GetTargetTag().(*RoutingRule_Tag); ok {
		return x.Tag
	}
	return ""
}

func (m *RoutingRule) GetBalancingTag() string {
	if x, ok := m.GetTargetTag().(*RoutingRule_BalancingTag); ok {
		return x.BalancingTag
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

// Deprecated: Do not use.
func (m *RoutingRule) GetPortRange() *net.PortRange {
	if m != nil {
		return m.PortRange
	}
	return nil
}

func (m *RoutingRule) GetPortList() *net.PortList {
	if m != nil {
		return m.PortList
	}
	return nil
}

// Deprecated: Do not use.
func (m *RoutingRule) GetNetworkList() *net.NetworkList {
	if m != nil {
		return m.NetworkList
	}
	return nil
}

func (m *RoutingRule) GetNetworks() []net.Network {
	if m != nil {
		return m.Networks
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

// XXX_OneofWrappers is for the internal use of the proto package.
func (*RoutingRule) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*RoutingRule_Tag)(nil),
		(*RoutingRule_BalancingTag)(nil),
	}
}

type BalancingRule struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	OutboundSelector     []string `protobuf:"bytes,2,rep,name=outbound_selector,json=outboundSelector,proto3" json:"outbound_selector,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BalancingRule) Reset()         { *m = BalancingRule{} }
func (m *BalancingRule) String() string { return proto.CompactTextString(m) }
func (*BalancingRule) ProtoMessage()    {}
func (*BalancingRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{7}
}

func (m *BalancingRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BalancingRule.Unmarshal(m, b)
}
func (m *BalancingRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BalancingRule.Marshal(b, m, deterministic)
}
func (m *BalancingRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalancingRule.Merge(m, src)
}
func (m *BalancingRule) XXX_Size() int {
	return xxx_messageInfo_BalancingRule.Size(m)
}
func (m *BalancingRule) XXX_DiscardUnknown() {
	xxx_messageInfo_BalancingRule.DiscardUnknown(m)
}

var xxx_messageInfo_BalancingRule proto.InternalMessageInfo

func (m *BalancingRule) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *BalancingRule) GetOutboundSelector() []string {
	if m != nil {
		return m.OutboundSelector
	}
	return nil
}

type Config struct {
	DomainStrategy       Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,proto3,enum=v2ray.core.app.router.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	Rule                 []*RoutingRule        `protobuf:"bytes,2,rep,name=rule,proto3" json:"rule,omitempty"`
	BalancingRule        []*BalancingRule      `protobuf:"bytes,3,rep,name=balancing_rule,json=balancingRule,proto3" json:"balancing_rule,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b1608360690c5fc, []int{8}
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

func (m *Config) GetBalancingRule() []*BalancingRule {
	if m != nil {
		return m.BalancingRule
	}
	return nil
}

func init() {
	proto.RegisterEnum("v2ray.core.app.router.Domain_Type", Domain_Type_name, Domain_Type_value)
	proto.RegisterEnum("v2ray.core.app.router.Config_DomainStrategy", Config_DomainStrategy_name, Config_DomainStrategy_value)
	proto.RegisterType((*Domain)(nil), "v2ray.core.app.router.Domain")
	proto.RegisterType((*Domain_Attribute)(nil), "v2ray.core.app.router.Domain.Attribute")
	proto.RegisterType((*CIDR)(nil), "v2ray.core.app.router.CIDR")
	proto.RegisterType((*GeoIP)(nil), "v2ray.core.app.router.GeoIP")
	proto.RegisterType((*GeoIPList)(nil), "v2ray.core.app.router.GeoIPList")
	proto.RegisterType((*GeoSite)(nil), "v2ray.core.app.router.GeoSite")
	proto.RegisterType((*GeoSiteList)(nil), "v2ray.core.app.router.GeoSiteList")
	proto.RegisterType((*RoutingRule)(nil), "v2ray.core.app.router.RoutingRule")
	proto.RegisterType((*BalancingRule)(nil), "v2ray.core.app.router.BalancingRule")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.router.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/router/config.proto", fileDescriptor_6b1608360690c5fc)
}

var fileDescriptor_6b1608360690c5fc = []byte{
	// 897 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0xae, 0xed, 0x24, 0x1b, 0x1f, 0x27, 0xc1, 0x8c, 0x58, 0x64, 0x0a, 0xa5, 0xc1, 0x5a, 0xd8,
	0x48, 0x20, 0x47, 0xca, 0x02, 0x17, 0x08, 0xb4, 0x34, 0xe9, 0xd2, 0x46, 0x40, 0xa9, 0xa6, 0xbb,
	0x7b, 0x01, 0x17, 0x91, 0xe3, 0x4c, 0x8d, 0xb5, 0xce, 0xcc, 0x68, 0x3c, 0x5e, 0x36, 0xaf, 0x84,
	0xc4, 0x33, 0xf0, 0x28, 0xbc, 0x0a, 0x9a, 0x9f, 0x34, 0x0d, 0x34, 0x25, 0xe2, 0x6e, 0xce, 0x39,
	0xdf, 0x39, 0xf3, 0xcd, 0xf9, 0x1b, 0xf8, 0xe4, 0xf5, 0x48, 0xa4, 0xab, 0x24, 0x63, 0xcb, 0x61,
	0xc6, 0x04, 0x19, 0xa6, 0x9c, 0x0f, 0x05, 0xab, 0x25, 0x11, 0xc3, 0x8c, 0xd1, 0xeb, 0x22, 0x4f,
	0xb8, 0x60, 0x92, 0xa1, 0x87, 0x6b, 0x9c, 0x20, 0x49, 0xca, 0x79, 0x62, 0x30, 0x87, 0x8f, 0xfe,
	0xe1, 0x9e, 0xb1, 0xe5, 0x92, 0xd1, 0x21, 0x25, 0x72, 0xc8, 0x99, 0x90, 0xc6, 0xf9, 0xf0, 0xf1,
	0x6e, 0x14, 0x25, 0xf2, 0x37, 0x26, 0x5e, 0x19, 0x60, 0xfc, 0xa7, 0x0b, 0xad, 0x53, 0xb6, 0x4c,
	0x0b, 0x8a, 0xbe, 0x84, 0x86, 0x5c, 0x71, 0x12, 0x39, 0x7d, 0x67, 0xd0, 0x1b, 0xc5, 0xc9, 0x9d,
	0xf7, 0x27, 0x06, 0x9c, 0x3c, 0x5f, 0x71, 0x82, 0x35, 0x1e, 0xbd, 0x03, 0xcd, 0xd7, 0x69, 0x59,
	0x93, 0xc8, 0xed, 0x3b, 0x03, 0x1f, 0x1b, 0x01, 0x3d, 0x03, 0x3f, 0x95, 0x52, 0x14, 0xf3, 0x5a,
	0x92, 0xc8, 0xeb, 0x7b, 0x83, 0x60, 0xf4, 0xf8, 0xfe, 0x90, 0x27, 0x6b, 0x38, 0xde, 0x78, 0x1e,
	0x96, 0xe0, 0xdf, 0xe8, 0x51, 0x08, 0xde, 0x2b, 0xb2, 0xd2, 0x04, 0x7d, 0xac, 0x8e, 0xe8, 0x18,
	0x60, 0xce, 0x58, 0x39, 0xdb, 0x10, 0x68, 0x9f, 0x1f, 0x60, 0x5f, 0xe9, 0x5e, 0x6a, 0x1a, 0x47,
	0xe0, 0x17, 0x54, 0x5a, 0xbb, 0xd7, 0x77, 0x06, 0xde, 0xf9, 0x01, 0x6e, 0x17, 0x54, 0x6a, 0xf3,
	0xb8, 0x0b, 0x81, 0x7a, 0xc3, 0xc2, 0x00, 0xe2, 0x11, 0x34, 0xd4, 0xc3, 0x90, 0x0f, 0xcd, 0xcb,
	0x32, 0x2d, 0x68, 0x78, 0xa0, 0x8e, 0x98, 0xe4, 0xe4, 0x4d, 0xe8, 0x20, 0x58, 0xa7, 0x2a, 0x74,
	0x51, 0x1b, 0x1a, 0xdf, 0xd5, 0x65, 0x19, 0x7a, 0x71, 0x02, 0x8d, 0xc9, 0xf4, 0x14, 0xa3, 0x1e,
	0xb8, 0x05, 0xd7, 0xdc, 0x3a, 0xd8, 0x2d, 0x38, 0x7a, 0x17, 0x5a, 0x5c, 0x90, 0xeb, 0xe2, 0x8d,
	0xa6, 0xd5, 0xc5, 0x56, 0x8a, 0x7f, 0x81, 0xe6, 0x19, 0x61, 0xd3, 0x4b, 0xf4, 0x11, 0x74, 0x32,
	0x56, 0x53, 0x29, 0x56, 0xb3, 0x8c, 0x2d, 0x88, 0x7d, 0x56, 0x60, 0x75, 0x13, 0xb6, 0x20, 0x68,
	0x08, 0x8d, 0xac, 0x58, 0x88, 0xc8, 0xd5, 0xf9, 0x7b, 0x7f, 0x47, 0xfe, 0xd4, 0xf5, 0x58, 0x03,
	0xe3, 0xa7, 0xe0, 0xeb, 0xe0, 0x3f, 0x14, 0x95, 0x44, 0x23, 0x68, 0x12, 0x15, 0x2a, 0x72, 0xb4,
	0xfb, 0x07, 0x3b, 0xdc, 0xb5, 0x03, 0x36, 0xd0, 0x38, 0x83, 0x07, 0x67, 0x84, 0x5d, 0x15, 0x92,
	0xec, 0xc3, 0xef, 0x0b, 0x68, 0x2d, 0x74, 0x46, 0x2c, 0xc3, 0xa3, 0x7b, 0x2b, 0x8c, 0x2d, 0x38,
	0x9e, 0x40, 0x60, 0x2f, 0xd1, 0x3c, 0x3f, 0xdf, 0xe6, 0xf9, 0xe1, 0x6e, 0x9e, 0xca, 0x65, 0xcd,
	0xf4, 0xaf, 0x26, 0x04, 0x98, 0xd5, 0xb2, 0xa0, 0x39, 0xae, 0x4b, 0x82, 0x10, 0x78, 0x32, 0xcd,
	0x0d, 0xcb, 0xf3, 0x03, 0xac, 0x04, 0xf4, 0x31, 0x74, 0xe7, 0x69, 0x99, 0xd2, 0xac, 0xa0, 0xf9,
	0x4c, 0x59, 0x3b, 0xd6, 0xda, 0xb9, 0x51, 0x3f, 0x4f, 0xf3, 0xff, 0xf9, 0x0c, 0xf4, 0xc4, 0x56,
	0xc7, 0xfb, 0xcf, 0xea, 0x8c, 0xdd, 0xc8, 0x31, 0x15, 0x52, 0x45, 0xc9, 0x09, 0x2b, 0x78, 0x04,
	0xfb, 0x14, 0x45, 0x43, 0xd1, 0x04, 0x40, 0xcd, 0xf6, 0x4c, 0xa4, 0x34, 0x27, 0x51, 0xa3, 0xef,
	0x0c, 0x82, 0x51, 0xff, 0xb6, 0xa3, 0x19, 0xef, 0x84, 0x12, 0x99, 0x5c, 0x32, 0x21, 0xb1, 0xc2,
	0xe9, 0x3b, 0x7d, 0xbe, 0x16, 0xd1, 0xd7, 0xa0, 0x85, 0x59, 0x59, 0x54, 0x32, 0xea, 0xe9, 0x18,
	0xc7, 0xf7, 0xc4, 0x50, 0x95, 0xc1, 0x6d, 0x6e, 0x4f, 0x68, 0x0a, 0x1d, 0xbb, 0x38, 0x4c, 0x80,
	0xa6, 0x0e, 0x10, 0xef, 0x08, 0x70, 0x61, 0xa0, 0xca, 0x53, 0xd3, 0x08, 0xe8, 0x46, 0x81, 0xbe,
	0x82, 0xb6, 0x15, 0xab, 0xa8, 0xdb, 0xf7, 0x06, 0xbd, 0xed, 0x8a, 0xff, 0x3b, 0x0c, 0xbe, 0xc1,
	0xa3, 0x6f, 0x21, 0xa8, 0x58, 0x2d, 0x32, 0x32, 0xd3, 0x99, 0x6f, 0xed, 0x97, 0x79, 0x30, 0x3e,
	0x13, 0x95, 0xff, 0xa7, 0xd0, 0xb1, 0x11, 0x4c, 0x19, 0x82, 0x3d, 0xca, 0x60, 0xef, 0x3c, 0xd3,
	0xc5, 0x38, 0x02, 0xa8, 0x2b, 0x22, 0x66, 0x64, 0x99, 0x16, 0x65, 0xf4, 0xa0, 0xef, 0x0d, 0x7c,
	0xec, 0x2b, 0xcd, 0x33, 0xa5, 0x40, 0xc7, 0x10, 0x14, 0x74, 0xce, 0x6a, 0xba, 0xd0, 0x0d, 0xd7,
	0xd6, 0x76, 0xb0, 0x2a, 0xd5, 0x6c, 0x87, 0xd0, 0xd6, 0xab, 0x37, 0x63, 0x65, 0xe4, 0x6b, 0xeb,
	0x8d, 0x3c, 0xee, 0x00, 0xc8, 0x54, 0xe4, 0x44, 0x2a, 0xdf, 0xf8, 0x02, 0xba, 0xe3, 0x75, 0x9b,
	0xea, 0x16, 0x0f, 0x6f, 0xb5, 0xb8, 0x69, 0xf0, 0x4f, 0xe1, 0x6d, 0x56, 0x4b, 0x73, 0x5d, 0x45,
	0x4a, 0x92, 0x49, 0x66, 0xb6, 0x85, 0x8f, 0xc3, 0xb5, 0xe1, 0xca, 0xea, 0xe3, 0x3f, 0x5c, 0x68,
	0x4d, 0xf4, 0x17, 0x83, 0x5e, 0xc0, 0x5b, 0xa6, 0x89, 0x67, 0x95, 0x14, 0xa9, 0x24, 0xf9, 0xca,
	0xae, 0xfd, 0xcf, 0x76, 0xe5, 0xd2, 0x7c, 0x4d, 0x66, 0x02, 0xae, 0xac, 0x0f, 0xee, 0x2d, 0xb6,
	0x64, 0xf5, 0x85, 0x88, 0xba, 0x24, 0x76, 0x8c, 0x76, 0x7d, 0x21, 0xb7, 0xa6, 0x16, 0x6b, 0x3c,
	0xfa, 0x1e, 0x7a, 0x9b, 0x39, 0xd5, 0x11, 0xcc, 0x4c, 0x3d, 0xda, 0x11, 0x61, 0x2b, 0x2d, 0x78,
	0x33, 0xe3, 0x4a, 0x8c, 0xcf, 0xa0, 0xb7, 0x4d, 0x53, 0x2d, 0xeb, 0x93, 0x6a, 0x5a, 0x99, 0x6d,
	0xfe, 0xa2, 0x22, 0x53, 0x1e, 0x3a, 0x28, 0x84, 0xce, 0x94, 0x4f, 0xaf, 0x2f, 0x18, 0xfd, 0x31,
	0x95, 0xd9, 0xaf, 0xa1, 0x8b, 0x7a, 0x00, 0x53, 0xfe, 0x13, 0x3d, 0x25, 0xcb, 0x94, 0x2e, 0x42,
	0x6f, 0xfc, 0x0d, 0xbc, 0x97, 0xb1, 0xe5, 0xdd, 0x14, 0x2e, 0x9d, 0x9f, 0x5b, 0xe6, 0xf4, 0xbb,
	0xfb, 0xf0, 0xe5, 0x08, 0xa7, 0xab, 0x64, 0xa2, 0x10, 0x27, 0x9c, 0xeb, 0xf7, 0x11, 0x31, 0x6f,
	0xe9, 0xb2, 0x3e, 0xf9, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xb2, 0x28, 0x0d, 0x41, 0xf1, 0x07, 0x00,
	0x00,
}
