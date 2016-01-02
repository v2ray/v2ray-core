package json

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
	"github.com/v2ray/v2ray-core/proxy/internal/config/json"
)

type HttpProxyConfig struct {
}

func init() {
	config.RegisterInboundConnectionConfig("http", json.JsonConfigLoader(func() interface{} {
		return new(HttpProxyConfig)
	}))
}
