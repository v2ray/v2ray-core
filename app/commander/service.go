package commander

import (
	"google.golang.org/grpc"
)

// Service is a Commander service.
type Service interface {
	// Register registers the service itself to a gRPC server.
	Register(*grpc.Server)
}
