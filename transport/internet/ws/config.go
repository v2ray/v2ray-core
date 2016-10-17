package ws

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func (this *ConnectionReuse) IsEnabled() bool {
	if this == nil {
		return false
	}
	return this.Enable
}

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_WebSocket, func() interface{} {
		return new(Config)
	})
}
