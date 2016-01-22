package testing

import (
	"github.com/v2ray/v2ray-core/common/dice"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

func PickPort() v2net.Port {
	return v2net.Port(30000 + dice.Roll(10000))
}
