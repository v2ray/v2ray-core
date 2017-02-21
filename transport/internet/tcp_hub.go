package internet

import (
	"net"

	"v2ray.com/core/app/log"
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
	listener     Listener
	connCallback ConnectionHandler
	closed       chan bool
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
	defer func() {
		recover()
	}()

	select {
	case <-v.closed:
		return
	default:
		v.listener.Close()
		close(v.closed)
	}
}

func (v *TCPHub) start() {
	for {
		select {
		case <-v.closed:
			return
		default:
		}
		var newConn Connection
		err := retry.ExponentialBackoff(10, 500).On(func() error {
			select {
			case <-v.closed:
				return nil
			default:
				conn, err := v.listener.Accept()
				if err != nil {
					return errors.Base(err).RequireUserAction().Message("Internet|Listener: Failed to accept new TCP connection.")
				}
				newConn = conn
				return nil
			}
		})
		if err != nil {
			if errors.IsActionRequired(err) {
				log.Warning(err)
			} else {
				log.Info(err)
			}
			continue
		}
		if newConn != nil {
			go v.connCallback(newConn)
		}
	}
}
