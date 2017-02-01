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
	GetHandler(ctx context.Context, tag string) (InboundHandler, error)
	AddHandler(ctx context.Context, config *InboundHandlerConfig) error
	Start() error
	Close()
}

type InboundHandler interface {
	Start() error
	Close()

	// For migration
	GetRandomInboundProxy() (proxy.Inbound, net.Port, int)
}

type OutboundHandlerManager interface {
	GetHandler(tag string) OutboundHandler
	GetDefaultHandler() OutboundHandler
	AddHandler(ctx context.Context, config *OutboundHandlerConfig) error
	Start() error
	Close()
}

type OutboundHandler interface {
	Dispatch(ctx context.Context, outboundRay ray.OutboundRay)
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
