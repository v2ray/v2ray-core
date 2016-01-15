// +build json

package freedom

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

func init() {
	config.RegisterOutboundConnectionConfig("freedom",
		func(data []byte) (interface{}, error) {
			return new(Config), nil
		})
}
