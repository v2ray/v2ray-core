package dokodemo

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

func init() {
	if err := proxy.RegisterInboundConnectionHandlerFactory("dokodemo-door", func(space app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
		config := rawConfig.(Config)
		return NewDokodemoDoor(space, config), nil
	}); err != nil {
		panic(err)
	}
}
