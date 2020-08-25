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
	core "v2ray.com/core"
	protocol "v2ray.com/core/common/protocol"
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

type AddUserOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *protocol.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *AddUserOperation) Reset() {
	*x = AddUserOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddUserOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddUserOperation) ProtoMessage() {}

func (x *AddUserOperation) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddUserOperation.ProtoReflect.Descriptor instead.
func (*AddUserOperation) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{0}
}

func (x *AddUserOperation) GetUser() *protocol.User {
	if x != nil {
		return x.User
	}
	return nil
}

type RemoveUserOperation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *RemoveUserOperation) Reset() {
	*x = RemoveUserOperation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveUserOperation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveUserOperation) ProtoMessage() {}

func (x *RemoveUserOperation) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveUserOperation.ProtoReflect.Descriptor instead.
func (*RemoveUserOperation) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{1}
}

func (x *RemoveUserOperation) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type AddInboundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inbound *core.InboundHandlerConfig `protobuf:"bytes,1,opt,name=inbound,proto3" json:"inbound,omitempty"`
}

func (x *AddInboundRequest) Reset() {
	*x = AddInboundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddInboundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddInboundRequest) ProtoMessage() {}

func (x *AddInboundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddInboundRequest.ProtoReflect.Descriptor instead.
func (*AddInboundRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{2}
}

func (x *AddInboundRequest) GetInbound() *core.InboundHandlerConfig {
	if x != nil {
		return x.Inbound
	}
	return nil
}

type AddInboundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddInboundResponse) Reset() {
	*x = AddInboundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddInboundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddInboundResponse) ProtoMessage() {}

func (x *AddInboundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddInboundResponse.ProtoReflect.Descriptor instead.
func (*AddInboundResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{3}
}

type RemoveInboundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (x *RemoveInboundRequest) Reset() {
	*x = RemoveInboundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveInboundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveInboundRequest) ProtoMessage() {}

func (x *RemoveInboundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveInboundRequest.ProtoReflect.Descriptor instead.
func (*RemoveInboundRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{4}
}

func (x *RemoveInboundRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

type RemoveInboundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveInboundResponse) Reset() {
	*x = RemoveInboundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveInboundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveInboundResponse) ProtoMessage() {}

func (x *RemoveInboundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveInboundResponse.ProtoReflect.Descriptor instead.
func (*RemoveInboundResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{5}
}

type AlterInboundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag       string               `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Operation *serial.TypedMessage `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
}

func (x *AlterInboundRequest) Reset() {
	*x = AlterInboundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlterInboundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlterInboundRequest) ProtoMessage() {}

func (x *AlterInboundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlterInboundRequest.ProtoReflect.Descriptor instead.
func (*AlterInboundRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{6}
}

func (x *AlterInboundRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *AlterInboundRequest) GetOperation() *serial.TypedMessage {
	if x != nil {
		return x.Operation
	}
	return nil
}

type AlterInboundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AlterInboundResponse) Reset() {
	*x = AlterInboundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlterInboundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlterInboundResponse) ProtoMessage() {}

func (x *AlterInboundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlterInboundResponse.ProtoReflect.Descriptor instead.
func (*AlterInboundResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{7}
}

type AddOutboundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Outbound *core.OutboundHandlerConfig `protobuf:"bytes,1,opt,name=outbound,proto3" json:"outbound,omitempty"`
}

func (x *AddOutboundRequest) Reset() {
	*x = AddOutboundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddOutboundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddOutboundRequest) ProtoMessage() {}

func (x *AddOutboundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddOutboundRequest.ProtoReflect.Descriptor instead.
func (*AddOutboundRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{8}
}

func (x *AddOutboundRequest) GetOutbound() *core.OutboundHandlerConfig {
	if x != nil {
		return x.Outbound
	}
	return nil
}

type AddOutboundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddOutboundResponse) Reset() {
	*x = AddOutboundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddOutboundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddOutboundResponse) ProtoMessage() {}

func (x *AddOutboundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddOutboundResponse.ProtoReflect.Descriptor instead.
func (*AddOutboundResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{9}
}

type RemoveOutboundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag string `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (x *RemoveOutboundRequest) Reset() {
	*x = RemoveOutboundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveOutboundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveOutboundRequest) ProtoMessage() {}

func (x *RemoveOutboundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveOutboundRequest.ProtoReflect.Descriptor instead.
func (*RemoveOutboundRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{10}
}

func (x *RemoveOutboundRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

type RemoveOutboundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveOutboundResponse) Reset() {
	*x = RemoveOutboundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveOutboundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveOutboundResponse) ProtoMessage() {}

func (x *RemoveOutboundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveOutboundResponse.ProtoReflect.Descriptor instead.
func (*RemoveOutboundResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{11}
}

type AlterOutboundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag       string               `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Operation *serial.TypedMessage `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
}

func (x *AlterOutboundRequest) Reset() {
	*x = AlterOutboundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlterOutboundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlterOutboundRequest) ProtoMessage() {}

func (x *AlterOutboundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlterOutboundRequest.ProtoReflect.Descriptor instead.
func (*AlterOutboundRequest) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{12}
}

func (x *AlterOutboundRequest) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *AlterOutboundRequest) GetOperation() *serial.TypedMessage {
	if x != nil {
		return x.Operation
	}
	return nil
}

type AlterOutboundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AlterOutboundResponse) Reset() {
	*x = AlterOutboundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlterOutboundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlterOutboundResponse) ProtoMessage() {}

func (x *AlterOutboundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlterOutboundResponse.ProtoReflect.Descriptor instead.
func (*AlterOutboundResponse) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{13}
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[14]
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
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP(), []int{14}
}

var File_v2ray_com_core_app_proxyman_command_command_proto protoreflect.FileDescriptor

var file_v2ray_com_core_app_proxyman_command_command_proto_rawDesc = []byte{
	0x0a, 0x31, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x61, 0x70, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x29, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x30, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x64, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1b, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48,
	0x0a, 0x10, 0x41, 0x64, 0x64, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x2b, 0x0a, 0x13, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x4f, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x62, 0x6f,
	0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x07, 0x69, 0x6e,
	0x62, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x76, 0x32,
	0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x07, 0x69,
	0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x62,
	0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x28, 0x0a, 0x14,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x22, 0x17, 0x0a, 0x15, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x6d, 0x0a, 0x13, 0x41, 0x6c, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x44, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x76, 0x32,
	0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x16,
	0x0a, 0x14, 0x41, 0x6c, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x53, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x4f, 0x75, 0x74,
	0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x08,
	0x6f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4f, 0x75, 0x74, 0x62,
	0x6f, 0x75, 0x6e, 0x64, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x08, 0x6f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x22, 0x15, 0x0a, 0x13, 0x41,
	0x64, 0x64, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x29, 0x0a, 0x15, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4f, 0x75, 0x74, 0x62,
	0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x74,
	0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x22, 0x18, 0x0a,
	0x16, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x6e, 0x0a, 0x14, 0x41, 0x6c, 0x74, 0x65, 0x72,
	0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61,
	0x67, 0x12, 0x44, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x2e,
	0x54, 0x79, 0x70, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x09, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x17, 0x0a, 0x15, 0x41, 0x6c, 0x74, 0x65, 0x72,
	0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x08, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x32, 0x90, 0x06, 0x0a, 0x0e, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x77, 0x0a,
	0x0a, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x32, 0x2e, 0x76, 0x32,
	0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x41, 0x64,
	0x64, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x33, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x80, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x35, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d,
	0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x36, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x7d, 0x0a, 0x0c, 0x41, 0x6c, 0x74,
	0x65, 0x72, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x34, 0x2e, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79,
	0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x41, 0x6c, 0x74, 0x65,
	0x72, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x35, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x41, 0x6c, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x7a, 0x0a, 0x0b, 0x41, 0x64, 0x64, 0x4f,
	0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x33, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61,
	0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x41, 0x64, 0x64, 0x4f, 0x75, 0x74,
	0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x76,
	0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x41,
	0x64, 0x64, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x83, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4f,
	0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x36, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61,
	0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x37, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x80, 0x01, 0x0a, 0x0d, 0x41,
	0x6c, 0x74, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x35, 0x2e, 0x76,
	0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x41,
	0x6c, 0x74, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x41, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x62, 0x6f,
	0x75, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x52, 0x0a,
	0x23, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x61, 0x70, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x50, 0x01, 0x5a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0xaa,
	0x02, 0x1f, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x70, 0x70,
	0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x6d, 0x61, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v2ray_com_core_app_proxyman_command_command_proto_rawDescOnce sync.Once
	file_v2ray_com_core_app_proxyman_command_command_proto_rawDescData = file_v2ray_com_core_app_proxyman_command_command_proto_rawDesc
)

func file_v2ray_com_core_app_proxyman_command_command_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_app_proxyman_command_command_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_app_proxyman_command_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_app_proxyman_command_command_proto_rawDescData)
	})
	return file_v2ray_com_core_app_proxyman_command_command_proto_rawDescData
}

var file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_v2ray_com_core_app_proxyman_command_command_proto_goTypes = []interface{}{
	(*AddUserOperation)(nil),           // 0: v2ray.core.app.proxyman.command.AddUserOperation
	(*RemoveUserOperation)(nil),        // 1: v2ray.core.app.proxyman.command.RemoveUserOperation
	(*AddInboundRequest)(nil),          // 2: v2ray.core.app.proxyman.command.AddInboundRequest
	(*AddInboundResponse)(nil),         // 3: v2ray.core.app.proxyman.command.AddInboundResponse
	(*RemoveInboundRequest)(nil),       // 4: v2ray.core.app.proxyman.command.RemoveInboundRequest
	(*RemoveInboundResponse)(nil),      // 5: v2ray.core.app.proxyman.command.RemoveInboundResponse
	(*AlterInboundRequest)(nil),        // 6: v2ray.core.app.proxyman.command.AlterInboundRequest
	(*AlterInboundResponse)(nil),       // 7: v2ray.core.app.proxyman.command.AlterInboundResponse
	(*AddOutboundRequest)(nil),         // 8: v2ray.core.app.proxyman.command.AddOutboundRequest
	(*AddOutboundResponse)(nil),        // 9: v2ray.core.app.proxyman.command.AddOutboundResponse
	(*RemoveOutboundRequest)(nil),      // 10: v2ray.core.app.proxyman.command.RemoveOutboundRequest
	(*RemoveOutboundResponse)(nil),     // 11: v2ray.core.app.proxyman.command.RemoveOutboundResponse
	(*AlterOutboundRequest)(nil),       // 12: v2ray.core.app.proxyman.command.AlterOutboundRequest
	(*AlterOutboundResponse)(nil),      // 13: v2ray.core.app.proxyman.command.AlterOutboundResponse
	(*Config)(nil),                     // 14: v2ray.core.app.proxyman.command.Config
	(*protocol.User)(nil),              // 15: v2ray.core.common.protocol.User
	(*core.InboundHandlerConfig)(nil),  // 16: v2ray.core.InboundHandlerConfig
	(*serial.TypedMessage)(nil),        // 17: v2ray.core.common.serial.TypedMessage
	(*core.OutboundHandlerConfig)(nil), // 18: v2ray.core.OutboundHandlerConfig
}
var file_v2ray_com_core_app_proxyman_command_command_proto_depIdxs = []int32{
	15, // 0: v2ray.core.app.proxyman.command.AddUserOperation.user:type_name -> v2ray.core.common.protocol.User
	16, // 1: v2ray.core.app.proxyman.command.AddInboundRequest.inbound:type_name -> v2ray.core.InboundHandlerConfig
	17, // 2: v2ray.core.app.proxyman.command.AlterInboundRequest.operation:type_name -> v2ray.core.common.serial.TypedMessage
	18, // 3: v2ray.core.app.proxyman.command.AddOutboundRequest.outbound:type_name -> v2ray.core.OutboundHandlerConfig
	17, // 4: v2ray.core.app.proxyman.command.AlterOutboundRequest.operation:type_name -> v2ray.core.common.serial.TypedMessage
	2,  // 5: v2ray.core.app.proxyman.command.HandlerService.AddInbound:input_type -> v2ray.core.app.proxyman.command.AddInboundRequest
	4,  // 6: v2ray.core.app.proxyman.command.HandlerService.RemoveInbound:input_type -> v2ray.core.app.proxyman.command.RemoveInboundRequest
	6,  // 7: v2ray.core.app.proxyman.command.HandlerService.AlterInbound:input_type -> v2ray.core.app.proxyman.command.AlterInboundRequest
	8,  // 8: v2ray.core.app.proxyman.command.HandlerService.AddOutbound:input_type -> v2ray.core.app.proxyman.command.AddOutboundRequest
	10, // 9: v2ray.core.app.proxyman.command.HandlerService.RemoveOutbound:input_type -> v2ray.core.app.proxyman.command.RemoveOutboundRequest
	12, // 10: v2ray.core.app.proxyman.command.HandlerService.AlterOutbound:input_type -> v2ray.core.app.proxyman.command.AlterOutboundRequest
	3,  // 11: v2ray.core.app.proxyman.command.HandlerService.AddInbound:output_type -> v2ray.core.app.proxyman.command.AddInboundResponse
	5,  // 12: v2ray.core.app.proxyman.command.HandlerService.RemoveInbound:output_type -> v2ray.core.app.proxyman.command.RemoveInboundResponse
	7,  // 13: v2ray.core.app.proxyman.command.HandlerService.AlterInbound:output_type -> v2ray.core.app.proxyman.command.AlterInboundResponse
	9,  // 14: v2ray.core.app.proxyman.command.HandlerService.AddOutbound:output_type -> v2ray.core.app.proxyman.command.AddOutboundResponse
	11, // 15: v2ray.core.app.proxyman.command.HandlerService.RemoveOutbound:output_type -> v2ray.core.app.proxyman.command.RemoveOutboundResponse
	13, // 16: v2ray.core.app.proxyman.command.HandlerService.AlterOutbound:output_type -> v2ray.core.app.proxyman.command.AlterOutboundResponse
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_app_proxyman_command_command_proto_init() }
func file_v2ray_com_core_app_proxyman_command_command_proto_init() {
	if File_v2ray_com_core_app_proxyman_command_command_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddUserOperation); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveUserOperation); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddInboundRequest); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddInboundResponse); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveInboundRequest); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveInboundResponse); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlterInboundRequest); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlterInboundResponse); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddOutboundRequest); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddOutboundResponse); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveOutboundRequest); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveOutboundResponse); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlterOutboundRequest); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlterOutboundResponse); i {
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
		file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_v2ray_com_core_app_proxyman_command_command_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v2ray_com_core_app_proxyman_command_command_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_app_proxyman_command_command_proto_depIdxs,
		MessageInfos:      file_v2ray_com_core_app_proxyman_command_command_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_app_proxyman_command_command_proto = out.File
	file_v2ray_com_core_app_proxyman_command_command_proto_rawDesc = nil
	file_v2ray_com_core_app_proxyman_command_command_proto_goTypes = nil
	file_v2ray_com_core_app_proxyman_command_command_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// HandlerServiceClient is the client API for HandlerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HandlerServiceClient interface {
	AddInbound(ctx context.Context, in *AddInboundRequest, opts ...grpc.CallOption) (*AddInboundResponse, error)
	RemoveInbound(ctx context.Context, in *RemoveInboundRequest, opts ...grpc.CallOption) (*RemoveInboundResponse, error)
	AlterInbound(ctx context.Context, in *AlterInboundRequest, opts ...grpc.CallOption) (*AlterInboundResponse, error)
	AddOutbound(ctx context.Context, in *AddOutboundRequest, opts ...grpc.CallOption) (*AddOutboundResponse, error)
	RemoveOutbound(ctx context.Context, in *RemoveOutboundRequest, opts ...grpc.CallOption) (*RemoveOutboundResponse, error)
	AlterOutbound(ctx context.Context, in *AlterOutboundRequest, opts ...grpc.CallOption) (*AlterOutboundResponse, error)
}

type handlerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHandlerServiceClient(cc grpc.ClientConnInterface) HandlerServiceClient {
	return &handlerServiceClient{cc}
}

func (c *handlerServiceClient) AddInbound(ctx context.Context, in *AddInboundRequest, opts ...grpc.CallOption) (*AddInboundResponse, error) {
	out := new(AddInboundResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/AddInbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *handlerServiceClient) RemoveInbound(ctx context.Context, in *RemoveInboundRequest, opts ...grpc.CallOption) (*RemoveInboundResponse, error) {
	out := new(RemoveInboundResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/RemoveInbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *handlerServiceClient) AlterInbound(ctx context.Context, in *AlterInboundRequest, opts ...grpc.CallOption) (*AlterInboundResponse, error) {
	out := new(AlterInboundResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/AlterInbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *handlerServiceClient) AddOutbound(ctx context.Context, in *AddOutboundRequest, opts ...grpc.CallOption) (*AddOutboundResponse, error) {
	out := new(AddOutboundResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/AddOutbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *handlerServiceClient) RemoveOutbound(ctx context.Context, in *RemoveOutboundRequest, opts ...grpc.CallOption) (*RemoveOutboundResponse, error) {
	out := new(RemoveOutboundResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/RemoveOutbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *handlerServiceClient) AlterOutbound(ctx context.Context, in *AlterOutboundRequest, opts ...grpc.CallOption) (*AlterOutboundResponse, error) {
	out := new(AlterOutboundResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/AlterOutbound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HandlerServiceServer is the server API for HandlerService service.
type HandlerServiceServer interface {
	AddInbound(context.Context, *AddInboundRequest) (*AddInboundResponse, error)
	RemoveInbound(context.Context, *RemoveInboundRequest) (*RemoveInboundResponse, error)
	AlterInbound(context.Context, *AlterInboundRequest) (*AlterInboundResponse, error)
	AddOutbound(context.Context, *AddOutboundRequest) (*AddOutboundResponse, error)
	RemoveOutbound(context.Context, *RemoveOutboundRequest) (*RemoveOutboundResponse, error)
	AlterOutbound(context.Context, *AlterOutboundRequest) (*AlterOutboundResponse, error)
}

// UnimplementedHandlerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedHandlerServiceServer struct {
}

func (*UnimplementedHandlerServiceServer) AddInbound(context.Context, *AddInboundRequest) (*AddInboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddInbound not implemented")
}
func (*UnimplementedHandlerServiceServer) RemoveInbound(context.Context, *RemoveInboundRequest) (*RemoveInboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveInbound not implemented")
}
func (*UnimplementedHandlerServiceServer) AlterInbound(context.Context, *AlterInboundRequest) (*AlterInboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AlterInbound not implemented")
}
func (*UnimplementedHandlerServiceServer) AddOutbound(context.Context, *AddOutboundRequest) (*AddOutboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOutbound not implemented")
}
func (*UnimplementedHandlerServiceServer) RemoveOutbound(context.Context, *RemoveOutboundRequest) (*RemoveOutboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveOutbound not implemented")
}
func (*UnimplementedHandlerServiceServer) AlterOutbound(context.Context, *AlterOutboundRequest) (*AlterOutboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AlterOutbound not implemented")
}

func RegisterHandlerServiceServer(s *grpc.Server, srv HandlerServiceServer) {
	s.RegisterService(&_HandlerService_serviceDesc, srv)
}

func _HandlerService_AddInbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddInboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).AddInbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/AddInbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).AddInbound(ctx, req.(*AddInboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HandlerService_RemoveInbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveInboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).RemoveInbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/RemoveInbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).RemoveInbound(ctx, req.(*RemoveInboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HandlerService_AlterInbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlterInboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).AlterInbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/AlterInbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).AlterInbound(ctx, req.(*AlterInboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HandlerService_AddOutbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddOutboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).AddOutbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/AddOutbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).AddOutbound(ctx, req.(*AddOutboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HandlerService_RemoveOutbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveOutboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).RemoveOutbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/RemoveOutbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).RemoveOutbound(ctx, req.(*RemoveOutboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HandlerService_AlterOutbound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlterOutboundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).AlterOutbound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/AlterOutbound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).AlterOutbound(ctx, req.(*AlterOutboundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HandlerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v2ray.core.app.proxyman.command.HandlerService",
	HandlerType: (*HandlerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddInbound",
			Handler:    _HandlerService_AddInbound_Handler,
		},
		{
			MethodName: "RemoveInbound",
			Handler:    _HandlerService_RemoveInbound_Handler,
		},
		{
			MethodName: "AlterInbound",
			Handler:    _HandlerService_AlterInbound_Handler,
		},
		{
			MethodName: "AddOutbound",
			Handler:    _HandlerService_AddOutbound_Handler,
		},
		{
			MethodName: "RemoveOutbound",
			Handler:    _HandlerService_RemoveOutbound_Handler,
		},
		{
			MethodName: "AlterOutbound",
			Handler:    _HandlerService_AlterOutbound_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v2ray.com/core/app/proxyman/command/command.proto",
}
