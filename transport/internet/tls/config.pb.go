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
	return fileDescriptor_config_094edd58cf826f66, []int{0, 0}
}

type Certificate struct {
	// TLS certificate in x509 format.
	Certificate []byte `protobuf:"bytes,1,opt,name=Certificate,proto3" json:"Certificate,omitempty"`
	// TLS key in x509 format.
	Key                  []byte            `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Usage                Certificate_Usage `protobuf:"varint,3,opt,name=usage,enum=v2ray.core.transport.internet.tls.Certificate_Usage" json:"usage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Certificate) Reset()         { *m = Certificate{} }
func (m *Certificate) String() string { return proto.CompactTextString(m) }
func (*Certificate) ProtoMessage()    {}
func (*Certificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_094edd58cf826f66, []int{0}
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
	AllowInsecure bool `protobuf:"varint,1,opt,name=allow_insecure,json=allowInsecure" json:"allow_insecure,omitempty"`
	// List of certificates to be served on server.
	Certificate []*Certificate `protobuf:"bytes,2,rep,name=certificate" json:"certificate,omitempty"`
	// Override server name.
	ServerName string `protobuf:"bytes,3,opt,name=server_name,json=serverName" json:"server_name,omitempty"`
	// Lists of string as ALPN values.
	NextProtocol         []string `protobuf:"bytes,4,rep,name=next_protocol,json=nextProtocol" json:"next_protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_094edd58cf826f66, []int{1}
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
	proto.RegisterFile("v2ray.com/core/transport/internet/tls/config.proto", fileDescriptor_config_094edd58cf826f66)
}

var fileDescriptor_config_094edd58cf826f66 = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xd1, 0x6e, 0xda, 0x30,
	0x14, 0x40, 0x97, 0x64, 0xa0, 0xe1, 0x00, 0x8b, 0xbc, 0x3d, 0xe4, 0x6d, 0x81, 0x09, 0x29, 0x4f,
	0x8e, 0x94, 0xed, 0x07, 0xb6, 0x34, 0x15, 0x69, 0x55, 0x1a, 0x99, 0x80, 0x44, 0x5f, 0x22, 0xd7,
	0x32, 0x28, 0x52, 0x12, 0x23, 0xdb, 0xd0, 0xf2, 0x4b, 0xfd, 0x91, 0x7e, 0x40, 0x7f, 0xa8, 0x4a,
	0x02, 0x14, 0x9e, 0x50, 0xdf, 0x7c, 0x8f, 0xcf, 0xbd, 0xd7, 0xf7, 0x1a, 0xf8, 0x5b, 0x5f, 0x90,
	0x1d, 0xa2, 0xbc, 0xf0, 0x28, 0x17, 0xcc, 0x53, 0x82, 0x94, 0x72, 0xcd, 0x85, 0xf2, 0xb2, 0x52,
	0x31, 0x51, 0x32, 0xe5, 0xa9, 0x5c, 0x7a, 0x94, 0x97, 0xcb, 0x6c, 0x85, 0xd6, 0x82, 0x2b, 0x0e,
	0x07, 0x87, 0x1c, 0xc1, 0xd0, 0xd1, 0x47, 0x07, 0x1f, 0xa9, 0x5c, 0x0e, 0xdf, 0x34, 0x60, 0x06,
	0x4c, 0xa8, 0x6c, 0x99, 0x51, 0xa2, 0x18, 0x74, 0xce, 0x42, 0x5b, 0x73, 0x34, 0xb7, 0x8b, 0xcf,
	0x0c, 0x0b, 0x18, 0xb7, 0x6c, 0x67, 0xeb, 0xf5, 0x4d, 0x75, 0x84, 0x37, 0xa0, 0xb5, 0x91, 0x64,
	0xc5, 0x6c, 0xc3, 0xd1, 0xdc, 0xbe, 0xff, 0x17, 0x5d, 0x6c, 0x8b, 0x4e, 0x0a, 0xa2, 0x59, 0x95,
	0x8b, 0x9b, 0x12, 0xc3, 0x2b, 0xd0, 0xaa, 0x63, 0x68, 0x81, 0x6e, 0x38, 0x09, 0xa2, 0x78, 0x1c,
	0xe2, 0xbb, 0x70, 0x92, 0x58, 0x5f, 0xe0, 0x4f, 0x60, 0xfd, 0x9b, 0x25, 0xe3, 0x7b, 0x1c, 0x25,
	0x8b, 0x74, 0x1e, 0xe2, 0xe8, 0x7a, 0x61, 0x69, 0xf0, 0x07, 0xf8, 0xfe, 0x41, 0xa3, 0xe9, 0x74,
	0x16, 0x5a, 0xfa, 0xf0, 0x55, 0x03, 0xed, 0xa0, 0xde, 0x04, 0x1c, 0x81, 0x3e, 0xc9, 0x73, 0xfe,
	0x94, 0x66, 0xa5, 0x64, 0x74, 0x23, 0x9a, 0x99, 0xbe, 0xe1, 0x5e, 0x4d, 0xa3, 0x3d, 0x84, 0x31,
	0x30, 0xe9, 0xc9, 0xdc, 0xba, 0x63, 0xb8, 0xa6, 0x8f, 0x3e, 0x37, 0x09, 0x3e, 0x2d, 0x01, 0x7f,
	0x01, 0x53, 0x32, 0xb1, 0x65, 0x22, 0x2d, 0x49, 0xd1, 0xec, 0xa6, 0x83, 0x41, 0x83, 0x26, 0xa4,
	0x60, 0xf0, 0x37, 0xe8, 0x95, 0xec, 0x59, 0xa5, 0xf5, 0x5f, 0x51, 0x9e, 0xdb, 0x5f, 0x1d, 0xc3,
	0xed, 0xe0, 0x6e, 0x05, 0xe3, 0x3d, 0xfb, 0x8f, 0xc1, 0x88, 0xf2, 0xe2, 0xf2, 0x3b, 0x62, 0xed,
	0xc1, 0x50, 0xb9, 0x7c, 0xd1, 0x07, 0x73, 0x1f, 0x93, 0x1d, 0x0a, 0x2a, 0x35, 0x39, 0xaa, 0xd1,
	0x41, 0x4d, 0x72, 0xf9, 0xd8, 0xae, 0x3b, 0xfe, 0x79, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xcd,
	0x2a, 0x68, 0x53, 0x02, 0x00, 0x00,
}
