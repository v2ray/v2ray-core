package testing

import (
	"sync/atomic"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	port = int32(30000)
)

func PickPort() v2net.Port {
	return v2net.Port(atomic.AddInt32(&port, 1))
}
