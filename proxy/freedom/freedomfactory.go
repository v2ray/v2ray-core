package freedom

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

func init() {
	if err := proxy.RegisterOutboundConnectionHandlerFactory("freedom", func(space app.Space, config interface{}) (connhandler.OutboundConnectionHandler, error) {
		return &FreedomConnection{space: space}, nil
	}); err != nil {
		panic(err)
	}
}
