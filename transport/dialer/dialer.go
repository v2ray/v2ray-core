package dialer

import (
	"errors"
	"net"

	"github.com/v2ray/v2ray-core/common/dice"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	ErrorInvalidHost = errors.New("Invalid Host.")
)

func Dial(dest v2net.Destination) (net.Conn, error) {
	var ip net.IP
	if dest.Address().IsIPv4() || dest.Address().IsIPv6() {
		ip = dest.Address().IP()
	} else {
		ips, err := net.LookupIP(dest.Address().Domain())
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, ErrorInvalidHost
		}
		ip = ips[dice.Roll(len(ips))]
	}
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
