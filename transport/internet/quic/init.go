// +build !confonly

package quic

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, Listen))
}
