package freedom

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func init() {
	internal.MustRegisterOutboundConnectionHandlerCreator("freedom",
		func(space app.Space, config interface{}) (proxy.OutboundConnectionHandler, error) {
			return &FreedomConnection{space: space}, nil
		})
}
