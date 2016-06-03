package repo

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func CreateInboundHandler(name string, space app.Space, rawConfig []byte, listen v2net.Address, port v2net.Port) (proxy.InboundHandler, error) {
	return internal.CreateInboundHandler(name, space, rawConfig, listen, port)
}

func CreateOutboundHandler(name string, space app.Space, rawConfig []byte, sendThrough v2net.Address) (proxy.OutboundHandler, error) {
	return internal.CreateOutboundHandler(name, space, rawConfig, sendThrough)
}
