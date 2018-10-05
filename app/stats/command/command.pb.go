package command

import (
	"context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type GetStatsRequest struct {
	// Name of the stat counter.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Whether or not to reset the counter to fetching its value.
	Reset_               bool     `protobuf:"varint,2,opt,name=reset,proto3" json:"reset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStatsRequest) Reset()         { *m = GetStatsRequest{} }
func (m *GetStatsRequest) String() string { return proto.CompactTextString(m) }
func (*GetStatsRequest) ProtoMessage()    {}
func (*GetStatsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c902411c4948f26b, []int{0}
}

func (m *GetStatsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStatsRequest.Unmarshal(m, b)
}
func (m *GetStatsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStatsRequest.Marshal(b, m, deterministic)
}
func (m *GetStatsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStatsRequest.Merge(m, src)
}
func (m *GetStatsRequest) XXX_Size() int {
	return xxx_messageInfo_GetStatsRequest.Size(m)
}
func (m *GetStatsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStatsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetStatsRequest proto.InternalMessageInfo

func (m *GetStatsRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetStatsRequest) GetReset_() bool {
	if m != nil {
		return m.Reset_
	}
	return false
}

type Stat struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                int64    `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stat) Reset()         { *m = Stat{} }
func (m *Stat) String() string { return proto.CompactTextString(m) }
func (*Stat) ProtoMessage()    {}
func (*Stat) Descriptor() ([]byte, []int) {
	return fileDescriptor_c902411c4948f26b, []int{1}
}

func (m *Stat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stat.Unmarshal(m, b)
}
func (m *Stat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stat.Marshal(b, m, deterministic)
}
func (m *Stat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stat.Merge(m, src)
}
func (m *Stat) XXX_Size() int {
	return xxx_messageInfo_Stat.Size(m)
}
func (m *Stat) XXX_DiscardUnknown() {
	xxx_messageInfo_Stat.DiscardUnknown(m)
}

var xxx_messageInfo_Stat proto.InternalMessageInfo

func (m *Stat) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Stat) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type GetStatsResponse struct {
	Stat                 *Stat    `protobuf:"bytes,1,opt,name=stat,proto3" json:"stat,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStatsResponse) Reset()         { *m = GetStatsResponse{} }
func (m *GetStatsResponse) String() string { return proto.CompactTextString(m) }
func (*GetStatsResponse) ProtoMessage()    {}
func (*GetStatsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c902411c4948f26b, []int{2}
}

func (m *GetStatsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStatsResponse.Unmarshal(m, b)
}
func (m *GetStatsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStatsResponse.Marshal(b, m, deterministic)
}
func (m *GetStatsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStatsResponse.Merge(m, src)
}
func (m *GetStatsResponse) XXX_Size() int {
	return xxx_messageInfo_GetStatsResponse.Size(m)
}
func (m *GetStatsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStatsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetStatsResponse proto.InternalMessageInfo

func (m *GetStatsResponse) GetStat() *Stat {
	if m != nil {
		return m.Stat
	}
	return nil
}

type QueryStatsRequest struct {
	Pattern              string   `protobuf:"bytes,1,opt,name=pattern,proto3" json:"pattern,omitempty"`
	Reset_               bool     `protobuf:"varint,2,opt,name=reset,proto3" json:"reset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryStatsRequest) Reset()         { *m = QueryStatsRequest{} }
func (m *QueryStatsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryStatsRequest) ProtoMessage()    {}
func (*QueryStatsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c902411c4948f26b, []int{3}
}

func (m *QueryStatsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryStatsRequest.Unmarshal(m, b)
}
func (m *QueryStatsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryStatsRequest.Marshal(b, m, deterministic)
}
func (m *QueryStatsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryStatsRequest.Merge(m, src)
}
func (m *QueryStatsRequest) XXX_Size() int {
	return xxx_messageInfo_QueryStatsRequest.Size(m)
}
func (m *QueryStatsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryStatsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryStatsRequest proto.InternalMessageInfo

func (m *QueryStatsRequest) GetPattern() string {
	if m != nil {
		return m.Pattern
	}
	return ""
}

func (m *QueryStatsRequest) GetReset_() bool {
	if m != nil {
		return m.Reset_
	}
	return false
}

type QueryStatsResponse struct {
	Stat                 []*Stat  `protobuf:"bytes,1,rep,name=stat,proto3" json:"stat,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryStatsResponse) Reset()         { *m = QueryStatsResponse{} }
func (m *QueryStatsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryStatsResponse) ProtoMessage()    {}
func (*QueryStatsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c902411c4948f26b, []int{4}
}

func (m *QueryStatsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryStatsResponse.Unmarshal(m, b)
}
func (m *QueryStatsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryStatsResponse.Marshal(b, m, deterministic)
}
func (m *QueryStatsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryStatsResponse.Merge(m, src)
}
func (m *QueryStatsResponse) XXX_Size() int {
	return xxx_messageInfo_QueryStatsResponse.Size(m)
}
func (m *QueryStatsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryStatsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryStatsResponse proto.InternalMessageInfo

func (m *QueryStatsResponse) GetStat() []*Stat {
	if m != nil {
		return m.Stat
	}
	return nil
}

type Config struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_c902411c4948f26b, []int{5}
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

func init() {
	proto.RegisterType((*GetStatsRequest)(nil), "v2ray.core.app.stats.command.GetStatsRequest")
	proto.RegisterType((*Stat)(nil), "v2ray.core.app.stats.command.Stat")
	proto.RegisterType((*GetStatsResponse)(nil), "v2ray.core.app.stats.command.GetStatsResponse")
	proto.RegisterType((*QueryStatsRequest)(nil), "v2ray.core.app.stats.command.QueryStatsRequest")
	proto.RegisterType((*QueryStatsResponse)(nil), "v2ray.core.app.stats.command.QueryStatsResponse")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.stats.command.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/stats/command/command.proto", fileDescriptor_c902411c4948f26b)
}

var fileDescriptor_c902411c4948f26b = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x49, 0x5b, 0xda, 0x72, 0x20, 0x01, 0x16, 0x43, 0x55, 0x75, 0x88, 0x3c, 0x75, 0xc1,
	0xa9, 0x82, 0xc4, 0xc2, 0x04, 0x19, 0x90, 0x50, 0x07, 0x70, 0x25, 0x06, 0x36, 0x13, 0x0e, 0x54,
	0x41, 0x62, 0xd7, 0x76, 0x22, 0xe5, 0x95, 0x78, 0x38, 0x9e, 0x01, 0xc5, 0x49, 0x54, 0xa0, 0x6a,
	0x54, 0xa6, 0xdc, 0xc5, 0xff, 0x77, 0xf7, 0xdf, 0xd9, 0xc0, 0xf2, 0x50, 0x8b, 0x82, 0xc5, 0x32,
	0x09, 0x62, 0xa9, 0x31, 0x10, 0x4a, 0x05, 0xc6, 0x0a, 0x6b, 0x82, 0x58, 0x26, 0x89, 0x48, 0x5f,
	0x9a, 0x2f, 0x53, 0x5a, 0x5a, 0x49, 0x26, 0x8d, 0x5e, 0x23, 0x13, 0x4a, 0x31, 0xa7, 0x65, 0xb5,
	0x86, 0x5e, 0xc1, 0xf1, 0x2d, 0xda, 0x45, 0xf9, 0x8f, 0xe3, 0x2a, 0x43, 0x63, 0x09, 0x81, 0x5e,
	0x2a, 0x12, 0x1c, 0x79, 0xbe, 0x37, 0x3d, 0xe0, 0x2e, 0x26, 0x67, 0xb0, 0xaf, 0xd1, 0xa0, 0x1d,
	0x75, 0x7c, 0x6f, 0x3a, 0xe4, 0x55, 0x42, 0x67, 0xd0, 0x2b, 0xc9, 0x6d, 0x44, 0x2e, 0x3e, 0x32,
	0x74, 0x44, 0x97, 0x57, 0x09, 0xbd, 0x83, 0x93, 0x75, 0x3b, 0xa3, 0x64, 0x6a, 0x90, 0x5c, 0x42,
	0xaf, 0xf4, 0xe4, 0xe8, 0xc3, 0x90, 0xb2, 0x36, 0xbf, 0xac, 0x44, 0xb9, 0xd3, 0xd3, 0x08, 0x4e,
	0x1f, 0x32, 0xd4, 0xc5, 0x2f, 0xf3, 0x23, 0x18, 0x28, 0x61, 0x2d, 0xea, 0xb4, 0x76, 0xd3, 0xa4,
	0x5b, 0x46, 0x98, 0x03, 0xf9, 0x59, 0x64, 0xc3, 0x52, 0xf7, 0x5f, 0x96, 0x86, 0xd0, 0x8f, 0x64,
	0xfa, 0xba, 0x7c, 0x0b, 0xbf, 0x3c, 0x38, 0x72, 0x35, 0x17, 0xa8, 0xf3, 0x65, 0x8c, 0xe4, 0x1d,
	0x86, 0xcd, 0xe4, 0xe4, 0xbc, 0xbd, 0xe0, 0x9f, 0x0b, 0x19, 0xb3, 0x5d, 0xe5, 0x95, 0x7b, 0xba,
	0x47, 0x56, 0x00, 0xeb, 0xa9, 0x48, 0xd0, 0xce, 0x6f, 0x2c, 0x71, 0x3c, 0xdb, 0x1d, 0x68, 0x5a,
	0xde, 0xcc, 0xc1, 0x8f, 0x65, 0xd2, 0x0a, 0xde, 0x7b, 0x4f, 0x83, 0x3a, 0xfc, 0xec, 0x4c, 0x1e,
	0x43, 0x2e, 0x0a, 0x16, 0x95, 0xca, 0x6b, 0xa5, 0xdc, 0x16, 0x0d, 0x8b, 0xaa, 0xe3, 0xe7, 0xbe,
	0x7b, 0xbb, 0x17, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x9e, 0xb8, 0xeb, 0xed, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StatsServiceClient is the client API for StatsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StatsServiceClient interface {
	GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error)
	QueryStats(ctx context.Context, in *QueryStatsRequest, opts ...grpc.CallOption) (*QueryStatsResponse, error)
}

type statsServiceClient struct {
	cc *grpc.ClientConn
}

func NewStatsServiceClient(cc *grpc.ClientConn) StatsServiceClient {
	return &statsServiceClient{cc}
}

func (c *statsServiceClient) GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.stats.command.StatsService/GetStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statsServiceClient) QueryStats(ctx context.Context, in *QueryStatsRequest, opts ...grpc.CallOption) (*QueryStatsResponse, error) {
	out := new(QueryStatsResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.stats.command.StatsService/QueryStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatsServiceServer is the server API for StatsService service.
type StatsServiceServer interface {
	GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error)
	QueryStats(context.Context, *QueryStatsRequest) (*QueryStatsResponse, error)
}

func RegisterStatsServiceServer(s *grpc.Server, srv StatsServiceServer) {
	s.RegisterService(&_StatsService_serviceDesc, srv)
}

func _StatsService_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServiceServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.stats.command.StatsService/GetStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServiceServer).GetStats(ctx, req.(*GetStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatsService_QueryStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServiceServer).QueryStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.stats.command.StatsService/QueryStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServiceServer).QueryStats(ctx, req.(*QueryStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StatsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v2ray.core.app.stats.command.StatsService",
	HandlerType: (*StatsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStats",
			Handler:    _StatsService_GetStats_Handler,
		},
		{
			MethodName: "QueryStats",
			Handler:    _StatsService_QueryStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v2ray.com/core/app/stats/command/command.proto",
}
