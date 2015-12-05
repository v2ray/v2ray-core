package connhandler

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

// A InboundConnectionHandlerFactory creates InboundConnectionHandler on demand.
type InboundConnectionHandlerFactory interface {
	// Create creates a new InboundConnectionHandler with given configuration.
	Create(space *app.Space, config interface{}) (InboundConnectionHandler, error)
}

// A InboundConnectionHandler handles inbound network connections to V2Ray.
type InboundConnectionHandler interface {
	// Listen starts a InboundConnectionHandler by listen on a specific port. This method is called
	// exactly once during runtime.
	Listen(port v2net.Port) error
}
