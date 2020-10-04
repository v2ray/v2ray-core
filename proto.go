package core

//go:generate go install -v google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install -v github.com/gogo/protobuf/protoc-gen-gofast
//go:generate go run ./infra/vprotogen/main.go
