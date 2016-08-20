package dokodemo

import (
	v2net "v2ray.com/core/common/net"
)

type Config struct {
	FollowRedirect bool
	Address        v2net.Address
	Port           v2net.Port
	Network        *v2net.NetworkList
	Timeout        int
}
