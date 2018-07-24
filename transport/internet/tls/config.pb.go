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
	return fileDescriptor_config_257f0631bcf0ff71, []int{0, 0}
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
	return fileDescriptor_config_257f0631bcf0ff71, []int{0}
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
	NextProtocol []string `protobuf:"bytes,4,rep,name=next_protocol,json=nextProtocol,proto3" json:"next_protocol,omitempty"`
	// Whether or not to disable session (ticket) resumption.
	DisableSessionResumption bool     `protobuf:"varint,6,opt,name=disable_session_resumption,json=disableSessionResumption,proto3" json:"disable_session_resumption,omitempty"`
	XXX_NoUnkeyedLiteral     struct{} `json:"-"`
	XXX_unrecognized         []byte   `json:"-"`
	XXX_sizecache            int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_257f0631bcf0ff71, []int{1}
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

func (m *Config) GetDisableSessionResumption() bool {
	if m != nil {
		return m.DisableSessionResumption
	}
	return false
}

func init() {
	proto.RegisterType((*Certificate)(nil), "v2ray.core.transport.internet.tls.Certificate")
	proto.RegisterType((*Config)(nil), "v2ray.core.transport.internet.tls.Config")
	proto.RegisterEnum("v2ray.core.transport.internet.tls.Certificate_Usage", Certificate_Usage_name, Certificate_Usage_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/tls/config.proto", fileDescriptor_config_257f0631bcf0ff71)
}

var fileDescriptor_config_257f0631bcf0ff71 = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x86, 0x49, 0x42, 0x2b, 0x76, 0xda, 0x8d, 0xc8, 0x4c, 0x28, 0xe2, 0x86, 0xac, 0x68, 0x52,
	0xaf, 0x1c, 0x29, 0xec, 0x92, 0x1b, 0x08, 0x41, 0x0b, 0x88, 0x52, 0xb9, 0xe9, 0xa4, 0x71, 0x13,
	0x79, 0xe6, 0x6c, 0x58, 0x4a, 0xec, 0xca, 0x76, 0x07, 0x7d, 0x25, 0x5e, 0x81, 0xc7, 0xe0, 0x85,
	0x50, 0x93, 0xb6, 0xb4, 0x57, 0x13, 0x77, 0x3e, 0xff, 0xf9, 0xce, 0xb1, 0xff, 0xdf, 0x90, 0xde,
	0xa7, 0x86, 0xaf, 0xa8, 0xd0, 0x4d, 0x22, 0xb4, 0xc1, 0xc4, 0x19, 0xae, 0xec, 0x42, 0x1b, 0x97,
	0x48, 0xe5, 0xd0, 0x28, 0x74, 0x89, 0xab, 0x6d, 0x22, 0xb4, 0xba, 0x95, 0x77, 0x74, 0x61, 0xb4,
	0xd3, 0xe4, 0x6c, 0x3b, 0x63, 0x90, 0xee, 0x78, 0xba, 0xe5, 0xa9, 0xab, 0xed, 0xe8, 0x8f, 0x07,
	0x83, 0x0c, 0x8d, 0x93, 0xb7, 0x52, 0x70, 0x87, 0x24, 0x3e, 0x28, 0x23, 0x2f, 0xf6, 0xc6, 0x43,
	0x76, 0x40, 0x84, 0x10, 0x7c, 0xc2, 0x55, 0xe4, 0xb7, 0x9d, 0xf5, 0x91, 0x7c, 0x84, 0xde, 0xd2,
	0xf2, 0x3b, 0x8c, 0x82, 0xd8, 0x1b, 0x9f, 0xa4, 0x17, 0xf4, 0xc1, 0x6b, 0xe9, 0xde, 0x42, 0x3a,
	0x5f, 0xcf, 0xb2, 0x6e, 0xc5, 0xe8, 0x3d, 0xf4, 0xda, 0x9a, 0x84, 0x30, 0xcc, 0x27, 0x59, 0x31,
	0xbd, 0xcc, 0xd9, 0xe7, 0x7c, 0x52, 0x86, 0x8f, 0xc8, 0x29, 0x84, 0x6f, 0xe7, 0xe5, 0xe5, 0x17,
	0x56, 0x94, 0xd7, 0xd5, 0x55, 0xce, 0x8a, 0x0f, 0xd7, 0xa1, 0x47, 0x9e, 0xc1, 0xd3, 0x7f, 0x6a,
	0x31, 0x9b, 0xcd, 0xf3, 0xd0, 0x1f, 0xfd, 0xf6, 0xa1, 0x9f, 0xb5, 0x49, 0x90, 0x73, 0x38, 0xe1,
	0x75, 0xad, 0x7f, 0x54, 0x52, 0x59, 0x14, 0x4b, 0xd3, 0x79, 0x7a, 0xc2, 0x8e, 0x5b, 0xb5, 0xd8,
	0x88, 0xe4, 0x02, 0x9e, 0x1f, 0x62, 0x95, 0x90, 0x8b, 0xef, 0x68, 0x6c, 0xd4, 0x6b, 0xf1, 0xd3,
	0x03, 0x3c, 0xeb, 0x7a, 0x64, 0x0a, 0x03, 0xb1, 0x97, 0x96, 0x1f, 0x07, 0xe3, 0x41, 0x4a, 0xff,
	0xcf, 0x3f, 0xdb, 0x5f, 0x41, 0x5e, 0xc2, 0xc0, 0xa2, 0xb9, 0x47, 0x53, 0x29, 0xde, 0x74, 0x89,
	0x1e, 0x31, 0xe8, 0xa4, 0x09, 0x6f, 0x90, 0xbc, 0x82, 0x63, 0x85, 0x3f, 0x5d, 0xd5, 0xfe, 0xb0,
	0xd0, 0x75, 0xf4, 0x38, 0x0e, 0xc6, 0x47, 0x6c, 0xb8, 0x16, 0xa7, 0x1b, 0x8d, 0xbc, 0x81, 0x17,
	0xdf, 0xa4, 0xe5, 0x37, 0x35, 0x56, 0x16, 0xad, 0x95, 0x5a, 0x55, 0x06, 0xed, 0xb2, 0x59, 0x38,
	0xa9, 0x55, 0xd4, 0x6f, 0x1d, 0x45, 0x1b, 0x62, 0xd6, 0x01, 0x6c, 0xd7, 0x7f, 0xc7, 0xe0, 0x5c,
	0xe8, 0xe6, 0x61, 0x17, 0x53, 0xef, 0x6b, 0xe0, 0x6a, 0xfb, 0xcb, 0x3f, 0xbb, 0x4a, 0x19, 0x5f,
	0xd1, 0x6c, 0x8d, 0x96, 0x3b, 0xb4, 0xd8, 0xa2, 0x65, 0x6d, 0x6f, 0xfa, 0xed, 0x7b, 0x5f, 0xff,
	0x0d, 0x00, 0x00, 0xff, 0xff, 0x80, 0x63, 0x12, 0xa7, 0xc7, 0x02, 0x00, 0x00,
}
