package internet

import (
	"crypto/tls"
	"errors"
	"net"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2tls "github.com/v2ray/v2ray-core/transport/internet/tls"
)

var (
	ErrUnsupportedStreamType = errors.New("Unsupported stream type.")
)

type Dialer func(src v2net.Address, dest v2net.Destination) (Connection, error)

var (
	TCPDialer    Dialer
	KCPDialer    Dialer
	RawTCPDialer Dialer
	UDPDialer    Dialer
)

func Dial(src v2net.Address, dest v2net.Destination, settings *StreamSettings) (Connection, error) {
	var connection Connection
	var err error
	if dest.IsTCP() {
		switch {
		case settings.IsCapableOf(StreamConnectionTypeTCP):
			connection, err = TCPDialer(src, dest)
		case settings.IsCapableOf(StreamConnectionTypeKCP):
			connection, err = KCPDialer(src, dest)
		case settings.IsCapableOf(StreamConnectionTypeRawTCP):
			connection, err = RawTCPDialer(src, dest)
		default:
			return nil, ErrUnsupportedStreamType
		}
		if err != nil {
			return nil, err
		}
		if settings.Security == StreamSecurityTypeNone {
			return connection, nil
		}

		config := settings.TLSSettings.GetTLSConfig()
		if dest.Address().IsDomain() {
			config.ServerName = dest.Address().Domain()
		}
		tlsConn := tls.Client(connection, config)
		return v2tls.NewConnection(tlsConn), nil
	}

	return UDPDialer(src, dest)
}

func DialToDest(src v2net.Address, dest v2net.Destination) (net.Conn, error) {

	if isDefaultDialerSubstituted() {
		dialer := v2AlternativeDialer
		return (*dialer).Dial(dest.Network().String(), dest.NetAddr())
	} else {
		dialer := &net.Dialer{
			Timeout:   time.Second * 60,
			DualStack: true,
		}
		if src != nil && src != v2net.AnyIP {
			var addr net.Addr
			if dest.IsTCP() {
				addr = &net.TCPAddr{
					IP:   src.IP(),
					Port: 0,
				}
			} else {
				addr = &net.UDPAddr{
					IP:   src.IP(),
					Port: 0,
				}
			}
			dialer.LocalAddr = addr
		}
		return dialer.Dial(dest.Network().String(), dest.NetAddr())
	}
}
