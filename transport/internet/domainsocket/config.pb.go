package domainsocket

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

// Flag Array
type DomainSocketSecurityMode int32

const (
	DomainSocketSecurityMode_Danger DomainSocketSecurityMode = 0
	// Verfify is Dialer have a pid or ppid match pid file
	DomainSocketSecurityMode_VerifyPid DomainSocketSecurityMode = 1
	// Do not tolerance ppid match
	DomainSocketSecurityMode_DisqualifyPPIDMatch DomainSocketSecurityMode = 2
	// Enforce Uid Verify On euid
	DomainSocketSecurityMode_VerifyEUID DomainSocketSecurityMode = 4
	// Enforce Uid Verify On ruid
	DomainSocketSecurityMode_VerifyRUID DomainSocketSecurityMode = 8
	// Does not allow same user exception
	DomainSocketSecurityMode_DisqualifySameUser DomainSocketSecurityMode = 16
	// Does not allow root user exception
	DomainSocketSecurityMode_DisqualifyRootUser DomainSocketSecurityMode = 32
)

var DomainSocketSecurityMode_name = map[int32]string{
	0:  "Danger",
	1:  "VerifyPid",
	2:  "DisqualifyPPIDMatch",
	4:  "VerifyEUID",
	8:  "VerifyRUID",
	16: "DisqualifySameUser",
	32: "DisqualifyRootUser",
}
var DomainSocketSecurityMode_value = map[string]int32{
	"Danger":              0,
	"VerifyPid":           1,
	"DisqualifyPPIDMatch": 2,
	"VerifyEUID":          4,
	"VerifyRUID":          8,
	"DisqualifySameUser":  16,
	"DisqualifyRootUser":  32,
}

func (x DomainSocketSecurityMode) String() string {
	return proto.EnumName(DomainSocketSecurityMode_name, int32(x))
}
func (DomainSocketSecurityMode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type DomainSocketSecurity struct {
	// Flag Array Type, User can set an integer to define various configure
	Mode DomainSocketSecurityMode `protobuf:"varint,1,opt,name=Mode,enum=v2ray.core.internet.domainsocket.DomainSocketSecurityMode" json:"Mode,omitempty"`
	// Set pid files to be allowed
	AllowedPid []string `protobuf:"bytes,2,rep,name=AllowedPid" json:"AllowedPid,omitempty"`
	// Set uids to be allowed, either euid or ruid should match one of following
	// uids AllowedUid, or user that v2ray is running or root.
	AllowedUid []uint64 `protobuf:"varint,3,rep,packed,name=AllowedUid" json:"AllowedUid,omitempty"`
}

func (m *DomainSocketSecurity) Reset()                    { *m = DomainSocketSecurity{} }
func (m *DomainSocketSecurity) String() string            { return proto.CompactTextString(m) }
func (*DomainSocketSecurity) ProtoMessage()               {}
func (*DomainSocketSecurity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DomainSocketSecurity) GetMode() DomainSocketSecurityMode {
	if m != nil {
		return m.Mode
	}
	return DomainSocketSecurityMode_Danger
}

func (m *DomainSocketSecurity) GetAllowedPid() []string {
	if m != nil {
		return m.AllowedPid
	}
	return nil
}

func (m *DomainSocketSecurity) GetAllowedUid() []uint64 {
	if m != nil {
		return m.AllowedUid
	}
	return nil
}

type DomainSocketSettings struct {
	// Path we should listen/dial
	Path     string                `protobuf:"bytes,1,opt,name=Path" json:"Path,omitempty"`
	Security *DomainSocketSecurity `protobuf:"bytes,2,opt,name=Security" json:"Security,omitempty"`
}

func (m *DomainSocketSettings) Reset()                    { *m = DomainSocketSettings{} }
func (m *DomainSocketSettings) String() string            { return proto.CompactTextString(m) }
func (*DomainSocketSettings) ProtoMessage()               {}
func (*DomainSocketSettings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DomainSocketSettings) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *DomainSocketSettings) GetSecurity() *DomainSocketSecurity {
	if m != nil {
		return m.Security
	}
	return nil
}

func init() {
	proto.RegisterType((*DomainSocketSecurity)(nil), "v2ray.core.internet.domainsocket.DomainSocketSecurity")
	proto.RegisterType((*DomainSocketSettings)(nil), "v2ray.core.internet.domainsocket.DomainSocketSettings")
	proto.RegisterEnum("v2ray.core.internet.domainsocket.DomainSocketSecurityMode", DomainSocketSecurityMode_name, DomainSocketSecurityMode_value)
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/domainsocket/config.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x41, 0x4b, 0xf3, 0x40,
	0x10, 0x86, 0xbf, 0xb4, 0xa1, 0xb4, 0xf3, 0x69, 0x09, 0xab, 0x68, 0x4e, 0x12, 0x7a, 0x0a, 0x1e,
	0x12, 0xa8, 0xe0, 0x41, 0x4f, 0x4a, 0x3c, 0xf4, 0x50, 0x09, 0x5b, 0xea, 0xc1, 0xdb, 0x9a, 0xdd,
	0xb6, 0x8b, 0xed, 0x4e, 0xdd, 0x4c, 0x95, 0x5e, 0xfc, 0x29, 0x1e, 0xfd, 0x9d, 0xd2, 0x95, 0xd8,
	0x28, 0x8a, 0xe0, 0x2d, 0xf3, 0xe4, 0x9d, 0x87, 0x77, 0x60, 0xe1, 0xfc, 0xb1, 0x6f, 0xc5, 0x3a,
	0x29, 0x70, 0x91, 0x16, 0x68, 0x55, 0x4a, 0x56, 0x98, 0x72, 0x89, 0x96, 0x52, 0x6d, 0x48, 0x59,
	0xa3, 0x28, 0x95, 0xb8, 0x10, 0xda, 0x94, 0x58, 0xdc, 0x2b, 0x4a, 0x0b, 0x34, 0x13, 0x3d, 0x4d,
	0x96, 0x16, 0x09, 0x59, 0x54, 0x2d, 0x5b, 0x95, 0x54, 0xf1, 0xa4, 0x1e, 0xef, 0xbd, 0x7a, 0xb0,
	0x9f, 0x39, 0x30, 0x72, 0x60, 0xa4, 0x8a, 0x95, 0xd5, 0xb4, 0x66, 0xd7, 0xe0, 0x0f, 0x51, 0xaa,
	0xd0, 0x8b, 0xbc, 0xb8, 0xdb, 0x3f, 0x4b, 0x7e, 0x33, 0x25, 0xdf, 0x59, 0x36, 0x06, 0xee, 0x3c,
	0xec, 0x08, 0xe0, 0x62, 0x3e, 0xc7, 0x27, 0x25, 0x73, 0x2d, 0xc3, 0x46, 0xd4, 0x8c, 0x3b, 0xbc,
	0x46, 0x6a, 0xff, 0xc7, 0x5a, 0x86, 0xcd, 0xa8, 0x19, 0xfb, 0xbc, 0x46, 0x7a, 0xcf, 0x5f, 0x7b,
	0x12, 0x69, 0x33, 0x2d, 0x19, 0x03, 0x3f, 0x17, 0x34, 0x73, 0x3d, 0x3b, 0xdc, 0x7d, 0x33, 0x0e,
	0xed, 0xaa, 0x41, 0xd8, 0x88, 0xbc, 0xf8, 0x7f, 0xff, 0xf4, 0x6f, 0xfd, 0xf9, 0x87, 0xe7, 0xf8,
	0xc5, 0x83, 0xf0, 0xa7, 0x13, 0x19, 0x40, 0x2b, 0x13, 0x66, 0xaa, 0x6c, 0xf0, 0x8f, 0xed, 0x42,
	0xe7, 0x46, 0x59, 0x3d, 0x59, 0xe7, 0x5a, 0x06, 0x1e, 0x3b, 0x84, 0xbd, 0x4c, 0x97, 0x0f, 0x2b,
	0x31, 0xdf, 0xa0, 0x7c, 0x90, 0x0d, 0x05, 0x15, 0xb3, 0xa0, 0xc1, 0xba, 0x00, 0xef, 0xb9, 0xab,
	0xf1, 0x20, 0x0b, 0xfc, 0xed, 0xcc, 0x37, 0x73, 0x9b, 0x1d, 0x00, 0xdb, 0x2e, 0x8e, 0xc4, 0x42,
	0x8d, 0x4b, 0x65, 0x83, 0xe0, 0x33, 0xe7, 0x88, 0xe4, 0x78, 0x74, 0xd9, 0xbd, 0xdd, 0xa9, 0xdf,
	0x73, 0xd7, 0x72, 0x4f, 0xe0, 0xe4, 0x2d, 0x00, 0x00, 0xff, 0xff, 0x30, 0xd5, 0x61, 0x90, 0x41,
	0x02, 0x00, 0x00,
}
