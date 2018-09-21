// +build linux

package udp

import (
	"syscall"

	"v2ray.com/core/common/net"
)

func RetrieveOriginalDest(oob []byte) net.Destination {
	msgs, err := syscall.ParseSocketControlMessage(oob)
	if err != nil {
		return net.Destination{}
	}
	for _, msg := range msgs {
		if msg.Header.Level == syscall.SOL_IP && msg.Header.Type == syscall.IP_RECVORIGDSTADDR {
			ip := net.IPAddress(msg.Data[4:8])
			port := net.PortFromBytes(msg.Data[2:4])
			return net.UDPDestination(ip, port)
		} else if msg.Header.Level == syscall.SOL_IPV6 && msg.Header.Type == syscall.IP_RECVORIGDSTADDR {
			ip := net.IPAddress(msg.Data[8:24])
			port := net.PortFromBytes(msg.Data[2:4])
			return net.UDPDestination(ip, port)
		}
	}
	return net.Destination{}
}

func ReadUDPMsg(conn *net.UDPConn, payload []byte, oob []byte) (int, int, int, *net.UDPAddr, error) {
	return conn.ReadMsgUDP(payload, oob)
}
