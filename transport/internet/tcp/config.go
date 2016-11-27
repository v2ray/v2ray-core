package tcp

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func (v *ConnectionReuse) IsEnabled() bool {
	if v == nil {
		return true
	}
	return v.Enable
}

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_TCP, func() interface{} {
		return new(Config)
	})
}
