package socks

import (
	"v2ray.com/core/common/protocol"
)

type ClientConfig struct {
	Servers []*protocol.ServerSpec
}
