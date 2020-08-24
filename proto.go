package core

//go:generate go get -u "google.golang.org/protobuf/proto"
//go:generate go get -u "google.golang.org/protobuf/cmd/protoc-gen-go"
//go:generate go get -u "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
//go:generate go get -u "github.com/gogo/protobuf/protoc-gen-gofast"
//go:generate go install "v2ray.com/core/infra/vprotogen"
//go:generate vprotogen -repo v2ray.com/core
