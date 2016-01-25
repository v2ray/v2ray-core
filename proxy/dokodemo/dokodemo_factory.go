package dokodemo

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func init() {
	internal.MustRegisterInboundHandlerCreator("dokodemo-door",
		func(space app.Space, rawConfig interface{}) (proxy.InboundHandler, error) {
			config := rawConfig.(*Config)
			return NewDokodemoDoor(space, config), nil
		})
}
