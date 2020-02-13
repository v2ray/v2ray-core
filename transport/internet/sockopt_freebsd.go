package internet

import (
	"encoding/binary"
	"net"
	"os"
	"syscall"
	"unsafe"
)

const (
	sysAF_INET       = 0x2
	sysAF_INET6      = 0x1c
	sysPF_INOUT      = 0x0
	sysPF_IN         = 0x1
	sysPF_OUT        = 0x2
	sysPF_FWD        = 0x3
	sysDIOCNATLOOK   = 0xc04c4417
	ianaProtocolIP   = 0x0
	ianaProtocolTCP  = 0x6
	ianaProtocolUDP  = 0x11
	ianaProtocolIPv6 = 0x29
)

type pfiocNatlook struct {
	Saddr     [16]byte /* pf_addr */
	Daddr     [16]byte /* pf_addr */
	Rsaddr    [16]byte /* pf_addr */
	Rdaddr    [16]byte /* pf_addr */
	Sport     uint16
	Dport     uint16
	Rsport    uint16
	Rdport    uint16
	Af        uint8
	Proto     uint8
	Direction uint8
	Pad_cgo_0 [1]byte
}

const (
	sizeofPfiocNatlook = 0x4c
	SO_REUSEPORT_LB    = 0x00010000
	IP_RECVORIGDSTADDR = 27
	TCP_FASTOPEN       = 0x401
	SO_REUSEADDR       = 0x4
)

func ioctl(s uintptr, ioc int, b []byte) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s, uintptr(ioc), uintptr(unsafe.Pointer(&b[0]))); errno != 0 {
		return error(errno)
	}
	return nil
}
func (nl *pfiocNatlook) rdPort() int {
	return int(binary.BigEndian.Uint16((*[2]byte)(unsafe.Pointer(&nl.Rdport))[:]))
}

func (nl *pfiocNatlook) setPort(remote, local int) {
	binary.BigEndian.PutUint16((*[2]byte)(unsafe.Pointer(&nl.Sport))[:], uint16(remote))
	binary.BigEndian.PutUint16((*[2]byte)(unsafe.Pointer(&nl.Dport))[:], uint16(local))
}
func OriginalDst(la, ra net.Addr) (net.IP, int, error) {
	f, err := os.Open("/dev/pf")
	if err != nil {
		return net.IP{}, -1, newError("failed to open device /dev/pf").Base(err)
	}
	defer f.Close()
	fd := f.Fd()
	b := make([]byte, sizeofPfiocNatlook)
	nl := (*pfiocNatlook)(unsafe.Pointer(&b[0]))
	var raIP, laIP net.IP
	var raPort, laPort int
	switch la.(type) {
	case *net.TCPAddr:
		raIP = ra.(*net.TCPAddr).IP
		laIP = la.(*net.TCPAddr).IP
		raPort = ra.(*net.TCPAddr).Port
		laPort = la.(*net.TCPAddr).Port
		nl.Proto = ianaProtocolTCP
	case *net.UDPAddr:
		raIP = ra.(*net.UDPAddr).IP
		laIP = la.(*net.UDPAddr).IP
		raPort = ra.(*net.UDPAddr).Port
		laPort = la.(*net.UDPAddr).Port
		nl.Proto = ianaProtocolUDP
	}
	if raIP.To4() != nil {
		if laIP.IsUnspecified() {
			laIP = net.ParseIP("127.0.0.1")
		}
		copy(nl.Saddr[:net.IPv4len], raIP.To4())
		copy(nl.Daddr[:net.IPv4len], laIP.To4())
		nl.Af = sysAF_INET
	}
	if raIP.To16() != nil && raIP.To4() == nil {
		if laIP.IsUnspecified() {
			laIP = net.ParseIP("::1")
		}
		copy(nl.Saddr[:], raIP)
		copy(nl.Daddr[:], laIP)
		nl.Af = sysAF_INET6
	}
	nl.setPort(raPort, laPort)
	ioc := uintptr(sysDIOCNATLOOK)
	for _, dir := range []byte{sysPF_OUT, sysPF_IN} {
		nl.Direction = dir
		err = ioctl(fd, int(ioc), b)
		if err == nil || err != syscall.ENOENT {
			break
		}
	}
	if err != nil {
		return net.IP{}, -1, os.NewSyscallError("ioctl", err)
	}

	odPort := nl.rdPort()
	var odIP net.IP
	switch nl.Af {
	case sysAF_INET:
		odIP = make(net.IP, net.IPv4len)
		copy(odIP, nl.Rdaddr[:net.IPv4len])
	case sysAF_INET6:
		odIP = make(net.IP, net.IPv6len)
		copy(odIP, nl.Rdaddr[:])
	}
	return odIP, odPort, nil
}

func applyOutboundSocketOptions(network string, address string, fd uintptr, config *SocketConfig) error {
	if config.Mark != 0 {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_USER_COOKIE, int(config.Mark)); err != nil {
			return newError("failed to set SO_USER_COOKIE").Base(err)
		}
	}

	if isTCPSocket(network) {
		switch config.Tfo {
		case SocketConfig_Enable:
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolTCP, TCP_FASTOPEN, 1); err != nil {
				return newError("failed to set TCP_FASTOPEN_CONNECT=1").Base(err)
			}
		case SocketConfig_Disable:
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolTCP, TCP_FASTOPEN, 0); err != nil {
				return newError("failed to set TCP_FASTOPEN_CONNECT=0").Base(err)
			}
		}
	}

	if config.Tproxy.IsEnabled() {
		ip, _, _ := net.SplitHostPort(address)
		if net.ParseIP(ip).To4() != nil {
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolIP, syscall.IP_BINDANY, 1); err != nil {
				return newError("failed to set outbound IP_BINDANY").Base(err)
			}
		} else {
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolIPv6, syscall.IPV6_BINDANY, 1); err != nil {
				return newError("failed to set outbound IPV6_BINDANY").Base(err)
			}
		}
	}
	return nil
}

func applyInboundSocketOptions(network string, fd uintptr, config *SocketConfig) error {
	if config.Mark != 0 {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_USER_COOKIE, int(config.Mark)); err != nil {
			return newError("failed to set SO_USER_COOKIE").Base(err)
		}
	}
	if isTCPSocket(network) {
		switch config.Tfo {
		case SocketConfig_Enable:
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolTCP, TCP_FASTOPEN, 1); err != nil {
				return newError("failed to set TCP_FASTOPEN=1").Base(err)
			}
		case SocketConfig_Disable:
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolTCP, TCP_FASTOPEN, 0); err != nil {
				return newError("failed to set TCP_FASTOPEN=0").Base(err)
			}
		}
	}

	if config.Tproxy.IsEnabled() {
		if err := syscall.SetsockoptInt(int(fd), ianaProtocolIPv6, syscall.IPV6_BINDANY, 1); err != nil {
			if err := syscall.SetsockoptInt(int(fd), ianaProtocolIP, syscall.IP_BINDANY, 1); err != nil {
				return newError("failed to set inbound IP_BINDANY").Base(err)
			}
		}
	}

	return nil
}

func bindAddr(fd uintptr, ip []byte, port uint32) error {
	if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, SO_REUSEADDR, 1); err != nil {
		return newError("failed to set resuse_addr").Base(err).AtWarning()
	}

	if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, SO_REUSEPORT_LB, 1); err != nil {
		return newError("failed to set resuse_port").Base(err).AtWarning()
	}

	var sockaddr syscall.Sockaddr

	switch len(ip) {
	case net.IPv4len:
		a4 := &syscall.SockaddrInet4{
			Port: int(port),
		}
		copy(a4.Addr[:], ip)
		sockaddr = a4
	case net.IPv6len:
		a6 := &syscall.SockaddrInet6{
			Port: int(port),
		}
		copy(a6.Addr[:], ip)
		sockaddr = a6
	default:
		return newError("unexpected length of ip")
	}

	return syscall.Bind(int(fd), sockaddr)
}
