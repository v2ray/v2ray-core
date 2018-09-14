package internet

import (
	"context"

	"v2ray.com/core/common/net"
)

var (
	transportListenerCache = make(map[string]ListenFunc)
)

func RegisterTransportListener(protocol string, listener ListenFunc) error {
	if _, found := transportListenerCache[protocol]; found {
		return newError(protocol, " listener already registered.").AtError()
	}
	transportListenerCache[protocol] = listener
	return nil
}

type ConnHandler func(Connection)

type ListenFunc func(ctx context.Context, address net.Address, port net.Port, handler ConnHandler) (Listener, error)

type Listener interface {
	Close() error
	Addr() net.Addr
}

func ListenTCP(ctx context.Context, address net.Address, port net.Port, handler ConnHandler) (Listener, error) {
	settings := StreamSettingsFromContext(ctx)
	if settings == nil {
		s, err := ToMemoryStreamConfig(nil)
		if err != nil {
			return nil, newError("failed to create default stream settings").Base(err)
		}
		settings = s
		ctx = ContextWithStreamSettings(ctx, settings)
	}

	if address.Family().IsDomain() && address.Domain() == "localhost" {
		address = net.LocalHostIP
	}

	protocol := settings.ProtocolName
	listenFunc := transportListenerCache[protocol]
	if listenFunc == nil {
		return nil, newError(protocol, " listener not registered.").AtError()
	}
	listener, err := listenFunc(ctx, address, port, handler)
	if err != nil {
		return nil, newError("failed to listen on address: ", address, ":", port).Base(err)
	}
	return listener, nil
}

func ListenSystemTCP(ctx context.Context, addr *net.TCPAddr) (*net.TCPListener, error) {
	return effectiveTCPListener.Listen(ctx, addr)
}
