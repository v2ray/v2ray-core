package tcp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(internet.TransportProtocol_TCP, func() interface{} {
		return new(Config)
	}))
}
