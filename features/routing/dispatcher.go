package routing

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/features"
	"v2ray.com/core/transport"
)

// Dispatcher is a feature that dispatches inbound requests to outbound handlers based on rules.
// Dispatcher is required to be registered in a V2Ray instance to make V2Ray function properly.
//
// v2ray:api:stable
type Dispatcher interface {
	features.Feature

	// Dispatch returns a Ray for transporting data for the given request.
	Dispatch(ctx context.Context, dest net.Destination) (*transport.Link, error)
}

// DispatcherType returns the type of Dispatcher interface. Can be used to implement common.HasType.
//
// v2ray:api:stable
func DispatcherType() interface{} {
	return (*Dispatcher)(nil)
}
