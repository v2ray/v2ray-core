// +build linux

package tcp

import (
	"syscall"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

const SO_ORIGINAL_DST = 80

func GetOriginalDestination(conn internet.Connection) net.Destination {
	tcpConn, ok := conn.(internet.SysFd)
	if !ok {
		log.Trace(errors.New("failed to get sys fd").Path("Transport", "Internet", "TCP"))
		return net.Destination{}
	}
	fd, err := tcpConn.SysFd()
	if err != nil {
		log.Trace(errors.New("failed to get original destination").Base(err).Path("Transport", "Internet", "TCP"))
		return net.Destination{}
	}

	addr, err := syscall.GetsockoptIPv6Mreq(fd, syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if err != nil {
		log.Trace(errors.New("failed to call getsockopt").Base(err).Path("Transport", "Internet", "TCP"))
		return net.Destination{}
	}
	ip := net.IPAddress(addr.Multiaddr[4:8])
	port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
	return net.TCPDestination(ip, net.Port(port))
}
