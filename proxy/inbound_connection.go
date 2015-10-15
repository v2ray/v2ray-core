package proxy

import (
	"github.com/v2ray/v2ray-core/app"
)

// A InboundConnectionHandlerFactory creates InboundConnectionHandler on demand.
type InboundConnectionHandlerFactory interface {
	// Create creates a new InboundConnectionHandler with given configuration.
	Create(dispatch app.PacketDispatcher, config interface{}) (InboundConnectionHandler, error)
}

// A InboundConnectionHandler handles inbound network connections to V2Ray.
type InboundConnectionHandler interface {
	// Listen starts a InboundConnectionHandler by listen on a specific port. This method is called
	// exactly once during runtime.
	Listen(port uint16) error
}
