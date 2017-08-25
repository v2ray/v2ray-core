// +build linux

package tcp

import (
	"syscall"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

const SO_ORIGINAL_DST = 80

func GetOriginalDestination(conn internet.Connection) (v2net.Destination, error) {
	sysrawconn, f := conn.(syscall.Conn)
	if !f {
		return v2net.Destination{}, newError("unable to get syscall.Conn")
	}
	rawConn, err := sysrawconn.SyscallConn()
	if err != nil {
		return v2net.Destination{}, newError("failed to get sys fd").Base(err)
	}
	var dest v2net.Destination
	err := rawConn.Control(func(fd uintptr) {
		addr, err := syscall.GetsockoptIPv6Mreq(int(fd), syscall.IPPROTO_IP, SO_ORIGINAL_DST)
		if err != nil {
			return v2net.Destination{}, newError("failed to call getsockopt").Base(err)
		}
		ip := v2net.IPAddress(addr.Multiaddr[4:8])
		port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
		addr = v2net.TCPDestination(ip, v2net.Port(port))
	})
	if err != nil {
		return v2net.Destination{}, newError("failed to get original destination").Base(err)
	}
	return dest, nil
}
