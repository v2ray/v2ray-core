package core

//go:generate go install "google.golang.org/protobuf/proto"
//go:generate go install "google.golang.org/protobuf/cmd/protoc-gen-go"
//go:generate go install "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
//go:generate go install "github.com/gogo/protobuf/protoc-gen-gofast"
//go:generate go install "v2ray.com/core/infra/vprotogen"
//go:generate vprotogen -repo v2ray.com/core
