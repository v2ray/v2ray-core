package core

//go:generate go install google.golang.org/protobuf/proto
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@v0.0.0-20200825170228-39ef2aaf62df
//go:generate go install github.com/gogo/protobuf/protoc-gen-gofast
//go:generate go run v2ray.com/core/infra/vprotogen
