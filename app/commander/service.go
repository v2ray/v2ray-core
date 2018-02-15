package commander

import (
	"google.golang.org/grpc"
)

type Service interface {
	Register(*grpc.Server)
}
