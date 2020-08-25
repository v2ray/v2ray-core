package socks

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	net "v2ray.com/core/common/net"
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

// AuthType is the authentication type of Socks proxy.
type AuthType int32

const (
	// NO_AUTH is for anounymous authentication.
	AuthType_NO_AUTH AuthType = 0
	// PASSWORD is for username/password authentication.
	AuthType_PASSWORD AuthType = 1
)

// Enum value maps for AuthType.
var (
	AuthType_name = map[int32]string{
		0: "NO_AUTH",
		1: "PASSWORD",
	}
	AuthType_value = map[string]int32{
		"NO_AUTH":  0,
		"PASSWORD": 1,
	}
)

func (x AuthType) Enum() *AuthType {
	p := new(AuthType)
	*p = x
	return p
}

func (x AuthType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AuthType) Descriptor() protoreflect.EnumDescriptor {
	return file_v2ray_com_core_proxy_socks_config_proto_enumTypes[0].Descriptor()
}

func (AuthType) Type() protoreflect.EnumType {
	return &file_v2ray_com_core_proxy_socks_config_proto_enumTypes[0]
}

func (x AuthType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AuthType.Descriptor instead.
func (AuthType) EnumDescriptor() ([]byte, []int) {
	return file_v2ray_com_core_proxy_socks_config_proto_rawDescGZIP(), []int{0}
}

// Account represents a Socks account.
type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_proxy_socks_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_proxy_socks_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_proxy_socks_config_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Account) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// ServerConfig is the protobuf config for Socks server.
type ServerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthType   AuthType          `protobuf:"varint,1,opt,name=auth_type,json=authType,proto3,enum=v2ray.core.proxy.socks.AuthType" json:"auth_type,omitempty"`
	Accounts   map[string]string `protobuf:"bytes,2,rep,name=accounts,proto3" json:"accounts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Address    *net.IPOrDomain   `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	UdpEnabled bool              `protobuf:"varint,4,opt,name=udp_enabled,json=udpEnabled,proto3" json:"udp_enabled,omitempty"`
	// Deprecated: Do not use.
	Timeout   uint32 `protobuf:"varint,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	UserLevel uint32 `protobuf:"varint,6,opt,name=user_level,json=userLevel,proto3" json:"user_level,omitempty"`
}

func (x *ServerConfig) Reset() {
	*x = ServerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_proxy_socks_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerConfig) ProtoMessage() {}

func (x *ServerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_proxy_socks_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerConfig.ProtoReflect.Descriptor instead.
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_proxy_socks_config_proto_rawDescGZIP(), []int{1}
}

func (x *ServerConfig) GetAuthType() AuthType {
	if x != nil {
		return x.AuthType
	}
	return AuthType_NO_AUTH
}

func (x *ServerConfig) GetAccounts() map[string]string {
	if x != nil {
		return x.Accounts
	}
	return nil
}

func (x *ServerConfig) GetAddress() *net.IPOrDomain {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *ServerConfig) GetUdpEnabled() bool {
	if x != nil {
		return x.UdpEnabled
	}
	return false
}

// Deprecated: Do not use.
func (x *ServerConfig) GetTimeout() uint32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *ServerConfig) GetUserLevel() uint32 {
	if x != nil {
		return x.UserLevel
	}
	return 0
}

// ClientConfig is the protobuf config for Socks client.
type ClientConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Sever is a list of Socks server addresses.
	Server []*protocol.ServerEndpoint `protobuf:"bytes,1,rep,name=server,proto3" json:"server,omitempty"`
}

func (x *ClientConfig) Reset() {
	*x = ClientConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v2ray_com_core_proxy_socks_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientConfig) ProtoMessage() {}

func (x *ClientConfig) ProtoReflect() protoreflect.Message {
	mi := &file_v2ray_com_core_proxy_socks_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientConfig.ProtoReflect.Descriptor instead.
func (*ClientConfig) Descriptor() ([]byte, []int) {
	return file_v2ray_com_core_proxy_socks_config_proto_rawDescGZIP(), []int{2}
}

func (x *ClientConfig) GetServer() []*protocol.ServerEndpoint {
	if x != nil {
		return x.Server
	}
	return nil
}

var File_v2ray_com_core_proxy_socks_config_proto protoreflect.FileDescriptor

var file_v2ray_com_core_proxy_socks_config_proto_rawDesc = []byte{
	0x0a, 0x27, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x73, 0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x76, 0x32, 0x72, 0x61, 0x79,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x73, 0x6f, 0x63, 0x6b,
	0x73, 0x1a, 0x27, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x30, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x07,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0xf5, 0x02, 0x0a, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x3d, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x73, 0x2e, 0x41, 0x75, 0x74,
	0x68, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x4e, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x32, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x78, 0x79, 0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12,
	0x3b, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x49, 0x50, 0x4f, 0x72, 0x44, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x0b,
	0x75, 0x64, 0x70, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0a, 0x75, 0x64, 0x70, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1c, 0x0a,
	0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x02,
	0x18, 0x01, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x1a, 0x3b, 0x0a, 0x0d, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x52, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x42, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2a, 0x25, 0x0a, 0x08, 0x41,
	0x75, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x5f, 0x41, 0x55,
	0x54, 0x48, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x41, 0x53, 0x53, 0x57, 0x4f, 0x52, 0x44,
	0x10, 0x01, 0x42, 0x3e, 0x0a, 0x1a, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x73,
	0x50, 0x01, 0x5a, 0x05, 0x73, 0x6f, 0x63, 0x6b, 0x73, 0xaa, 0x02, 0x16, 0x56, 0x32, 0x52, 0x61,
	0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x53, 0x6f, 0x63,
	0x6b, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v2ray_com_core_proxy_socks_config_proto_rawDescOnce sync.Once
	file_v2ray_com_core_proxy_socks_config_proto_rawDescData = file_v2ray_com_core_proxy_socks_config_proto_rawDesc
)

func file_v2ray_com_core_proxy_socks_config_proto_rawDescGZIP() []byte {
	file_v2ray_com_core_proxy_socks_config_proto_rawDescOnce.Do(func() {
		file_v2ray_com_core_proxy_socks_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_v2ray_com_core_proxy_socks_config_proto_rawDescData)
	})
	return file_v2ray_com_core_proxy_socks_config_proto_rawDescData
}

var file_v2ray_com_core_proxy_socks_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v2ray_com_core_proxy_socks_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_v2ray_com_core_proxy_socks_config_proto_goTypes = []interface{}{
	(AuthType)(0),                   // 0: v2ray.core.proxy.socks.AuthType
	(*Account)(nil),                 // 1: v2ray.core.proxy.socks.Account
	(*ServerConfig)(nil),            // 2: v2ray.core.proxy.socks.ServerConfig
	(*ClientConfig)(nil),            // 3: v2ray.core.proxy.socks.ClientConfig
	nil,                             // 4: v2ray.core.proxy.socks.ServerConfig.AccountsEntry
	(*net.IPOrDomain)(nil),          // 5: v2ray.core.common.net.IPOrDomain
	(*protocol.ServerEndpoint)(nil), // 6: v2ray.core.common.protocol.ServerEndpoint
}
var file_v2ray_com_core_proxy_socks_config_proto_depIdxs = []int32{
	0, // 0: v2ray.core.proxy.socks.ServerConfig.auth_type:type_name -> v2ray.core.proxy.socks.AuthType
	4, // 1: v2ray.core.proxy.socks.ServerConfig.accounts:type_name -> v2ray.core.proxy.socks.ServerConfig.AccountsEntry
	5, // 2: v2ray.core.proxy.socks.ServerConfig.address:type_name -> v2ray.core.common.net.IPOrDomain
	6, // 3: v2ray.core.proxy.socks.ClientConfig.server:type_name -> v2ray.core.common.protocol.ServerEndpoint
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_v2ray_com_core_proxy_socks_config_proto_init() }
func file_v2ray_com_core_proxy_socks_config_proto_init() {
	if File_v2ray_com_core_proxy_socks_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v2ray_com_core_proxy_socks_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
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
		file_v2ray_com_core_proxy_socks_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerConfig); i {
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
		file_v2ray_com_core_proxy_socks_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientConfig); i {
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
			RawDescriptor: file_v2ray_com_core_proxy_socks_config_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v2ray_com_core_proxy_socks_config_proto_goTypes,
		DependencyIndexes: file_v2ray_com_core_proxy_socks_config_proto_depIdxs,
		EnumInfos:         file_v2ray_com_core_proxy_socks_config_proto_enumTypes,
		MessageInfos:      file_v2ray_com_core_proxy_socks_config_proto_msgTypes,
	}.Build()
	File_v2ray_com_core_proxy_socks_config_proto = out.File
	file_v2ray_com_core_proxy_socks_config_proto_rawDesc = nil
	file_v2ray_com_core_proxy_socks_config_proto_goTypes = nil
	file_v2ray_com_core_proxy_socks_config_proto_depIdxs = nil
}
