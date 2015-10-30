package json

import (
	"github.com/v2ray/v2ray-core/proxy/common/config/json"
)

type HttpProxyConfig struct {
}

func init() {
	json.RegisterInboundConnectionConfig("http", func() interface{} {
		return new(HttpProxyConfig)
	})
}
