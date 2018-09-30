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

type Network int32

const (
	Network_Unknown Network = 0
	Network_RawTCP  Network = 1 // Deprecated: Do not use.
	Network_TCP     Network = 2
	Network_UDP     Network = 3
)

var Network_name = map[int32]string{
	0: "Unknown",
	1: "RawTCP",
	2: "TCP",
	3: "UDP",
}

var Network_value = map[string]int32{
	"Unknown": 0,
	"RawTCP":  1,
	"TCP":     2,
	"UDP":     3,
}

func (x Network) String() string {
	return proto.EnumName(Network_name, int32(x))
}

func (Network) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6a103d5ccb9e785e, []int{0}
}

// NetworkList is a list of Networks.
type NetworkList struct {
	Network              []Network `protobuf:"varint,1,rep,packed,name=network,proto3,enum=v2ray.core.common.net.Network" json:"network,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *NetworkList) Reset()         { *m = NetworkList{} }
func (m *NetworkList) String() string { return proto.CompactTextString(m) }
func (*NetworkList) ProtoMessage()    {}
func (*NetworkList) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a103d5ccb9e785e, []int{0}
}

func (m *NetworkList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkList.Unmarshal(m, b)
}
func (m *NetworkList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkList.Marshal(b, m, deterministic)
}
func (m *NetworkList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkList.Merge(m, src)
}
func (m *NetworkList) XXX_Size() int {
	return xxx_messageInfo_NetworkList.Size(m)
}
func (m *NetworkList) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkList.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkList proto.InternalMessageInfo

func (m *NetworkList) GetNetwork() []Network {
	if m != nil {
		return m.Network
	}
	return nil
}

func init() {
	proto.RegisterEnum("v2ray.core.common.net.Network", Network_name, Network_value)
	proto.RegisterType((*NetworkList)(nil), "v2ray.core.common.net.NetworkList")
}

func init() {
	proto.RegisterFile("v2ray.com/core/common/net/network.proto", fileDescriptor_6a103d5ccb9e785e)
}

var fileDescriptor_6a103d5ccb9e785e = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2f, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0xce, 0xcf, 0xcd, 0xcd,
	0xcf, 0xd3, 0xcf, 0x4b, 0x2d, 0x01, 0xe1, 0xf2, 0xfc, 0xa2, 0x6c, 0xbd, 0x82, 0xa2, 0xfc, 0x92,
	0x7c, 0x21, 0x51, 0x98, 0xc2, 0xa2, 0x54, 0x3d, 0x88, 0x22, 0xbd, 0xbc, 0xd4, 0x12, 0x25, 0x77,
	0x2e, 0x6e, 0x3f, 0x88, 0x3a, 0x9f, 0xcc, 0xe2, 0x12, 0x21, 0x0b, 0x2e, 0x76, 0xa8, 0x36, 0x09,
	0x46, 0x05, 0x66, 0x0d, 0x3e, 0x23, 0x39, 0x3d, 0xac, 0xfa, 0xf4, 0xa0, 0x9a, 0x82, 0x60, 0xca,
	0xb5, 0x2c, 0xb8, 0xd8, 0xa1, 0x62, 0x42, 0xdc, 0x5c, 0xec, 0xa1, 0x79, 0xd9, 0x79, 0xf9, 0xe5,
	0x79, 0x02, 0x0c, 0x42, 0x7c, 0x5c, 0x6c, 0x41, 0x89, 0xe5, 0x21, 0xce, 0x01, 0x02, 0x8c, 0x52,
	0x4c, 0x1c, 0x8c, 0x42, 0xec, 0x5c, 0xcc, 0x20, 0x0e, 0x13, 0x88, 0x11, 0xea, 0x12, 0x20, 0xc0,
	0xec, 0x64, 0xc5, 0x25, 0x99, 0x9c, 0x9f, 0x8b, 0xdd, 0x9e, 0x00, 0xc6, 0x28, 0xe6, 0xbc, 0xd4,
	0x92, 0x55, 0x4c, 0xa2, 0x61, 0x46, 0x41, 0x89, 0x95, 0x7a, 0xce, 0x20, 0x69, 0x67, 0x88, 0xb4,
	0x5f, 0x6a, 0x49, 0x12, 0x1b, 0xd8, 0x73, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xef, 0x75,
	0xd9, 0x5b, 0x07, 0x01, 0x00, 0x00,
}
