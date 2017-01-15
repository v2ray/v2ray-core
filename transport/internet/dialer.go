package internet

import (
	"net"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

type DialerOptions struct {
	Stream *StreamConfig
	Proxy  *ProxyConfig
}

type Dialer func(src v2net.Address, dest v2net.Destination, options DialerOptions) (Connection, error)

var (
	transportDialerCache = make(map[TransportProtocol]Dialer)

	ProxyDialer Dialer
)

func RegisterTransportDialer(protocol TransportProtocol, dialer Dialer) error {
	if _, found := transportDialerCache[protocol]; found {
		return errors.New("Internet|Dialer: ", protocol, " dialer already registered.")
	}
	transportDialerCache[protocol] = dialer
	return nil
}

func Dial(src v2net.Address, dest v2net.Destination, options DialerOptions) (Connection, error) {
	if options.Proxy.HasTag() && ProxyDialer != nil {
		log.Info("Internet: Proxying outbound connection through: ", options.Proxy.Tag)
		return ProxyDialer(src, dest, options)
	}

	if dest.Network == v2net.Network_TCP {
		protocol := options.Stream.GetEffectiveProtocol()
		dialer := transportDialerCache[protocol]
		if dialer == nil {
			return nil, errors.New("Internet|Dialer: ", options.Stream.Protocol, " dialer not registered.")
		}
		return dialer(src, dest, options)
	}

	udpDialer := transportDialerCache[TransportProtocol_UDP]
	if udpDialer == nil {
		return nil, errors.New("Internet|Dialer: UDP dialer not registered.")
	}
	return udpDialer(src, dest, options)
}

// DialSystem calls system dialer to create a network connection.
func DialSystem(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	return effectiveSystemDialer.Dial(src, dest)
}
