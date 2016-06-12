package transport

import "github.com/v2ray/v2ray-core/transport/hub/kcpv"

// Config for V2Ray transport layer.
type Config struct {
	ConnectionReuse bool
	enableKcp       bool
	kcpConfig       *kcpv.Config
}

// Apply applies this Config.
func (this *Config) Apply() error {
	if this.ConnectionReuse {
		connectionReuse = true
	}
	enableKcp = this.enableKcp
	if enableKcp {
		KcpConfig = this.kcpConfig
	}
	return nil
}
