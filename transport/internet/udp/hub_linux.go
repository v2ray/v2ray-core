// +build linux

package udp

import (
	"syscall"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func SetOriginalDestOptions(fd int) error {
	if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
		return err
	}
	if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_RECVORIGDSTADDR, 1); err != nil {
		return err
	}
	return nil
}

func RetrieveOriginalDest(oob []byte) v2net.Destination {
	msgs, err := syscall.ParseSocketControlMessage(oob)
	if err != nil {
		return nil
	}
	for _, msg := range msgs {
		if msg.Header.Level == syscall.SOL_IP && msg.Header.Type == syscall.IP_ORIGDSTADDR {
			ip := v2net.IPAddress(msg.Data[4:8])
			port := v2net.PortFromBytes(msg.Data[2:4])
			return v2net.UDPDestination(ip, port)
		}
	}
	return nil
}
