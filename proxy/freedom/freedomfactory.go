package freedom

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func init() {
	if err := internal.RegisterOutboundConnectionHandlerFactory("freedom", func(space app.Space, config interface{}) (connhandler.OutboundConnectionHandler, error) {
		return &FreedomConnection{space: space}, nil
	}); err != nil {
		panic(err)
	}
}
