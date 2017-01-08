package dokodemo

import (
	v2net "v2ray.com/core/common/net"
)

// GetPredefinedAddress returns the defined address from proto config. Null if address is not valid.
func (v *Config) GetPredefinedAddress() v2net.Address {
	addr := v.Address.AsAddress()
	if addr == nil {
		return nil
	}
	return addr
}
