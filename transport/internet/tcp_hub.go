package internet

import (
	"net"
	"sync"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
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

func (v *TCPHub) Close() {
	v.accepting = false
	v.listener.Close()
}

func (v *TCPHub) start() {
	v.accepting = true
	for v.accepting {
		var newConn Connection
		err := retry.ExponentialBackoff(10, 200).On(func() error {
			if !v.accepting {
				return nil
			}
			conn, err := v.listener.Accept()
			if err != nil {
				if v.accepting {
					log.Warning("Internet|Listener: Failed to accept new TCP connection: ", err)
				}
				return err
			}
			newConn = conn
			return nil
		})
		if err == nil && newConn != nil {
			go v.connCallback(newConn)
		}
	}
}
