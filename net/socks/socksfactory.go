package socks

import (
	"github.com/v2ray/v2ray-core"
)

type SocksServerFactory struct {
}

func (factory *SocksServerFactory) Create(vp *core.VPoint) *SocksServer {
	return NewSocksServer(vp)
}
