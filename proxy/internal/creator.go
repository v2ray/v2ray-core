package internal

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
)

type InboundHandlerCreator func(space app.Space, config interface{}, listenOn v2net.Address, port v2net.Port) (proxy.InboundHandler, error)
type OutboundHandlerCreator func(space app.Space, config interface{}, sendThrough v2net.Address) (proxy.OutboundHandler, error)
