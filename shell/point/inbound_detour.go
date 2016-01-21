package point

import (
	"github.com/v2ray/v2ray-core/proxy"
)

type InboundDetourHandler interface {
	Start() error
	Close()
	GetConnectionHandler() (proxy.InboundConnectionHandler, int)
}
