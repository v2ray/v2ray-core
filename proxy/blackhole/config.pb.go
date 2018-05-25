package blackhole

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import serial "v2ray.com/core/common/serial"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NoneResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NoneResponse) Reset()         { *m = NoneResponse{} }
func (m *NoneResponse) String() string { return proto.CompactTextString(m) }
func (*NoneResponse) ProtoMessage()    {}
func (*NoneResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_9e2968d90a4a62fe, []int{0}
}
func (m *NoneResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NoneResponse.Unmarshal(m, b)
}
func (m *NoneResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NoneResponse.Marshal(b, m, deterministic)
}
func (dst *NoneResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoneResponse.Merge(dst, src)
}
func (m *NoneResponse) XXX_Size() int {
	return xxx_messageInfo_NoneResponse.Size(m)
}
func (m *NoneResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NoneResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NoneResponse proto.InternalMessageInfo

type HTTPResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HTTPResponse) Reset()         { *m = HTTPResponse{} }
func (m *HTTPResponse) String() string { return proto.CompactTextString(m) }
func (*HTTPResponse) ProtoMessage()    {}
func (*HTTPResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_9e2968d90a4a62fe, []int{1}
}
func (m *HTTPResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPResponse.Unmarshal(m, b)
}
func (m *HTTPResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPResponse.Marshal(b, m, deterministic)
}
func (dst *HTTPResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPResponse.Merge(dst, src)
}
func (m *HTTPResponse) XXX_Size() int {
	return xxx_messageInfo_HTTPResponse.Size(m)
}
func (m *HTTPResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPResponse proto.InternalMessageInfo

type Config struct {
	Response             *serial.TypedMessage `protobuf:"bytes,1,opt,name=response" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_9e2968d90a4a62fe, []int{2}
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

func (m *Config) GetResponse() *serial.TypedMessage {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*NoneResponse)(nil), "v2ray.core.proxy.blackhole.NoneResponse")
	proto.RegisterType((*HTTPResponse)(nil), "v2ray.core.proxy.blackhole.HTTPResponse")
	proto.RegisterType((*Config)(nil), "v2ray.core.proxy.blackhole.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/proxy/blackhole/config.proto", fileDescriptor_config_9e2968d90a4a62fe)
}

var fileDescriptor_config_9e2968d90a4a62fe = []byte{
	// 217 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2e, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x2f, 0x28, 0xca, 0xaf, 0xa8,
	0xd4, 0x4f, 0xca, 0x49, 0x4c, 0xce, 0xce, 0xc8, 0xcf, 0x49, 0xd5, 0x4f, 0xce, 0xcf, 0x4b, 0xcb,
	0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x82, 0x29, 0x2e, 0x4a, 0xd5, 0x03, 0x2b,
	0xd4, 0x83, 0x2b, 0x94, 0x32, 0x40, 0x33, 0x28, 0x39, 0x3f, 0x37, 0x37, 0x3f, 0x4f, 0xbf, 0x38,
	0xb5, 0x28, 0x33, 0x31, 0x47, 0xbf, 0xa4, 0xb2, 0x20, 0x35, 0x25, 0x3e, 0x37, 0xb5, 0xb8, 0x38,
	0x31, 0x3d, 0x15, 0x62, 0x9a, 0x12, 0x1f, 0x17, 0x8f, 0x5f, 0x7e, 0x5e, 0x6a, 0x50, 0x6a, 0x71,
	0x41, 0x7e, 0x5e, 0x71, 0x2a, 0x88, 0xef, 0x11, 0x12, 0x12, 0x00, 0xe7, 0xfb, 0x70, 0xb1, 0x39,
	0x83, 0x6d, 0x17, 0x72, 0xe2, 0xe2, 0x28, 0x82, 0x8a, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b,
	0xa9, 0xe9, 0x21, 0x39, 0x05, 0x62, 0x95, 0x1e, 0xc4, 0x2a, 0xbd, 0x10, 0x90, 0x55, 0xbe, 0x10,
	0x9b, 0x82, 0xe0, 0xfa, 0x9c, 0xbc, 0xb8, 0xe4, 0x92, 0xf3, 0x73, 0xf5, 0x70, 0xfb, 0x20, 0x80,
	0x31, 0x8a, 0x13, 0xce, 0x59, 0xc5, 0x24, 0x15, 0x66, 0x14, 0x94, 0x58, 0xa9, 0xe7, 0x0c, 0x52,
	0x19, 0x00, 0x56, 0xe9, 0x04, 0x93, 0x4c, 0x62, 0x03, 0x7b, 0xc0, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0xb6, 0x3c, 0xef, 0x4f, 0x3d, 0x01, 0x00, 0x00,
}
