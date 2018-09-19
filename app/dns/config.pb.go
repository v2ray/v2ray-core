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
	// 523 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0xd1, 0x6e, 0xd3, 0x30,
	0x18, 0x85, 0x49, 0xd2, 0x76, 0xeb, 0x9f, 0x52, 0x15, 0x5f, 0x4c, 0x55, 0x91, 0xa0, 0x0c, 0x31,
	0x2a, 0x10, 0x8e, 0x14, 0x90, 0x80, 0xdd, 0x4c, 0x6c, 0x2b, 0xa2, 0x42, 0x83, 0xca, 0x43, 0x5c,
	0x00, 0x52, 0xe5, 0x25, 0xa6, 0xb3, 0x68, 0x6c, 0x63, 0xbb, 0x85, 0xf0, 0x06, 0x88, 0x37, 0xe1,
	0x29, 0x51, 0xed, 0x8e, 0x16, 0xd6, 0x01, 0x37, 0xdc, 0xd5, 0xee, 0xf9, 0xfe, 0x73, 0x72, 0xfe,
	0x04, 0x6e, 0xce, 0x52, 0x4d, 0x4b, 0x9c, 0xc9, 0x22, 0xc9, 0xa4, 0x66, 0x09, 0x55, 0x2a, 0xc9,
	0x85, 0x49, 0x32, 0x29, 0xde, 0xf3, 0x31, 0x56, 0x5a, 0x5a, 0x89, 0xd0, 0x99, 0x48, 0x33, 0x4c,
	0x95, 0xc2, 0xb9, 0x30, 0x9d, 0xdb, 0xbf, 0x81, 0x99, 0x2c, 0x0a, 0x29, 0x12, 0xc1, 0x6c, 0x42,
	0xf3, 0x5c, 0x33, 0x63, 0x3c, 0xdc, 0xb9, 0x7b, 0xb1, 0x30, 0x67, 0xc6, 0x72, 0x41, 0x2d, 0x97,
	0xc2, 0x8b, 0xb7, 0xbf, 0x85, 0x00, 0x2f, 0x68, 0xc1, 0x8e, 0x99, 0x9e, 0x31, 0x8d, 0x1e, 0xc3,
	0xc6, 0x62, 0x58, 0x3b, 0xe8, 0x06, 0xbd, 0x38, 0xbd, 0x8e, 0x57, 0xa2, 0xf8, 0x49, 0x58, 0x30,
	0x8b, 0xfb, 0x22, 0x57, 0x92, 0x0b, 0x4b, 0xce, 0xf4, 0xe8, 0x1d, 0x20, 0xa5, 0xb9, 0xd4, 0xdc,
	0xf2, 0x2f, 0x2c, 0x1f, 0xe5, 0xb2, 0xa0, 0x5c, 0xb4, 0xc3, 0x6e, 0xd4, 0x8b, 0xd3, 0x7b, 0xf8,
	0xfc, 0x03, 0xe1, 0xa5, 0x2d, 0x1e, 0x7a, 0xb0, 0x3c, 0x74, 0x10, 0xb9, 0xb2, 0x32, 0xc8, 0x5f,
	0x75, 0x72, 0x68, 0xfe, 0x2a, 0x42, 0xbb, 0x50, 0xb1, 0xa5, 0x62, 0x2e, 0x67, 0x33, 0xdd, 0x59,
	0xe7, 0xe0, 0x95, 0x47, 0xd4, 0x66, 0xa7, 0x5c, 0x8c, 0x5f, 0x95, 0x8a, 0x11, 0xc7, 0xa0, 0x2d,
	0xa8, 0xfd, 0xcc, 0x17, 0xf4, 0xea, 0x64, 0x71, 0xda, 0xfe, 0x5a, 0x81, 0xda, 0x81, 0x5b, 0x04,
	0xea, 0x43, 0xbc, 0x0c, 0x38, 0x6f, 0x23, 0xfa, 0x87, 0x36, 0xf6, 0xc3, 0x76, 0x40, 0x56, 0x39,
	0xb4, 0x07, 0xb1, 0xa0, 0x05, 0x1b, 0x19, 0x77, 0x6e, 0x57, 0xdd, 0x98, 0x6b, 0x7f, 0xae, 0x83,
	0x80, 0x58, 0x6e, 0x64, 0x0f, 0xaa, 0xcf, 0xa4, 0xb1, 0x66, 0xd1, 0xe4, 0xad, 0x75, 0xa8, 0x8f,
	0x8c, 0x9d, 0xae, 0x2f, 0xac, 0x2e, 0x5d, 0x0e, 0xcf, 0xa1, 0xab, 0x50, 0xcf, 0x26, 0x9c, 0x09,
	0x3b, 0xe2, 0xaa, 0x1d, 0x75, 0x83, 0x5e, 0x83, 0x6c, 0xfa, 0x8b, 0x81, 0x42, 0x03, 0x68, 0x18,
	0x4b, 0x2d, 0xcf, 0x46, 0xa7, 0xce, 0xa4, 0xe2, 0x4c, 0x76, 0xfe, 0x62, 0x72, 0x44, 0x95, 0xe2,
	0x62, 0x4c, 0x62, 0xcf, 0x3a, 0x9f, 0xce, 0x5b, 0x80, 0x65, 0x00, 0xd4, 0x82, 0xe8, 0x03, 0x2b,
	0xdd, 0x72, 0xea, 0x64, 0xfe, 0x13, 0x3d, 0x84, 0xea, 0x8c, 0x4e, 0xa6, 0xcc, 0x55, 0x1e, 0xa7,
	0x37, 0x2e, 0xa8, 0x72, 0x30, 0x7c, 0xa9, 0x17, 0xaf, 0x81, 0xd7, 0xef, 0x86, 0x8f, 0x82, 0xce,
	0x47, 0x88, 0x57, 0x8c, 0xff, 0xc7, 0xee, 0x51, 0x13, 0x42, 0x57, 0x50, 0xd4, 0x6b, 0x90, 0x90,
	0xab, 0x3b, 0x7d, 0x40, 0xe7, 0x67, 0xa0, 0x4d, 0xa8, 0x3c, 0x9d, 0x4e, 0x26, 0xad, 0x4b, 0xe8,
	0x32, 0xd4, 0x8f, 0xa7, 0x27, 0x1e, 0x6e, 0x05, 0x28, 0x86, 0x8d, 0xe7, 0xac, 0xfc, 0x24, 0x75,
	0xde, 0x0a, 0x51, 0x1d, 0xaa, 0x84, 0x8d, 0xd9, 0xe7, 0x56, 0xb4, 0xff, 0x00, 0xb6, 0x32, 0x59,
	0xac, 0x49, 0x38, 0x0c, 0xde, 0x44, 0xb9, 0x30, 0xdf, 0x43, 0xf4, 0x3a, 0x25, 0xb4, 0xc4, 0x07,
	0xf3, 0xff, 0x9e, 0x28, 0x85, 0x0f, 0x85, 0x39, 0xa9, 0xb9, 0xaf, 0xf3, 0xfe, 0x8f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x32, 0xd5, 0xea, 0x6d, 0x2e, 0x04, 0x00, 0x00,
}
