package transport

import (
	"v2ray.com/core/transport/internet"
)

// Apply applies this Config.
func (this *Config) Apply() error {
	internet.ApplyGlobalNetworkSettings(this.NetworkSettings)
	return nil
}
