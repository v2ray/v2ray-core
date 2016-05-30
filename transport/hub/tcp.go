package hub

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	ErrorClosedConnection = errors.New("Connection already closed.")
)

type TCPHub struct {
	sync.Mutex
	listener     net.Listener
	connCallback ConnectionHandler
	accepting    bool
}

func ListenTCP(address v2net.Address, port v2net.Port, callback ConnectionHandler, tlsConfig *tls.Config) (*TCPHub, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	var hub *TCPHub
	if tlsConfig != nil {
		tlsListener := tls.NewListener(listener, tlsConfig)
		hub = &TCPHub{
			listener:     tlsListener,
			connCallback: callback,
		}
	} else {
		hub = &TCPHub{
			listener:     listener,
			connCallback: callback,
		}
	}

	go hub.start()
	return hub, nil
}

func (this *TCPHub) Close() {
	this.accepting = false
	this.listener.Close()
}

func (this *TCPHub) start() {
	this.accepting = true
	for this.accepting {
		conn, err := this.listener.Accept()

		if err != nil {
			if this.accepting {
				log.Warning("Listener: Failed to accept new TCP connection: ", err)
			}
			continue
		}
		go this.connCallback(&Connection{
			dest:     conn.RemoteAddr().String(),
			conn:     conn,
			listener: this,
		})
	}
}

// @Private
func (this *TCPHub) Recycle(dest string, conn net.Conn) {
	if this.accepting {
		go this.connCallback(&Connection{
			dest:     dest,
			conn:     conn,
			listener: this,
		})
	}
}
