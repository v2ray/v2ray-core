package command

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_log_command_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_log_command_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_log_command_config_proto_rawDescGZIP(), []int{0}
}

type RestartLoggerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RestartLoggerRequest) Reset() {
	*x = RestartLoggerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_log_command_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestartLoggerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestartLoggerRequest) ProtoMessage() {}

func (x *RestartLoggerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_log_command_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestartLoggerRequest.ProtoReflect.Descriptor instead.
func (*RestartLoggerRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_log_command_config_proto_rawDescGZIP(), []int{1}
}

type RestartLoggerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RestartLoggerResponse) Reset() {
	*x = RestartLoggerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_log_command_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestartLoggerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestartLoggerResponse) ProtoMessage() {}

func (x *RestartLoggerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_log_command_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestartLoggerResponse.ProtoReflect.Descriptor instead.
func (*RestartLoggerResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_log_command_config_proto_rawDescGZIP(), []int{2}
}

var File_v2ray_com_core_app_log_command_config_proto protoreflect.FileDescriptor

var file_v2ray_com_core_app_log_command_config_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x61, 0x70, 0x70, 0x2f, 0x6c, 0x6f, 0x67, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x76,
	0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6c, 0x6f,
	0x67, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x08, 0x0a, 0x06, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x22, 0x16, 0x0a, 0x14, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4c, 0x6f,
	0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x17, 0x0a, 0x15, 0x52,
	0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x87, 0x01, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x76, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x12, 0x30, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4c, 0x6f, 0x67, 0x67,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4c, 0x6f,
	0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x48,
	0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x6c, 0x6f, 0x67, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x50, 0x01, 0x5a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0xaa, 0x02, 0x1a, 0x56, 0x32,
	0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x2e, 0x4c, 0x6f, 0x67,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v2ray_com_core_app_log_command_config_proto_rawDescOnce sync.Once
	file_v2ray_com_core_app_log_command_config_proto_rawDescData = file_v2ray_com_core_app_log_command_config_proto_rawDesc
)

func file_v2ray_com_core_app_log_command_config_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_app_log_command_config_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_app_log_command_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_app_log_command_config_proto_rawDescData)
	})
	return file_v2ray_com_core_app_log_command_config_proto_rawDescData
}

var file_v2ray_com_core_app_log_command_config_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_v2ray_com_core_app_log_command_config_proto_goTypes = []interface{}{
	(*Config)(nil),                // 0: v2ray.core.app.log.command.Config
	(*RestartLoggerRequest)(nil),  // 1: v2ray.core.app.log.command.RestartLoggerRequest
	(*RestartLoggerResponse)(nil), // 2: v2ray.core.app.log.command.RestartLoggerResponse
}
var file_v2ray_com_core_app_log_command_config_proto_depIdxs = []int32{
	1, // 0: v2ray.core.app.log.command.LoggerService.RestartLogger:input_type -> v2ray.core.app.log.command.RestartLoggerRequest
	2, // 1: v2ray.core.app.log.command.LoggerService.RestartLogger:output_type -> v2ray.core.app.log.command.RestartLoggerResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_app_log_command_config_proto_init() }
func file_v2ray_com_core_app_log_command_config_proto_init() {
	if File_v2ray_com_core_app_log_command_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_app_log_command_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v2ray_com_core_app_log_command_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestartLoggerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_v2ray_com_core_app_log_command_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestartLoggerResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v2ray_com_core_app_log_command_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v2ray_com_core_app_log_command_config_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_app_log_command_config_proto_depIdxs,
		MessageInfos:      file_v2ray_com_core_app_log_command_config_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_app_log_command_config_proto = out.File
	file_v2ray_com_core_app_log_command_config_proto_rawDesc = nil
	file_v2ray_com_core_app_log_command_config_proto_goTypes = nil
	file_v2ray_com_core_app_log_command_config_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LoggerServiceClient is the client API for LoggerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LoggerServiceClient interface {
	RestartLogger(ctx context.Context, in *RestartLoggerRequest, opts ...grpc.CallOption) (*RestartLoggerResponse, error)
}

type loggerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoggerServiceClient(cc grpc.ClientConnInterface) LoggerServiceClient {
	return &loggerServiceClient{cc}
}

func (c *loggerServiceClient) RestartLogger(ctx context.Context, in *RestartLoggerRequest, opts ...grpc.CallOption) (*RestartLoggerResponse, error) {
	out := new(RestartLoggerResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.log.command.LoggerService/RestartLogger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoggerServiceServer is the server API for LoggerService service.
type LoggerServiceServer interface {
	RestartLogger(context.Context, *RestartLoggerRequest) (*RestartLoggerResponse, error)
}

// UnimplementedLoggerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLoggerServiceServer struct {
}

func (*UnimplementedLoggerServiceServer) RestartLogger(context.Context, *RestartLoggerRequest) (*RestartLoggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartLogger not implemented")
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
