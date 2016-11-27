package dokodemo

import (
	v2net "v2ray.com/core/common/net"
)

func (v *Config) GetPredefinedAddress() v2net.Address {
	addr := v.Address.AsAddress()
	if addr == nil {
		return nil
	}
	return addr
}
