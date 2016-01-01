package connhandler

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

// A InboundConnectionHandler handles inbound network connections to V2Ray.
type InboundConnectionHandler interface {
	// Listen starts a InboundConnectionHandler by listen on a specific port. This method is called
	// exactly once during runtime.
	Listen(port v2net.Port) error
}
