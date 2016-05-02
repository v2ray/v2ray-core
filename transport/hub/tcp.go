package hub

import (
	"errors"
	"net"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	ErrorClosedConnection = errors.New("Connection already closed.")
)

type TCPHub struct {
	listener     *net.TCPListener
	connCallback ConnectionHandler
	accepting    bool
}

func ListenTCP(port v2net.Port, callback ConnectionHandler) (*TCPHub, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	tcpListener := &TCPHub{
		listener:     listener,
		connCallback: callback,
	}
	go tcpListener.start()
	return tcpListener, nil
}

func (this *TCPHub) Close() {
	this.accepting = false
	this.listener.Close()
	this.listener = nil
}

func (this *TCPHub) start() {
	this.accepting = true
	for this.accepting {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
			if this.accepting {
				log.Warning("Listener: Failed to accept new TCP connection: ", err)
			}
			continue
		}
		go this.connCallback(&Connection{
			conn:     conn,
			listener: this,
		})
	}
}

func (this *TCPHub) recycle(conn *net.TCPConn) {

}
