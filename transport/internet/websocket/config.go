package websocket

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func (c *Config) IsConnectionReuse() bool {
	if c == nil || c.ConnectionReuse == nil {
		return true
	}
	return c.ConnectionReuse.Enable
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(internet.TransportProtocol_WebSocket, func() interface{} {
		return new(Config)
	}))
}
