// +build !linux

package udp

import (
	"net"
)

func TransmitSocket(src net.Addr, dst net.Addr) (net.Conn, error) {
	return nil, newError("forging source address is not supported on non-Linux platform.").AtWarning()
}
