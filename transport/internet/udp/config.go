package udp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreatorByName(protocolName, func() interface{} {
		return new(Config)
	}))
}
