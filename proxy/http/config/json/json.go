package json

import (
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
)

type HttpProxyConfig struct {
}

func init() {
	json.RegisterConfigType("http", config.TypeInbound, func() interface{} {
		return new(HttpProxyConfig)
	})
}
