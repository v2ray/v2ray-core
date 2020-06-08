package net

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

// Endpoint of a network connection.
type Endpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network Network     `protobuf:"varint,1,opt,name=network,proto3,enum=v2ray.core.common.net.Network" json:"network,omitempty"`
	Address *IPOrDomain `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Port    uint32      `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Endpoint) Reset() {
	*x = Endpoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_common_net_destination_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Endpoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Endpoint) ProtoMessage() {}

func (x *Endpoint) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_common_net_destination_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Endpoint.ProtoReflect.Descriptor instead.
func (*Endpoint) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_common_net_destination_proto_rawDescGZIP(), []int{0}
}

func (x *Endpoint) GetNetwork() Network {
	if x != nil {
		return x.Network
	}
	return Network_Unknown
}

func (x *Endpoint) GetAddress() *IPOrDomain {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Endpoint) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

var File_v2ray_com_core_common_net_destination_proto protoreflect.FileDescriptor

var file_v2ray_com_core_common_net_destination_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x64, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x76,
	0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x6e, 0x65, 0x74, 0x1a, 0x27, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x27, 0x76,
	0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01, 0x0a, 0x08, 0x45, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x3b, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x49, 0x50, 0x4f, 0x72, 0x44, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x3a,
	0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x50, 0x01, 0x5a, 0x03, 0x6e,
	0x65, 0x74, 0xaa, 0x02, 0x15, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_v2ray_com_core_common_net_destination_proto_rawDescOnce sync.Once
	file_v2ray_com_core_common_net_destination_proto_rawDescData = file_v2ray_com_core_common_net_destination_proto_rawDesc
)

func file_v2ray_com_core_common_net_destination_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_common_net_destination_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_common_net_destination_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_common_net_destination_proto_rawDescData)
	})
	return file_v2ray_com_core_common_net_destination_proto_rawDescData
}

var file_v2ray_com_core_common_net_destination_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_v2ray_com_core_common_net_destination_proto_goTypes = []interface{}{
	(*Endpoint)(nil),   // 0: v2ray.core.common.net.Endpoint
	(Network)(0),       // 1: v2ray.core.common.net.Network
	(*IPOrDomain)(nil), // 2: v2ray.core.common.net.IPOrDomain
}
var file_v2ray_com_core_common_net_destination_proto_depIdxs = []int32{
	1, // 0: v2ray.core.common.net.Endpoint.network:type_name -> v2ray.core.common.net.Network
	2, // 1: v2ray.core.common.net.Endpoint.address:type_name -> v2ray.core.common.net.IPOrDomain
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_common_net_destination_proto_init() }
func file_v2ray_com_core_common_net_destination_proto_init() {
	if File_v2ray_com_core_common_net_destination_proto != nil {
		return
	}
	file_v2ray_com_core_common_net_network_proto_init()
	file_v2ray_com_core_common_net_address_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_common_net_destination_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Endpoint); i {
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
			RawDescriptor: file_v2ray_com_core_common_net_destination_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v2ray_com_core_common_net_destination_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_common_net_destination_proto_depIdxs,
		MessageInfos:      file_v2ray_com_core_common_net_destination_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_common_net_destination_proto = out.File
	file_v2ray_com_core_common_net_destination_proto_rawDesc = nil
	file_v2ray_com_core_common_net_destination_proto_goTypes = nil
	file_v2ray_com_core_common_net_destination_proto_depIdxs = nil
}
