package commander

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	serial "v2ray.com/core/common/serial"
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

type ApiCertTuple struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ca  string `protobuf:"bytes,1,opt,name=ca,proto3" json:"ca,omitempty"`
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Pub string `protobuf:"bytes,3,opt,name=pub,proto3" json:"pub,omitempty"`
}

func (x *ApiCertTuple) Reset() {
	*x = ApiCertTuple{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_commander_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiCertTuple) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiCertTuple) ProtoMessage() {}

func (x *ApiCertTuple) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_commander_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiCertTuple.ProtoReflect.Descriptor instead.
func (*ApiCertTuple) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_commander_config_proto_rawDescGZIP(), []int{0}
}

func (x *ApiCertTuple) GetCa() string {
	if x != nil {
		return x.Ca
	}
	return ""
}

func (x *ApiCertTuple) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ApiCertTuple) GetPub() string {
	if x != nil {
		return x.Pub
	}
	return ""
}

// Config is the settings for Commander.
type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Tag of the outbound handler that handles grpc connections.
	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	// Services that supported by this server. All services must implement Service interface.
	Service []*serial.TypedMessage `protobuf:"bytes,2,rep,name=service,proto3" json:"service,omitempty"`
	// Cert
	Cert *ApiCertTuple `protobuf:"bytes,3,opt,name=cert,proto3" json:"cert,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_commander_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_commander_config_proto_msgTypes[1]
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
	return file_v2ray_com_core_app_commander_config_proto_rawDescGZIP(), []int{1}
}

func (x *Config) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Config) GetService() []*serial.TypedMessage {
	if x != nil {
		return x.Service
	}
	return nil
}

func (x *Config) GetCert() *ApiCertTuple {
	if x != nil {
		return x.Cert
	}
	return nil
}

var File_v2ray_com_core_app_commander_config_proto protoreflect.FileDescriptor

var file_v2ray_com_core_app_commander_config_proto_rawDesc = []byte{
	0x0a, 0x29, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x61, 0x70, 0x70, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x76, 0x32, 0x72,
	0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x65, 0x72, 0x1a, 0x30, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72,
	0x69, 0x61, 0x6c, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42, 0x0a, 0x0c, 0x41, 0x70, 0x69, 0x43, 0x65,
	0x72, 0x74, 0x54, 0x75, 0x70, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x63, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x63, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x75, 0x62,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x75, 0x62, 0x22, 0x98, 0x01, 0x0a, 0x06,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x40, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x73, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x04, 0x63, 0x65,
	0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x65, 0x72, 0x2e, 0x41, 0x70, 0x69, 0x43, 0x65, 0x72, 0x74, 0x54, 0x75, 0x70, 0x6c, 0x65,
	0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x42, 0x46, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32,
	0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x50, 0x01, 0x5a, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x65, 0x72, 0xaa, 0x02, 0x18, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65,
	0x2e, 0x41, 0x70, 0x70, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v2ray_com_core_app_commander_config_proto_rawDescOnce sync.Once
	file_v2ray_com_core_app_commander_config_proto_rawDescData = file_v2ray_com_core_app_commander_config_proto_rawDesc
)

func file_v2ray_com_core_app_commander_config_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_app_commander_config_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_app_commander_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_app_commander_config_proto_rawDescData)
	})
	return file_v2ray_com_core_app_commander_config_proto_rawDescData
}

var file_v2ray_com_core_app_commander_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v2ray_com_core_app_commander_config_proto_goTypes = []interface{}{
	(*ApiCertTuple)(nil),        // 0: v2ray.core.app.commander.ApiCertTuple
	(*Config)(nil),              // 1: v2ray.core.app.commander.Config
	(*serial.TypedMessage)(nil), // 2: v2ray.core.common.serial.TypedMessage
}
var file_v2ray_com_core_app_commander_config_proto_depIdxs = []int32{
	2, // 0: v2ray.core.app.commander.Config.service:type_name -> v2ray.core.common.serial.TypedMessage
	0, // 1: v2ray.core.app.commander.Config.cert:type_name -> v2ray.core.app.commander.ApiCertTuple
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_app_commander_config_proto_init() }
func file_v2ray_com_core_app_commander_config_proto_init() {
	if File_v2ray_com_core_app_commander_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_app_commander_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiCertTuple); i {
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
		file_v2ray_com_core_app_commander_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v2ray_com_core_app_commander_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v2ray_com_core_app_commander_config_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_app_commander_config_proto_depIdxs,
		MessageInfos:      file_v2ray_com_core_app_commander_config_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_app_commander_config_proto = out.File
	file_v2ray_com_core_app_commander_config_proto_rawDesc = nil
	file_v2ray_com_core_app_commander_config_proto_goTypes = nil
	file_v2ray_com_core_app_commander_config_proto_depIdxs = nil
}
