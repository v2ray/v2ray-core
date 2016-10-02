package internet

import (
	"errors"
	"net"
	"sync"

	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

var (
	ErrClosedConnection = errors.New("Connection already closed.")

	KCPListenFunc    ListenFunc
	TCPListenFunc    ListenFunc
	RawTCPListenFunc ListenFunc
	WSListenFunc     ListenFunc
)

type ListenFunc func(address v2net.Address, port v2net.Port, options ListenOptions) (Listener, error)
type ListenOptions struct {
	Stream *StreamConfig
}

type Listener interface {
	Accept() (Connection, error)
	Close() error
	Addr() net.Addr
}

type TCPHub struct {
	sync.Mutex
	listener     Listener
	connCallback ConnectionHandler
	accepting    bool
}

func ListenTCP(address v2net.Address, port v2net.Port, callback ConnectionHandler, settings *StreamConfig) (*TCPHub, error) {
	var listener Listener
	var err error
	options := ListenOptions{
		Stream: settings,
	}
	switch settings.Network {
	case v2net.Network_TCP:
		listener, err = TCPListenFunc(address, port, options)
	case v2net.Network_KCP:
		listener, err = KCPListenFunc(address, port, options)
	case v2net.Network_WebSocket:
		listener, err = WSListenFunc(address, port, options)
	case v2net.Network_RawTCP:
		listener, err = RawTCPListenFunc(address, port, options)
	default:
		log.Error("Internet|Listener: Unknown stream type: ", settings.Network)
		err = ErrUnsupportedStreamType
	}

	if err != nil {
		log.Warning("Internet|Listener: Failed to listen on ", address, ":", port)
		return nil, err
	}

	hub := &TCPHub{
		listener:     listener,
		connCallback: callback,
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
				log.Warning("Internet|Listener: Failed to accept new TCP connection: ", err)
			}
			continue
		}
		go this.connCallback(conn)
	}
}
