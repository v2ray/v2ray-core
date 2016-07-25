package socks

import (
	"github.com/v2ray/v2ray-core/common/protocol"
)

type ClientConfig struct {
	Servers []*protocol.ServerSpec
}
