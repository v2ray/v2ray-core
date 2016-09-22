package dokodemo

import (
	v2net "v2ray.com/core/common/net"
)

func (this *Config) GetPredefinedAddress() v2net.Address {
	return this.Address.AsAddress()
}
