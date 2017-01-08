package internet

import (
	"net"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

var (
	ErrUnsupportedStreamType = errors.New("Unsupported stream type.")
)

type DialerOptions struct {
	Stream *StreamConfig
	Proxy  *ProxyConfig
}

type Dialer func(src v2net.Address, dest v2net.Destination, options DialerOptions) (Connection, error)

var (
	networkDialerCache = make(map[v2net.Network]Dialer)

	ProxyDialer Dialer
)

func RegisterNetworkDialer(network v2net.Network, dialer Dialer) error {
	if _, found := networkDialerCache[network]; found {
		return errors.New("Internet|Dialer: ", network, " dialer already registered.")
	}
	networkDialerCache[network] = dialer
	return nil
}

func Dial(src v2net.Address, dest v2net.Destination, options DialerOptions) (Connection, error) {
	if options.Proxy.HasTag() && ProxyDialer != nil {
		log.Info("Internet: Proxying outbound connection through: ", options.Proxy.Tag)
		return ProxyDialer(src, dest, options)
	}

	if dest.Network == v2net.Network_TCP {
		dialer := networkDialerCache[options.Stream.Network]
		if dialer == nil {
			return nil, errors.New("Internet|Dialer: ", options.Stream.Network, " dialer not registered.")
		}
		return dialer(src, dest, options)
	}

	udpDialer := networkDialerCache[v2net.Network_UDP]
	if udpDialer == nil {
		return nil, errors.New("Internet|Dialer: UDP dialer not registered.")
	}
	return udpDialer(src, dest, options)
}

// DialSystem calls system dialer to create a network connection.
func DialSystem(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	return effectiveSystemDialer.Dial(src, dest)
}
