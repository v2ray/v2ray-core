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

// PortRange represents a range of ports.
type PortRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The port that this range starts from.
	From uint32 `protobuf:"varint,1,opt,name=From,proto3" json:"From,omitempty"`
	// The port that this range ends with (inclusive).
	To uint32 `protobuf:"varint,2,opt,name=To,proto3" json:"To,omitempty"`
}

func (x *PortRange) Reset() {
	*x = PortRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_common_net_port_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortRange) ProtoMessage() {}

func (x *PortRange) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_common_net_port_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortRange.ProtoReflect.Descriptor instead.
func (*PortRange) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_common_net_port_proto_rawDescGZIP(), []int{0}
}

func (x *PortRange) GetFrom() uint32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *PortRange) GetTo() uint32 {
	if x != nil {
		return x.To
	}
	return 0
}

// PortList is a list of ports.
type PortList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Range []*PortRange `protobuf:"bytes,1,rep,name=range,proto3" json:"range,omitempty"`
}

func (x *PortList) Reset() {
	*x = PortList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_common_net_port_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortList) ProtoMessage() {}

func (x *PortList) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_common_net_port_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortList.ProtoReflect.Descriptor instead.
func (*PortList) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_common_net_port_proto_rawDescGZIP(), []int{1}
}

func (x *PortList) GetRange() []*PortRange {
	if x != nil {
		return x.Range
	}
	return nil
}

var File_v2ray_com_core_common_net_port_proto protoreflect.FileDescriptor

var file_v2ray_com_core_common_net_port_proto_rawDesc = []byte{
	0x0a, 0x24, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x22, 0x2f, 0x0a,
	0x09, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x46, 0x72,
	0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x0e,
	0x0a, 0x02, 0x54, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x54, 0x6f, 0x22, 0x42,
	0x0a, 0x08, 0x50, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x05, 0x72, 0x61,
	0x6e, 0x67, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65,
	0x74, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x05, 0x72, 0x61, 0x6e,
	0x67, 0x65, 0x42, 0x3a, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x50,
	0x01, 0x5a, 0x03, 0x6e, 0x65, 0x74, 0xaa, 0x02, 0x15, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43,
	0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v2ray_com_core_common_net_port_proto_rawDescOnce sync.Once
	file_v2ray_com_core_common_net_port_proto_rawDescData = file_v2ray_com_core_common_net_port_proto_rawDesc
)

func file_v2ray_com_core_common_net_port_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_common_net_port_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_common_net_port_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_common_net_port_proto_rawDescData)
	})
	return file_v2ray_com_core_common_net_port_proto_rawDescData
}

var file_v2ray_com_core_common_net_port_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_v2ray_com_core_common_net_port_proto_goTypes = []interface{}{
	(*PortRange)(nil), // 0: v2ray.core.common.net.PortRange
	(*PortList)(nil),  // 1: v2ray.core.common.net.PortList
}
var file_v2ray_com_core_common_net_port_proto_depIdxs = []int32{
	0, // 0: v2ray.core.common.net.PortList.range:type_name -> v2ray.core.common.net.PortRange
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_common_net_port_proto_init() }
func file_v2ray_com_core_common_net_port_proto_init() {
	if File_v2ray_com_core_common_net_port_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_common_net_port_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortRange); i {
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
		file_v2ray_com_core_common_net_port_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortList); i {
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
			RawDescriptor: file_v2ray_com_core_common_net_port_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v2ray_com_core_common_net_port_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_common_net_port_proto_depIdxs,
		MessageInfos:      file_v2ray_com_core_common_net_port_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_common_net_port_proto = out.File
	file_v2ray_com_core_common_net_port_proto_rawDesc = nil
	file_v2ray_com_core_common_net_port_proto_goTypes = nil
	file_v2ray_com_core_common_net_port_proto_depIdxs = nil
}
