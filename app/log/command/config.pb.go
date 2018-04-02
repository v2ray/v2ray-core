package command

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	"context"

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

type Config struct {
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type RestartLoggerRequest struct {
}

func (m *RestartLoggerRequest) Reset()                    { *m = RestartLoggerRequest{} }
func (m *RestartLoggerRequest) String() string            { return proto.CompactTextString(m) }
func (*RestartLoggerRequest) ProtoMessage()               {}
func (*RestartLoggerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type RestartLoggerResponse struct {
}

func (m *RestartLoggerResponse) Reset()                    { *m = RestartLoggerResponse{} }
func (m *RestartLoggerResponse) String() string            { return proto.CompactTextString(m) }
func (*RestartLoggerResponse) ProtoMessage()               {}
func (*RestartLoggerResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*Config)(nil), "v2ray.core.app.log.command.Config")
	proto.RegisterType((*RestartLoggerRequest)(nil), "v2ray.core.app.log.command.RestartLoggerRequest")
	proto.RegisterType((*RestartLoggerResponse)(nil), "v2ray.core.app.log.command.RestartLoggerResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for LoggerService service

type LoggerServiceClient interface {
	RestartLogger(ctx context.Context, in *RestartLoggerRequest, opts ...grpc.CallOption) (*RestartLoggerResponse, error)
}

type loggerServiceClient struct {
	cc *grpc.ClientConn
}

func NewLoggerServiceClient(cc *grpc.ClientConn) LoggerServiceClient {
	return &loggerServiceClient{cc}
}

func (c *loggerServiceClient) RestartLogger(ctx context.Context, in *RestartLoggerRequest, opts ...grpc.CallOption) (*RestartLoggerResponse, error) {
	out := new(RestartLoggerResponse)
	err := grpc.Invoke(ctx, "/v2ray.core.app.log.command.LoggerService/RestartLogger", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LoggerService service

type LoggerServiceServer interface {
	RestartLogger(context.Context, *RestartLoggerRequest) (*RestartLoggerResponse, error)
}

func RegisterLoggerServiceServer(s *grpc.Server, srv LoggerServiceServer) {
	s.RegisterService(&_LoggerService_serviceDesc, srv)
}

func _LoggerService_RestartLogger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestartLoggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggerServiceServer).RestartLogger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.log.command.LoggerService/RestartLogger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggerServiceServer).RestartLogger(ctx, req.(*RestartLoggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LoggerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v2ray.core.app.log.command.LoggerService",
	HandlerType: (*LoggerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RestartLogger",
			Handler:    _LoggerService_RestartLogger_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v2ray.com/core/app/log/command/config.proto",
}

func init() { proto.RegisterFile("v2ray.com/core/app/log/command/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2e, 0x33, 0x2a, 0x4a,
	0xac, 0xd4, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0x2c, 0x28, 0xd0, 0xcf,
	0xc9, 0x4f, 0xd7, 0x4f, 0xce, 0xcf, 0xcd, 0x4d, 0xcc, 0x4b, 0xd1, 0x4f, 0xce, 0xcf, 0x4b, 0xcb,
	0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x82, 0x29, 0x2e, 0x4a, 0xd5, 0x4b, 0x2c,
	0x28, 0xd0, 0xcb, 0xc9, 0x4f, 0xd7, 0x83, 0x2a, 0x54, 0xe2, 0xe0, 0x62, 0x73, 0x06, 0xab, 0x55,
	0x12, 0xe3, 0x12, 0x09, 0x4a, 0x2d, 0x2e, 0x49, 0x2c, 0x2a, 0xf1, 0xc9, 0x4f, 0x4f, 0x4f, 0x2d,
	0x0a, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x51, 0x12, 0xe7, 0x12, 0x45, 0x13, 0x2f, 0x2e, 0xc8,
	0xcf, 0x2b, 0x4e, 0x35, 0x6a, 0x67, 0xe4, 0xe2, 0x85, 0x08, 0x05, 0xa7, 0x16, 0x95, 0x65, 0x26,
	0xa7, 0x0a, 0x95, 0x71, 0xf1, 0xa2, 0x28, 0x15, 0x32, 0xd0, 0xc3, 0x6d, 0xb5, 0x1e, 0x36, 0xdb,
	0xa4, 0x0c, 0x49, 0xd0, 0x01, 0x71, 0x87, 0x12, 0x83, 0x93, 0x07, 0x97, 0x5c, 0x72, 0x7e, 0x2e,
	0x1e, 0x9d, 0x01, 0x8c, 0x51, 0xec, 0x50, 0xe6, 0x2a, 0x26, 0xa9, 0x30, 0xa3, 0xa0, 0xc4, 0x4a,
	0x3d, 0x67, 0x90, 0x3a, 0xc7, 0x82, 0x02, 0x3d, 0x9f, 0xfc, 0x74, 0x3d, 0x67, 0x88, 0x64, 0x12,
	0x1b, 0x38, 0xc4, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x37, 0xc7, 0xfc, 0xda, 0x60, 0x01,
	0x00, 0x00,
}
