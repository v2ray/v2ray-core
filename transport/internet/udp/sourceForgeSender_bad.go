// +build !linux

package udp

import (
	"errors"
	"net"
)

func TransmitionSocket(src net.Addr, dst net.Addr) (net.Conn, error) {
	return nil, errors.New("Using an Linux only functionality on an non-Linux OS.")
}
