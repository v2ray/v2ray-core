package tcp

import (
	"net"

	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

var (
	globalCache = NewConnectionCache()
)

func Dial(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	log.Info("Dailing TCP to ", dest)
	if src == nil {
		src = v2net.AnyIP
	}
	id := src.String() + "-" + dest.NetAddr()
	var conn net.Conn
	if dest.Network() == v2net.Network_TCP && effectiveConfig.ConnectionReuse {
		conn = globalCache.Get(id)
	}
	if conn == nil {
		var err error
		conn, err = internet.DialToDest(src, dest)
		if err != nil {
			return nil, err
		}
	}
	return NewConnection(id, conn, globalCache), nil
}

func DialRaw(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
	log.Info("Dailing Raw TCP to ", dest)
	conn, err := internet.DialToDest(src, dest)
	if err != nil {
		return nil, err
	}
	return &RawConnection{
		TCPConn: *conn.(*net.TCPConn),
	}, nil
}

func init() {
	internet.TCPDialer = Dial
	internet.RawTCPDialer = DialRaw
}
