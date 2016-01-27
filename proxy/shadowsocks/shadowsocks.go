// R.I.P Shadowsocks

package shadowsocks

import (
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/listener"
)

type Shadowsocks struct {
	config      *Config
	port        v2net.Port
	accepting   bool
	tcpListener *listener.TCPListener
}

func (this *Shadowsocks) Port() v2net.Port {
	return this.port
}

func (this *Shadowsocks) Close() {
	this.accepting = false
	this.tcpListener.Close()
	this.tcpListener = nil
}

func (this *Shadowsocks) Listen(port v2net.Port) error {
	if this.accepting {
		if this.port == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}

	tcpListener, err := listener.ListenTCP(port, this.handleConnection)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen on port ", port, ": ", err)
		return err
	}
	this.tcpListener = tcpListener
	this.accepting = true
	return nil
}

func (this *Shadowsocks) handleConnection(conn *listener.TCPConn) {
	defer conn.Close()

}
