package proxyman

import (
	"v2ray.com/core/app"
	"v2ray.com/core/proxy"
)

const (
	APP_ID_INBOUND_MANAGER  = app.ID(4)
	APP_ID_OUTBOUND_MANAGER = app.ID(6)
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
