package http

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(internet.TransportProtocol_HTTP, func() interface{} {
		return new(Config)
	}))
}
