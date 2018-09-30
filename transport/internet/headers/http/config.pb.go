package http

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
	// "Accept", "Cookie", etc
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Each entry must be valid in one piece. Random entry will be chosen if multiple entries present.
	Value                []string `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{0}
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

func (m *Header) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Header) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

// HTTP version. Default value "1.1".
type Version struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{1}
}
func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// HTTP method. Default value "GET".
type Method struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Method) Reset()         { *m = Method{} }
func (m *Method) String() string { return proto.CompactTextString(m) }
func (*Method) ProtoMessage()    {}
func (*Method) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{2}
}
func (m *Method) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Method.Unmarshal(m, b)
}
func (m *Method) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Method.Marshal(b, m, deterministic)
}
func (m *Method) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Method.Merge(m, src)
}
func (m *Method) XXX_Size() int {
	return xxx_messageInfo_Method.Size(m)
}
func (m *Method) XXX_DiscardUnknown() {
	xxx_messageInfo_Method.DiscardUnknown(m)
}

var xxx_messageInfo_Method proto.InternalMessageInfo

func (m *Method) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type RequestConfig struct {
	// Full HTTP version like "1.1".
	Version *Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// GET, POST, CONNECT etc
	Method *Method `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	// URI like "/login.php"
	Uri                  []string  `protobuf:"bytes,3,rep,name=uri,proto3" json:"uri,omitempty"`
	Header               []*Header `protobuf:"bytes,4,rep,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *RequestConfig) Reset()         { *m = RequestConfig{} }
func (m *RequestConfig) String() string { return proto.CompactTextString(m) }
func (*RequestConfig) ProtoMessage()    {}
func (*RequestConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{3}
}
func (m *RequestConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestConfig.Unmarshal(m, b)
}
func (m *RequestConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestConfig.Marshal(b, m, deterministic)
}
func (m *RequestConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestConfig.Merge(m, src)
}
func (m *RequestConfig) XXX_Size() int {
	return xxx_messageInfo_RequestConfig.Size(m)
}
func (m *RequestConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestConfig.DiscardUnknown(m)
}

var xxx_messageInfo_RequestConfig proto.InternalMessageInfo

func (m *RequestConfig) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *RequestConfig) GetMethod() *Method {
	if m != nil {
		return m.Method
	}
	return nil
}

func (m *RequestConfig) GetUri() []string {
	if m != nil {
		return m.Uri
	}
	return nil
}

func (m *RequestConfig) GetHeader() []*Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type Status struct {
	// Status code. Default "200".
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	// Statue reason. Default "OK".
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{4}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Status) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type ResponseConfig struct {
	Version              *Version  `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Status               *Status   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Header               []*Header `protobuf:"bytes,3,rep,name=header,proto3" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ResponseConfig) Reset()         { *m = ResponseConfig{} }
func (m *ResponseConfig) String() string { return proto.CompactTextString(m) }
func (*ResponseConfig) ProtoMessage()    {}
func (*ResponseConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{5}
}
func (m *ResponseConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseConfig.Unmarshal(m, b)
}
func (m *ResponseConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseConfig.Marshal(b, m, deterministic)
}
func (m *ResponseConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseConfig.Merge(m, src)
}
func (m *ResponseConfig) XXX_Size() int {
	return xxx_messageInfo_ResponseConfig.Size(m)
}
func (m *ResponseConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseConfig proto.InternalMessageInfo

func (m *ResponseConfig) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *ResponseConfig) GetStatus() *Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *ResponseConfig) GetHeader() []*Header {
	if m != nil {
		return m.Header
	}
	return nil
}

type Config struct {
	// Settings for authenticating requests. If not set, client side will not send authenication header, and server side will bypass authentication.
	Request *RequestConfig `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	// Settings for authenticating responses. If not set, client side will bypass authentication, and server side will not send authentication header.
	Response             *ResponseConfig `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2685d0b4b039e80, []int{6}
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

func (m *Config) GetRequest() *RequestConfig {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *Config) GetResponse() *ResponseConfig {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*Header)(nil), "v2ray.core.transport.internet.headers.http.Header")
	proto.RegisterType((*Version)(nil), "v2ray.core.transport.internet.headers.http.Version")
	proto.RegisterType((*Method)(nil), "v2ray.core.transport.internet.headers.http.Method")
	proto.RegisterType((*RequestConfig)(nil), "v2ray.core.transport.internet.headers.http.RequestConfig")
	proto.RegisterType((*Status)(nil), "v2ray.core.transport.internet.headers.http.Status")
	proto.RegisterType((*ResponseConfig)(nil), "v2ray.core.transport.internet.headers.http.ResponseConfig")
	proto.RegisterType((*Config)(nil), "v2ray.core.transport.internet.headers.http.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/headers/http/config.proto", fileDescriptor_e2685d0b4b039e80)
}

var fileDescriptor_e2685d0b4b039e80 = []byte{
	// 394 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xbf, 0x8a, 0xdb, 0x40,
	0x10, 0xc6, 0x91, 0xe4, 0xc8, 0xf1, 0x84, 0x84, 0xb0, 0x84, 0xa0, 0x2a, 0x31, 0xaa, 0x8c, 0x8b,
	0x15, 0xc8, 0x69, 0x92, 0x74, 0x71, 0xe3, 0x04, 0x0c, 0x61, 0x1d, 0x5c, 0xa4, 0xdb, 0xc8, 0x93,
	0x58, 0x10, 0xef, 0x2a, 0xbb, 0x2b, 0x83, 0xdf, 0x20, 0xcf, 0x72, 0xfd, 0x3d, 0xdb, 0xb5, 0xc7,
	0xfe, 0x91, 0xce, 0x57, 0x1c, 0x9c, 0xee, 0xb8, 0x4a, 0x33, 0x68, 0xbe, 0x1f, 0xf3, 0x7d, 0x5a,
	0x2d, 0x7c, 0x3e, 0x96, 0x8a, 0x9f, 0x68, 0x25, 0x0f, 0x45, 0x25, 0x15, 0x16, 0x46, 0x71, 0xa1,
	0x1b, 0xa9, 0x4c, 0x51, 0x0b, 0x83, 0x4a, 0xa0, 0x29, 0xf6, 0xc8, 0x77, 0xa8, 0x74, 0xb1, 0x37,
	0xa6, 0x29, 0x2a, 0x29, 0x7e, 0xd7, 0x7f, 0x68, 0xa3, 0xa4, 0x91, 0x64, 0xde, 0x89, 0x15, 0xd2,
	0x5e, 0x48, 0x3b, 0x21, 0x0d, 0x42, 0x6a, 0x85, 0x79, 0x09, 0xe9, 0xca, 0xf5, 0x84, 0xc0, 0x48,
	0xf0, 0x03, 0x66, 0xd1, 0x34, 0x9a, 0x4d, 0x98, 0xab, 0xc9, 0x1b, 0x78, 0x76, 0xe4, 0x7f, 0x5b,
	0xcc, 0xe2, 0x69, 0x32, 0x9b, 0x30, 0xdf, 0xe4, 0xef, 0x61, 0xbc, 0x45, 0xa5, 0x6b, 0x29, 0x6e,
	0x06, 0xbc, 0x2a, 0x0c, 0xbc, 0x83, 0x74, 0x8d, 0x66, 0x2f, 0x77, 0x77, 0xbc, 0xff, 0x1f, 0xc3,
	0x4b, 0x86, 0xff, 0x5a, 0xd4, 0x66, 0xe9, 0x16, 0x27, 0x6b, 0x18, 0x1f, 0x3d, 0xd2, 0x4d, 0xbe,
	0x28, 0x17, 0xf4, 0xfe, 0x26, 0x68, 0xd8, 0x86, 0x75, 0x0c, 0xf2, 0x0d, 0xd2, 0x83, 0x5b, 0x20,
	0x8b, 0x1d, 0xad, 0x1c, 0x42, 0xf3, 0xab, 0xb3, 0x40, 0x20, 0xaf, 0x21, 0x69, 0x55, 0x9d, 0x25,
	0x2e, 0x01, 0x5b, 0x5a, 0xba, 0x17, 0x64, 0xa3, 0x69, 0x32, 0x94, 0xee, 0xd3, 0x66, 0x81, 0x90,
	0x7f, 0x80, 0x74, 0x63, 0xb8, 0x69, 0xb5, 0xcd, 0xbf, 0x92, 0xbb, 0x3e, 0x7f, 0x5b, 0x93, 0xb7,
	0x90, 0x2a, 0xe4, 0x5a, 0x0a, 0xe7, 0x63, 0xc2, 0x42, 0x97, 0x5f, 0x45, 0xf0, 0x8a, 0xa1, 0x6e,
	0xa4, 0xd0, 0xf8, 0x64, 0x09, 0x6a, 0xb7, 0xd7, 0x43, 0x12, 0xf4, 0x8e, 0x58, 0x20, 0x9c, 0xe5,
	0x95, 0x3c, 0x3a, 0xaf, 0xcb, 0x08, 0xd2, 0xe0, 0x78, 0x03, 0x63, 0xe5, 0x0f, 0x51, 0x70, 0xfc,
	0x71, 0x08, 0xf7, 0xd6, 0xf9, 0x63, 0x1d, 0x89, 0x6c, 0xe1, 0xb9, 0x0a, 0xc1, 0x06, 0xe7, 0x9f,
	0x86, 0x51, 0xcf, 0x3f, 0x0a, 0xeb, 0x59, 0x5f, 0x10, 0xec, 0xcf, 0x3c, 0x00, 0xf5, 0x3d, 0xfa,
	0x39, 0xb2, 0xcf, 0x8b, 0x78, 0xbe, 0x2d, 0x19, 0x3f, 0xd1, 0xa5, 0x15, 0xfd, 0xe8, 0x45, 0x5f,
	0x3b, 0xd1, 0x2a, 0x88, 0x56, 0xc6, 0x34, 0xbf, 0x52, 0x77, 0x03, 0x2c, 0xae, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x10, 0xea, 0x22, 0x27, 0x40, 0x04, 0x00, 0x00,
}
