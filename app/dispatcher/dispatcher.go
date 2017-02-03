package dispatcher

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/ray"
)

// Interface dispatch a packet and possibly further network payload to its destination.
type Interface interface {
	Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error)
	Start() error
	Close()
}

func FromSpace(space app.Space) Interface {
	if app := space.GetApplication((*Interface)(nil)); app != nil {
		return app.(Interface)
	}
	return nil
}
