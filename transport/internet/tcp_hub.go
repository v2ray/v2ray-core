package internet

import (
	"context"
	"time"

	"v2ray.com/core/common/net"
)

var (
	transportListenerCache = make(map[TransportProtocol]ListenFunc)
)

func RegisterTransportListener(protocol TransportProtocol, listener ListenFunc) error {
	if _, found := transportListenerCache[protocol]; found {
		return newError(protocol, " listener already registered.").AtError()
	}
	transportListenerCache[protocol] = listener
	return nil
}

type AddConnection func(context.Context, Connection) bool

type ListenFunc func(ctx context.Context, address net.Address, port net.Port, addConn AddConnection) (Listener, error)

type Listener interface {
	Close() error
	Addr() net.Addr
}

func ListenTCP(ctx context.Context, address net.Address, port net.Port, conns chan<- Connection) (Listener, error) {
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
	listener, err := listenFunc(ctx, address, port, func(ctx context.Context, conn Connection) bool {
		select {
		case <-ctx.Done():
			conn.Close()
			return false
		case conns <- conn:
			return true
		default:
			select {
			case <-ctx.Done():
				conn.Close()
				return false
			case conns <- conn:
				return true
			case <-time.After(time.Second * 5):
				conn.Close()
				return false
			}
		}
	})
	if err != nil {
		return nil, newError("failed to listen on address: ", address, ":", port).Base(err)
	}
	return listener, nil
}
