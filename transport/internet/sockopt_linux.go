package internet

import (
	"syscall"

	"v2ray.com/core/common/net"
)

const (
	// For incoming connections.
	TCP_FASTOPEN = 23
	// For out-going connections.
	TCP_FASTOPEN_CONNECT = 30
)

func bindAddr(fd uintptr, address net.Address, port net.Port) error {
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return newError("failed to set resuse_addr").Base(err).AtWarning()
	}

	var sockaddr syscall.Sockaddr

	switch address.Family() {
	case net.AddressFamilyIPv4:
		a4 := &syscall.SockaddrInet4{
			Port: int(port),
		}
		copy(a4.Addr[:], address.IP())
		sockaddr = a4
	case net.AddressFamilyIPv6:
		a6 := &syscall.SockaddrInet6{
			Port: int(port),
		}
		copy(a6.Addr[:], address.IP())
		sockaddr = a6
	default:
		return newError("unsupported address family: ", address.Family())
	}

	return syscall.Bind(int(fd), sockaddr)
}

func applyOutboundSocketOptions(network string, address string, fd uintptr, config *SocketConfig) error {
	if config.Mark != 0 {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(config.Mark)); err != nil {
			return newError("failed to set SO_MARK").Base(err)
		}
	}

	if isTCPSocket(network) {
		switch config.Tfo {
		case SocketConfig_Enable:
			if err := syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, TCP_FASTOPEN_CONNECT, 1); err != nil {
				return newError("failed to set TCP_FASTOPEN_CONNECT=1").Base(err)
			}
		case SocketConfig_Disable:
			if err := syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, TCP_FASTOPEN_CONNECT, 0); err != nil {
				return newError("failed to set TCP_FASTOPEN_CONNECT=0").Base(err)
			}
		}
	}

	if config.Tproxy.IsEnabled() {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			return newError("failed to set IP_TRANSPARENT").Base(err)
		}
	}

	return nil
}

func applyInboundSocketOptions(network string, fd uintptr, config *SocketConfig) error {
	if isTCPSocket(network) {
		switch config.Tfo {
		case SocketConfig_Enable:
			if err := syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, TCP_FASTOPEN, 1); err != nil {
				return newError("failed to set TCP_FASTOPEN=1").Base(err)
			}
		case SocketConfig_Disable:
			if err := syscall.SetsockoptInt(int(fd), syscall.SOL_TCP, TCP_FASTOPEN, 0); err != nil {
				return newError("failed to set TCP_FASTOPEN=0").Base(err)
			}
		}
	}

	if config.Tproxy.IsEnabled() {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			return newError("failed to set IP_TRANSPARENT").Base(err)
		}
	}

	if config.ReceiveOriginalDestAddress && isUDPSocket(network) {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_IP, syscall.IP_RECVORIGDSTADDR, 1); err != nil {
			return err
		}
	}

	return nil
}
