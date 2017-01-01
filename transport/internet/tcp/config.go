package tcp

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func (v *Config) IsConnectionReuse() bool {
	if v == nil || v.ConnectionReuse == nil {
		return true
	}
	return v.ConnectionReuse.Enable
}

func init() {
	internet.RegisterNetworkConfigCreator(v2net.Network_TCP, func() interface{} {
		return new(Config)
	})
}
