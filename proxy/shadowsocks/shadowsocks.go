// R.I.P Shadowsocks

package shadowsocks

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Shadowsocks struct {
	config *Config
	port   v2net.Port
}

func (this *Shadowsocks) Port() v2net.Port {
	return this.port
}

func (this *Shadowsocks) Listen(port v2net.Port) error {
	return nil
}
