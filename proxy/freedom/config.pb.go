package freedom

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	protocol "v2ray.com/core/common/protocol"
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

type Config_DomainStrategy int32

const (
	Config_AS_IS   Config_DomainStrategy = 0
	Config_USE_IP  Config_DomainStrategy = 1
	Config_USE_IP4 Config_DomainStrategy = 2
	Config_USE_IP6 Config_DomainStrategy = 3
)

// Enum value maps for Config_DomainStrategy.
var (
	Config_DomainStrategy_name = map[int32]string{
		0: "AS_IS",
		1: "USE_IP",
		2: "USE_IP4",
		3: "USE_IP6",
	}
	Config_DomainStrategy_value = map[string]int32{
		"AS_IS":   0,
		"USE_IP":  1,
		"USE_IP4": 2,
		"USE_IP6": 3,
	}
)

func (x Config_DomainStrategy) Enum() *Config_DomainStrategy {
	p := new(Config_DomainStrategy)
	*p = x
	return p
}

func (x Config_DomainStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Config_DomainStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_v2ray_com_core_proxy_freedom_config_proto_enumTypes[0].Descriptor()
}

func (Config_DomainStrategy) Type() protoreflect.EnumType {
	return &file_v2ray_com_core_proxy_freedom_config_proto_enumTypes[0]
}

func (x Config_DomainStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Config_DomainStrategy.Descriptor instead.
func (Config_DomainStrategy) EnumDescriptor() ([]byte, []int) {
	return file_v2ray_com_core_proxy_freedom_config_proto_rawDescGZIP(), []int{1, 0}
}

type DestinationOverride struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Server *protocol.ServerEndpoint `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
}

func (x *DestinationOverride) Reset() {
	*x = DestinationOverride{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_proxy_freedom_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DestinationOverride) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestinationOverride) ProtoMessage() {}

func (x *DestinationOverride) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_proxy_freedom_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestinationOverride.ProtoReflect.Descriptor instead.
func (*DestinationOverride) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_proxy_freedom_config_proto_rawDescGZIP(), []int{0}
}

func (x *DestinationOverride) GetServer() *protocol.ServerEndpoint {
	if x != nil {
		return x.Server
	}
	return nil
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DomainStrategy Config_DomainStrategy `protobuf:"varint,1,opt,name=domain_strategy,json=domainStrategy,proto3,enum=v2ray.core.proxy.freedom.Config_DomainStrategy" json:"domain_strategy,omitempty"`
	// Deprecated: Do not use.
	Timeout             uint32               `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	DestinationOverride *DestinationOverride `protobuf:"bytes,3,opt,name=destination_override,json=destinationOverride,proto3" json:"destination_override,omitempty"`
	UserLevel           uint32               `protobuf:"varint,4,opt,name=user_level,json=userLevel,proto3" json:"user_level,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_proxy_freedom_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_proxy_freedom_config_proto_msgTypes[1]
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
	return file_v2ray_com_core_proxy_freedom_config_proto_rawDescGZIP(), []int{1}
}

func (x *Config) GetDomainStrategy() Config_DomainStrategy {
	if x != nil {
		return x.DomainStrategy
	}
	return Config_AS_IS
}

// Deprecated: Do not use.
func (x *Config) GetTimeout() uint32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *Config) GetDestinationOverride() *DestinationOverride {
	if x != nil {
		return x.DestinationOverride
	}
	return nil
}

func (x *Config) GetUserLevel() uint32 {
	if x != nil {
		return x.UserLevel
	}
	return 0
}

var File_v2ray_com_core_proxy_freedom_config_proto protoreflect.FileDescriptor

var file_v2ray_com_core_proxy_freedom_config_proto_rawDesc = []byte{
	0x0a, 0x29, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x66, 0x72, 0x65, 0x65, 0x64, 0x6f, 0x6d, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x76, 0x32, 0x72,
	0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x66, 0x72,
	0x65, 0x65, 0x64, 0x6f, 0x6d, 0x1a, 0x30, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x70, 0x65,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x13, 0x44, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x12, 0x42,
	0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a,
	0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x22, 0xc4, 0x02, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x58, 0x0a,
	0x0f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x66, 0x72, 0x65, 0x65, 0x64, 0x6f,
	0x6d, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x53,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x0e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x53,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x1c, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x02, 0x18, 0x01, 0x52, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x60, 0x0a, 0x14, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x66, 0x72, 0x65, 0x65, 0x64, 0x6f, 0x6d, 0x2e, 0x44,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x76, 0x65, 0x72, 0x72, 0x69,
	0x64, 0x65, 0x52, 0x13, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f,
	0x76, 0x65, 0x72, 0x72, 0x69, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x75, 0x73, 0x65,
	0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x22, 0x41, 0x0a, 0x0e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x53, 0x5f, 0x49,
	0x53, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x53, 0x45, 0x5f, 0x49, 0x50, 0x10, 0x01, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x53, 0x45, 0x5f, 0x49, 0x50, 0x34, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x53, 0x45, 0x5f, 0x49, 0x50, 0x36, 0x10, 0x03, 0x42, 0x44, 0x0a, 0x1c, 0x63, 0x6f, 0x6d,
	0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x78,
	0x79, 0x2e, 0x66, 0x72, 0x65, 0x65, 0x64, 0x6f, 0x6d, 0x50, 0x01, 0x5a, 0x07, 0x66, 0x72, 0x65,
	0x65, 0x64, 0x6f, 0x6d, 0xaa, 0x02, 0x18, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72,
	0x65, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x46, 0x72, 0x65, 0x65, 0x64, 0x6f, 0x6d, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v2ray_com_core_proxy_freedom_config_proto_rawDescOnce sync.Once
	file_v2ray_com_core_proxy_freedom_config_proto_rawDescData = file_v2ray_com_core_proxy_freedom_config_proto_rawDesc
)

func file_v2ray_com_core_proxy_freedom_config_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_proxy_freedom_config_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_proxy_freedom_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_proxy_freedom_config_proto_rawDescData)
	})
	return file_v2ray_com_core_proxy_freedom_config_proto_rawDescData
}

var file_v2ray_com_core_proxy_freedom_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v2ray_com_core_proxy_freedom_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v2ray_com_core_proxy_freedom_config_proto_goTypes = []interface{}{
	(Config_DomainStrategy)(0),      // 0: v2ray.core.proxy.freedom.Config.DomainStrategy
	(*DestinationOverride)(nil),     // 1: v2ray.core.proxy.freedom.DestinationOverride
	(*Config)(nil),                  // 2: v2ray.core.proxy.freedom.Config
	(*protocol.ServerEndpoint)(nil), // 3: v2ray.core.common.protocol.ServerEndpoint
}
var file_v2ray_com_core_proxy_freedom_config_proto_depIdxs = []int32{
	3, // 0: v2ray.core.proxy.freedom.DestinationOverride.server:type_name -> v2ray.core.common.protocol.ServerEndpoint
	0, // 1: v2ray.core.proxy.freedom.Config.domain_strategy:type_name -> v2ray.core.proxy.freedom.Config.DomainStrategy
	1, // 2: v2ray.core.proxy.freedom.Config.destination_override:type_name -> v2ray.core.proxy.freedom.DestinationOverride
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_proxy_freedom_config_proto_init() }
func file_v2ray_com_core_proxy_freedom_config_proto_init() {
	if File_v2ray_com_core_proxy_freedom_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_proxy_freedom_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DestinationOverride); i {
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
		file_v2ray_com_core_proxy_freedom_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_v2ray_com_core_proxy_freedom_config_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v2ray_com_core_proxy_freedom_config_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_proxy_freedom_config_proto_depIdxs,
		EnumInfos:         file_v2ray_com_core_proxy_freedom_config_proto_enumTypes,
		MessageInfos:      file_v2ray_com_core_proxy_freedom_config_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_proxy_freedom_config_proto = out.File
	file_v2ray_com_core_proxy_freedom_config_proto_rawDesc = nil
	file_v2ray_com_core_proxy_freedom_config_proto_goTypes = nil
	file_v2ray_com_core_proxy_freedom_config_proto_depIdxs = nil
}
