package command

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
	core "v2ray.com/core"
	protocol "v2ray.com/core/common/protocol"
	serial "v2ray.com/core/common/serial"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AddUserOperation struct {
	User                 *protocol.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *AddUserOperation) Reset()         { *m = AddUserOperation{} }
func (m *AddUserOperation) String() string { return proto.CompactTextString(m) }
func (*AddUserOperation) ProtoMessage()    {}
func (*AddUserOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{0}
}

func (m *AddUserOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddUserOperation.Unmarshal(m, b)
}
func (m *AddUserOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddUserOperation.Marshal(b, m, deterministic)
}
func (m *AddUserOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddUserOperation.Merge(m, src)
}
func (m *AddUserOperation) XXX_Size() int {
	return xxx_messageInfo_AddUserOperation.Size(m)
}
func (m *AddUserOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_AddUserOperation.DiscardUnknown(m)
}

var xxx_messageInfo_AddUserOperation proto.InternalMessageInfo

func (m *AddUserOperation) GetUser() *protocol.User {
	if m != nil {
		return m.User
	}
	return nil
}

type RemoveUserOperation struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveUserOperation) Reset()         { *m = RemoveUserOperation{} }
func (m *RemoveUserOperation) String() string { return proto.CompactTextString(m) }
func (*RemoveUserOperation) ProtoMessage()    {}
func (*RemoveUserOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{1}
}

func (m *RemoveUserOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveUserOperation.Unmarshal(m, b)
}
func (m *RemoveUserOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveUserOperation.Marshal(b, m, deterministic)
}
func (m *RemoveUserOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveUserOperation.Merge(m, src)
}
func (m *RemoveUserOperation) XXX_Size() int {
	return xxx_messageInfo_RemoveUserOperation.Size(m)
}
func (m *RemoveUserOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveUserOperation.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveUserOperation proto.InternalMessageInfo

func (m *RemoveUserOperation) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type AddInboundRequest struct {
	Inbound              *core.InboundHandlerConfig `protobuf:"bytes,1,opt,name=inbound,proto3" json:"inbound,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *AddInboundRequest) Reset()         { *m = AddInboundRequest{} }
func (m *AddInboundRequest) String() string { return proto.CompactTextString(m) }
func (*AddInboundRequest) ProtoMessage()    {}
func (*AddInboundRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{2}
}

func (m *AddInboundRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddInboundRequest.Unmarshal(m, b)
}
func (m *AddInboundRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddInboundRequest.Marshal(b, m, deterministic)
}
func (m *AddInboundRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddInboundRequest.Merge(m, src)
}
func (m *AddInboundRequest) XXX_Size() int {
	return xxx_messageInfo_AddInboundRequest.Size(m)
}
func (m *AddInboundRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddInboundRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddInboundRequest proto.InternalMessageInfo

func (m *AddInboundRequest) GetInbound() *core.InboundHandlerConfig {
	if m != nil {
		return m.Inbound
	}
	return nil
}

type AddInboundResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddInboundResponse) Reset()         { *m = AddInboundResponse{} }
func (m *AddInboundResponse) String() string { return proto.CompactTextString(m) }
func (*AddInboundResponse) ProtoMessage()    {}
func (*AddInboundResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{3}
}

func (m *AddInboundResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddInboundResponse.Unmarshal(m, b)
}
func (m *AddInboundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddInboundResponse.Marshal(b, m, deterministic)
}
func (m *AddInboundResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddInboundResponse.Merge(m, src)
}
func (m *AddInboundResponse) XXX_Size() int {
	return xxx_messageInfo_AddInboundResponse.Size(m)
}
func (m *AddInboundResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddInboundResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddInboundResponse proto.InternalMessageInfo

type RemoveInboundRequest struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveInboundRequest) Reset()         { *m = RemoveInboundRequest{} }
func (m *RemoveInboundRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveInboundRequest) ProtoMessage()    {}
func (*RemoveInboundRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{4}
}

func (m *RemoveInboundRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveInboundRequest.Unmarshal(m, b)
}
func (m *RemoveInboundRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveInboundRequest.Marshal(b, m, deterministic)
}
func (m *RemoveInboundRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveInboundRequest.Merge(m, src)
}
func (m *RemoveInboundRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveInboundRequest.Size(m)
}
func (m *RemoveInboundRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveInboundRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveInboundRequest proto.InternalMessageInfo

func (m *RemoveInboundRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type RemoveInboundResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveInboundResponse) Reset()         { *m = RemoveInboundResponse{} }
func (m *RemoveInboundResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveInboundResponse) ProtoMessage()    {}
func (*RemoveInboundResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{5}
}

func (m *RemoveInboundResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveInboundResponse.Unmarshal(m, b)
}
func (m *RemoveInboundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveInboundResponse.Marshal(b, m, deterministic)
}
func (m *RemoveInboundResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveInboundResponse.Merge(m, src)
}
func (m *RemoveInboundResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveInboundResponse.Size(m)
}
func (m *RemoveInboundResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveInboundResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveInboundResponse proto.InternalMessageInfo

type AlterInboundRequest struct {
	Tag                  string               `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Operation            *serial.TypedMessage `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *AlterInboundRequest) Reset()         { *m = AlterInboundRequest{} }
func (m *AlterInboundRequest) String() string { return proto.CompactTextString(m) }
func (*AlterInboundRequest) ProtoMessage()    {}
func (*AlterInboundRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{6}
}

func (m *AlterInboundRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlterInboundRequest.Unmarshal(m, b)
}
func (m *AlterInboundRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlterInboundRequest.Marshal(b, m, deterministic)
}
func (m *AlterInboundRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlterInboundRequest.Merge(m, src)
}
func (m *AlterInboundRequest) XXX_Size() int {
	return xxx_messageInfo_AlterInboundRequest.Size(m)
}
func (m *AlterInboundRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AlterInboundRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AlterInboundRequest proto.InternalMessageInfo

func (m *AlterInboundRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *AlterInboundRequest) GetOperation() *serial.TypedMessage {
	if m != nil {
		return m.Operation
	}
	return nil
}

type AlterInboundResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AlterInboundResponse) Reset()         { *m = AlterInboundResponse{} }
func (m *AlterInboundResponse) String() string { return proto.CompactTextString(m) }
func (*AlterInboundResponse) ProtoMessage()    {}
func (*AlterInboundResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{7}
}

func (m *AlterInboundResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlterInboundResponse.Unmarshal(m, b)
}
func (m *AlterInboundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlterInboundResponse.Marshal(b, m, deterministic)
}
func (m *AlterInboundResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlterInboundResponse.Merge(m, src)
}
func (m *AlterInboundResponse) XXX_Size() int {
	return xxx_messageInfo_AlterInboundResponse.Size(m)
}
func (m *AlterInboundResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AlterInboundResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AlterInboundResponse proto.InternalMessageInfo

type ListInboundUserRequest struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListInboundUserRequest) Reset()         { *m = ListInboundUserRequest{} }
func (m *ListInboundUserRequest) String() string { return proto.CompactTextString(m) }
func (*ListInboundUserRequest) ProtoMessage()    {}
func (*ListInboundUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{8}
}

func (m *ListInboundUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListInboundUserRequest.Unmarshal(m, b)
}
func (m *ListInboundUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListInboundUserRequest.Marshal(b, m, deterministic)
}
func (m *ListInboundUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListInboundUserRequest.Merge(m, src)
}
func (m *ListInboundUserRequest) XXX_Size() int {
	return xxx_messageInfo_ListInboundUserRequest.Size(m)
}
func (m *ListInboundUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListInboundUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListInboundUserRequest proto.InternalMessageInfo

func (m *ListInboundUserRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type ListInboundUserResponse struct {
	User                 []*protocol.User `protobuf:"bytes,1,rep,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ListInboundUserResponse) Reset()         { *m = ListInboundUserResponse{} }
func (m *ListInboundUserResponse) String() string { return proto.CompactTextString(m) }
func (*ListInboundUserResponse) ProtoMessage()    {}
func (*ListInboundUserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{9}
}

func (m *ListInboundUserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListInboundUserResponse.Unmarshal(m, b)
}
func (m *ListInboundUserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListInboundUserResponse.Marshal(b, m, deterministic)
}
func (m *ListInboundUserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListInboundUserResponse.Merge(m, src)
}
func (m *ListInboundUserResponse) XXX_Size() int {
	return xxx_messageInfo_ListInboundUserResponse.Size(m)
}
func (m *ListInboundUserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListInboundUserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListInboundUserResponse proto.InternalMessageInfo

func (m *ListInboundUserResponse) GetUser() []*protocol.User {
	if m != nil {
		return m.User
	}
	return nil
}

type AddOutboundRequest struct {
	Outbound             *core.OutboundHandlerConfig `protobuf:"bytes,1,opt,name=outbound,proto3" json:"outbound,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *AddOutboundRequest) Reset()         { *m = AddOutboundRequest{} }
func (m *AddOutboundRequest) String() string { return proto.CompactTextString(m) }
func (*AddOutboundRequest) ProtoMessage()    {}
func (*AddOutboundRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{10}
}

func (m *AddOutboundRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddOutboundRequest.Unmarshal(m, b)
}
func (m *AddOutboundRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddOutboundRequest.Marshal(b, m, deterministic)
}
func (m *AddOutboundRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddOutboundRequest.Merge(m, src)
}
func (m *AddOutboundRequest) XXX_Size() int {
	return xxx_messageInfo_AddOutboundRequest.Size(m)
}
func (m *AddOutboundRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddOutboundRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddOutboundRequest proto.InternalMessageInfo

func (m *AddOutboundRequest) GetOutbound() *core.OutboundHandlerConfig {
	if m != nil {
		return m.Outbound
	}
	return nil
}

type AddOutboundResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddOutboundResponse) Reset()         { *m = AddOutboundResponse{} }
func (m *AddOutboundResponse) String() string { return proto.CompactTextString(m) }
func (*AddOutboundResponse) ProtoMessage()    {}
func (*AddOutboundResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{11}
}

func (m *AddOutboundResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddOutboundResponse.Unmarshal(m, b)
}
func (m *AddOutboundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddOutboundResponse.Marshal(b, m, deterministic)
}
func (m *AddOutboundResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddOutboundResponse.Merge(m, src)
}
func (m *AddOutboundResponse) XXX_Size() int {
	return xxx_messageInfo_AddOutboundResponse.Size(m)
}
func (m *AddOutboundResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddOutboundResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddOutboundResponse proto.InternalMessageInfo

type RemoveOutboundRequest struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveOutboundRequest) Reset()         { *m = RemoveOutboundRequest{} }
func (m *RemoveOutboundRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveOutboundRequest) ProtoMessage()    {}
func (*RemoveOutboundRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{12}
}

func (m *RemoveOutboundRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveOutboundRequest.Unmarshal(m, b)
}
func (m *RemoveOutboundRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveOutboundRequest.Marshal(b, m, deterministic)
}
func (m *RemoveOutboundRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveOutboundRequest.Merge(m, src)
}
func (m *RemoveOutboundRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveOutboundRequest.Size(m)
}
func (m *RemoveOutboundRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveOutboundRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveOutboundRequest proto.InternalMessageInfo

func (m *RemoveOutboundRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type RemoveOutboundResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveOutboundResponse) Reset()         { *m = RemoveOutboundResponse{} }
func (m *RemoveOutboundResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveOutboundResponse) ProtoMessage()    {}
func (*RemoveOutboundResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{13}
}

func (m *RemoveOutboundResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveOutboundResponse.Unmarshal(m, b)
}
func (m *RemoveOutboundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveOutboundResponse.Marshal(b, m, deterministic)
}
func (m *RemoveOutboundResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveOutboundResponse.Merge(m, src)
}
func (m *RemoveOutboundResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveOutboundResponse.Size(m)
}
func (m *RemoveOutboundResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveOutboundResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveOutboundResponse proto.InternalMessageInfo

type AlterOutboundRequest struct {
	Tag                  string               `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Operation            *serial.TypedMessage `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *AlterOutboundRequest) Reset()         { *m = AlterOutboundRequest{} }
func (m *AlterOutboundRequest) String() string { return proto.CompactTextString(m) }
func (*AlterOutboundRequest) ProtoMessage()    {}
func (*AlterOutboundRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{14}
}

func (m *AlterOutboundRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlterOutboundRequest.Unmarshal(m, b)
}
func (m *AlterOutboundRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlterOutboundRequest.Marshal(b, m, deterministic)
}
func (m *AlterOutboundRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlterOutboundRequest.Merge(m, src)
}
func (m *AlterOutboundRequest) XXX_Size() int {
	return xxx_messageInfo_AlterOutboundRequest.Size(m)
}
func (m *AlterOutboundRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AlterOutboundRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AlterOutboundRequest proto.InternalMessageInfo

func (m *AlterOutboundRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *AlterOutboundRequest) GetOperation() *serial.TypedMessage {
	if m != nil {
		return m.Operation
	}
	return nil
}

type AlterOutboundResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AlterOutboundResponse) Reset()         { *m = AlterOutboundResponse{} }
func (m *AlterOutboundResponse) String() string { return proto.CompactTextString(m) }
func (*AlterOutboundResponse) ProtoMessage()    {}
func (*AlterOutboundResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{15}
}

func (m *AlterOutboundResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlterOutboundResponse.Unmarshal(m, b)
}
func (m *AlterOutboundResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlterOutboundResponse.Marshal(b, m, deterministic)
}
func (m *AlterOutboundResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlterOutboundResponse.Merge(m, src)
}
func (m *AlterOutboundResponse) XXX_Size() int {
	return xxx_messageInfo_AlterOutboundResponse.Size(m)
}
func (m *AlterOutboundResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AlterOutboundResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AlterOutboundResponse proto.InternalMessageInfo

type Config struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2c30a70a48636a0, []int{16}
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
	proto.RegisterType((*AddUserOperation)(nil), "v2ray.core.app.proxyman.command.AddUserOperation")
	proto.RegisterType((*RemoveUserOperation)(nil), "v2ray.core.app.proxyman.command.RemoveUserOperation")
	proto.RegisterType((*AddInboundRequest)(nil), "v2ray.core.app.proxyman.command.AddInboundRequest")
	proto.RegisterType((*AddInboundResponse)(nil), "v2ray.core.app.proxyman.command.AddInboundResponse")
	proto.RegisterType((*RemoveInboundRequest)(nil), "v2ray.core.app.proxyman.command.RemoveInboundRequest")
	proto.RegisterType((*RemoveInboundResponse)(nil), "v2ray.core.app.proxyman.command.RemoveInboundResponse")
	proto.RegisterType((*AlterInboundRequest)(nil), "v2ray.core.app.proxyman.command.AlterInboundRequest")
	proto.RegisterType((*AlterInboundResponse)(nil), "v2ray.core.app.proxyman.command.AlterInboundResponse")
	proto.RegisterType((*ListInboundUserRequest)(nil), "v2ray.core.app.proxyman.command.ListInboundUserRequest")
	proto.RegisterType((*ListInboundUserResponse)(nil), "v2ray.core.app.proxyman.command.ListInboundUserResponse")
	proto.RegisterType((*AddOutboundRequest)(nil), "v2ray.core.app.proxyman.command.AddOutboundRequest")
	proto.RegisterType((*AddOutboundResponse)(nil), "v2ray.core.app.proxyman.command.AddOutboundResponse")
	proto.RegisterType((*RemoveOutboundRequest)(nil), "v2ray.core.app.proxyman.command.RemoveOutboundRequest")
	proto.RegisterType((*RemoveOutboundResponse)(nil), "v2ray.core.app.proxyman.command.RemoveOutboundResponse")
	proto.RegisterType((*AlterOutboundRequest)(nil), "v2ray.core.app.proxyman.command.AlterOutboundRequest")
	proto.RegisterType((*AlterOutboundResponse)(nil), "v2ray.core.app.proxyman.command.AlterOutboundResponse")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.proxyman.command.Config")
}

func init() {
	proto.RegisterFile("v2ray.com/core/app/proxyman/command/command.proto", fileDescriptor_e2c30a70a48636a0)
}

var fileDescriptor_e2c30a70a48636a0 = []byte{
	// 603 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x25, 0x1b, 0xac, 0xdb, 0x1d, 0x8c, 0xe1, 0x7e, 0x2a, 0x3c, 0x6c, 0x04, 0x09, 0x6d, 0x20,
	0x39, 0xd0, 0x75, 0x1b, 0x42, 0xe2, 0xa1, 0x94, 0x87, 0x21, 0x81, 0x3a, 0x65, 0xc0, 0x03, 0x2f,
	0xc8, 0x4b, 0x4c, 0x15, 0xa9, 0x89, 0x8d, 0x93, 0x16, 0x8a, 0x84, 0x84, 0x84, 0xc4, 0x7f, 0xe0,
	0x2f, 0xf0, 0x2b, 0x51, 0x62, 0xa7, 0x4b, 0xdc, 0xa2, 0x34, 0x12, 0x4f, 0x4d, 0x9d, 0x73, 0xce,
	0x3d, 0xf7, 0x1e, 0xdf, 0xc0, 0x93, 0x69, 0x57, 0x90, 0x19, 0x76, 0x59, 0x60, 0xbb, 0x4c, 0x50,
	0x9b, 0x70, 0x6e, 0x73, 0xc1, 0xbe, 0xce, 0x02, 0x12, 0xda, 0x2e, 0x0b, 0x02, 0x12, 0x7a, 0xd9,
	0x2f, 0xe6, 0x82, 0xc5, 0x0c, 0xed, 0x65, 0x14, 0x41, 0x31, 0xe1, 0x1c, 0x67, 0x70, 0xac, 0x60,
	0xe6, 0xa1, 0xa6, 0x99, 0x9c, 0xb3, 0xd0, 0x4e, 0xd9, 0x2e, 0x1b, 0xdb, 0x93, 0x88, 0x0a, 0xa9,
	0x65, 0x3e, 0x5e, 0x0e, 0x8d, 0xa8, 0xf0, 0xc9, 0xd8, 0x8e, 0x67, 0x9c, 0x7a, 0x1f, 0x03, 0x1a,
	0x45, 0x64, 0x44, 0x15, 0xe3, 0xee, 0x02, 0x23, 0xfc, 0xe4, 0x8f, 0xe4, 0x4b, 0xeb, 0x0c, 0x76,
	0xfb, 0x9e, 0xf7, 0x2e, 0xa2, 0x62, 0xc8, 0xa9, 0x20, 0xb1, 0xcf, 0x42, 0xd4, 0x83, 0xeb, 0x49,
	0xc1, 0x8e, 0xb1, 0x6f, 0x1c, 0x6c, 0x77, 0xf7, 0x71, 0xce, 0xbd, 0xac, 0x86, 0x33, 0x63, 0x38,
	0x21, 0x3a, 0x29, 0xda, 0x7a, 0x04, 0x75, 0x87, 0x06, 0x6c, 0x4a, 0x8b, 0x62, 0x0d, 0xb8, 0x41,
	0x03, 0xe2, 0x8f, 0x53, 0xb5, 0x2d, 0x47, 0xfe, 0xb1, 0x86, 0x70, 0xa7, 0xef, 0x79, 0xaf, 0xc2,
	0x4b, 0x36, 0x09, 0x3d, 0x87, 0x7e, 0x9e, 0xd0, 0x28, 0x46, 0xcf, 0xa0, 0xe6, 0xcb, 0x93, 0x65,
	0xa5, 0x15, 0xf8, 0x8c, 0x84, 0xde, 0x98, 0x8a, 0x41, 0xda, 0x84, 0x93, 0x11, 0xac, 0x06, 0xa0,
	0xbc, 0x60, 0xc4, 0x59, 0x18, 0x51, 0xeb, 0x00, 0x1a, 0xd2, 0x93, 0x56, 0x69, 0x17, 0xd6, 0x63,
	0x32, 0x52, 0x96, 0x92, 0x47, 0xab, 0x0d, 0x4d, 0x0d, 0xa9, 0x24, 0x02, 0xa8, 0xf7, 0xc7, 0x31,
	0x15, 0x65, 0x0a, 0xe8, 0x25, 0x6c, 0xb1, 0xac, 0xeb, 0xce, 0x5a, 0xea, 0xff, 0xc1, 0x92, 0xd1,
	0xc9, 0xa0, 0xf0, 0xdb, 0x24, 0xa8, 0x37, 0x32, 0x27, 0xe7, 0x8a, 0x68, 0xb5, 0xa0, 0x51, 0x2c,
	0xa7, 0x6c, 0x3c, 0x84, 0xd6, 0x6b, 0x3f, 0x8a, 0xd5, 0x71, 0x3a, 0xf6, 0x7f, 0xf6, 0x32, 0x84,
	0xf6, 0x02, 0x56, 0xca, 0xe4, 0xa2, 0x5d, 0xaf, 0x10, 0xed, 0x45, 0x3a, 0xdc, 0xe1, 0x24, 0x2e,
	0x8c, 0xe0, 0x39, 0x6c, 0x32, 0x75, 0xa4, 0xf2, 0xba, 0x97, 0xd7, 0xcb, 0xe0, 0xc5, 0xc0, 0xe6,
	0x14, 0xab, 0x09, 0xf5, 0x82, 0xa8, 0x6a, 0xf4, 0x30, 0x0b, 0x42, 0x2f, 0xb7, 0xd8, 0x67, 0x07,
	0x5a, 0x3a, 0x54, 0x89, 0x84, 0x6a, 0x8a, 0xa5, 0x1a, 0xff, 0x29, 0xb5, 0x36, 0x34, 0xb5, 0x7a,
	0xca, 0xc8, 0x26, 0x6c, 0xc8, 0xc6, 0xbb, 0xbf, 0x6b, 0xb0, 0xa3, 0x46, 0x71, 0x41, 0xc5, 0xd4,
	0x77, 0x29, 0xfa, 0x02, 0x70, 0x75, 0x67, 0x51, 0x17, 0x97, 0x7c, 0x25, 0xf0, 0xc2, 0xc6, 0x98,
	0x47, 0x95, 0x38, 0xca, 0xd3, 0x35, 0xf4, 0xc3, 0x80, 0x5b, 0x85, 0xdb, 0x8e, 0x8e, 0x4b, 0x85,
	0x96, 0xed, 0x91, 0x79, 0x52, 0x95, 0x36, 0xb7, 0xf0, 0x1d, 0x6e, 0xe6, 0xef, 0x39, 0xea, 0x95,
	0x77, 0xb2, 0xb8, 0x85, 0xe6, 0x71, 0x45, 0xd6, 0xbc, 0xfc, 0x2f, 0x03, 0x6e, 0x6b, 0x3b, 0x82,
	0x4e, 0x4b, 0xc5, 0x96, 0x6f, 0xa0, 0xf9, 0xb4, 0x3a, 0x71, 0x6e, 0xe4, 0x1b, 0x6c, 0xe7, 0xb6,
	0x00, 0xad, 0x14, 0xa8, 0x76, 0xab, 0xcd, 0x5e, 0x35, 0xd2, 0xbc, 0xf6, 0x4f, 0x03, 0x76, 0x8a,
	0x0b, 0x84, 0x56, 0x0d, 0x54, 0xb7, 0x70, 0x5a, 0x99, 0x57, 0xb8, 0x8c, 0x85, 0xe5, 0x41, 0x2b,
	0xa6, 0xaa, 0x7b, 0x38, 0xa9, 0x4a, 0xcb, 0x2c, 0xbc, 0x70, 0xe0, 0xbe, 0xcb, 0x82, 0x32, 0xfa,
	0xb9, 0xf1, 0xa1, 0xa6, 0x1e, 0xff, 0xac, 0xed, 0xbd, 0xef, 0x3a, 0x64, 0x86, 0x07, 0x09, 0xb8,
	0xcf, 0x39, 0x3e, 0xcf, 0xc0, 0x03, 0x89, 0xb8, 0xdc, 0x48, 0x3f, 0xa4, 0x47, 0x7f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0xf7, 0x36, 0x57, 0x86, 0x2f, 0x08, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HandlerServiceClient is the client API for HandlerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HandlerServiceClient interface {
	AddInbound(ctx context.Context, in *AddInboundRequest, opts ...grpc.CallOption) (*AddInboundResponse, error)
	RemoveInbound(ctx context.Context, in *RemoveInboundRequest, opts ...grpc.CallOption) (*RemoveInboundResponse, error)
	AlterInbound(ctx context.Context, in *AlterInboundRequest, opts ...grpc.CallOption) (*AlterInboundResponse, error)
	ListInboundUser(ctx context.Context, in *ListInboundUserRequest, opts ...grpc.CallOption) (*ListInboundUserResponse, error)
	AddOutbound(ctx context.Context, in *AddOutboundRequest, opts ...grpc.CallOption) (*AddOutboundResponse, error)
	RemoveOutbound(ctx context.Context, in *RemoveOutboundRequest, opts ...grpc.CallOption) (*RemoveOutboundResponse, error)
	AlterOutbound(ctx context.Context, in *AlterOutboundRequest, opts ...grpc.CallOption) (*AlterOutboundResponse, error)
}

type handlerServiceClient struct {
	cc *grpc.ClientConn
}

func NewHandlerServiceClient(cc *grpc.ClientConn) HandlerServiceClient {
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

func (c *handlerServiceClient) ListInboundUser(ctx context.Context, in *ListInboundUserRequest, opts ...grpc.CallOption) (*ListInboundUserResponse, error) {
	out := new(ListInboundUserResponse)
	err := c.cc.Invoke(ctx, "/v2ray.core.app.proxyman.command.HandlerService/ListInboundUser", in, out, opts...)
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
	ListInboundUser(context.Context, *ListInboundUserRequest) (*ListInboundUserResponse, error)
	AddOutbound(context.Context, *AddOutboundRequest) (*AddOutboundResponse, error)
	RemoveOutbound(context.Context, *RemoveOutboundRequest) (*RemoveOutboundResponse, error)
	AlterOutbound(context.Context, *AlterOutboundRequest) (*AlterOutboundResponse, error)
}

// UnimplementedHandlerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedHandlerServiceServer struct {
}

func (*UnimplementedHandlerServiceServer) AddInbound(ctx context.Context, req *AddInboundRequest) (*AddInboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddInbound not implemented")
}
func (*UnimplementedHandlerServiceServer) RemoveInbound(ctx context.Context, req *RemoveInboundRequest) (*RemoveInboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveInbound not implemented")
}
func (*UnimplementedHandlerServiceServer) AlterInbound(ctx context.Context, req *AlterInboundRequest) (*AlterInboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AlterInbound not implemented")
}
func (*UnimplementedHandlerServiceServer) ListInboundUser(ctx context.Context, req *ListInboundUserRequest) (*ListInboundUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInboundUser not implemented")
}
func (*UnimplementedHandlerServiceServer) AddOutbound(ctx context.Context, req *AddOutboundRequest) (*AddOutboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOutbound not implemented")
}
func (*UnimplementedHandlerServiceServer) RemoveOutbound(ctx context.Context, req *RemoveOutboundRequest) (*RemoveOutboundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveOutbound not implemented")
}
func (*UnimplementedHandlerServiceServer) AlterOutbound(ctx context.Context, req *AlterOutboundRequest) (*AlterOutboundResponse, error) {
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

func _HandlerService_ListInboundUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInboundUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServiceServer).ListInboundUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2ray.core.app.proxyman.command.HandlerService/ListInboundUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServiceServer).ListInboundUser(ctx, req.(*ListInboundUserRequest))
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
			MethodName: "ListInboundUser",
			Handler:    _HandlerService_ListInboundUser_Handler,
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
