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
)

var DomainMatchingType_name = map[int32]string{
	0: "Full",
	1: "Subdomain",
	2: "Keyword",
	3: "Regex",
}

var DomainMatchingType_value = map[string]int32{
	"Full":      0,
	"Subdomain": 1,
	"Keyword":   2,
	"Regex":     3,
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
	Type                 DomainMatchingType `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.dns.DomainMatchingType" json:"type,omitempty"`
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
	StaticHosts []*Config_HostMapping `protobuf:"bytes,4,rep,name=static_hosts,json=staticHosts,proto3" json:"static_hosts,omitempty"` // Deprecated: Do not use.
	// Tag is the inbound tag of DNS client.
	Tag                  string                     `protobuf:"bytes,6,opt,name=tag,proto3" json:"tag,omitempty"`
	UseFake              []string                   `protobuf:"bytes,7,rep,name=use_fake,json=useFake,proto3" json:"use_fake,omitempty"`
	HostRules            []*Config_HostMapping      `protobuf:"bytes,8,rep,name=host_rules,json=hostRules,proto3" json:"host_rules,omitempty"`
	ExternalRules        map[string]*ConfigPatterns `protobuf:"bytes,9,rep,name=external_rules,json=externalRules,proto3" json:"external_rules,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
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

// Deprecated: Do not use.
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

func (m *Config) GetUseFake() []string {
	if m != nil {
		return m.UseFake
	}
	return nil
}

func (m *Config) GetHostRules() []*Config_HostMapping {
	if m != nil {
		return m.HostRules
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
	Type    DomainMatchingType `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.dns.DomainMatchingType" json:"type,omitempty"` // Deprecated: Do not use.
	Pattern string             `protobuf:"bytes,2,opt,name=pattern,proto3" json:"pattern,omitempty"`
	Ip      [][]byte           `protobuf:"bytes,3,rep,name=ip,proto3" json:"ip,omitempty"`
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

func (m *Config_HostMapping) GetPattern() string {
	if m != nil {
		return m.Pattern
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
	return fileDescriptor_ed5695198e3def8f, []int{1, 2}
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
	proto.RegisterType((*ConfigPatterns)(nil), "v2ray.core.app.dns.Config.patterns")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dns/config.proto", fileDescriptor_ed5695198e3def8f)
}

var fileDescriptor_ed5695198e3def8f = []byte{
	// 704 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x5d, 0x4f, 0x1b, 0x47,
	0x14, 0xed, 0xae, 0x3f, 0xf7, 0x1a, 0x2c, 0x77, 0x1e, 0xd0, 0x76, 0x5b, 0xb5, 0x14, 0x84, 0x6b,
	0xb5, 0xea, 0x5a, 0x72, 0x2b, 0xb5, 0xf0, 0x50, 0x54, 0xc0, 0xb4, 0x28, 0x22, 0xb1, 0x06, 0x94,
	0x87, 0x24, 0x92, 0x35, 0xec, 0x5e, 0xcc, 0x0a, 0x7b, 0x66, 0x34, 0x3b, 0x26, 0x6c, 0x7e, 0x47,
	0x9e, 0xf3, 0x03, 0xf2, 0x27, 0xf2, 0xd7, 0xa2, 0x9d, 0x59, 0x63, 0x03, 0x86, 0x7c, 0xbc, 0xcd,
	0x8c, 0xef, 0x39, 0xe7, 0x9e, 0xb3, 0xf7, 0x1a, 0x36, 0xaf, 0x7a, 0x8a, 0x65, 0x61, 0x24, 0x26,
	0xdd, 0x48, 0x28, 0xec, 0x32, 0x29, 0xbb, 0x31, 0x4f, 0xbb, 0x91, 0xe0, 0xe7, 0xc9, 0x28, 0x94,
	0x4a, 0x68, 0x41, 0xc8, 0xac, 0x48, 0x61, 0xc8, 0xa4, 0x0c, 0x63, 0x9e, 0x06, 0xbf, 0xdc, 0x01,
	0x46, 0x62, 0x32, 0x11, 0xbc, 0xcb, 0x51, 0x77, 0x59, 0x1c, 0x2b, 0x4c, 0x53, 0x0b, 0x0e, 0x7e,
	0x7b, 0xb8, 0x30, 0xc6, 0x54, 0x27, 0x9c, 0xe9, 0x44, 0xf0, 0xa2, 0xb8, 0xbd, 0xa4, 0x1d, 0x25,
	0xa6, 0x1a, 0xd5, 0xad, 0x8e, 0x36, 0x3e, 0xb8, 0x00, 0x4f, 0xd9, 0x04, 0x4f, 0x50, 0x5d, 0xa1,
	0x22, 0xdb, 0x50, 0x2b, 0x44, 0x7d, 0x67, 0xdd, 0xe9, 0x34, 0x7a, 0x3f, 0x85, 0x0b, 0x2d, 0x5b,
	0xc5, 0x90, 0xa3, 0x0e, 0xfb, 0x3c, 0x96, 0x22, 0xe1, 0x9a, 0xce, 0xea, 0xc9, 0x2b, 0x20, 0x52,
	0x25, 0x42, 0x25, 0x3a, 0x79, 0x83, 0xf1, 0x30, 0x16, 0x13, 0x96, 0x70, 0xdf, 0x5d, 0x2f, 0x75,
	0x1a, 0xbd, 0xdf, 0xc3, 0xfb, 0xc6, 0xc3, 0xb9, 0x6c, 0x38, 0xb0, 0xc0, 0xec, 0xc0, 0x80, 0xe8,
	0xb7, 0x0b, 0x44, 0xf6, 0x89, 0xf4, 0xa0, 0x32, 0x42, 0x91, 0x48, 0xbf, 0x64, 0x08, 0x7f, 0xb8,
	0x4b, 0x68, 0xbd, 0x85, 0xff, 0xa1, 0x38, 0x1a, 0x50, 0x5b, 0x1a, 0xc4, 0xd0, 0xbc, 0x4d, 0x4c,
	0x76, 0xa0, 0xac, 0x33, 0x89, 0xc6, 0x5b, 0xb3, 0xd7, 0x5e, 0xd6, 0x95, 0xad, 0x3c, 0x66, 0x3a,
	0xba, 0x48, 0xf8, 0xe8, 0x34, 0x93, 0x48, 0x0d, 0x86, 0xac, 0x41, 0xf5, 0xc6, 0x93, 0xd3, 0xf1,
	0x68, 0x71, 0xdb, 0x78, 0x5b, 0x83, 0xea, 0xbe, 0x89, 0x94, 0xf4, 0xa1, 0x31, 0x37, 0x95, 0x27,
	0x58, 0xfa, 0x8c, 0x04, 0xf7, 0x5c, 0xdf, 0xa1, 0x8b, 0x38, 0xb2, 0x0b, 0x0d, 0xce, 0x26, 0x38,
	0x4c, 0xcd, 0xdd, 0xaf, 0x18, 0x9a, 0x1f, 0x1f, 0x8f, 0x90, 0x02, 0x9f, 0x7f, 0xc5, 0x5d, 0xa8,
	0xfc, 0x2f, 0x52, 0x9d, 0x16, 0xe9, 0x6f, 0x2d, 0x83, 0xda, 0x96, 0x43, 0x53, 0xd7, 0xe7, 0x5a,
	0x65, 0xa6, 0x0f, 0x8b, 0x23, 0xdf, 0x83, 0x17, 0x8d, 0x13, 0xe4, 0x7a, 0x68, 0x12, 0x77, 0x3a,
	0x2b, 0xb4, 0x6e, 0x1f, 0x8e, 0x24, 0x39, 0x86, 0x95, 0x54, 0x33, 0x9d, 0x44, 0xc3, 0x0b, 0x23,
	0x52, 0x36, 0x22, 0xed, 0x4f, 0x88, 0x1c, 0x33, 0x29, 0x13, 0x3e, 0xb2, 0x6e, 0x2d, 0xde, 0x6a,
	0xb5, 0xa0, 0xa4, 0xd9, 0xc8, 0xaf, 0x9a, 0x50, 0xf3, 0x23, 0xf9, 0x0e, 0xea, 0xd3, 0x14, 0x87,
	0xe7, 0xec, 0x12, 0xfd, 0xda, 0x7a, 0xa9, 0xe3, 0xd1, 0xda, 0x34, 0xc5, 0x43, 0x76, 0x89, 0xa4,
	0x0f, 0x90, 0x8b, 0x0e, 0xd5, 0x74, 0x8c, 0xa9, 0x5f, 0xff, 0x12, 0x65, 0xea, 0xe5, 0x48, 0x9a,
	0x03, 0xc9, 0x29, 0x34, 0xf1, 0x5a, 0xa3, 0xe2, 0x6c, 0x5c, 0x50, 0x79, 0x0f, 0xcf, 0x69, 0x41,
	0xd5, 0x2f, 0x00, 0x86, 0xc1, 0x24, 0x46, 0x57, 0x71, 0xf1, 0x2d, 0x78, 0x09, 0x30, 0x8f, 0x33,
	0xf7, 0x75, 0x89, 0x99, 0x19, 0x35, 0x8f, 0xe6, 0x47, 0xf2, 0x17, 0x54, 0xae, 0xd8, 0x78, 0x8a,
	0x66, 0x80, 0x1a, 0xbd, 0x9f, 0x1f, 0x18, 0x8c, 0xa3, 0xc1, 0x33, 0x55, 0x2c, 0x82, 0xad, 0xdf,
	0x71, 0xff, 0x76, 0x82, 0x77, 0x0e, 0x34, 0x16, 0xdc, 0x90, 0x7f, 0xbe, 0x66, 0x94, 0x4d, 0xfa,
	0x76, 0x9c, 0x7d, 0xa8, 0x49, 0xa6, 0xf3, 0xf6, 0x8b, 0x79, 0x9e, 0x5d, 0x49, 0x13, 0xdc, 0x62,
	0xcf, 0x56, 0xa8, 0x9b, 0x48, 0xb2, 0x05, 0x4d, 0xa9, 0xc4, 0x75, 0x32, 0x5f, 0xea, 0xb2, 0x01,
	0xac, 0x16, 0xaf, 0x56, 0x26, 0x68, 0x43, 0xbd, 0x60, 0x48, 0x49, 0x30, 0x3f, 0x9b, 0x2d, 0xf0,
	0xe8, 0xcd, 0x3d, 0x40, 0x20, 0xf7, 0xa3, 0x5c, 0x92, 0xd6, 0xf6, 0xed, 0xb4, 0x36, 0x1f, 0xf9,
	0x34, 0x33, 0xee, 0x85, 0xbc, 0x7e, 0xed, 0x03, 0xb9, 0xef, 0x9f, 0xd4, 0xa1, 0x7c, 0x38, 0x1d,
	0x8f, 0x5b, 0xdf, 0x90, 0x55, 0xf0, 0x4e, 0xa6, 0x67, 0xd6, 0x50, 0xcb, 0x21, 0x0d, 0xa8, 0x3d,
	0xc1, 0xec, 0xb5, 0x50, 0x71, 0xcb, 0x25, 0x1e, 0x54, 0x28, 0x8e, 0xf0, 0xba, 0x55, 0xda, 0xfb,
	0x13, 0xd6, 0x22, 0x31, 0x59, 0xa2, 0x3d, 0x70, 0x5e, 0x94, 0x62, 0x9e, 0xbe, 0x77, 0xc9, 0xf3,
	0x1e, 0x65, 0x59, 0xb8, 0x9f, 0xff, 0xf6, 0xaf, 0x94, 0xe1, 0x01, 0x4f, 0xcf, 0xaa, 0xe6, 0xcf,
	0xf5, 0x8f, 0x8f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x30, 0xe2, 0x0a, 0xbf, 0x15, 0x06, 0x00, 0x00,
}
