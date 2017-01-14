// Package proxyman defines applications for manageing inbound and outbound proxies.
package proxyman

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type InboundHandlerManager interface {
	GetHandler(tag string) (proxy.InboundHandler, int)
}

type InboundHandler interface {
}

type OutboundHandlerManager interface {
	GetHandler(tag string) proxy.OutboundHandler
	GetDefaultHandler() proxy.OutboundHandler
	SetDefaultHandler(handler proxy.OutboundHandler) error
	SetHandler(tag string, handler proxy.OutboundHandler) error
}

type OutboundHandler interface {
	Dispatch(ctx context.Context, destination net.Destination, outboundRay ray.OutboundRay)
}

func InboundHandlerManagerFromSpace(space app.Space) InboundHandlerManager {
	app := space.GetApplication((*InboundHandlerManager)(nil))
	if app == nil {
		return nil
	}
	return app.(InboundHandlerManager)
}

func OutboundHandlerManagerFromSpace(space app.Space) OutboundHandlerManager {
	app := space.GetApplication((*OutboundHandlerManager)(nil))
	if app == nil {
		return nil
	}
	return app.(OutboundHandlerManager)
}
