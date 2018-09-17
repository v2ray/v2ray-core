// +build js dragonfly freebsd netbsd openbsd

package internet

import "v2ray.com/core/common/net"

func applyOutboundSocketOptions(network string, address string, fd uintptr, config *SocketConfig) error {
	return nil
}

func applyInboundSocketOptions(network string, fd uintptr, config *SocketConfig) error {
	return nil
}

func bindAddr(fd uintptr, address net.Address, port net.Port) error {
	return nil
}
