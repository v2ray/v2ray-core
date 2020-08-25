package protocol

import (
	proto "github.com/golang/protobuf/proto"
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

type SecurityType int32

const (
	SecurityType_UNKNOWN           SecurityType = 0
	SecurityType_LEGACY            SecurityType = 1
	SecurityType_AUTO              SecurityType = 2
	SecurityType_AES128_GCM        SecurityType = 3
	SecurityType_CHACHA20_POLY1305 SecurityType = 4
	SecurityType_NONE              SecurityType = 5
)

// Enum value maps for SecurityType.
var (
	SecurityType_name = map[int32]string{
		0: "UNKNOWN",
		1: "LEGACY",
		2: "AUTO",
		3: "AES128_GCM",
		4: "CHACHA20_POLY1305",
		5: "NONE",
	}
	SecurityType_value = map[string]int32{
		"UNKNOWN":           0,
		"LEGACY":            1,
		"AUTO":              2,
		"AES128_GCM":        3,
		"CHACHA20_POLY1305": 4,
		"NONE":              5,
	}
)

func (x SecurityType) Enum() *SecurityType {
	p := new(SecurityType)
	*p = x
	return p
}

func (x SecurityType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecurityType) Descriptor() protoreflect.EnumDescriptor {
	return file_v2ray_com_core_common_protocol_headers_proto_enumTypes[0].Descriptor()
}

func (SecurityType) Type() protoreflect.EnumType {
	return &file_v2ray_com_core_common_protocol_headers_proto_enumTypes[0]
}

func (x SecurityType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecurityType.Descriptor instead.
func (SecurityType) EnumDescriptor() ([]byte, []int) {
	return file_v2ray_com_core_common_protocol_headers_proto_rawDescGZIP(), []int{0}
}

type SecurityConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type SecurityType `protobuf:"varint,1,opt,name=type,proto3,enum=v2ray.core.common.protocol.SecurityType" json:"type,omitempty"`
}

func (x *SecurityConfig) Reset() {
	*x = SecurityConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_common_protocol_headers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecurityConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecurityConfig) ProtoMessage() {}

func (x *SecurityConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_common_protocol_headers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecurityConfig.ProtoReflect.Descriptor instead.
func (*SecurityConfig) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_common_protocol_headers_proto_rawDescGZIP(), []int{0}
}

func (x *SecurityConfig) GetType() SecurityType {
	if x != nil {
		return x.Type
	}
	return SecurityType_UNKNOWN
}

var File_v2ray_com_core_common_protocol_headers_proto protoreflect.FileDescriptor

var file_v2ray_com_core_common_protocol_headers_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a,
	0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x4e, 0x0a, 0x0e, 0x53, 0x65,
	0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x76, 0x32, 0x72,
	0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x62, 0x0a, 0x0c, 0x53, 0x65,
	0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x45, 0x47, 0x41, 0x43,
	0x59, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x55, 0x54, 0x4f, 0x10, 0x02, 0x12, 0x0e, 0x0a,
	0x0a, 0x41, 0x45, 0x53, 0x31, 0x32, 0x38, 0x5f, 0x47, 0x43, 0x4d, 0x10, 0x03, 0x12, 0x15, 0x0a,
	0x11, 0x43, 0x48, 0x41, 0x43, 0x48, 0x41, 0x32, 0x30, 0x5f, 0x50, 0x4f, 0x4c, 0x59, 0x31, 0x33,
	0x30, 0x35, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x05, 0x42, 0x49,
	0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x50, 0x01, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0xaa, 0x02, 0x1a, 0x56,
	0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_v2ray_com_core_common_protocol_headers_proto_rawDescOnce sync.Once
	file_v2ray_com_core_common_protocol_headers_proto_rawDescData = file_v2ray_com_core_common_protocol_headers_proto_rawDesc
)

func file_v2ray_com_core_common_protocol_headers_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_common_protocol_headers_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_common_protocol_headers_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_common_protocol_headers_proto_rawDescData)
	})
	return file_v2ray_com_core_common_protocol_headers_proto_rawDescData
}

var file_v2ray_com_core_common_protocol_headers_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v2ray_com_core_common_protocol_headers_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_v2ray_com_core_common_protocol_headers_proto_goTypes = []interface{}{
	(SecurityType)(0),      // 0: v2ray.core.common.protocol.SecurityType
	(*SecurityConfig)(nil), // 1: v2ray.core.common.protocol.SecurityConfig
}
var file_v2ray_com_core_common_protocol_headers_proto_depIdxs = []int32{
	0, // 0: v2ray.core.common.protocol.SecurityConfig.type:type_name -> v2ray.core.common.protocol.SecurityType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_common_protocol_headers_proto_init() }
func file_v2ray_com_core_common_protocol_headers_proto_init() {
	if File_v2ray_com_core_common_protocol_headers_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_common_protocol_headers_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecurityConfig); i {
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
			RawDescriptor: file_v2ray_com_core_common_protocol_headers_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v2ray_com_core_common_protocol_headers_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_common_protocol_headers_proto_depIdxs,
		EnumInfos:         file_v2ray_com_core_common_protocol_headers_proto_enumTypes,
		MessageInfos:      file_v2ray_com_core_common_protocol_headers_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_common_protocol_headers_proto = out.File
	file_v2ray_com_core_common_protocol_headers_proto_rawDesc = nil
	file_v2ray_com_core_common_protocol_headers_proto_goTypes = nil
	file_v2ray_com_core_common_protocol_headers_proto_depIdxs = nil
}
