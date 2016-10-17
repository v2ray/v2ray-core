package conf

import (
	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/http"
)

type HttpServerConfig struct {
	Timeout uint32 `json:"timeout"`
}

func (this *HttpServerConfig) Build() (*loader.TypedSettings, error) {
	config := &http.ServerConfig{
		Timeout: this.Timeout,
	}

	return loader.NewTypedSettings(config), nil
}
