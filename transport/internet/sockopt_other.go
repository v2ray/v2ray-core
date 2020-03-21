// +build js dragonfly netbsd openbsd solaris

package internet

func applyOutboundSocketOptions(network string, address string, fd uintptr, config *SocketConfig) error {
	return nil
}

func applyInboundSocketOptions(network string, fd uintptr, config *SocketConfig) error {
	return nil
}

func bindAddr(fd uintptr, ip []byte, port uint32) error {
	return nil
}
