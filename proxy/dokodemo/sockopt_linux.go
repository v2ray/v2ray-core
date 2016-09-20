// +build linux

package dokodemo

import (
	"syscall"

	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

const SO_ORIGINAL_DST = 80

func GetOriginalDestination(conn internet.Connection) v2net.Destination {
	tcpConn, ok := conn.(internet.SysFd)
	if !ok {
		log.Info("Dokodemo: Failed to get sys fd.")
		return v2net.Destination{}
	}
	fd, err := tcpConn.SysFd()
	if err != nil {
		log.Info("Dokodemo: Failed to get original destination: ", err)
		return v2net.Destination{}
	}

	addr, err := syscall.GetsockoptIPv6Mreq(fd, syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if err != nil {
		log.Info("Dokodemo: Failed to call getsockopt: ", err)
		return v2net.Destination{}
	}
	ip := v2net.IPAddress(addr.Multiaddr[4:8])
	port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
	return v2net.TCPDestination(ip, v2net.Port(port))
}
