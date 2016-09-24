package dokodemo

import (
	v2net "v2ray.com/core/common/net"
)

func (this *Config) GetPredefinedAddress() v2net.Address {
	addr := this.Address.AsAddress()
	if addr == nil {
		return nil
	}
	return addr
}
