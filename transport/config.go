package transport

import (
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/ws"
)

// Config for V2Ray transport layer.
type Config struct {
	tcpConfig *tcp.Config
	kcpConfig kcp.Config
	wsConfig  *ws.Config
}

// Apply applies this Config.
func (this *Config) Apply() error {
	if this.tcpConfig != nil {
		this.tcpConfig.Apply()
	}
	this.kcpConfig.Apply()
	if this.wsConfig != nil {
		this.wsConfig.Apply()
	}
	return nil
}
