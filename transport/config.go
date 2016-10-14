package transport

import (
	"v2ray.com/core/transport/internet"
)

// Apply applies this Config.
func (this *Config) Apply() error {
	if this == nil {
		return nil
	}
	if err := internet.ApplyGlobalNetworkSettings(this.NetworkSettings); err != nil {
		return err
	}
	return nil
}
