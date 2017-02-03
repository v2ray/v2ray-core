package kcp

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_serial "v2ray.com/core/common/serial"

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
	Value uint32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *MTU) Reset()                    { *m = MTU{} }
func (m *MTU) String() string            { return proto.CompactTextString(m) }
func (*MTU) ProtoMessage()               {}
func (*MTU) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MTU) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Transmission Time Interview, in milli-sec.
type TTI struct {
	Value uint32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *TTI) Reset()                    { *m = TTI{} }
func (m *TTI) String() string            { return proto.CompactTextString(m) }
func (*TTI) ProtoMessage()               {}
func (*TTI) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TTI) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Uplink capacity, in MB.
type UplinkCapacity struct {
	Value uint32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *UplinkCapacity) Reset()                    { *m = UplinkCapacity{} }
func (m *UplinkCapacity) String() string            { return proto.CompactTextString(m) }
func (*UplinkCapacity) ProtoMessage()               {}
func (*UplinkCapacity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UplinkCapacity) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Downlink capacity, in MB.
type DownlinkCapacity struct {
	Value uint32 `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
}

func (m *DownlinkCapacity) Reset()                    { *m = DownlinkCapacity{} }
func (m *DownlinkCapacity) String() string            { return proto.CompactTextString(m) }
func (*DownlinkCapacity) ProtoMessage()               {}
func (*DownlinkCapacity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DownlinkCapacity) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type WriteBuffer struct {
	// Buffer size in bytes.
	Size uint32 `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
}

func (m *WriteBuffer) Reset()                    { *m = WriteBuffer{} }
func (m *WriteBuffer) String() string            { return proto.CompactTextString(m) }
func (*WriteBuffer) ProtoMessage()               {}
func (*WriteBuffer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *WriteBuffer) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

type ReadBuffer struct {
	// Buffer size in bytes.
	Size uint32 `protobuf:"varint,1,opt,name=size" json:"size,omitempty"`
}

func (m *ReadBuffer) Reset()                    { *m = ReadBuffer{} }
func (m *ReadBuffer) String() string            { return proto.CompactTextString(m) }
func (*ReadBuffer) ProtoMessage()               {}
func (*ReadBuffer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ReadBuffer) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

type ConnectionReuse struct {
	Enable bool `protobuf:"varint,1,opt,name=enable" json:"enable,omitempty"`
}

func (m *ConnectionReuse) Reset()                    { *m = ConnectionReuse{} }
func (m *ConnectionReuse) String() string            { return proto.CompactTextString(m) }
func (*ConnectionReuse) ProtoMessage()               {}
func (*ConnectionReuse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ConnectionReuse) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

type Config struct {
	Mtu              *MTU                                   `protobuf:"bytes,1,opt,name=mtu" json:"mtu,omitempty"`
	Tti              *TTI                                   `protobuf:"bytes,2,opt,name=tti" json:"tti,omitempty"`
	UplinkCapacity   *UplinkCapacity                        `protobuf:"bytes,3,opt,name=uplink_capacity,json=uplinkCapacity" json:"uplink_capacity,omitempty"`
	DownlinkCapacity *DownlinkCapacity                      `protobuf:"bytes,4,opt,name=downlink_capacity,json=downlinkCapacity" json:"downlink_capacity,omitempty"`
	Congestion       bool                                   `protobuf:"varint,5,opt,name=congestion" json:"congestion,omitempty"`
	WriteBuffer      *WriteBuffer                           `protobuf:"bytes,6,opt,name=write_buffer,json=writeBuffer" json:"write_buffer,omitempty"`
	ReadBuffer       *ReadBuffer                            `protobuf:"bytes,7,opt,name=read_buffer,json=readBuffer" json:"read_buffer,omitempty"`
	HeaderConfig     *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,8,opt,name=header_config,json=headerConfig" json:"header_config,omitempty"`
	ConnectionReuse  *ConnectionReuse                       `protobuf:"bytes,9,opt,name=connection_reuse,json=connectionReuse" json:"connection_reuse,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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

func (m *Config) GetHeaderConfig() *v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.HeaderConfig
	}
	return nil
}

func (m *Config) GetConnectionReuse() *ConnectionReuse {
	if m != nil {
		return m.ConnectionReuse
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

func init() { proto.RegisterFile("v2ray.com/core/transport/internet/kcp/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x94, 0x5f, 0x6f, 0xd3, 0x30,
	0x14, 0xc5, 0xd5, 0x75, 0x2d, 0xe3, 0x76, 0x5b, 0x4b, 0x84, 0x50, 0x04, 0x12, 0x5a, 0x2b, 0x31,
	0x8d, 0x07, 0x1c, 0xc8, 0x5e, 0x78, 0x5e, 0x79, 0xa9, 0xa6, 0x22, 0xb0, 0x52, 0x90, 0x26, 0xa1,
	0xe0, 0x3a, 0xb7, 0x25, 0x6a, 0x63, 0x47, 0x8e, 0xb3, 0xaa, 0x7c, 0x23, 0xf8, 0x94, 0xc8, 0x4e,
	0xd3, 0x7f, 0x68, 0x6b, 0xde, 0x6a, 0xdf, 0x73, 0x7f, 0xae, 0xce, 0x3d, 0x37, 0xe0, 0xdf, 0xfb,
	0x8a, 0x2d, 0x09, 0x97, 0x89, 0xc7, 0xa5, 0x42, 0x4f, 0x2b, 0x26, 0xb2, 0x54, 0x2a, 0xed, 0xc5,
	0x42, 0xa3, 0x12, 0xa8, 0xbd, 0x19, 0x4f, 0x3d, 0x2e, 0xc5, 0x24, 0x9e, 0x92, 0x54, 0x49, 0x2d,
	0x9d, 0x6e, 0xd9, 0xa3, 0x90, 0xac, 0xf5, 0xa4, 0xd4, 0x93, 0x19, 0x4f, 0x5f, 0xbe, 0xdf, 0xc3,
	0x72, 0x99, 0x24, 0x52, 0x78, 0x19, 0xaa, 0x98, 0xcd, 0x3d, 0xbd, 0x4c, 0x31, 0x0a, 0x13, 0xcc,
	0x32, 0x36, 0xc5, 0x02, 0xda, 0x7b, 0x05, 0xf5, 0x61, 0x30, 0x72, 0x9e, 0x43, 0xe3, 0x9e, 0xcd,
	0x73, 0x74, 0x6b, 0x17, 0xb5, 0xab, 0x33, 0x5a, 0x1c, 0x4c, 0x31, 0x08, 0x06, 0x0f, 0x14, 0x2f,
	0xe1, 0x7c, 0x94, 0xce, 0x63, 0x31, 0xeb, 0xb3, 0x94, 0xf1, 0x58, 0x2f, 0x1f, 0xd0, 0x5d, 0x41,
	0xe7, 0x93, 0x5c, 0x88, 0x0a, 0xca, 0x2e, 0xb4, 0xbe, 0xab, 0x58, 0xe3, 0x4d, 0x3e, 0x99, 0xa0,
	0x72, 0x1c, 0x38, 0xce, 0xe2, 0xdf, 0xa5, 0xc6, 0xfe, 0xee, 0x5d, 0x00, 0x50, 0x64, 0xd1, 0x23,
	0x8a, 0xb7, 0xd0, 0xee, 0x4b, 0x21, 0x90, 0xeb, 0x58, 0x0a, 0x8a, 0x79, 0x86, 0xce, 0x0b, 0x68,
	0xa2, 0x60, 0xe3, 0x79, 0x21, 0x3c, 0xa1, 0xab, 0x53, 0xef, 0x4f, 0x03, 0x9a, 0x7d, 0xeb, 0xb0,
	0xf3, 0x11, 0xea, 0x89, 0xce, 0x6d, 0xbd, 0xe5, 0x5f, 0x92, 0x83, 0x4e, 0x93, 0x61, 0x30, 0xa2,
	0xa6, 0xc5, 0x74, 0x6a, 0x1d, 0xbb, 0x47, 0x95, 0x3b, 0x83, 0x60, 0x40, 0x4d, 0x8b, 0x73, 0x07,
	0xed, 0xdc, 0x1a, 0x18, 0xf2, 0x95, 0x2f, 0x6e, 0xdd, 0x52, 0x3e, 0x54, 0xa0, 0xec, 0x5a, 0x4f,
	0xcf, 0xf3, 0xdd, 0x51, 0xfc, 0x84, 0x67, 0xd1, 0xca, 0xf4, 0x0d, 0xfd, 0xd8, 0xd2, 0xaf, 0x2b,
	0xd0, 0xf7, 0x07, 0x46, 0x3b, 0xd1, 0xfe, 0x08, 0x5f, 0x03, 0x70, 0x29, 0xa6, 0x98, 0x19, 0x9f,
	0xdd, 0x86, 0x35, 0x76, 0xeb, 0xc6, 0xf9, 0x0a, 0xa7, 0x0b, 0x33, 0xcc, 0x70, 0x6c, 0x67, 0xe5,
	0x36, 0xed, 0xe3, 0xa4, 0xc2, 0xe3, 0x5b, 0x19, 0xa0, 0xad, 0xc5, 0x56, 0x20, 0x3e, 0x43, 0x4b,
	0x21, 0x8b, 0x4a, 0xe2, 0x13, 0x4b, 0x7c, 0x57, 0x81, 0xb8, 0x89, 0x0c, 0x05, 0xb5, 0x89, 0xcf,
	0x2d, 0x9c, 0xfd, 0x42, 0x16, 0xa1, 0x0a, 0x8b, 0x3d, 0x73, 0x4f, 0xfe, 0x1f, 0x62, 0xb1, 0x41,
	0xa4, 0xd8, 0x20, 0x12, 0x98, 0x0d, 0x1a, 0x16, 0x0b, 0x44, 0x4f, 0x8b, 0xe6, 0x55, 0x82, 0x7e,
	0x40, 0x87, 0xaf, 0x73, 0x17, 0x2a, 0x13, 0x3c, 0xf7, 0xa9, 0xe5, 0xf9, 0x15, 0xfe, 0xe1, 0x5e,
	0x64, 0x69, 0x9b, 0xef, 0x5e, 0xdc, 0x50, 0x78, 0xc3, 0x65, 0x72, 0x98, 0xf4, 0xa5, 0x76, 0x57,
	0x9f, 0xf1, 0xf4, 0xef, 0x51, 0xf7, 0x9b, 0x4f, 0xd9, 0x92, 0xf4, 0x8d, 0x34, 0x58, 0x4b, 0x07,
	0xa5, 0xf4, 0x96, 0xa7, 0xe3, 0xa6, 0xfd, 0x04, 0x5c, 0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x1d,
	0x04, 0x18, 0x2b, 0x8d, 0x04, 0x00, 0x00,
}
