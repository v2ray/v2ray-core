package core

//go:generate go install -v google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install -v github.com/gogo/protobuf/protoc-gen-gofast
//go:generate go run ./infra/vprotogen/main.go

import "path/filepath"

// ProtoFilesUsingProtocGenGoFast is the map of Proto files
// that use `protoc-gen-gofast` to generate pb.go files
var ProtoFilesUsingProtocGenGoFast = map[string]bool{"proxy/vless/encoding/addons.proto": true}

// ProtocMap is the map of paths to `protoc` binary excutable files of specific platform
var ProtocMap = map[string]string{
	"windows": filepath.Join(".dev", "protoc", "windows", "protoc.exe"),
	"darwin":  filepath.Join(".dev", "protoc", "macos", "protoc"),
	"linux":   filepath.Join(".dev", "protoc", "linux", "protoc"),
}
