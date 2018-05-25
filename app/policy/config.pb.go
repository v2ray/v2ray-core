package policy

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

type Second struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Second) Reset()         { *m = Second{} }
func (m *Second) String() string { return proto.CompactTextString(m) }
func (*Second) ProtoMessage()    {}
func (*Second) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{0}
}
func (m *Second) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Second.Unmarshal(m, b)
}
func (m *Second) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Second.Marshal(b, m, deterministic)
}
func (dst *Second) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Second.Merge(dst, src)
}
func (m *Second) XXX_Size() int {
	return xxx_messageInfo_Second.Size(m)
}
func (m *Second) XXX_DiscardUnknown() {
	xxx_messageInfo_Second.DiscardUnknown(m)
}

var xxx_messageInfo_Second proto.InternalMessageInfo

func (m *Second) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type Policy struct {
	Timeout              *Policy_Timeout `protobuf:"bytes,1,opt,name=timeout" json:"timeout,omitempty"`
	Stats                *Policy_Stats   `protobuf:"bytes,2,opt,name=stats" json:"stats,omitempty"`
	Buffer               *Policy_Buffer  `protobuf:"bytes,3,opt,name=buffer" json:"buffer,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Policy) Reset()         { *m = Policy{} }
func (m *Policy) String() string { return proto.CompactTextString(m) }
func (*Policy) ProtoMessage()    {}
func (*Policy) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{1}
}
func (m *Policy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Policy.Unmarshal(m, b)
}
func (m *Policy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Policy.Marshal(b, m, deterministic)
}
func (dst *Policy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Policy.Merge(dst, src)
}
func (m *Policy) XXX_Size() int {
	return xxx_messageInfo_Policy.Size(m)
}
func (m *Policy) XXX_DiscardUnknown() {
	xxx_messageInfo_Policy.DiscardUnknown(m)
}

var xxx_messageInfo_Policy proto.InternalMessageInfo

func (m *Policy) GetTimeout() *Policy_Timeout {
	if m != nil {
		return m.Timeout
	}
	return nil
}

func (m *Policy) GetStats() *Policy_Stats {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *Policy) GetBuffer() *Policy_Buffer {
	if m != nil {
		return m.Buffer
	}
	return nil
}

// Timeout is a message for timeout settings in various stages, in seconds.
type Policy_Timeout struct {
	Handshake            *Second  `protobuf:"bytes,1,opt,name=handshake" json:"handshake,omitempty"`
	ConnectionIdle       *Second  `protobuf:"bytes,2,opt,name=connection_idle,json=connectionIdle" json:"connection_idle,omitempty"`
	UplinkOnly           *Second  `protobuf:"bytes,3,opt,name=uplink_only,json=uplinkOnly" json:"uplink_only,omitempty"`
	DownlinkOnly         *Second  `protobuf:"bytes,4,opt,name=downlink_only,json=downlinkOnly" json:"downlink_only,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Policy_Timeout) Reset()         { *m = Policy_Timeout{} }
func (m *Policy_Timeout) String() string { return proto.CompactTextString(m) }
func (*Policy_Timeout) ProtoMessage()    {}
func (*Policy_Timeout) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{1, 0}
}
func (m *Policy_Timeout) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Policy_Timeout.Unmarshal(m, b)
}
func (m *Policy_Timeout) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Policy_Timeout.Marshal(b, m, deterministic)
}
func (dst *Policy_Timeout) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Policy_Timeout.Merge(dst, src)
}
func (m *Policy_Timeout) XXX_Size() int {
	return xxx_messageInfo_Policy_Timeout.Size(m)
}
func (m *Policy_Timeout) XXX_DiscardUnknown() {
	xxx_messageInfo_Policy_Timeout.DiscardUnknown(m)
}

var xxx_messageInfo_Policy_Timeout proto.InternalMessageInfo

func (m *Policy_Timeout) GetHandshake() *Second {
	if m != nil {
		return m.Handshake
	}
	return nil
}

func (m *Policy_Timeout) GetConnectionIdle() *Second {
	if m != nil {
		return m.ConnectionIdle
	}
	return nil
}

func (m *Policy_Timeout) GetUplinkOnly() *Second {
	if m != nil {
		return m.UplinkOnly
	}
	return nil
}

func (m *Policy_Timeout) GetDownlinkOnly() *Second {
	if m != nil {
		return m.DownlinkOnly
	}
	return nil
}

type Policy_Stats struct {
	UserUplink           bool     `protobuf:"varint,1,opt,name=user_uplink,json=userUplink" json:"user_uplink,omitempty"`
	UserDownlink         bool     `protobuf:"varint,2,opt,name=user_downlink,json=userDownlink" json:"user_downlink,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Policy_Stats) Reset()         { *m = Policy_Stats{} }
func (m *Policy_Stats) String() string { return proto.CompactTextString(m) }
func (*Policy_Stats) ProtoMessage()    {}
func (*Policy_Stats) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{1, 1}
}
func (m *Policy_Stats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Policy_Stats.Unmarshal(m, b)
}
func (m *Policy_Stats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Policy_Stats.Marshal(b, m, deterministic)
}
func (dst *Policy_Stats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Policy_Stats.Merge(dst, src)
}
func (m *Policy_Stats) XXX_Size() int {
	return xxx_messageInfo_Policy_Stats.Size(m)
}
func (m *Policy_Stats) XXX_DiscardUnknown() {
	xxx_messageInfo_Policy_Stats.DiscardUnknown(m)
}

var xxx_messageInfo_Policy_Stats proto.InternalMessageInfo

func (m *Policy_Stats) GetUserUplink() bool {
	if m != nil {
		return m.UserUplink
	}
	return false
}

func (m *Policy_Stats) GetUserDownlink() bool {
	if m != nil {
		return m.UserDownlink
	}
	return false
}

type Policy_Buffer struct {
	Enabled              bool     `protobuf:"varint,1,opt,name=enabled" json:"enabled,omitempty"`
	Size                 uint32   `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Policy_Buffer) Reset()         { *m = Policy_Buffer{} }
func (m *Policy_Buffer) String() string { return proto.CompactTextString(m) }
func (*Policy_Buffer) ProtoMessage()    {}
func (*Policy_Buffer) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{1, 2}
}
func (m *Policy_Buffer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Policy_Buffer.Unmarshal(m, b)
}
func (m *Policy_Buffer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Policy_Buffer.Marshal(b, m, deterministic)
}
func (dst *Policy_Buffer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Policy_Buffer.Merge(dst, src)
}
func (m *Policy_Buffer) XXX_Size() int {
	return xxx_messageInfo_Policy_Buffer.Size(m)
}
func (m *Policy_Buffer) XXX_DiscardUnknown() {
	xxx_messageInfo_Policy_Buffer.DiscardUnknown(m)
}

var xxx_messageInfo_Policy_Buffer proto.InternalMessageInfo

func (m *Policy_Buffer) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *Policy_Buffer) GetSize() uint32 {
	if m != nil {
		return m.Size
	}
	return 0
}

type SystemPolicy struct {
	Stats                *SystemPolicy_Stats `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *SystemPolicy) Reset()         { *m = SystemPolicy{} }
func (m *SystemPolicy) String() string { return proto.CompactTextString(m) }
func (*SystemPolicy) ProtoMessage()    {}
func (*SystemPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{2}
}
func (m *SystemPolicy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemPolicy.Unmarshal(m, b)
}
func (m *SystemPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemPolicy.Marshal(b, m, deterministic)
}
func (dst *SystemPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemPolicy.Merge(dst, src)
}
func (m *SystemPolicy) XXX_Size() int {
	return xxx_messageInfo_SystemPolicy.Size(m)
}
func (m *SystemPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_SystemPolicy proto.InternalMessageInfo

func (m *SystemPolicy) GetStats() *SystemPolicy_Stats {
	if m != nil {
		return m.Stats
	}
	return nil
}

type SystemPolicy_Stats struct {
	InboundUplink        bool     `protobuf:"varint,1,opt,name=inbound_uplink,json=inboundUplink" json:"inbound_uplink,omitempty"`
	InboundDownlink      bool     `protobuf:"varint,2,opt,name=inbound_downlink,json=inboundDownlink" json:"inbound_downlink,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemPolicy_Stats) Reset()         { *m = SystemPolicy_Stats{} }
func (m *SystemPolicy_Stats) String() string { return proto.CompactTextString(m) }
func (*SystemPolicy_Stats) ProtoMessage()    {}
func (*SystemPolicy_Stats) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{2, 0}
}
func (m *SystemPolicy_Stats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemPolicy_Stats.Unmarshal(m, b)
}
func (m *SystemPolicy_Stats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemPolicy_Stats.Marshal(b, m, deterministic)
}
func (dst *SystemPolicy_Stats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemPolicy_Stats.Merge(dst, src)
}
func (m *SystemPolicy_Stats) XXX_Size() int {
	return xxx_messageInfo_SystemPolicy_Stats.Size(m)
}
func (m *SystemPolicy_Stats) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemPolicy_Stats.DiscardUnknown(m)
}

var xxx_messageInfo_SystemPolicy_Stats proto.InternalMessageInfo

func (m *SystemPolicy_Stats) GetInboundUplink() bool {
	if m != nil {
		return m.InboundUplink
	}
	return false
}

func (m *SystemPolicy_Stats) GetInboundDownlink() bool {
	if m != nil {
		return m.InboundDownlink
	}
	return false
}

type Config struct {
	Level                map[uint32]*Policy `protobuf:"bytes,1,rep,name=level" json:"level,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	System               *SystemPolicy      `protobuf:"bytes,2,opt,name=system" json:"system,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_505638f2092d854e, []int{3}
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

func (m *Config) GetLevel() map[uint32]*Policy {
	if m != nil {
		return m.Level
	}
	return nil
}

func (m *Config) GetSystem() *SystemPolicy {
	if m != nil {
		return m.System
	}
	return nil
}

func init() {
	proto.RegisterType((*Second)(nil), "v2ray.core.app.policy.Second")
	proto.RegisterType((*Policy)(nil), "v2ray.core.app.policy.Policy")
	proto.RegisterType((*Policy_Timeout)(nil), "v2ray.core.app.policy.Policy.Timeout")
	proto.RegisterType((*Policy_Stats)(nil), "v2ray.core.app.policy.Policy.Stats")
	proto.RegisterType((*Policy_Buffer)(nil), "v2ray.core.app.policy.Policy.Buffer")
	proto.RegisterType((*SystemPolicy)(nil), "v2ray.core.app.policy.SystemPolicy")
	proto.RegisterType((*SystemPolicy_Stats)(nil), "v2ray.core.app.policy.SystemPolicy.Stats")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.policy.Config")
	proto.RegisterMapType((map[uint32]*Policy)(nil), "v2ray.core.app.policy.Config.LevelEntry")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/policy/config.proto", fileDescriptor_config_505638f2092d854e)
}

var fileDescriptor_config_505638f2092d854e = []byte{
	// 523 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xeb, 0x6a, 0x13, 0x41,
	0x14, 0xc7, 0xd9, 0x5c, 0x36, 0xf5, 0x24, 0xdb, 0x96, 0xc1, 0xc2, 0xba, 0xa0, 0x96, 0xd4, 0x4a,
	0xfa, 0x65, 0x03, 0x29, 0x88, 0x5a, 0xad, 0x18, 0x2f, 0x20, 0x28, 0x96, 0x89, 0x17, 0xf4, 0x4b,
	0xd8, 0xec, 0x9e, 0xd8, 0x25, 0x93, 0x99, 0x65, 0x2f, 0x91, 0xf5, 0x31, 0x7c, 0x8c, 0x3e, 0x54,
	0x9f, 0x45, 0x76, 0x2e, 0xa6, 0x95, 0x26, 0xf1, 0xdb, 0xcc, 0xe1, 0xf7, 0xff, 0x33, 0xff, 0xb3,
	0xe7, 0x2c, 0x3c, 0x5c, 0x0c, 0xd2, 0xa0, 0xf4, 0x43, 0x31, 0xef, 0x87, 0x22, 0xc5, 0x7e, 0x90,
	0x24, 0xfd, 0x44, 0xb0, 0x38, 0x2c, 0xfb, 0xa1, 0xe0, 0xd3, 0xf8, 0x87, 0x9f, 0xa4, 0x22, 0x17,
	0x64, 0xcf, 0x70, 0x29, 0xfa, 0x41, 0x92, 0xf8, 0x8a, 0xe9, 0xde, 0x03, 0x7b, 0x84, 0xa1, 0xe0,
	0x11, 0xb9, 0x0d, 0xcd, 0x45, 0xc0, 0x0a, 0x74, 0xad, 0x7d, 0xab, 0xe7, 0x50, 0x75, 0xe9, 0x5e,
	0x36, 0xc0, 0x3e, 0x93, 0x28, 0x79, 0x01, 0xad, 0x3c, 0x9e, 0xa3, 0x28, 0x72, 0x89, 0xb4, 0x07,
	0x87, 0xfe, 0x8d, 0x9e, 0xbe, 0xe2, 0xfd, 0x4f, 0x0a, 0xa6, 0x46, 0x45, 0x9e, 0x40, 0x33, 0xcb,
	0x83, 0x3c, 0x73, 0x6b, 0x52, 0x7e, 0xb0, 0x5e, 0x3e, 0xaa, 0x50, 0xaa, 0x14, 0xe4, 0x19, 0xd8,
	0x93, 0x62, 0x3a, 0xc5, 0xd4, 0xad, 0x4b, 0xed, 0x83, 0xf5, 0xda, 0xa1, 0x64, 0xa9, 0xd6, 0x78,
	0xbf, 0x6b, 0xd0, 0xd2, 0xaf, 0x21, 0x27, 0x70, 0xeb, 0x3c, 0xe0, 0x51, 0x76, 0x1e, 0xcc, 0x50,
	0xe7, 0xb8, 0xbb, 0xc2, 0x4c, 0x35, 0x86, 0x2e, 0x79, 0xf2, 0x16, 0x76, 0x42, 0xc1, 0x39, 0x86,
	0x79, 0x2c, 0xf8, 0x38, 0x8e, 0x18, 0xea, 0x2c, 0x1b, 0x2c, 0xb6, 0x97, 0xaa, 0x77, 0x11, 0x43,
	0x72, 0x0a, 0xed, 0x22, 0x61, 0x31, 0x9f, 0x8d, 0x05, 0x67, 0xa5, 0xce, 0xb4, 0xc1, 0x03, 0x94,
	0xe2, 0x23, 0x67, 0x25, 0x19, 0x82, 0x13, 0x89, 0x9f, 0x7c, 0xe9, 0xd0, 0xf8, 0x1f, 0x87, 0x8e,
	0xd1, 0x54, 0x1e, 0xde, 0x07, 0x68, 0xca, 0x16, 0x93, 0xfb, 0xd0, 0x2e, 0x32, 0x4c, 0xc7, 0xca,
	0x5f, 0xf6, 0x64, 0x8b, 0x42, 0x55, 0xfa, 0x2c, 0x2b, 0xe4, 0x00, 0x1c, 0x09, 0x18, 0xb9, 0xcc,
	0xbc, 0x45, 0x3b, 0x55, 0xf1, 0xb5, 0xae, 0x79, 0x8f, 0xc0, 0x56, 0x5d, 0x27, 0x2e, 0xb4, 0x90,
	0x07, 0x13, 0x86, 0x91, 0xf6, 0x32, 0x57, 0x42, 0xa0, 0x91, 0xc5, 0xbf, 0x54, 0xcf, 0x1c, 0x2a,
	0xcf, 0xdd, 0x0b, 0x0b, 0x3a, 0xa3, 0x32, 0xcb, 0x71, 0xfe, 0x77, 0xcc, 0xf4, 0x94, 0xa8, 0x8f,
	0x73, 0xb4, 0x2a, 0xd3, 0x15, 0xcd, 0xb5, 0x59, 0xf1, 0xbe, 0x99, 0x60, 0x87, 0xb0, 0x1d, 0xf3,
	0x89, 0x28, 0x78, 0x74, 0x3d, 0x9b, 0xa3, 0xab, 0x3a, 0xde, 0x11, 0xec, 0x1a, 0xec, 0x9f, 0x84,
	0x3b, 0xba, 0x6e, 0x42, 0x76, 0x2f, 0x2d, 0xb0, 0x5f, 0xc9, 0xad, 0x22, 0xa7, 0xd0, 0x64, 0xb8,
	0x40, 0xe6, 0x5a, 0xfb, 0xf5, 0x5e, 0x7b, 0xd0, 0x5b, 0xf1, 0x4c, 0x45, 0xfb, 0xef, 0x2b, 0xf4,
	0x0d, 0xcf, 0xd3, 0x92, 0x2a, 0x19, 0x39, 0x01, 0x3b, 0x93, 0x11, 0x36, 0x6c, 0xc3, 0xd5, 0x9c,
	0x54, 0x4b, 0xbc, 0xaf, 0x00, 0x4b, 0x47, 0xb2, 0x0b, 0xf5, 0x19, 0x96, 0x7a, 0x6f, 0xab, 0x23,
	0x39, 0x36, 0xbb, 0xbc, 0x7e, 0x3a, 0xb5, 0xab, 0x62, 0x9f, 0xd6, 0x1e, 0x5b, 0xc3, 0xe7, 0x70,
	0x27, 0x14, 0xf3, 0x9b, 0xf1, 0x33, 0xeb, 0xbb, 0xad, 0x4e, 0x17, 0xb5, 0xbd, 0x2f, 0x03, 0x1a,
	0x54, 0xe9, 0x52, 0xf4, 0x5f, 0x26, 0x89, 0x76, 0x9a, 0xd8, 0xf2, 0x5f, 0x73, 0xfc, 0x27, 0x00,
	0x00, 0xff, 0xff, 0x2c, 0x2e, 0xe6, 0xcf, 0x95, 0x04, 0x00, 0x00,
}
