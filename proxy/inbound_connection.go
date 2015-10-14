package proxy

import (
	"github.com/v2ray/v2ray-core/app"
)

type InboundConnectionHandlerFactory interface {
	Create(dispatch app.PacketDispatcher, config interface{}) (InboundConnectionHandler, error)
}

type InboundConnectionHandler interface {
	Listen(port uint16) error
}
