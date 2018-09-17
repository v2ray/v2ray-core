package dispatcher

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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

type SessionConfig struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionConfig) Reset()         { *m = SessionConfig{} }
func (m *SessionConfig) String() string { return proto.CompactTextString(m) }
func (*SessionConfig) ProtoMessage()    {}
func (*SessionConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_90b7c41cce355532, []int{0}
}
func (m *SessionConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionConfig.Unmarshal(m, b)
}
func (m *SessionConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionConfig.Marshal(b, m, deterministic)
}
func (m *SessionConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionConfig.Merge(m, src)
}
func (m *SessionConfig) XXX_Size() int {
	return xxx_messageInfo_SessionConfig.Size(m)
}
func (m *SessionConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SessionConfig proto.InternalMessageInfo

type Config struct {
	Settings             *SessionConfig `protobuf:"bytes,1,opt,name=settings,proto3" json:"settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_90b7c41cce355532, []int{1}
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

func (m *Config) GetSettings() *SessionConfig {
	if m != nil {
		return m.Settings
	}
	return nil
}

func init() {
	proto.RegisterType((*SessionConfig)(nil), "v2ray.core.app.dispatcher.SessionConfig")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.dispatcher.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/dispatcher/config.proto", fileDescriptor_90b7c41cce355532)
}

var fileDescriptor_90b7c41cce355532 = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2a, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0x2c, 0x28, 0xd0, 0x4f,
	0xc9, 0x2c, 0x2e, 0x48, 0x2c, 0x49, 0xce, 0x48, 0x2d, 0xd2, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x84, 0xa9, 0x2d, 0x4a, 0xd5, 0x4b, 0x2c, 0x28,
	0xd0, 0x43, 0xa8, 0x53, 0x12, 0xe5, 0xe2, 0x0d, 0x4e, 0x2d, 0x2e, 0xce, 0xcc, 0xcf, 0x73, 0x06,
	0xeb, 0xf0, 0x62, 0xe1, 0x60, 0x14, 0x60, 0x52, 0xf2, 0xe3, 0x62, 0x83, 0xf0, 0x85, 0x5c, 0xb8,
	0x38, 0x8a, 0x53, 0x4b, 0x4a, 0x32, 0xf3, 0xd2, 0x8b, 0x25, 0x18, 0x15, 0x18, 0x35, 0xb8, 0x8d,
	0x34, 0xf4, 0x70, 0x1a, 0xa7, 0x87, 0x62, 0x56, 0x10, 0x5c, 0xa7, 0x93, 0x27, 0x97, 0x6c, 0x72,
	0x7e, 0x2e, 0x6e, 0x8d, 0x01, 0x8c, 0x51, 0x5c, 0x08, 0xde, 0x2a, 0x26, 0xc9, 0x30, 0xa3, 0xa0,
	0xc4, 0x4a, 0x3d, 0x67, 0x90, 0x4a, 0xc7, 0x82, 0x02, 0x3d, 0x17, 0xb8, 0x5c, 0x12, 0x1b, 0xd8,
	0x4f, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x47, 0xc5, 0xd6, 0x01, 0x01, 0x00, 0x00,
}
