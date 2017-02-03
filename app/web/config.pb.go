package web

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v2ray_core_common_serial "v2ray.com/core/common/serial"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FileServer struct {
	Entry []*FileServer_Entry `protobuf:"bytes,1,rep,name=entry" json:"entry,omitempty"`
}

func (m *FileServer) Reset()                    { *m = FileServer{} }
func (m *FileServer) String() string            { return proto.CompactTextString(m) }
func (*FileServer) ProtoMessage()               {}
func (*FileServer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *FileServer) GetEntry() []*FileServer_Entry {
	if m != nil {
		return m.Entry
	}
	return nil
}

type FileServer_Entry struct {
	// Types that are valid to be assigned to FileOrDir:
	//	*FileServer_Entry_File
	//	*FileServer_Entry_Directory
	FileOrDir isFileServer_Entry_FileOrDir `protobuf_oneof:"FileOrDir"`
	Path      string                       `protobuf:"bytes,3,opt,name=path" json:"path,omitempty"`
}

func (m *FileServer_Entry) Reset()                    { *m = FileServer_Entry{} }
func (m *FileServer_Entry) String() string            { return proto.CompactTextString(m) }
func (*FileServer_Entry) ProtoMessage()               {}
func (*FileServer_Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type isFileServer_Entry_FileOrDir interface {
	isFileServer_Entry_FileOrDir()
}

type FileServer_Entry_File struct {
	File string `protobuf:"bytes,1,opt,name=File,oneof"`
}
type FileServer_Entry_Directory struct {
	Directory string `protobuf:"bytes,2,opt,name=Directory,oneof"`
}

func (*FileServer_Entry_File) isFileServer_Entry_FileOrDir()      {}
func (*FileServer_Entry_Directory) isFileServer_Entry_FileOrDir() {}

func (m *FileServer_Entry) GetFileOrDir() isFileServer_Entry_FileOrDir {
	if m != nil {
		return m.FileOrDir
	}
	return nil
}

func (m *FileServer_Entry) GetFile() string {
	if x, ok := m.GetFileOrDir().(*FileServer_Entry_File); ok {
		return x.File
	}
	return ""
}

func (m *FileServer_Entry) GetDirectory() string {
	if x, ok := m.GetFileOrDir().(*FileServer_Entry_Directory); ok {
		return x.Directory
	}
	return ""
}

func (m *FileServer_Entry) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FileServer_Entry) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FileServer_Entry_OneofMarshaler, _FileServer_Entry_OneofUnmarshaler, _FileServer_Entry_OneofSizer, []interface{}{
		(*FileServer_Entry_File)(nil),
		(*FileServer_Entry_Directory)(nil),
	}
}

func _FileServer_Entry_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FileServer_Entry)
	// FileOrDir
	switch x := m.FileOrDir.(type) {
	case *FileServer_Entry_File:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.File)
	case *FileServer_Entry_Directory:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Directory)
	case nil:
	default:
		return fmt.Errorf("FileServer_Entry.FileOrDir has unexpected type %T", x)
	}
	return nil
}

func _FileServer_Entry_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FileServer_Entry)
	switch tag {
	case 1: // FileOrDir.File
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.FileOrDir = &FileServer_Entry_File{x}
		return true, err
	case 2: // FileOrDir.Directory
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.FileOrDir = &FileServer_Entry_Directory{x}
		return true, err
	default:
		return false, nil
	}
}

func _FileServer_Entry_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FileServer_Entry)
	// FileOrDir
	switch x := m.FileOrDir.(type) {
	case *FileServer_Entry_File:
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.File)))
		n += len(x.File)
	case *FileServer_Entry_Directory:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Directory)))
		n += len(x.Directory)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Server struct {
	Domain   []string                               `protobuf:"bytes,1,rep,name=domain" json:"domain,omitempty"`
	Settings *v2ray_core_common_serial.TypedMessage `protobuf:"bytes,2,opt,name=settings" json:"settings,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Server) GetDomain() []string {
	if m != nil {
		return m.Domain
	}
	return nil
}

func (m *Server) GetSettings() *v2ray_core_common_serial.TypedMessage {
	if m != nil {
		return m.Settings
	}
	return nil
}

type Config struct {
	Server []*Server `protobuf:"bytes,1,rep,name=server" json:"server,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Config) GetServer() []*Server {
	if m != nil {
		return m.Server
	}
	return nil
}

func init() {
	proto.RegisterType((*FileServer)(nil), "v2ray.core.app.web.FileServer")
	proto.RegisterType((*FileServer_Entry)(nil), "v2ray.core.app.web.FileServer.Entry")
	proto.RegisterType((*Server)(nil), "v2ray.core.app.web.Server")
	proto.RegisterType((*Config)(nil), "v2ray.core.app.web.Config")
}

func init() { proto.RegisterFile("v2ray.com/core/app/web/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x91, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0x87, 0x4d, 0xff, 0x04, 0x33, 0xbd, 0x2d, 0x52, 0x42, 0x0f, 0x52, 0xaa, 0x48, 0x4f, 0x1b,
	0x89, 0x9e, 0xc4, 0x8b, 0x6d, 0x15, 0x2f, 0xa2, 0xac, 0xa2, 0xe0, 0x41, 0xd9, 0xa4, 0x63, 0x5d,
	0x68, 0xb2, 0xcb, 0x64, 0x69, 0xc9, 0x1b, 0x89, 0x4f, 0x29, 0xd9, 0x44, 0x2b, 0xda, 0x5b, 0x26,
	0xf3, 0xfd, 0x76, 0x66, 0xbf, 0x85, 0x83, 0x55, 0x4c, 0xb2, 0xe4, 0xa9, 0xce, 0xa2, 0x54, 0x13,
	0x46, 0xd2, 0x98, 0x68, 0x8d, 0x49, 0x94, 0xea, 0xfc, 0x4d, 0x2d, 0xb8, 0x21, 0x6d, 0x35, 0x63,
	0xdf, 0x10, 0x21, 0x97, 0xc6, 0xf0, 0x35, 0x26, 0x83, 0xe3, 0x3f, 0xc1, 0x54, 0x67, 0x99, 0xce,
	0xa3, 0x02, 0x49, 0xc9, 0x65, 0x64, 0x4b, 0x83, 0xf3, 0xd7, 0x0c, 0x8b, 0x42, 0x2e, 0xb0, 0x3e,
	0x65, 0xf4, 0xe1, 0x01, 0x5c, 0xa9, 0x25, 0xde, 0x23, 0xad, 0x90, 0xd8, 0x19, 0x74, 0x31, 0xb7,
	0x54, 0x86, 0xde, 0xb0, 0x3d, 0xee, 0xc5, 0x87, 0xfc, 0xff, 0x10, 0xbe, 0xc1, 0xf9, 0x65, 0xc5,
	0x8a, 0x3a, 0x32, 0x78, 0x81, 0xae, 0xab, 0xd9, 0x1e, 0x74, 0x2a, 0x26, 0xf4, 0x86, 0xde, 0x38,
	0xb8, 0xde, 0x11, 0xae, 0x62, 0xfb, 0x10, 0xcc, 0x14, 0x61, 0x6a, 0x35, 0x95, 0x61, 0xab, 0x69,
	0x6d, 0x7e, 0x31, 0x06, 0x1d, 0x23, 0xed, 0x7b, 0xd8, 0xae, 0x5a, 0xc2, 0x7d, 0x4f, 0x7a, 0x10,
	0x54, 0xd9, 0x5b, 0x9a, 0x29, 0x1a, 0xcd, 0xc1, 0x6f, 0xb6, 0xec, 0x83, 0x3f, 0xd7, 0x99, 0x54,
	0xb9, 0x5b, 0x33, 0x10, 0x4d, 0xc5, 0x26, 0xb0, 0x5b, 0xa0, 0xb5, 0x2a, 0x5f, 0x14, 0x6e, 0x42,
	0x2f, 0x3e, 0xfa, 0x7d, 0x81, 0xda, 0x06, 0xaf, 0x6d, 0xf0, 0x87, 0xca, 0xc6, 0x4d, 0x2d, 0x43,
	0xfc, 0xe4, 0x46, 0xe7, 0xe0, 0x4f, 0x9d, 0x66, 0x16, 0x83, 0x5f, 0xb8, 0x79, 0x8d, 0x8c, 0xc1,
	0x36, 0x19, 0xf5, 0x46, 0xa2, 0x21, 0x27, 0xa7, 0xd0, 0x4f, 0x75, 0xb6, 0x05, 0xbc, 0xf3, 0x9e,
	0xdb, 0x6b, 0x4c, 0x3e, 0x5b, 0xec, 0x31, 0x16, 0xb2, 0xe4, 0xd3, 0xaa, 0x77, 0x61, 0x0c, 0x7f,
	0xc2, 0x24, 0xf1, 0xdd, 0x5b, 0x9c, 0x7c, 0x05, 0x00, 0x00, 0xff, 0xff, 0xf9, 0x2a, 0xe1, 0x59,
	0xf8, 0x01, 0x00, 0x00,
}
