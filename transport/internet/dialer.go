package internet

import (
	"context"

	"v2ray.com/core/common/session"

	"v2ray.com/core/common/net"
)

type Dialer func(ctx context.Context, dest net.Destination) (Connection, error)

var (
	transportDialerCache = make(map[string]Dialer)
)

func RegisterTransportDialer(protocol string, dialer Dialer) error {
	if _, found := transportDialerCache[protocol]; found {
		return newError(protocol, " dialer already registered").AtError()
	}
	transportDialerCache[protocol] = dialer
	return nil
}

// Dial dials a internet connection towards the given destination.
func Dial(ctx context.Context, dest net.Destination) (Connection, error) {
	if dest.Network == net.Network_TCP {
		streamSettings := StreamSettingsFromContext(ctx)
		if streamSettings == nil {
			s, err := ToMemoryStreamConfig(nil)
			if err != nil {
				return nil, newError("failed to create default stream settings").Base(err)
			}
			streamSettings = s
			ctx = ContextWithStreamSettings(ctx, streamSettings)
		}

		protocol := streamSettings.ProtocolName
		dialer := transportDialerCache[protocol]
		if dialer == nil {
			return nil, newError(protocol, " dialer not registered").AtError()
		}
		return dialer(ctx, dest)
	}

	if dest.Network == net.Network_UDP {
		udpDialer := transportDialerCache["udp"]
		if udpDialer == nil {
			return nil, newError("UDP dialer not registered").AtError()
		}
		return udpDialer(ctx, dest)
	}

	return nil, newError("unknown network ", dest.Network)
}

// DialSystem calls system dialer to create a network connection.
func DialSystem(ctx context.Context, dest net.Destination) (net.Conn, error) {
	var src net.Address
	if outbound := session.OutboundFromContext(ctx); outbound != nil {
		src = outbound.Gateway
	}
	return effectiveSystemDialer.Dial(ctx, src, dest)
}
