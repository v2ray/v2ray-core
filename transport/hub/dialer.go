package hub

import (
	"errors"
	"net"
	"time"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport"
)

var (
	ErrorInvalidHost = errors.New("Invalid Host.")

	globalCache = NewConnectionCache()
)

func Dial(dest v2net.Destination) (*Connection, error) {
	destStr := dest.String()
	var conn net.Conn
	if transport.IsConnectionReusable() {
		conn = globalCache.Get(destStr)
	}
	if conn == nil {
		var err error
		log.Debug("Hub: Dialling new connection to ", dest)
		conn, err = DialWithoutCache(dest)
		if err != nil {
			return nil, err
		}
	} else {
		log.Debug("Hub: Reusing connection to ", dest)
	}
	return &Connection{
		dest:     destStr,
		conn:     conn,
		listener: globalCache,
	}, nil
}

func DialWithoutCache(dest v2net.Destination) (net.Conn, error) {
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
	}

	return net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ip,
		Port: int(dest.Port()),
	})
}
