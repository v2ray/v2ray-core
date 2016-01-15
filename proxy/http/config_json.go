// +build json

package http

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

func init() {
	config.RegisterInboundConnectionConfig("http",
		func(data []byte) (interface{}, error) {
			return new(Config), nil
		})
}
