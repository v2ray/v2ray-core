package net

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

// PortRange represents a range of ports.
type PortRange struct {
	// The port that this range starts from.
	From uint32 `protobuf:"varint,1,opt,name=From,proto3" json:"From,omitempty"`
	// The port that this range ends with (inclusive).
	To                   uint32   `protobuf:"varint,2,opt,name=To,proto3" json:"To,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PortRange) Reset()         { *m = PortRange{} }
func (m *PortRange) String() string { return proto.CompactTextString(m) }
func (*PortRange) ProtoMessage()    {}
func (*PortRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_166067e37a39f913, []int{0}
}

func (m *PortRange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PortRange.Unmarshal(m, b)
}
func (m *PortRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PortRange.Marshal(b, m, deterministic)
}
func (m *PortRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PortRange.Merge(m, src)
}
func (m *PortRange) XXX_Size() int {
	return xxx_messageInfo_PortRange.Size(m)
}
func (m *PortRange) XXX_DiscardUnknown() {
	xxx_messageInfo_PortRange.DiscardUnknown(m)
}

var xxx_messageInfo_PortRange proto.InternalMessageInfo

func (m *PortRange) GetFrom() uint32 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *PortRange) GetTo() uint32 {
	if m != nil {
		return m.To
	}
	return 0
}

func init() {
	proto.RegisterType((*PortRange)(nil), "v2ray.core.common.net.PortRange")
}

func init() {
	proto.RegisterFile("v2ray.com/core/common/net/port.proto", fileDescriptor_166067e37a39f913)
}

var fileDescriptor_166067e37a39f913 = []byte{
	// 158 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x29, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0xce, 0xcf, 0xcd, 0xcd,
	0xcf, 0xd3, 0xcf, 0x4b, 0x2d, 0xd1, 0x2f, 0xc8, 0x2f, 0x2a, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x12, 0x85, 0xa9, 0x2a, 0x4a, 0xd5, 0x83, 0xa8, 0xd0, 0xcb, 0x4b, 0x2d, 0x51, 0xd2, 0xe7,
	0xe2, 0x0c, 0xc8, 0x2f, 0x2a, 0x09, 0x4a, 0xcc, 0x4b, 0x4f, 0x15, 0x12, 0xe2, 0x62, 0x71, 0x2b,
	0xca, 0xcf, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0x02, 0xb3, 0x85, 0xf8, 0xb8, 0x98, 0x42,
	0xf2, 0x25, 0x98, 0xc0, 0x22, 0x4c, 0x21, 0xf9, 0x4e, 0x56, 0x5c, 0x92, 0xc9, 0xf9, 0xb9, 0x7a,
	0x58, 0x4d, 0x0b, 0x60, 0x8c, 0x62, 0xce, 0x4b, 0x2d, 0x59, 0xc5, 0x24, 0x1a, 0x66, 0x14, 0x94,
	0x58, 0xa9, 0xe7, 0x0c, 0x92, 0x76, 0x86, 0x48, 0xfb, 0xa5, 0x96, 0x24, 0xb1, 0x81, 0x9d, 0x62,
	0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x0f, 0x6a, 0xc0, 0xb2, 0x00, 0x00, 0x00,
}
