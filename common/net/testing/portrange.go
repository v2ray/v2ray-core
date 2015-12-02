package testing

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type PortRange struct {
	FromValue v2net.Port
	ToValue   v2net.Port
}

func (this *PortRange) From() v2net.Port {
	return this.FromValue
}

func (this *PortRange) To() v2net.Port {
	return this.ToValue
}
