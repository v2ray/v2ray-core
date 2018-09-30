package websocket

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

type Header struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_c4869c9c0fc9b72f, []int{0}
}

func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Header) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Config struct {
	// URL path to the WebSocket service. Empty value means root(/).
	Path                 string    `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Header               []*Header `protobuf:"bytes,3,rep,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_c4869c9c0fc9b72f, []int{1}
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

func (m *Config) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Config) GetHeader() []*Header {
	if m != nil {
		return m.Header
	}
	return nil
}

func init() {
	proto.RegisterType((*Header)(nil), "v2ray.core.transport.internet.websocket.Header")
	proto.RegisterType((*Config)(nil), "v2ray.core.transport.internet.websocket.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/websocket/config.proto", fileDescriptor_c4869c9c0fc9b72f)
}

var fileDescriptor_c4869c9c0fc9b72f = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0x28, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x2f, 0x29, 0x4a, 0xcc, 0x2b,
	0x2e, 0xc8, 0x2f, 0x2a, 0xd1, 0xcf, 0xcc, 0x2b, 0x49, 0x2d, 0xca, 0x4b, 0x2d, 0xd1, 0x2f, 0x4f,
	0x4d, 0x2a, 0xce, 0x4f, 0xce, 0x4e, 0x2d, 0xd1, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x52, 0x87, 0xe9, 0x2c, 0x4a, 0xd5, 0x83, 0xeb, 0xd2, 0x83, 0xe9,
	0xd2, 0x83, 0xeb, 0x52, 0x32, 0xe0, 0x62, 0xf3, 0x48, 0x4d, 0x4c, 0x49, 0x2d, 0x12, 0x12, 0xe0,
	0x62, 0xce, 0x4e, 0xad, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0x31, 0x85, 0x44, 0xb8,
	0x58, 0xcb, 0x12, 0x73, 0x4a, 0x53, 0x25, 0x98, 0xc0, 0x62, 0x10, 0x8e, 0x52, 0x36, 0x17, 0x9b,
	0x33, 0xd8, 0x2a, 0x21, 0x21, 0x2e, 0x96, 0x82, 0xc4, 0x92, 0x0c, 0xa8, 0x34, 0x98, 0x2d, 0xe4,
	0xce, 0xc5, 0x96, 0x01, 0x36, 0x4f, 0x82, 0x59, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x5f, 0x8f, 0x48,
	0x97, 0xe8, 0x41, 0x9c, 0x11, 0x04, 0xd5, 0xee, 0xc5, 0xc2, 0xc1, 0x28, 0xc0, 0xe4, 0x94, 0xc2,
	0xa5, 0x9d, 0x9c, 0x9f, 0x4b, 0xac, 0x19, 0x01, 0x8c, 0x51, 0x9c, 0x70, 0xce, 0x2a, 0x26, 0xf5,
	0x30, 0xa3, 0xa0, 0xc4, 0x4a, 0x3d, 0x67, 0x90, 0xb6, 0x10, 0xb8, 0x36, 0x4f, 0x98, 0xb6, 0x70,
	0x98, 0xca, 0x24, 0x36, 0x70, 0xa0, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xf5, 0x7e, 0x60,
	0xf9, 0x70, 0x01, 0x00, 0x00,
}
