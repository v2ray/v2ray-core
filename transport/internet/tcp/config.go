package tcp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/transport/internet"
)

func (v *Config) IsConnectionReuse() bool {
	if v == nil || v.ConnectionReuse == nil {
		return true
	}
	return v.ConnectionReuse.Enable
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(internet.TransportProtocol_TCP, func() interface{} {
		return new(Config)
	}))
}
