// Package proxyman defines applications for manageing inbound and outbound proxies.
package proxyman

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg proxyman -path App,Proxyman

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

type key int

const (
	protocolsKey key = iota
)

func ContextWithProtocolSniffers(ctx context.Context, list []KnownProtocols) context.Context {
	return context.WithValue(ctx, protocolsKey, list)
}

func ProtocoSniffersFromContext(ctx context.Context) []KnownProtocols {
	if list, ok := ctx.Value(protocolsKey).([]KnownProtocols); ok {
		return list
	}
	return nil
}
