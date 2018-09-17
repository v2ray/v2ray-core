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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DomainMatchingType int32

const (
	DomainMatchingType_Full      DomainMatchingType = 0
	DomainMatchingType_Subdomain DomainMatchingType = 1
	DomainMatchingType_Keyword   DomainMatchingType = 2
)

var DomainMatchingType_name = map[int32]string{
	0: "Full",
	1: "Subdomain",
	2: "Keyword",
}

var DomainMatchingType_value = map[string]int32{
	"Full":      0,
	"Subdomain": 1,
	"Keyword":   2,
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
	proto.RegisterType((*NameServer)(nil), "v2ray.core.app.dns.NameServer")
	proto.RegisterType((*NameServer_PriorityDomain)(nil), "v2ray.core.app.dns.NameServer.PriorityDomain")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.dns.Config")
	proto.RegisterMapType((map[string]*net.IPOrDomain)(nil), "v2ray.core.app.dns.Config.HostsEntry")
	proto.RegisterType((*Config_HostMapping)(nil), "v2ray.core.app.dns.Config.HostMapping")
	proto.RegisterEnum("v2ray.core.app.dns.DomainMatchingType", DomainMatchingType_name, DomainMatchingType_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dns/config.proto", fileDescriptor_ed5695198e3def8f)
}

var fileDescriptor_ed5695198e3def8f = []byte{
	// 516 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0xd1, 0x6e, 0xd3, 0x3e,
	0x18, 0xc5, 0xff, 0x49, 0xda, 0x6e, 0xfd, 0xd2, 0x7f, 0x55, 0x7c, 0x31, 0x55, 0x45, 0x82, 0x32,
	0xc4, 0xa8, 0x40, 0x38, 0x52, 0x40, 0x02, 0x7a, 0x33, 0xb1, 0xad, 0x88, 0x0a, 0x0d, 0xaa, 0x0c,
	0x71, 0x01, 0x48, 0x95, 0x17, 0x9b, 0xce, 0xa2, 0xb1, 0x8d, 0xed, 0x16, 0x85, 0x37, 0x40, 0xbc,
	0x09, 0x4f, 0x89, 0x6a, 0x77, 0xb4, 0xb0, 0x0e, 0xb8, 0xe1, 0xae, 0x76, 0xcf, 0xef, 0x3b, 0x27,
	0xe7, 0x4b, 0xe0, 0xe6, 0x3c, 0xd5, 0xa4, 0xc4, 0xb9, 0x2c, 0x92, 0x5c, 0x6a, 0x96, 0x10, 0xa5,
	0x12, 0x2a, 0x4c, 0x92, 0x4b, 0xf1, 0x9e, 0x4f, 0xb0, 0xd2, 0xd2, 0x4a, 0x84, 0xce, 0x45, 0x9a,
	0x61, 0xa2, 0x14, 0xa6, 0xc2, 0x74, 0x6e, 0xff, 0x02, 0xe6, 0xb2, 0x28, 0xa4, 0x48, 0x04, 0xb3,
	0x09, 0xa1, 0x54, 0x33, 0x63, 0x3c, 0xdc, 0xb9, 0x7b, 0xb9, 0x90, 0x32, 0x63, 0xb9, 0x20, 0x96,
	0x4b, 0xe1, 0xc5, 0xbb, 0x5f, 0x43, 0x80, 0x17, 0xa4, 0x60, 0x27, 0x4c, 0xcf, 0x99, 0x46, 0x8f,
	0x61, 0x6b, 0x39, 0xac, 0x1d, 0x74, 0x83, 0x5e, 0x9c, 0x5e, 0xc7, 0x6b, 0x51, 0xfc, 0x24, 0x2c,
	0x98, 0xc5, 0x03, 0x41, 0x95, 0xe4, 0xc2, 0x66, 0xe7, 0x7a, 0xf4, 0x0e, 0x90, 0xd2, 0x5c, 0x6a,
	0x6e, 0xf9, 0x67, 0x46, 0xc7, 0x54, 0x16, 0x84, 0x8b, 0x76, 0xd8, 0x8d, 0x7a, 0x71, 0x7a, 0x0f,
	0x5f, 0x7c, 0x20, 0xbc, 0xb2, 0xc5, 0x23, 0x0f, 0x96, 0x47, 0x0e, 0xca, 0xae, 0xac, 0x0d, 0xf2,
	0x57, 0x1d, 0x0a, 0xcd, 0x9f, 0x45, 0xa8, 0x0f, 0x15, 0x5b, 0x2a, 0xe6, 0x72, 0x36, 0xd3, 0xbd,
	0x4d, 0x0e, 0x5e, 0x79, 0x4c, 0x6c, 0x7e, 0xc6, 0xc5, 0xe4, 0x55, 0xa9, 0x58, 0xe6, 0x18, 0xb4,
	0x03, 0xb5, 0x1f, 0xf9, 0x82, 0x5e, 0x3d, 0x5b, 0x9e, 0x76, 0xbf, 0x54, 0xa0, 0x76, 0xe8, 0x16,
	0x81, 0x06, 0x10, 0xaf, 0x02, 0x2e, 0xda, 0x88, 0xfe, 0xa2, 0x8d, 0x83, 0xb0, 0x1d, 0x64, 0xeb,
	0x1c, 0xda, 0x87, 0x58, 0x90, 0x82, 0x8d, 0x8d, 0x3b, 0xb7, 0xab, 0x6e, 0xcc, 0xb5, 0xdf, 0xd7,
	0x91, 0x81, 0x58, 0x6d, 0x64, 0x1f, 0xaa, 0xcf, 0xa4, 0xb1, 0x66, 0xd9, 0xe4, 0xad, 0x4d, 0xa8,
	0x8f, 0x8c, 0x9d, 0x6e, 0x20, 0xac, 0x2e, 0x5d, 0x0e, 0xcf, 0xa1, 0xab, 0x50, 0xcf, 0xa7, 0x9c,
	0x09, 0x3b, 0xe6, 0xaa, 0x1d, 0x75, 0x83, 0x5e, 0x23, 0xdb, 0xf6, 0x17, 0x43, 0x85, 0x86, 0xd0,
	0x30, 0x96, 0x58, 0x9e, 0x8f, 0xcf, 0x9c, 0x49, 0xc5, 0x99, 0xec, 0xfd, 0xc1, 0xe4, 0x98, 0x28,
	0xc5, 0xc5, 0x24, 0x8b, 0x3d, 0xeb, 0x7c, 0x3a, 0x6f, 0x01, 0x56, 0x01, 0x50, 0x0b, 0xa2, 0x0f,
	0xac, 0x74, 0xcb, 0xa9, 0x67, 0x8b, 0x9f, 0xe8, 0x21, 0x54, 0xe7, 0x64, 0x3a, 0x63, 0xae, 0xf2,
	0x38, 0xbd, 0x71, 0x49, 0x95, 0xc3, 0xd1, 0x4b, 0xbd, 0x7c, 0x0d, 0xbc, 0xbe, 0x1f, 0x3e, 0x0a,
	0x3a, 0x1f, 0x21, 0x5e, 0x33, 0xfe, 0x17, 0xbb, 0x47, 0x4d, 0x08, 0x5d, 0x41, 0x51, 0xaf, 0x91,
	0x85, 0x5c, 0xdd, 0xe9, 0x03, 0xba, 0x38, 0x03, 0x6d, 0x43, 0xe5, 0xe9, 0x6c, 0x3a, 0x6d, 0xfd,
	0x87, 0xfe, 0x87, 0xfa, 0xc9, 0xec, 0xd4, 0xc3, 0xad, 0x00, 0xc5, 0xb0, 0xf5, 0x9c, 0x95, 0x9f,
	0xa4, 0xa6, 0xad, 0xf0, 0xe0, 0x01, 0xec, 0xe4, 0xb2, 0xd8, 0x10, 0x6b, 0x14, 0xbc, 0x89, 0xa8,
	0x30, 0xdf, 0x42, 0xf4, 0x3a, 0xcd, 0x48, 0x89, 0x0f, 0x17, 0xff, 0x3d, 0x51, 0x0a, 0x1f, 0x09,
	0x73, 0x5a, 0x73, 0x9f, 0xe4, 0xfd, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff, 0x78, 0xfa, 0x21, 0x21,
	0x23, 0x04, 0x00, 0x00,
}
