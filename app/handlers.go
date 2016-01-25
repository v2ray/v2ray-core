package app

import (
	"github.com/v2ray/v2ray-core/proxy"
)

type InboundHandlerManager interface {
	GetHandler(tag string) (proxy.InboundHandler, int)
}
