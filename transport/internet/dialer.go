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
	TCPDialer    Dialer
	KCPDialer    Dialer
	RawTCPDialer Dialer
	UDPDialer    Dialer
	WSDialer     Dialer
	ProxyDialer  Dialer
)

func Dial(src v2net.Address, dest v2net.Destination, options DialerOptions) (Connection, error) {
	if options.Proxy.HasTag() && ProxyDialer != nil {
		log.Info("Internet: Proxying outbound connection through: ", options.Proxy.Tag)
		return ProxyDialer(src, dest, options)
	}

	var connection Connection
	var err error
	if dest.Network == v2net.Network_TCP {
		switch options.Stream.Network {
		case v2net.Network_TCP:
			connection, err = TCPDialer(src, dest, options)
		case v2net.Network_KCP:
			connection, err = KCPDialer(src, dest, options)
		case v2net.Network_WebSocket:
			connection, err = WSDialer(src, dest, options)

			// This check has to be the last one.
		case v2net.Network_RawTCP:
			connection, err = RawTCPDialer(src, dest, options)
		default:
			return nil, ErrUnsupportedStreamType
		}
		if err != nil {
			return nil, err
		}

		return connection, nil
	}

	return UDPDialer(src, dest, options)
}

func DialToDest(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	return effectiveSystemDialer.Dial(src, dest)
}
