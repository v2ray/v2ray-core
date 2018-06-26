package tls

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Certificate_Usage int32

const (
	Certificate_ENCIPHERMENT     Certificate_Usage = 0
	Certificate_AUTHORITY_VERIFY Certificate_Usage = 1
	Certificate_AUTHORITY_ISSUE  Certificate_Usage = 2
)

var Certificate_Usage_name = map[int32]string{
	0: "ENCIPHERMENT",
	1: "AUTHORITY_VERIFY",
	2: "AUTHORITY_ISSUE",
}
var Certificate_Usage_value = map[string]int32{
	"ENCIPHERMENT":     0,
	"AUTHORITY_VERIFY": 1,
	"AUTHORITY_ISSUE":  2,
}

func (x Certificate_Usage) String() string {
	return proto.EnumName(Certificate_Usage_name, int32(x))
}
func (Certificate_Usage) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_e641ddfb53ec9d25, []int{0, 0}
}

type Certificate struct {
	// TLS certificate in x509 format.
	Certificate []byte `protobuf:"bytes,1,opt,name=Certificate,proto3" json:"Certificate,omitempty"`
	// TLS key in x509 format.
	Key                  []byte            `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Usage                Certificate_Usage `protobuf:"varint,3,opt,name=usage,proto3,enum=v2ray.core.transport.internet.tls.Certificate_Usage" json:"usage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Certificate) Reset()         { *m = Certificate{} }
func (m *Certificate) String() string { return proto.CompactTextString(m) }
func (*Certificate) ProtoMessage()    {}
func (*Certificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_e641ddfb53ec9d25, []int{0}
}
func (m *Certificate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Certificate.Unmarshal(m, b)
}
func (m *Certificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Certificate.Marshal(b, m, deterministic)
}
func (dst *Certificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Certificate.Merge(dst, src)
}
func (m *Certificate) XXX_Size() int {
	return xxx_messageInfo_Certificate.Size(m)
}
func (m *Certificate) XXX_DiscardUnknown() {
	xxx_messageInfo_Certificate.DiscardUnknown(m)
}

var xxx_messageInfo_Certificate proto.InternalMessageInfo

func (m *Certificate) GetCertificate() []byte {
	if m != nil {
		return m.Certificate
	}
	return nil
}

func (m *Certificate) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Certificate) GetUsage() Certificate_Usage {
	if m != nil {
		return m.Usage
	}
	return Certificate_ENCIPHERMENT
}

type Config struct {
	// Whether or not to allow self-signed certificates.
	AllowInsecure bool `protobuf:"varint,1,opt,name=allow_insecure,json=allowInsecure,proto3" json:"allow_insecure,omitempty"`
	// Whether or not to allow insecure cipher suites.
	AllowInsecureCiphers bool `protobuf:"varint,5,opt,name=allow_insecure_ciphers,json=allowInsecureCiphers,proto3" json:"allow_insecure_ciphers,omitempty"`
	// List of certificates to be served on server.
	Certificate []*Certificate `protobuf:"bytes,2,rep,name=certificate,proto3" json:"certificate,omitempty"`
	// Override server name.
	ServerName string `protobuf:"bytes,3,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`
	// Lists of string as ALPN values.
	NextProtocol         []string `protobuf:"bytes,4,rep,name=next_protocol,json=nextProtocol,proto3" json:"next_protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_e641ddfb53ec9d25, []int{1}
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

func (m *Config) GetAllowInsecure() bool {
	if m != nil {
		return m.AllowInsecure
	}
	return false
}

func (m *Config) GetAllowInsecureCiphers() bool {
	if m != nil {
		return m.AllowInsecureCiphers
	}
	return false
}

func (m *Config) GetCertificate() []*Certificate {
	if m != nil {
		return m.Certificate
	}
	return nil
}

func (m *Config) GetServerName() string {
	if m != nil {
		return m.ServerName
	}
	return ""
}

func (m *Config) GetNextProtocol() []string {
	if m != nil {
		return m.NextProtocol
	}
	return nil
}

func init() {
	proto.RegisterType((*Certificate)(nil), "v2ray.core.transport.internet.tls.Certificate")
	proto.RegisterType((*Config)(nil), "v2ray.core.transport.internet.tls.Config")
	proto.RegisterEnum("v2ray.core.transport.internet.tls.Certificate_Usage", Certificate_Usage_name, Certificate_Usage_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/tls/config.proto", fileDescriptor_config_e641ddfb53ec9d25)
}

var fileDescriptor_config_e641ddfb53ec9d25 = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x51, 0x8f, 0x93, 0x40,
	0x10, 0x80, 0x05, 0xec, 0xc5, 0x1b, 0x7a, 0x27, 0x59, 0x2f, 0x86, 0x37, 0xb9, 0x9a, 0x26, 0x3c,
	0x2d, 0x09, 0xf6, 0x0f, 0x28, 0x62, 0x8a, 0xc6, 0x4a, 0xb6, 0xb4, 0x49, 0x7d, 0x21, 0xeb, 0x66,
	0x5b, 0x49, 0x80, 0x6d, 0x76, 0xb7, 0xd5, 0xfe, 0x25, 0xff, 0x8e, 0xbf, 0xc7, 0xc4, 0x00, 0x6d,
	0x2d, 0x4f, 0x8d, 0x6f, 0xcc, 0x37, 0xdf, 0xcc, 0x30, 0xb3, 0x10, 0xee, 0x43, 0x49, 0x0f, 0x98,
	0x89, 0x2a, 0x60, 0x42, 0xf2, 0x40, 0x4b, 0x5a, 0xab, 0xad, 0x90, 0x3a, 0x28, 0x6a, 0xcd, 0x65,
	0xcd, 0x75, 0xa0, 0x4b, 0x15, 0x30, 0x51, 0xaf, 0x8b, 0x0d, 0xde, 0x4a, 0xa1, 0x05, 0x7a, 0x3c,
	0xd5, 0x48, 0x8e, 0xcf, 0x3e, 0x3e, 0xf9, 0x58, 0x97, 0x6a, 0xf4, 0xdb, 0x00, 0x3b, 0xe2, 0x52,
	0x17, 0xeb, 0x82, 0x51, 0xcd, 0x91, 0xd7, 0x0b, 0x5d, 0xc3, 0x33, 0xfc, 0x21, 0xe9, 0x19, 0x0e,
	0x58, 0x9f, 0xf8, 0xc1, 0x35, 0xdb, 0x4c, 0xf3, 0x89, 0x3e, 0xc2, 0x60, 0xa7, 0xe8, 0x86, 0xbb,
	0x96, 0x67, 0xf8, 0xf7, 0xe1, 0x04, 0x5f, 0x1d, 0x8b, 0x2f, 0x1a, 0xe2, 0x45, 0x53, 0x4b, 0xba,
	0x16, 0xa3, 0xf7, 0x30, 0x68, 0x63, 0xe4, 0xc0, 0x30, 0x9e, 0x45, 0x49, 0x3a, 0x8d, 0xc9, 0xe7,
	0x78, 0x96, 0x39, 0x4f, 0xd0, 0x03, 0x38, 0x6f, 0x17, 0xd9, 0xf4, 0x0b, 0x49, 0xb2, 0x55, 0xbe,
	0x8c, 0x49, 0xf2, 0x61, 0xe5, 0x18, 0xe8, 0x05, 0x3c, 0xff, 0x47, 0x93, 0xf9, 0x7c, 0x11, 0x3b,
	0xe6, 0xe8, 0x8f, 0x01, 0x37, 0x51, 0x7b, 0x09, 0x34, 0x86, 0x7b, 0x5a, 0x96, 0xe2, 0x47, 0x5e,
	0xd4, 0x8a, 0xb3, 0x9d, 0xec, 0x76, 0x7a, 0x46, 0xee, 0x5a, 0x9a, 0x1c, 0x21, 0x9a, 0xc0, 0xcb,
	0xbe, 0x96, 0xb3, 0x62, 0xfb, 0x9d, 0x4b, 0xe5, 0x0e, 0x5a, 0xfd, 0xa1, 0xa7, 0x47, 0x5d, 0x0e,
	0xa5, 0x60, 0xb3, 0x8b, 0x6b, 0x99, 0x9e, 0xe5, 0xdb, 0x21, 0xfe, 0xbf, 0xfd, 0xc9, 0x65, 0x0b,
	0xf4, 0x0a, 0x6c, 0xc5, 0xe5, 0x9e, 0xcb, 0xbc, 0xa6, 0x55, 0x77, 0xd1, 0x5b, 0x02, 0x1d, 0x9a,
	0xd1, 0x8a, 0xa3, 0xd7, 0x70, 0x57, 0xf3, 0x9f, 0x3a, 0x6f, 0x5f, 0x98, 0x89, 0xd2, 0x7d, 0xea,
	0x59, 0xfe, 0x2d, 0x19, 0x36, 0x30, 0x3d, 0xb2, 0x77, 0x04, 0xc6, 0x4c, 0x54, 0xd7, 0xff, 0x23,
	0x35, 0xbe, 0x5a, 0xba, 0x54, 0xbf, 0xcc, 0xc7, 0x65, 0x48, 0xe8, 0x01, 0x47, 0x8d, 0x9a, 0x9d,
	0xd5, 0xe4, 0xa4, 0x66, 0xa5, 0xfa, 0x76, 0xd3, 0x4e, 0x7c, 0xf3, 0x37, 0x00, 0x00, 0xff, 0xff,
	0xd1, 0x3b, 0xdd, 0x37, 0x89, 0x02, 0x00, 0x00,
}
