package dns

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	router "v2ray.com/core/app/router"
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

type DomainMatchingType int32

const (
	DomainMatchingType_Full      DomainMatchingType = 0
	DomainMatchingType_Subdomain DomainMatchingType = 1
	DomainMatchingType_Keyword   DomainMatchingType = 2
	DomainMatchingType_Regex     DomainMatchingType = 3
	DomainMatchingType_New       DomainMatchingType = 4
)

var DomainMatchingType_name = map[int32]string{
	0: "Full",
	1: "Subdomain",
	2: "Keyword",
	3: "Regex",
	4: "New",
}

var DomainMatchingType_value = map[string]int32{
	"Full":      0,
	"Subdomain": 1,
	"Keyword":   2,
	"Regex":     3,
	"New":       4,
}

func (x DomainMatchingType) String() string {
	return proto.EnumName(DomainMatchingType_name, int32(x))
}

func (DomainMatchingType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{0}
}

type NameServer struct {
	Address              *net.Endpoint                `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PrioritizedDomain    []*NameServer_PriorityDomain `protobuf:"bytes,2,rep,name=prioritized_domain,json=prioritizedDomain,proto3" json:"prioritized_domain,omitempty"`
	Geoip                []*router.GeoIP              `protobuf:"bytes,3,rep,name=geoip,proto3" json:"geoip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *NameServer) Reset()         { *m = NameServer{} }
func (m *NameServer) String() string { return proto.CompactTextString(m) }
func (*NameServer) ProtoMessage()    {}
func (*NameServer) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{0}
}

func (m *NameServer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameServer.Unmarshal(m, b)
}
func (m *NameServer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameServer.Marshal(b, m, deterministic)
}
func (m *NameServer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameServer.Merge(m, src)
}
func (m *NameServer) XXX_Size() int {
	return xxx_messageInfo_NameServer.Size(m)
}
func (m *NameServer) XXX_DiscardUnknown() {
	xxx_messageInfo_NameServer.DiscardUnknown(m)
}

var xxx_messageInfo_NameServer proto.InternalMessageInfo

func (m *NameServer) GetAddress() *net.Endpoint {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *NameServer) GetPrioritizedDomain() []*NameServer_PriorityDomain {
	if m != nil {
		return m.PrioritizedDomain
	}
	return nil
}

func (m *NameServer) GetGeoip() []*router.GeoIP {
	if m != nil {
		return m.Geoip
	}
	return nil
}

type NameServer_PriorityDomain struct {
	Type                 DomainMatchingType `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.dns.DomainMatchingType" json:"type,omitempty"` // Deprecated: Do not use.
	Domain               string             `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *NameServer_PriorityDomain) Reset()         { *m = NameServer_PriorityDomain{} }
func (m *NameServer_PriorityDomain) String() string { return proto.CompactTextString(m) }
func (*NameServer_PriorityDomain) ProtoMessage()    {}
func (*NameServer_PriorityDomain) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{0, 0}
}

func (m *NameServer_PriorityDomain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameServer_PriorityDomain.Unmarshal(m, b)
}
func (m *NameServer_PriorityDomain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameServer_PriorityDomain.Marshal(b, m, deterministic)
}
func (m *NameServer_PriorityDomain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameServer_PriorityDomain.Merge(m, src)
}
func (m *NameServer_PriorityDomain) XXX_Size() int {
	return xxx_messageInfo_NameServer_PriorityDomain.Size(m)
}
func (m *NameServer_PriorityDomain) XXX_DiscardUnknown() {
	xxx_messageInfo_NameServer_PriorityDomain.DiscardUnknown(m)
}

var xxx_messageInfo_NameServer_PriorityDomain proto.InternalMessageInfo

// Deprecated: Do not use.
func (m *NameServer_PriorityDomain) GetType() DomainMatchingType {
	if m != nil {
		return m.Type
	}
	return DomainMatchingType_Full
}

func (m *NameServer_PriorityDomain) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

type Config struct {
	// Nameservers used by this DNS. Only traditional UDP servers are support at the moment.
	// A special value 'localhost' as a domain address can be set to use DNS on local system.
	NameServers []*net.Endpoint `protobuf:"bytes,1,rep,name=NameServers,proto3" json:"NameServers,omitempty"` // Deprecated: Do not use.
	// NameServer list used by this DNS client.
	NameServer []*NameServer `protobuf:"bytes,5,rep,name=name_server,json=nameServer,proto3" json:"name_server,omitempty"`
	// Static hosts. Domain to IP.
	// Deprecated. Use static_hosts.
	Hosts map[string]*net.IPOrDomain `protobuf:"bytes,2,rep,name=Hosts,proto3" json:"Hosts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // Deprecated: Do not use.
	// Client IP for EDNS client subnet. Must be 4 bytes (IPv4) or 16 bytes (IPv6).
	ClientIp    []byte                `protobuf:"bytes,3,opt,name=client_ip,json=clientIp,proto3" json:"client_ip,omitempty"`
	StaticHosts []*Config_HostMapping `protobuf:"bytes,4,rep,name=static_hosts,json=staticHosts,proto3" json:"static_hosts,omitempty"`
	// Tag is the inbound tag of DNS client.
	Tag                  string                     `protobuf:"bytes,6,opt,name=tag,proto3" json:"tag,omitempty"`
	Fake                 *Config_Fake               `protobuf:"bytes,7,opt,name=fake,proto3" json:"fake,omitempty"`
	ExternalRules        map[string]*ConfigPatterns `protobuf:"bytes,8,rep,name=external_rules,json=externalRules,proto3" json:"external_rules,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{1}
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

// Deprecated: Do not use.
func (m *Config) GetNameServers() []*net.Endpoint {
	if m != nil {
		return m.NameServers
	}
	return nil
}

func (m *Config) GetNameServer() []*NameServer {
	if m != nil {
		return m.NameServer
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

func (m *Config) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *Config) GetFake() *Config_Fake {
	if m != nil {
		return m.Fake
	}
	return nil
}

func (m *Config) GetExternalRules() map[string]*ConfigPatterns {
	if m != nil {
		return m.ExternalRules
	}
	return nil
}

type Config_HostMapping struct {
	Type   DomainMatchingType `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.dns.DomainMatchingType" json:"type,omitempty"` // Deprecated: Do not use.
	Domain string             `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Ip     [][]byte           `protobuf:"bytes,3,rep,name=ip,proto3" json:"ip,omitempty"`
	// ProxiedDomain indicates the mapped domain has the same IP address on this domain. V2Ray will use this domain for IP queries.
	// This field is only effective if ip is empty.
	ProxiedDomain        string   `protobuf:"bytes,4,opt,name=proxied_domain,json=proxiedDomain,proto3" json:"proxied_domain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config_HostMapping) Reset()         { *m = Config_HostMapping{} }
func (m *Config_HostMapping) String() string { return proto.CompactTextString(m) }
func (*Config_HostMapping) ProtoMessage()    {}
func (*Config_HostMapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{1, 1}
}

func (m *Config_HostMapping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config_HostMapping.Unmarshal(m, b)
}
func (m *Config_HostMapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config_HostMapping.Marshal(b, m, deterministic)
}
func (m *Config_HostMapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config_HostMapping.Merge(m, src)
}
func (m *Config_HostMapping) XXX_Size() int {
	return xxx_messageInfo_Config_HostMapping.Size(m)
}
func (m *Config_HostMapping) XXX_DiscardUnknown() {
	xxx_messageInfo_Config_HostMapping.DiscardUnknown(m)
}

var xxx_messageInfo_Config_HostMapping proto.InternalMessageInfo

// Deprecated: Do not use.
func (m *Config_HostMapping) GetType() DomainMatchingType {
	if m != nil {
		return m.Type
	}
	return DomainMatchingType_Full
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

func (m *Config_HostMapping) GetProxiedDomain() string {
	if m != nil {
		return m.ProxiedDomain
	}
	return ""
}

type Config_Fake struct {
	FakeRules            []string `protobuf:"bytes,1,rep,name=fake_rules,json=fakeRules,proto3" json:"fake_rules,omitempty"`
	FakeNet              string   `protobuf:"bytes,2,opt,name=fake_net,json=fakeNet,proto3" json:"fake_net,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config_Fake) Reset()         { *m = Config_Fake{} }
func (m *Config_Fake) String() string { return proto.CompactTextString(m) }
func (*Config_Fake) ProtoMessage()    {}
func (*Config_Fake) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{1, 2}
}

func (m *Config_Fake) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config_Fake.Unmarshal(m, b)
}
func (m *Config_Fake) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config_Fake.Marshal(b, m, deterministic)
}
func (m *Config_Fake) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config_Fake.Merge(m, src)
}
func (m *Config_Fake) XXX_Size() int {
	return xxx_messageInfo_Config_Fake.Size(m)
}
func (m *Config_Fake) XXX_DiscardUnknown() {
	xxx_messageInfo_Config_Fake.DiscardUnknown(m)
}

var xxx_messageInfo_Config_Fake proto.InternalMessageInfo

func (m *Config_Fake) GetFakeRules() []string {
	if m != nil {
		return m.FakeRules
	}
	return nil
}

func (m *Config_Fake) GetFakeNet() string {
	if m != nil {
		return m.FakeNet
	}
	return ""
}

type ConfigPatterns struct {
	Patterns             []string `protobuf:"bytes,1,rep,name=patterns,proto3" json:"patterns,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigPatterns) Reset()         { *m = ConfigPatterns{} }
func (m *ConfigPatterns) String() string { return proto.CompactTextString(m) }
func (*ConfigPatterns) ProtoMessage()    {}
func (*ConfigPatterns) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed5695198e3def8f, []int{1, 3}
}

func (m *ConfigPatterns) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigPatterns.Unmarshal(m, b)
}
func (m *ConfigPatterns) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigPatterns.Marshal(b, m, deterministic)
}
func (m *ConfigPatterns) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigPatterns.Merge(m, src)
}
func (m *ConfigPatterns) XXX_Size() int {
	return xxx_messageInfo_ConfigPatterns.Size(m)
}
func (m *ConfigPatterns) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigPatterns.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigPatterns proto.InternalMessageInfo

func (m *ConfigPatterns) GetPatterns() []string {
	if m != nil {
		return m.Patterns
	}
	return nil
}

func init() {
	proto.RegisterEnum("v2ray.core.app.dns.DomainMatchingType", DomainMatchingType_name, DomainMatchingType_value)
	proto.RegisterType((*NameServer)(nil), "v2ray.core.app.dns.NameServer")
	proto.RegisterType((*NameServer_PriorityDomain)(nil), "v2ray.core.app.dns.NameServer.PriorityDomain")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.dns.Config")
	proto.RegisterMapType((map[string]*ConfigPatterns)(nil), "v2ray.core.app.dns.Config.ExternalRulesEntry")
	proto.RegisterMapType((map[string]*net.IPOrDomain)(nil), "v2ray.core.app.dns.Config.HostsEntry")
	proto.RegisterType((*Config_HostMapping)(nil), "v2ray.core.app.dns.Config.HostMapping")
	proto.RegisterType((*Config_Fake)(nil), "v2ray.core.app.dns.Config.Fake")
	proto.RegisterType((*ConfigPatterns)(nil), "v2ray.core.app.dns.Config.patterns")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dns/config.proto", fileDescriptor_ed5695198e3def8f)
}

var fileDescriptor_ed5695198e3def8f = []byte{
	// 712 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x5d, 0x4f, 0xdb, 0x48,
	0x14, 0x5d, 0x3b, 0xdf, 0xd7, 0x10, 0x65, 0xe7, 0x01, 0x79, 0xbd, 0x5f, 0x2c, 0x88, 0x6c, 0xb4,
	0xab, 0x75, 0xa4, 0xb0, 0xd2, 0x2e, 0xfb, 0xb0, 0xb4, 0x40, 0x68, 0xa3, 0x8a, 0x34, 0x1a, 0x50,
	0x1f, 0xda, 0x4a, 0xd1, 0x60, 0x5f, 0x82, 0x45, 0x32, 0x33, 0x1a, 0x4f, 0x00, 0xf7, 0xb7, 0xf4,
	0xad, 0x6f, 0xfd, 0x1b, 0xfd, 0x63, 0x95, 0xc7, 0x86, 0x04, 0x12, 0x68, 0x5f, 0xfa, 0x36, 0x73,
	0x7d, 0xce, 0x3d, 0xf7, 0x9e, 0x7b, 0xc7, 0xb0, 0x79, 0xd9, 0x51, 0x2c, 0xf1, 0x03, 0x31, 0x69,
	0x07, 0x42, 0x61, 0x9b, 0x49, 0xd9, 0x0e, 0x79, 0xdc, 0x0e, 0x04, 0x3f, 0x8b, 0x46, 0xbe, 0x54,
	0x42, 0x0b, 0x42, 0x6e, 0x40, 0x0a, 0x7d, 0x26, 0xa5, 0x1f, 0xf2, 0xd8, 0xfb, 0xfd, 0x1e, 0x31,
	0x10, 0x93, 0x89, 0xe0, 0x6d, 0x8e, 0xba, 0xcd, 0xc2, 0x50, 0x61, 0x1c, 0x67, 0x64, 0xef, 0xcf,
	0x87, 0x81, 0x21, 0xc6, 0x3a, 0xe2, 0x4c, 0x47, 0x82, 0xe7, 0xe0, 0xe6, 0x92, 0x72, 0x94, 0x98,
	0x6a, 0x54, 0x77, 0x2a, 0xda, 0xf8, 0x64, 0x03, 0xf4, 0xd9, 0x04, 0x8f, 0x51, 0x5d, 0xa2, 0x22,
	0x3b, 0x50, 0xc9, 0x45, 0x5d, 0x6b, 0xdd, 0x6a, 0x39, 0x9d, 0x5f, 0xfd, 0xb9, 0x92, 0x33, 0x45,
	0x9f, 0xa3, 0xf6, 0xbb, 0x3c, 0x94, 0x22, 0xe2, 0x9a, 0xde, 0xe0, 0xc9, 0x5b, 0x20, 0x52, 0x45,
	0x42, 0x45, 0x3a, 0x7a, 0x87, 0xe1, 0x30, 0x14, 0x13, 0x16, 0x71, 0xd7, 0x5e, 0x2f, 0xb4, 0x9c,
	0xce, 0x5f, 0xfe, 0x62, 0xe3, 0xfe, 0x4c, 0xd6, 0x1f, 0x64, 0xc4, 0xe4, 0xc0, 0x90, 0xe8, 0xf7,
	0x73, 0x89, 0xb2, 0x10, 0xe9, 0x40, 0x69, 0x84, 0x22, 0x92, 0x6e, 0xc1, 0x24, 0xfc, 0xe9, 0x7e,
	0xc2, 0xac, 0x37, 0xff, 0x19, 0x8a, 0xde, 0x80, 0x66, 0x50, 0xef, 0x1c, 0xea, 0x77, 0x13, 0x93,
	0xff, 0xa1, 0xa8, 0x13, 0x89, 0xa6, 0xb7, 0x7a, 0xa7, 0xb9, 0xac, 0xaa, 0x0c, 0x79, 0xc4, 0x74,
	0x70, 0x1e, 0xf1, 0xd1, 0x49, 0x22, 0x71, 0xcf, 0x76, 0x2d, 0x6a, 0x78, 0x64, 0x0d, 0xca, 0xb7,
	0x7d, 0x59, 0xad, 0x1a, 0xcd, 0x6f, 0x1b, 0x1f, 0x2a, 0x50, 0xde, 0x37, 0xb6, 0x92, 0x2e, 0x38,
	0xb3, 0xc6, 0x52, 0x17, 0x0b, 0x5f, 0xe1, 0xa2, 0x91, 0x98, 0xe7, 0x91, 0x5d, 0x70, 0x38, 0x9b,
	0xe0, 0x30, 0x36, 0x77, 0xb7, 0x64, 0xd2, 0xfc, 0xf2, 0xb8, 0x8d, 0x14, 0xf8, 0x6c, 0x92, 0xbb,
	0x50, 0x7a, 0x2e, 0x62, 0x1d, 0xe7, 0x13, 0xd8, 0x5a, 0x46, 0xcd, 0x4a, 0xf6, 0x0d, 0xae, 0xcb,
	0xb5, 0x4a, 0x4c, 0x1d, 0x19, 0x8f, 0xfc, 0x08, 0xb5, 0x60, 0x1c, 0x21, 0xd7, 0x43, 0xe3, 0xba,
	0xd5, 0x5a, 0xa1, 0xd5, 0x2c, 0xd0, 0x93, 0xa4, 0x07, 0x2b, 0xb1, 0x66, 0x3a, 0x0a, 0x86, 0xe7,
	0x46, 0xa4, 0x68, 0x44, 0x9a, 0x5f, 0x10, 0x39, 0x62, 0x52, 0x46, 0x7c, 0x44, 0x9d, 0x8c, 0x9b,
	0xe9, 0x34, 0xa0, 0xa0, 0xd9, 0xc8, 0x2d, 0x1b, 0x43, 0xd3, 0x23, 0xd9, 0x86, 0xe2, 0x19, 0xbb,
	0x40, 0xb7, 0xb2, 0xb8, 0x81, 0xf7, 0x92, 0x1e, 0xb2, 0x0b, 0xa4, 0x06, 0x4c, 0x4e, 0xa0, 0x8e,
	0xd7, 0x1a, 0x15, 0x67, 0xe3, 0xa1, 0x9a, 0x8e, 0x31, 0x76, 0xab, 0x0f, 0xaf, 0x5e, 0x4e, 0xef,
	0xe6, 0x04, 0x9a, 0xe2, 0x8d, 0x01, 0x74, 0x15, 0xe7, 0x63, 0xde, 0x1b, 0x80, 0x99, 0x3b, 0x69,
	0xa9, 0x17, 0x98, 0x98, 0xed, 0xa9, 0xd1, 0xf4, 0x48, 0xfe, 0x81, 0xd2, 0x25, 0x1b, 0x4f, 0xd1,
	0xec, 0x83, 0xd3, 0xf9, 0xed, 0x81, 0x39, 0xf7, 0x06, 0x2f, 0x55, 0xbe, 0xdb, 0x19, 0xfe, 0x3f,
	0xfb, 0x5f, 0xcb, 0x7b, 0x6f, 0x81, 0x33, 0x67, 0xcb, 0xb7, 0xda, 0x4e, 0x52, 0x07, 0x3b, 0x7f,
	0x38, 0x2b, 0xd4, 0x8e, 0x24, 0xd9, 0x82, 0xba, 0x54, 0xe2, 0x3a, 0x9a, 0xbd, 0xd2, 0xa2, 0xc1,
	0xaf, 0xe6, 0xd1, 0x4c, 0xc4, 0x7b, 0x02, 0xc5, 0xd4, 0x5f, 0xf2, 0x33, 0x40, 0xea, 0x70, 0xee,
	0x6a, 0xba, 0xd0, 0x35, 0x5a, 0x4b, 0x23, 0xc6, 0x22, 0xf2, 0x03, 0x54, 0xcd, 0x67, 0x8e, 0x3a,
	0xd7, 0xad, 0xa4, 0xf7, 0x3e, 0x6a, 0xaf, 0x09, 0x55, 0xc9, 0x74, 0xea, 0x67, 0x4c, 0xbc, 0xd9,
	0x39, 0xcf, 0x71, 0x7b, 0xf7, 0x10, 0xc8, 0xe2, 0x28, 0x96, 0xb8, 0xbd, 0x73, 0xd7, 0xed, 0xcd,
	0x47, 0x46, 0x7b, 0x93, 0x7b, 0xce, 0xef, 0x3f, 0xfa, 0x40, 0x16, 0xfd, 0x23, 0x55, 0x28, 0x1e,
	0x4e, 0xc7, 0xe3, 0xc6, 0x77, 0x64, 0x15, 0x6a, 0xc7, 0xd3, 0xd3, 0xcc, 0x92, 0x86, 0x45, 0x1c,
	0xa8, 0xbc, 0xc0, 0xe4, 0x4a, 0xa8, 0xb0, 0x61, 0x93, 0x1a, 0x94, 0x28, 0x8e, 0xf0, 0xba, 0x51,
	0x20, 0x15, 0x28, 0xf4, 0xf1, 0xaa, 0x51, 0xdc, 0xfb, 0x1b, 0xd6, 0x02, 0x31, 0x59, 0x52, 0xc4,
	0xc0, 0x7a, 0x5d, 0x08, 0x79, 0xfc, 0xd1, 0x26, 0xaf, 0x3a, 0x94, 0x25, 0xfe, 0x7e, 0xfa, 0xed,
	0xa9, 0x94, 0xfe, 0x01, 0x8f, 0x4f, 0xcb, 0xe6, 0xc7, 0xbb, 0xfd, 0x39, 0x00, 0x00, 0xff, 0xff,
	0x04, 0xa0, 0x0f, 0x70, 0x31, 0x06, 0x00, 0x00,
}
