// +build linux

package dokodemo

import (
	"syscall"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/hub"
)

const SO_ORIGINAL_DST = 80

func GetOriginalDestination(conn *hub.Connection) v2net.Destination {
	fd, err := conn.SysFd()
	if err != nil {
		log.Info("Dokodemo: Failed to get original destination: ", err)
		return nil
	}

	addr, err := syscall.GetsockoptIPv6Mreq(fd, syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if err != nil {
		log.Info("Dokodemo: Failed to call getsockopt: ", err)
		return nil
	}
	ip := v2net.IPAddress(addr.Multiaddr[4:8])
	port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])
	return v2net.TCPDestination(ip, v2net.Port(port))
}
