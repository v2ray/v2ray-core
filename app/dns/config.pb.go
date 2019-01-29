package dns

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
	StaticHosts []*Config_HostMapping `protobuf:"bytes,4,rep,name=static_hosts,json=staticHosts,proto3" json:"static_hosts,omitempty"`
	// Tag is the inbound tag of DNS client.
	Tag                  string   `protobuf:"bytes,6,opt,name=tag,proto3" json:"tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
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

type Config_HostMapping struct {
	Type                 DomainMatchingType `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.app.dns.DomainMatchingType" json:"type,omitempty"`
	Domain               string             `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Ip                   [][]byte           `protobuf:"bytes,3,rep,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
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

func init() {
	proto.RegisterEnum("v2ray.core.app.dns.DomainMatchingType", DomainMatchingType_name, DomainMatchingType_value)
	proto.RegisterType((*NameServer)(nil), "v2ray.core.app.dns.NameServer")
	proto.RegisterType((*NameServer_PriorityDomain)(nil), "v2ray.core.app.dns.NameServer.PriorityDomain")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.dns.Config")
	proto.RegisterMapType((map[string]*net.IPOrDomain)(nil), "v2ray.core.app.dns.Config.HostsEntry")
	proto.RegisterType((*Config_HostMapping)(nil), "v2ray.core.app.dns.Config.HostMapping")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dns/config.proto", fileDescriptor_ed5695198e3def8f)
}

var fileDescriptor_ed5695198e3def8f = []byte{
	// 530 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0xdb, 0x6e, 0xd3, 0x30,
	0x18, 0x26, 0x49, 0xdb, 0xad, 0x7f, 0xc6, 0x54, 0x7c, 0x31, 0x45, 0x45, 0x82, 0x32, 0xc4, 0xa8,
	0x40, 0x38, 0x52, 0x40, 0x02, 0x76, 0x33, 0xb1, 0xad, 0x88, 0x0a, 0x0d, 0x2a, 0x0f, 0x71, 0x01,
	0x48, 0x95, 0x97, 0x98, 0xce, 0xa2, 0xb1, 0x8d, 0xed, 0x16, 0xc2, 0x2b, 0xf0, 0x08, 0xbc, 0x01,
	0x4f, 0x89, 0x6a, 0x77, 0xb4, 0xb0, 0x0e, 0xb8, 0xe1, 0xce, 0x87, 0xef, 0x94, 0xef, 0x77, 0xe0,
	0xe6, 0x34, 0xd3, 0xb4, 0xc2, 0xb9, 0x2c, 0xd3, 0x5c, 0x6a, 0x96, 0x52, 0xa5, 0xd2, 0x42, 0x98,
	0x34, 0x97, 0xe2, 0x3d, 0x1f, 0x61, 0xa5, 0xa5, 0x95, 0x08, 0x9d, 0x81, 0x34, 0xc3, 0x54, 0x29,
	0x5c, 0x08, 0xd3, 0xbe, 0xfd, 0x1b, 0x31, 0x97, 0x65, 0x29, 0x45, 0x2a, 0x98, 0x4d, 0x69, 0x51,
	0x68, 0x66, 0x8c, 0x27, 0xb7, 0xef, 0x5e, 0x0c, 0x2c, 0x98, 0xb1, 0x5c, 0x50, 0xcb, 0xa5, 0xf0,
	0xe0, 0xed, 0xaf, 0x21, 0xc0, 0x0b, 0x5a, 0xb2, 0x63, 0xa6, 0xa7, 0x4c, 0xa3, 0xc7, 0xb0, 0x36,
	0x17, 0x4b, 0x82, 0x4e, 0xd0, 0x8d, 0xb3, 0xeb, 0x78, 0x29, 0x8a, 0x57, 0xc2, 0x82, 0x59, 0xdc,
	0x13, 0x85, 0x92, 0x5c, 0x58, 0x72, 0x86, 0x47, 0xef, 0x00, 0x29, 0xcd, 0xa5, 0xe6, 0x96, 0x7f,
	0x61, 0xc5, 0xb0, 0x90, 0x25, 0xe5, 0x22, 0x09, 0x3b, 0x51, 0x37, 0xce, 0xee, 0xe1, 0xf3, 0x1f,
	0x84, 0x17, 0xb6, 0x78, 0xe0, 0x89, 0xd5, 0xa1, 0x23, 0x91, 0x2b, 0x4b, 0x42, 0xfe, 0xa8, 0x5d,
	0xc0, 0xe6, 0xaf, 0x20, 0xb4, 0x0b, 0x35, 0x5b, 0x29, 0xe6, 0x72, 0x6e, 0x66, 0x3b, 0xab, 0x1c,
	0x3c, 0xf2, 0x88, 0xda, 0xfc, 0x94, 0x8b, 0xd1, 0xab, 0x4a, 0x31, 0xe2, 0x38, 0x68, 0x0b, 0x1a,
	0x3f, 0xf3, 0x05, 0xdd, 0x26, 0x99, 0xef, 0xb6, 0xbf, 0xd5, 0xa0, 0x71, 0xe0, 0x06, 0x81, 0x7a,
	0x10, 0x2f, 0x02, 0xce, 0xda, 0x88, 0xfe, 0xa1, 0x8d, 0xfd, 0x30, 0x09, 0xc8, 0x32, 0x0f, 0xed,
	0x41, 0x2c, 0x68, 0xc9, 0x86, 0xc6, 0xed, 0x93, 0xba, 0x93, 0xb9, 0xf6, 0xe7, 0x3a, 0x08, 0x88,
	0xc5, 0x44, 0xf6, 0xa0, 0xfe, 0x4c, 0x1a, 0x6b, 0xe6, 0x4d, 0xde, 0x5a, 0x45, 0xf5, 0x91, 0xb1,
	0xc3, 0xf5, 0x84, 0xd5, 0x95, 0xcb, 0xe1, 0x79, 0xe8, 0x2a, 0x34, 0xf3, 0x31, 0x67, 0xc2, 0x0e,
	0xb9, 0x4a, 0xa2, 0x4e, 0xd0, 0xdd, 0x20, 0xeb, 0xfe, 0xa0, 0xaf, 0x50, 0x1f, 0x36, 0x8c, 0xa5,
	0x96, 0xe7, 0xc3, 0x53, 0x67, 0x52, 0x73, 0x26, 0x3b, 0x7f, 0x31, 0x39, 0xa2, 0x4a, 0x71, 0x31,
	0x22, 0xb1, 0xe7, 0x7a, 0x9f, 0x16, 0x44, 0x96, 0x8e, 0x92, 0x86, 0x2b, 0x74, 0xb6, 0x6c, 0xbf,
	0x05, 0x58, 0x44, 0x9a, 0xdd, 0x7f, 0x60, 0x95, 0x1b, 0x57, 0x93, 0xcc, 0x96, 0xe8, 0x21, 0xd4,
	0xa7, 0x74, 0x3c, 0x61, 0x6e, 0x08, 0x71, 0x76, 0xe3, 0x82, 0x72, 0xfb, 0x83, 0x97, 0x7a, 0xfe,
	0x30, 0x3c, 0x7e, 0x37, 0x7c, 0x14, 0xb4, 0x3f, 0x42, 0xbc, 0x14, 0xe5, 0x7f, 0xbc, 0x06, 0xb4,
	0x09, 0xa1, 0xab, 0x2c, 0xea, 0x6e, 0x90, 0x90, 0xab, 0x3b, 0x3d, 0x40, 0xe7, 0x35, 0xd0, 0x3a,
	0xd4, 0x9e, 0x4e, 0xc6, 0xe3, 0xd6, 0x25, 0x74, 0x19, 0x9a, 0xc7, 0x93, 0x13, 0x4f, 0x6e, 0x05,
	0x28, 0x86, 0xb5, 0xe7, 0xac, 0xfa, 0x24, 0x75, 0xd1, 0x0a, 0x51, 0x13, 0xea, 0x84, 0x8d, 0xd8,
	0xe7, 0x56, 0xb4, 0xff, 0x00, 0xb6, 0x72, 0x59, 0xae, 0x48, 0x38, 0x08, 0xde, 0x44, 0x85, 0x30,
	0xdf, 0x43, 0xf4, 0x3a, 0x23, 0xb4, 0xc2, 0x07, 0xb3, 0xbb, 0x27, 0x4a, 0xe1, 0x43, 0x61, 0x4e,
	0x1a, 0xee, 0x7f, 0xbd, 0xff, 0x23, 0x00, 0x00, 0xff, 0xff, 0x83, 0x2b, 0x1c, 0x4c, 0x40, 0x04,
	0x00, 0x00,
}
