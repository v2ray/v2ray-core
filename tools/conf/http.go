package conf

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/http"
)

type HttpServerConfig struct {
	Timeout uint32 `json:"timeout"`
}

func (v *HttpServerConfig) Build() (*serial.TypedMessage, error) {
	config := &http.ServerConfig{
		Timeout: v.Timeout,
	}

	return serial.ToTypedMessage(config), nil
}
