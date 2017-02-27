package internet

import (
	"context"
	"net"

	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
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

type ListenFunc func(ctx context.Context, address v2net.Address, port v2net.Port, conns chan<- Connection) (Listener, error)

type Listener interface {
	Close() error
	Addr() net.Addr
}

func ListenTCP(ctx context.Context, address v2net.Address, port v2net.Port, conns chan<- Connection) (Listener, error) {
	settings := StreamSettingsFromContext(ctx)
	protocol := settings.GetEffectiveProtocol()
	transportSettings, err := settings.GetEffectiveTransportSettings()
	if err != nil {
		return nil, err
	}
	ctx = ContextWithTransportSettings(ctx, transportSettings)
	if settings != nil && settings.HasSecuritySettings() {
		securitySettings, err := settings.GetEffectiveSecuritySettings()
		if err != nil {
			return nil, err
		}
		ctx = ContextWithSecuritySettings(ctx, securitySettings)
	}
	listenFunc := transportListenerCache[protocol]
	if listenFunc == nil {
		return nil, errors.New("Internet|TCPHub: ", protocol, " listener not registered.")
	}
	listener, err := listenFunc(ctx, address, port, conns)
	if err != nil {
		return nil, errors.Base(err).Message("Internet|TCPHub: Failed to listen on address: ", address, ":", port)
	}
	return listener, nil
}
