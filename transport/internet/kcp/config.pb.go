package kcp

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

// Maximum Transmission Unit, in bytes.
type MTU struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MTU) Reset()         { *m = MTU{} }
func (m *MTU) String() string { return proto.CompactTextString(m) }
func (*MTU) ProtoMessage()    {}
func (*MTU) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{0}
}
func (m *MTU) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MTU.Unmarshal(m, b)
}
func (m *MTU) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MTU.Marshal(b, m, deterministic)
}
func (dst *MTU) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MTU.Merge(dst, src)
}
func (m *MTU) XXX_Size() int {
	return xxx_messageInfo_MTU.Size(m)
}
func (m *MTU) XXX_DiscardUnknown() {
	xxx_messageInfo_MTU.DiscardUnknown(m)
}

var xxx_messageInfo_MTU proto.InternalMessageInfo

func (m *MTU) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Transmission Time Interview, in milli-sec.
type TTI struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TTI) Reset()         { *m = TTI{} }
func (m *TTI) String() string { return proto.CompactTextString(m) }
func (*TTI) ProtoMessage()    {}
func (*TTI) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{1}
}
func (m *TTI) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TTI.Unmarshal(m, b)
}
func (m *TTI) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TTI.Marshal(b, m, deterministic)
}
func (dst *TTI) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TTI.Merge(dst, src)
}
func (m *TTI) XXX_Size() int {
	return xxx_messageInfo_TTI.Size(m)
}
func (m *TTI) XXX_DiscardUnknown() {
	xxx_messageInfo_TTI.DiscardUnknown(m)
}

var xxx_messageInfo_TTI proto.InternalMessageInfo

func (m *TTI) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Uplink capacity, in MB.
type UplinkCapacity struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UplinkCapacity) Reset()         { *m = UplinkCapacity{} }
func (m *UplinkCapacity) String() string { return proto.CompactTextString(m) }
func (*UplinkCapacity) ProtoMessage()    {}
func (*UplinkCapacity) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{2}
}
func (m *UplinkCapacity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UplinkCapacity.Unmarshal(m, b)
}
func (m *UplinkCapacity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UplinkCapacity.Marshal(b, m, deterministic)
}
func (dst *UplinkCapacity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UplinkCapacity.Merge(dst, src)
}
func (m *UplinkCapacity) XXX_Size() int {
	return xxx_messageInfo_UplinkCapacity.Size(m)
}
func (m *UplinkCapacity) XXX_DiscardUnknown() {
	xxx_messageInfo_UplinkCapacity.DiscardUnknown(m)
}

var xxx_messageInfo_UplinkCapacity proto.InternalMessageInfo

func (m *UplinkCapacity) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Downlink capacity, in MB.
type DownlinkCapacity struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownlinkCapacity) Reset()         { *m = DownlinkCapacity{} }
func (m *DownlinkCapacity) String() string { return proto.CompactTextString(m) }
func (*DownlinkCapacity) ProtoMessage()    {}
func (*DownlinkCapacity) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{3}
}
func (m *DownlinkCapacity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownlinkCapacity.Unmarshal(m, b)
}
func (m *DownlinkCapacity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownlinkCapacity.Marshal(b, m, deterministic)
}
func (dst *DownlinkCapacity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownlinkCapacity.Merge(dst, src)
}
func (m *DownlinkCapacity) XXX_Size() int {
	return xxx_messageInfo_DownlinkCapacity.Size(m)
}
func (m *DownlinkCapacity) XXX_DiscardUnknown() {
	xxx_messageInfo_DownlinkCapacity.DiscardUnknown(m)
}

var xxx_messageInfo_DownlinkCapacity proto.InternalMessageInfo

func (m *DownlinkCapacity) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type WriteBuffer struct {
	// Buffer size in bytes.
	Size                 uint32   `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WriteBuffer) Reset()         { *m = WriteBuffer{} }
func (m *WriteBuffer) String() string { return proto.CompactTextString(m) }
func (*WriteBuffer) ProtoMessage()    {}
func (*WriteBuffer) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{4}
}
func (m *WriteBuffer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WriteBuffer.Unmarshal(m, b)
}
func (m *WriteBuffer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WriteBuffer.Marshal(b, m, deterministic)
}
func (dst *WriteBuffer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteBuffer.Merge(dst, src)
}
func (m *WriteBuffer) XXX_Size() int {
	return xxx_messageInfo_WriteBuffer.Size(m)
}
func (m *WriteBuffer) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteBuffer.DiscardUnknown(m)
}

var xxx_messageInfo_WriteBuffer proto.InternalMessageInfo

func (m *WriteBuffer) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

type ReadBuffer struct {
	// Buffer size in bytes.
	Size                 uint32   `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReadBuffer) Reset()         { *m = ReadBuffer{} }
func (m *ReadBuffer) String() string { return proto.CompactTextString(m) }
func (*ReadBuffer) ProtoMessage()    {}
func (*ReadBuffer) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{5}
}
func (m *ReadBuffer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadBuffer.Unmarshal(m, b)
}
func (m *ReadBuffer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadBuffer.Marshal(b, m, deterministic)
}
func (dst *ReadBuffer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadBuffer.Merge(dst, src)
}
func (m *ReadBuffer) XXX_Size() int {
	return xxx_messageInfo_ReadBuffer.Size(m)
}
func (m *ReadBuffer) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadBuffer.DiscardUnknown(m)
}

var xxx_messageInfo_ReadBuffer proto.InternalMessageInfo

func (m *ReadBuffer) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

type ConnectionReuse struct {
	Enable               bool     `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectionReuse) Reset()         { *m = ConnectionReuse{} }
func (m *ConnectionReuse) String() string { return proto.CompactTextString(m) }
func (*ConnectionReuse) ProtoMessage()    {}
func (*ConnectionReuse) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{6}
}
func (m *ConnectionReuse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectionReuse.Unmarshal(m, b)
}
func (m *ConnectionReuse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectionReuse.Marshal(b, m, deterministic)
}
func (dst *ConnectionReuse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectionReuse.Merge(dst, src)
}
func (m *ConnectionReuse) XXX_Size() int {
	return xxx_messageInfo_ConnectionReuse.Size(m)
}
func (m *ConnectionReuse) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectionReuse.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectionReuse proto.InternalMessageInfo

func (m *ConnectionReuse) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

type Config struct {
	Mtu                  *MTU                 `protobuf:"bytes,1,opt,name=mtu,proto3" json:"mtu,omitempty"`
	Tti                  *TTI                 `protobuf:"bytes,2,opt,name=tti,proto3" json:"tti,omitempty"`
	UplinkCapacity       *UplinkCapacity      `protobuf:"bytes,3,opt,name=uplink_capacity,json=uplinkCapacity,proto3" json:"uplink_capacity,omitempty"`
	DownlinkCapacity     *DownlinkCapacity    `protobuf:"bytes,4,opt,name=downlink_capacity,json=downlinkCapacity,proto3" json:"downlink_capacity,omitempty"`
	Congestion           bool                 `protobuf:"varint,5,opt,name=congestion,proto3" json:"congestion,omitempty"`
	WriteBuffer          *WriteBuffer         `protobuf:"bytes,6,opt,name=write_buffer,json=writeBuffer,proto3" json:"write_buffer,omitempty"`
	ReadBuffer           *ReadBuffer          `protobuf:"bytes,7,opt,name=read_buffer,json=readBuffer,proto3" json:"read_buffer,omitempty"`
	HeaderConfig         *serial.TypedMessage `protobuf:"bytes,8,opt,name=header_config,json=headerConfig,proto3" json:"header_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_4bc2f043099e7e59, []int{7}
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

func (m *Config) GetMtu() *MTU {
	if m != nil {
		return m.Mtu
	}
	return nil
}

func (m *Config) GetTti() *TTI {
	if m != nil {
		return m.Tti
	}
	return nil
}

func (m *Config) GetUplinkCapacity() *UplinkCapacity {
	if m != nil {
		return m.UplinkCapacity
	}
	return nil
}

func (m *Config) GetDownlinkCapacity() *DownlinkCapacity {
	if m != nil {
		return m.DownlinkCapacity
	}
	return nil
}

func (m *Config) GetCongestion() bool {
	if m != nil {
		return m.Congestion
	}
	return false
}

func (m *Config) GetWriteBuffer() *WriteBuffer {
	if m != nil {
		return m.WriteBuffer
	}
	return nil
}

func (m *Config) GetReadBuffer() *ReadBuffer {
	if m != nil {
		return m.ReadBuffer
	}
	return nil
}

func (m *Config) GetHeaderConfig() *serial.TypedMessage {
	if m != nil {
		return m.HeaderConfig
	}
	return nil
}

func init() {
	proto.RegisterType((*MTU)(nil), "v2ray.core.transport.internet.kcp.MTU")
	proto.RegisterType((*TTI)(nil), "v2ray.core.transport.internet.kcp.TTI")
	proto.RegisterType((*UplinkCapacity)(nil), "v2ray.core.transport.internet.kcp.UplinkCapacity")
	proto.RegisterType((*DownlinkCapacity)(nil), "v2ray.core.transport.internet.kcp.DownlinkCapacity")
	proto.RegisterType((*WriteBuffer)(nil), "v2ray.core.transport.internet.kcp.WriteBuffer")
	proto.RegisterType((*ReadBuffer)(nil), "v2ray.core.transport.internet.kcp.ReadBuffer")
	proto.RegisterType((*ConnectionReuse)(nil), "v2ray.core.transport.internet.kcp.ConnectionReuse")
	proto.RegisterType((*Config)(nil), "v2ray.core.transport.internet.kcp.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/transport/internet/kcp/config.proto", fileDescriptor_config_4bc2f043099e7e59)
}

var fileDescriptor_config_4bc2f043099e7e59 = []byte{
	// 471 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x5f, 0x6f, 0xd3, 0x3e,
	0x14, 0x55, 0xd7, 0xae, 0xbf, 0xfe, 0x6e, 0xf7, 0xa7, 0x44, 0x08, 0x45, 0x20, 0xa1, 0xb5, 0x12,
	0xd3, 0x78, 0xc0, 0x81, 0xee, 0x85, 0xe7, 0x95, 0x97, 0x32, 0x15, 0x81, 0x95, 0x82, 0xb4, 0x97,
	0xe2, 0x3a, 0xb7, 0xc5, 0x6a, 0x63, 0x5b, 0x8e, 0xb3, 0xaa, 0x7c, 0x24, 0x3e, 0x0d, 0x1f, 0x09,
	0xc5, 0x6e, 0xd6, 0xae, 0x68, 0x2c, 0x6f, 0x71, 0xee, 0x39, 0xc7, 0xd6, 0x39, 0xf7, 0x40, 0xff,
	0xb6, 0x6f, 0xd8, 0x9a, 0x70, 0x95, 0x46, 0x5c, 0x19, 0x8c, 0xac, 0x61, 0x32, 0xd3, 0xca, 0xd8,
	0x48, 0x48, 0x8b, 0x46, 0xa2, 0x8d, 0x16, 0x5c, 0x47, 0x5c, 0xc9, 0x99, 0x98, 0x13, 0x6d, 0x94,
	0x55, 0x41, 0xb7, 0xe4, 0x18, 0x24, 0x77, 0x78, 0x52, 0xe2, 0xc9, 0x82, 0xeb, 0xe7, 0x6f, 0xf7,
	0x64, 0xb9, 0x4a, 0x53, 0x25, 0xa3, 0x0c, 0x8d, 0x60, 0xcb, 0xc8, 0xae, 0x35, 0x26, 0x93, 0x14,
	0xb3, 0x8c, 0xcd, 0xd1, 0x8b, 0xf6, 0x5e, 0x40, 0x7d, 0x14, 0x8f, 0x83, 0xa7, 0x70, 0x78, 0xcb,
	0x96, 0x39, 0x86, 0xb5, 0xb3, 0xda, 0xc5, 0x31, 0xf5, 0x87, 0x62, 0x18, 0xc7, 0xc3, 0x07, 0x86,
	0xe7, 0x70, 0x32, 0xd6, 0x4b, 0x21, 0x17, 0x03, 0xa6, 0x19, 0x17, 0x76, 0xfd, 0x00, 0xee, 0x02,
	0x3a, 0x1f, 0xd4, 0x4a, 0x56, 0x40, 0x76, 0xa1, 0xfd, 0xcd, 0x08, 0x8b, 0x57, 0xf9, 0x6c, 0x86,
	0x26, 0x08, 0xa0, 0x91, 0x89, 0x9f, 0x25, 0xc6, 0x7d, 0xf7, 0xce, 0x00, 0x28, 0xb2, 0xe4, 0x1f,
	0x88, 0xd7, 0x70, 0x3a, 0x50, 0x52, 0x22, 0xb7, 0x42, 0x49, 0x8a, 0x79, 0x86, 0xc1, 0x33, 0x68,
	0xa2, 0x64, 0xd3, 0xa5, 0x07, 0xb6, 0xe8, 0xe6, 0xd4, 0xfb, 0xdd, 0x80, 0xe6, 0xc0, 0x39, 0x1c,
	0xbc, 0x87, 0x7a, 0x6a, 0x73, 0x37, 0x6f, 0xf7, 0xcf, 0xc9, 0xa3, 0x4e, 0x93, 0x51, 0x3c, 0xa6,
	0x05, 0xa5, 0x60, 0x5a, 0x2b, 0xc2, 0x83, 0xca, 0xcc, 0x38, 0x1e, 0xd2, 0x82, 0x12, 0xdc, 0xc0,
	0x69, 0xee, 0x0c, 0x9c, 0xf0, 0x8d, 0x2f, 0x61, 0xdd, 0xa9, 0xbc, 0xab, 0xa0, 0x72, 0xdf, 0x7a,
	0x7a, 0x92, 0xdf, 0x8f, 0xe2, 0x3b, 0x3c, 0x49, 0x36, 0xa6, 0x6f, 0xd5, 0x1b, 0x4e, 0xfd, 0xb2,
	0x82, 0xfa, 0x7e, 0x60, 0xb4, 0x93, 0xec, 0x47, 0xf8, 0x12, 0x80, 0x2b, 0x39, 0xc7, 0xac, 0xf0,
	0x39, 0x3c, 0x74, 0xc6, 0xee, 0xfc, 0x09, 0xbe, 0xc0, 0xd1, 0xaa, 0x08, 0x73, 0x32, 0x75, 0x59,
	0x85, 0x4d, 0x77, 0x39, 0xa9, 0x70, 0xf9, 0xce, 0x0e, 0xd0, 0xf6, 0x6a, 0x67, 0x21, 0x3e, 0x41,
	0xdb, 0x20, 0x4b, 0x4a, 0xc5, 0xff, 0x9c, 0xe2, 0x9b, 0x0a, 0x8a, 0xdb, 0x95, 0xa1, 0x60, 0xb6,
	0xeb, 0x73, 0x0d, 0xc7, 0x3f, 0x90, 0x25, 0x68, 0x26, 0xbe, 0x67, 0x61, 0xeb, 0xef, 0x10, 0x7d,
	0x83, 0x88, 0x6f, 0x10, 0x89, 0x8b, 0x06, 0x8d, 0x7c, 0x81, 0xe8, 0x91, 0x27, 0xfb, 0x0d, 0xfa,
	0xd8, 0x68, 0xfd, 0xdf, 0x81, 0x2b, 0x0a, 0xaf, 0xb8, 0x4a, 0x1f, 0x7f, 0xd2, 0xe7, 0xda, 0x4d,
	0x7d, 0xc1, 0xf5, 0xaf, 0x83, 0xee, 0xd7, 0x3e, 0x65, 0x6b, 0x32, 0x28, 0xa0, 0xf1, 0x1d, 0x74,
	0x58, 0x42, 0xaf, 0xb9, 0x9e, 0x36, 0x5d, 0x53, 0x2f, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x54,
	0xdd, 0xba, 0xf9, 0x34, 0x04, 0x00, 0x00,
}
