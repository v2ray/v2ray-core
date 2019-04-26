package tls

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	os "os"
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
	return fileDescriptor_42ed70cad60a2736, []int{0, 0}
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
	return fileDescriptor_42ed70cad60a2736, []int{0}
}

func (m *Certificate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Certificate.Unmarshal(m, b)
}
func (m *Certificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Certificate.Marshal(b, m, deterministic)
}
func (m *Certificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Certificate.Merge(m, src)
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
	DisableSessionResumption bool `protobuf:"varint,6,opt,name=disable_session_resumption,json=disableSessionResumption,proto3" json:"disable_session_resumption,omitempty"`
	// If true, root certificates on the system will not be loaded for verification.
	DisableSystemRoot    bool     `protobuf:"varint,7,opt,name=disable_system_root,json=disableSystemRoot,proto3" json:"disable_system_root,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_42ed70cad60a2736, []int{1}
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

func (m *Config) GetDisableSystemRoot() bool {
	if m != nil {
		return m.DisableSystemRoot
	}
	return false
}

func init() {
	proto.RegisterEnum("v2ray.core.transport.internet.tls.Certificate_Usage", Certificate_Usage_name, Certificate_Usage_value)
	proto.RegisterType((*Certificate)(nil), "v2ray.core.transport.internet.tls.Certificate")
	proto.RegisterType((*Config)(nil), "v2ray.core.transport.internet.tls.Config")
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/tls/config.proto", fileDescriptor_42ed70cad60a2736)
}

var fileDescriptor_42ed70cad60a2736 = []byte{
	// 435 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x86, 0x49, 0x43, 0x0b, 0x3b, 0xed, 0x46, 0xf0, 0x26, 0x14, 0x71, 0x43, 0x56, 0x34, 0xa9,
	0x57, 0x8e, 0x14, 0x76, 0xc9, 0x0d, 0x84, 0xa0, 0x05, 0x44, 0xa9, 0xdc, 0x74, 0xd2, 0xb8, 0x89,
	0x32, 0x73, 0x36, 0x2c, 0x25, 0x76, 0x65, 0xbb, 0x83, 0xbe, 0x12, 0xaf, 0xc3, 0x63, 0xf0, 0x12,
	0xa8, 0x4e, 0x5b, 0xda, 0xab, 0x69, 0x77, 0x39, 0xff, 0xff, 0xfd, 0x27, 0xfa, 0x8f, 0x21, 0xb9,
	0x4b, 0x74, 0xb5, 0xa4, 0x5c, 0x35, 0x31, 0x57, 0x1a, 0x63, 0xab, 0x2b, 0x69, 0xe6, 0x4a, 0xdb,
	0x58, 0x48, 0x8b, 0x5a, 0xa2, 0x8d, 0x6d, 0x6d, 0x62, 0xae, 0xe4, 0x8d, 0xb8, 0xa5, 0x73, 0xad,
	0xac, 0x22, 0xa7, 0x9b, 0x8c, 0x46, 0xba, 0xe5, 0xe9, 0x86, 0xa7, 0xb6, 0x36, 0xc3, 0x3f, 0x1e,
	0xf4, 0x53, 0xd4, 0x56, 0xdc, 0x08, 0x5e, 0x59, 0x24, 0xd1, 0xde, 0x18, 0x7a, 0x91, 0x37, 0x1a,
	0xb0, 0x3d, 0x22, 0x00, 0xff, 0x33, 0x2e, 0xc3, 0x8e, 0x73, 0x56, 0x9f, 0xe4, 0x13, 0x74, 0x17,
	0xa6, 0xba, 0xc5, 0xd0, 0x8f, 0xbc, 0xd1, 0x51, 0x72, 0x4e, 0xef, 0xfd, 0x2d, 0xdd, 0x59, 0x48,
	0x67, 0xab, 0x2c, 0x6b, 0x57, 0x0c, 0x3f, 0x40, 0xd7, 0xcd, 0x24, 0x80, 0x41, 0x36, 0x4e, 0xf3,
	0xc9, 0x45, 0xc6, 0xbe, 0x64, 0xe3, 0x22, 0x78, 0x44, 0x4e, 0x20, 0x78, 0x37, 0x2b, 0x2e, 0xbe,
	0xb2, 0xbc, 0xb8, 0x2a, 0x2f, 0x33, 0x96, 0x7f, 0xbc, 0x0a, 0x3c, 0x72, 0x0c, 0xcf, 0xfe, 0xab,
	0xf9, 0x74, 0x3a, 0xcb, 0x82, 0xce, 0xf0, 0x6f, 0x07, 0x7a, 0xa9, 0xbb, 0x04, 0x39, 0x83, 0xa3,
	0xaa, 0xae, 0xd5, 0xcf, 0x52, 0x48, 0x83, 0x7c, 0xa1, 0xdb, 0x4e, 0x4f, 0xd9, 0xa1, 0x53, 0xf3,
	0xb5, 0x48, 0xce, 0xe1, 0xc5, 0x3e, 0x56, 0x72, 0x31, 0xff, 0x81, 0xda, 0x84, 0x5d, 0x87, 0x9f,
	0xec, 0xe1, 0x69, 0xeb, 0x91, 0x09, 0xf4, 0xf9, 0xce, 0xb5, 0x3a, 0x91, 0x3f, 0xea, 0x27, 0xf4,
	0x61, 0xfd, 0xd9, 0xee, 0x0a, 0xf2, 0x0a, 0xfa, 0x06, 0xf5, 0x1d, 0xea, 0x52, 0x56, 0x4d, 0x7b,
	0xd1, 0x03, 0x06, 0xad, 0x34, 0xae, 0x1a, 0x24, 0xaf, 0xe1, 0x50, 0xe2, 0x2f, 0x5b, 0xba, 0x17,
	0xe6, 0xaa, 0x0e, 0x1f, 0x47, 0xfe, 0xe8, 0x80, 0x0d, 0x56, 0xe2, 0x64, 0xad, 0x91, 0xb7, 0xf0,
	0xf2, 0xbb, 0x30, 0xd5, 0x75, 0x8d, 0xa5, 0x41, 0x63, 0x84, 0x92, 0xa5, 0x46, 0xb3, 0x68, 0xe6,
	0x56, 0x28, 0x19, 0xf6, 0x5c, 0xa3, 0x70, 0x4d, 0x4c, 0x5b, 0x80, 0x6d, 0x7d, 0x42, 0xe1, 0x78,
	0x9b, 0x5e, 0x1a, 0x8b, 0x4d, 0xa9, 0x95, 0xb2, 0xe1, 0x13, 0x17, 0x7b, 0xbe, 0x89, 0x39, 0x87,
	0x29, 0x65, 0xdf, 0x33, 0x38, 0xe3, 0xaa, 0xb9, 0xbf, 0xf5, 0xc4, 0xfb, 0xe6, 0xdb, 0xda, 0xfc,
	0xee, 0x9c, 0x5e, 0x26, 0xac, 0x5a, 0xd2, 0x74, 0x85, 0x16, 0x5b, 0x34, 0xdf, 0xa0, 0x45, 0x6d,
	0xae, 0x7b, 0xae, 0xdf, 0x9b, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb6, 0x58, 0x74, 0x54, 0xf7,
	0x02, 0x00, 0x00,
}
