package http

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Config struct {
	OwnHosts []v2net.Address
}

func (this *Config) IsOwnHost(host v2net.Address) bool {
	for _, ownHost := range this.OwnHosts {
		if ownHost.Equals(host) {
			return true
		}
	}
	return false
}
