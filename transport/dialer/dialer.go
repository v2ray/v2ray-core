package dialer

import (
	"errors"
	"net"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	ErrorInvalidHost = errors.New("Invalid Host.")
)

func Dial(dest v2net.Destination) (net.Conn, error) {
	if dest.Address().IsDomain() {
		dialer := &net.Dialer{
			Timeout:   time.Second * 60,
			DualStack: true,
		}
		network := "tcp"
		if dest.IsUDP() {
			network = "udp"
		}
		return dialer.Dial(network, dest.NetAddr())
	}

	ip := dest.Address().IP()
	if dest.IsTCP() {
		return net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   ip,
			Port: int(dest.Port()),
		})
	} else {
		return net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   ip,
			Port: int(dest.Port()),
		})
	}
}
