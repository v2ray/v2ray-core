package dokodemo

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func init() {
	if err := internal.RegisterInboundConnectionHandlerFactory("dokodemo-door", func(space app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
		config := rawConfig.(Config)
		return NewDokodemoDoor(space, config), nil
	}); err != nil {
		panic(err)
	}
}
