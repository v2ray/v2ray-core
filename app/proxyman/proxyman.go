// Package proxyman defines applications for manageing inbound and outbound proxies.
package proxyman

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

type InboundHandlerManager interface {
	GetHandler(tag string) (proxy.InboundHandler, int)
}

type OutboundHandlerManager interface {
	GetHandler(tag string) proxy.OutboundHandler
	GetDefaultHandler() proxy.OutboundHandler
	SetDefaultHandler(handler proxy.OutboundHandler) error
	SetHandler(tag string, handler proxy.OutboundHandler) error
}

func InboundHandlerManagerFromSpace(space app.Space) InboundHandlerManager {
	app := space.(app.AppGetter).GetApp(serial.GetMessageType((*InboundConfig)(nil)))
	if app == nil {
		return nil
	}
	return app.(InboundHandlerManager)
}

func OutboundHandlerManagerFromSpace(space app.Space) OutboundHandlerManager {
	app := space.(app.AppGetter).GetApp(serial.GetMessageType((*OutboundConfig)(nil)))
	if app == nil {
		return nil
	}
	return app.(OutboundHandlerManager)
}
