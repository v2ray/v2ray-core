package websocket

import (
	"net/http"

	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func (c *Config) GetNormalizedPath() string {
	path := c.Path
	if len(path) == 0 {
		return "/"
	}
	if path[0] != '/' {
		return "/" + path
	}
	return path
}

func (c *Config) GetRequestHeader() http.Header {
	header := http.Header{}
	for _, h := range c.Header {
		header.Add(h.Key, h.Value)
	}
	return header
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(internet.TransportProtocol_WebSocket, func() interface{} {
		return new(Config)
	}))
}
