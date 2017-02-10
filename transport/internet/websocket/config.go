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

func (c *Config) GetNormailzedPath() string {
	path := c.Path
	if len(path) == 0 {
		return "/"
	}
	if path[0] != '/' {
		return "/" + path
	}
	return path
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(internet.TransportProtocol_WebSocket, func() interface{} {
		return new(Config)
	}))
}
