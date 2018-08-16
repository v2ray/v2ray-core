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
		return nil, newError(protocol, " listener not registered.").AtError()
	}
	listener, err := listenFunc(ctx, address, port, handler)
	if err != nil {
		return nil, newError("failed to listen on address: ", address, ":", port).Base(err)
	}
	return listener, nil
}
