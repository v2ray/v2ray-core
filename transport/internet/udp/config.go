package udp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
