package router

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/ray"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error)
}
