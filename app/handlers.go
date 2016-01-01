package app

import (
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type InboundHandlerManager interface {
	GetHandler(tag string) (connhandler.InboundConnectionHandler, int)
}
