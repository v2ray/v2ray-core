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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

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

// PortList is a list of ports.
type PortList struct {
	Range                []*PortRange `protobuf:"bytes,1,rep,name=range,proto3" json:"range,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PortList) Reset()         { *m = PortList{} }
func (m *PortList) String() string { return proto.CompactTextString(m) }
func (*PortList) ProtoMessage()    {}
func (*PortList) Descriptor() ([]byte, []int) {
	return fileDescriptor_166067e37a39f913, []int{1}
}

func (m *PortList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PortList.Unmarshal(m, b)
}
func (m *PortList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PortList.Marshal(b, m, deterministic)
}
func (m *PortList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PortList.Merge(m, src)
}
func (m *PortList) XXX_Size() int {
	return xxx_messageInfo_PortList.Size(m)
}
func (m *PortList) XXX_DiscardUnknown() {
	xxx_messageInfo_PortList.DiscardUnknown(m)
}

var xxx_messageInfo_PortList proto.InternalMessageInfo

func (m *PortList) GetRange() []*PortRange {
	if m != nil {
		return m.Range
	}
	return nil
}

func init() {
	proto.RegisterType((*PortRange)(nil), "v2ray.core.common.net.PortRange")
	proto.RegisterType((*PortList)(nil), "v2ray.core.common.net.PortList")
}

func init() {
	proto.RegisterFile("v2ray.com/core/common/net/port.proto", fileDescriptor_166067e37a39f913)
}

var fileDescriptor_166067e37a39f913 = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x29, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0xce, 0xcf, 0xcd, 0xcd,
	0xcf, 0xd3, 0xcf, 0x4b, 0x2d, 0xd1, 0x2f, 0xc8, 0x2f, 0x2a, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x12, 0x85, 0xa9, 0x2a, 0x4a, 0xd5, 0x83, 0xa8, 0xd0, 0xcb, 0x4b, 0x2d, 0x51, 0xd2, 0xe7,
	0xe2, 0x0c, 0xc8, 0x2f, 0x2a, 0x09, 0x4a, 0xcc, 0x4b, 0x4f, 0x15, 0x12, 0xe2, 0x62, 0x71, 0x2b,
	0xca, 0xcf, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0x02, 0xb3, 0x85, 0xf8, 0xb8, 0x98, 0x42,
	0xf2, 0x25, 0x98, 0xc0, 0x22, 0x4c, 0x21, 0xf9, 0x4a, 0x4e, 0x5c, 0x1c, 0x20, 0x0d, 0x3e, 0x99,
	0xc5, 0x25, 0x42, 0x66, 0x5c, 0xac, 0x45, 0x20, 0x8d, 0x12, 0x8c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46,
	0x0a, 0x7a, 0x58, 0xed, 0xd0, 0x83, 0x5b, 0x10, 0x04, 0x51, 0xee, 0x64, 0xc5, 0x25, 0x99, 0x9c,
	0x9f, 0x8b, 0x5d, 0x75, 0x00, 0x63, 0x14, 0x73, 0x5e, 0x6a, 0xc9, 0x2a, 0x26, 0xd1, 0x30, 0xa3,
	0xa0, 0xc4, 0x4a, 0x3d, 0x67, 0x90, 0xb4, 0x33, 0x44, 0xda, 0x2f, 0xb5, 0x24, 0x89, 0x0d, 0xec,
	0x1d, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xba, 0xd0, 0x7b, 0xfa, 0xf6, 0x00, 0x00, 0x00,
}
