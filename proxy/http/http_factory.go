package http

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func init() {
	internal.MustRegisterInboundHandlerCreator("http",
		func(space app.Space, rawConfig interface{}) (proxy.InboundHandler, error) {
			if !space.HasApp(dispatcher.APP_ID) {
				return nil, internal.ErrorBadConfiguration
			}
			return NewHttpProxyServer(
				rawConfig.(*Config),
				space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)), nil
		})
}
