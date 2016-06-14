package internet

import (
	"errors"
	"net"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
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
	if dest.IsTCP() {
		switch {
		case settings.IsCapableOf(StreamConnectionTypeKCP):
			return KCPDialer(src, dest)
		case settings.IsCapableOf(StreamConnectionTypeTCP):
			return TCPDialer(src, dest)
		case settings.IsCapableOf(StreamConnectionTypeRawTCP):
			return RawTCPDialer(src, dest)
		}
		return nil, ErrUnsupportedStreamType
	}

	return UDPDialer(src, dest)
}

func DialToDest(src v2net.Address, dest v2net.Destination) (net.Conn, error) {
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
