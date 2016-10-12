package core

import (
	"v2ray.com/core/proxy"
)

type InboundDetourHandler interface {
	Start() error
	Close()
	GetConnectionHandler() (proxy.InboundHandler, int)
}
