// +build !linux

package udp

import (
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func SetOriginalDestOptions(fd int) error {
	return nil
}

func RetrieveOriginalDest(oob []byte) v2net.Destination {
	return nil
}

func ReadUDPMsg(conn *net.UDPConn, payload []byte, oob []byte) (int, int, int, *net.UDPAddr, error) {
	nBytes, addr, err := conn.ReadFromUDP(payload)
	return nBytes, 0, 0, addr, err
}
