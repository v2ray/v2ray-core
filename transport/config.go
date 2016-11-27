package transport

import (
	"v2ray.com/core/transport/internet"
)

// Apply applies this Config.
func (v *Config) Apply() error {
	if v == nil {
		return nil
	}
	if err := internet.ApplyGlobalNetworkSettings(v.NetworkSettings); err != nil {
		return err
	}
	return nil
}
