package dokodemo

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Config struct {
	Address v2net.Address
	Port    v2net.Port
	Network v2net.NetworkList
	Timeout int
}
