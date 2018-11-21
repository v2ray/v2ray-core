package quic

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

//go:generate errorgen

// Here is some modification needs to be done before update quic vendor.
// * use bytespool in buffer_pool.go
// * set MaxReceivePacketSize to 1452 - 32 (16 bytes auth, 16 bytes head)
//
//

const protocolName = "quic"

func init() {
	common.Must(internet.RegisterProtocolConfigCreatorByName(protocolName, func() interface{} {
		return new(Config)
	}))
}
