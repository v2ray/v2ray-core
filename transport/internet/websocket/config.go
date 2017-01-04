package websocket

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func (c *Config) IsConnectionReuse() bool {
	if c == nil || c.ConnectionReuse == nil {
		return true
	}
	return c.ConnectionReuse.Enable
}

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_WebSocket, func() interface{} {
		return new(Config)
	})
}
