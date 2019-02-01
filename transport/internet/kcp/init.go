// +build !confonly

package kcp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, DialKCP))
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, ListenKCP))
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
