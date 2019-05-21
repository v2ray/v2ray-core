package udp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
)

// PickPort returns an unused UDP port in the system. The port returned is highly likely to be unused, but not guaranteed.
func PickPort() net.Port {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.LocalHostIP.IP(),
		Port: 0,
	})
	common.Must(err)
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	return net.Port(addr.Port)
}
