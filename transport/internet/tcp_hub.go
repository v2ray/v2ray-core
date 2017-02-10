package internet

import (
	"net"
	"sync"

	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
)

var (
	transportListenerCache = make(map[TransportProtocol]ListenFunc)
)

func RegisterTransportListener(protocol TransportProtocol, listener ListenFunc) error {
	if _, found := transportListenerCache[protocol]; found {
		return errors.New("Internet|TCPHub: ", protocol, " listener already registered.")
	}
	transportListenerCache[protocol] = listener
	return nil
}

type ListenFunc func(address v2net.Address, port v2net.Port, options ListenOptions) (Listener, error)
type ListenOptions struct {
	Stream       *StreamConfig
	RecvOrigDest bool
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
	options := ListenOptions{
		Stream: settings,
	}
	protocol := settings.GetEffectiveProtocol()
	listenFunc := transportListenerCache[protocol]
	if listenFunc == nil {
		return nil, errors.New("Internet|TCPHub: ", protocol, " listener not registered.")
	}
	listener, err := listenFunc(address, port, options)
	if err != nil {
		return nil, errors.Base(err).Message("Internet|TCPHub: Failed to listen on address: ", address, ":", port)
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
		err := retry.ExponentialBackoff(10, 500).On(func() error {
			if !v.accepting {
				return nil
			}
			conn, err := v.listener.Accept()
			if err != nil {
				return errors.Base(err).Message("Internet|Listener: Failed to accept new TCP connection.")
			}
			newConn = conn
			return nil
		})
		if err == nil && newConn != nil {
			go v.connCallback(newConn)
		}
	}
}
