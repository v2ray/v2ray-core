// +build linux

package udp

import (
	"net"
	"os"
	"syscall"
)

//Currently, Only IPv4 Forge is supported
func TransmitSocket(src net.Addr, dst net.Addr) (net.Conn, error) {
	var fd int
	var err error
	fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, newError("failed to create fd").Base(err).AtWarning()
	}
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return nil, newError("failed to set resuse_addr").Base(err).AtWarning()
	}

	err = syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1)
	if err != nil {
		return nil, newError("failed to set transparent").Base(err).AtWarning()
	}

	ip := src.(*net.UDPAddr).IP.To4()
	var ip2 [4]byte
	copy(ip2[:], ip)
	srcaddr := syscall.SockaddrInet4{}
	srcaddr.Addr = ip2
	srcaddr.Port = src.(*net.UDPAddr).Port
	err = syscall.Bind(fd, &srcaddr)
	if err != nil {
		return nil, newError("failed to bind source address").Base(err).AtWarning()
	}
	ipd := dst.(*net.UDPAddr).IP.To4()
	var ip2d [4]byte
	copy(ip2d[:], ipd)
	dstaddr := syscall.SockaddrInet4{}
	dstaddr.Addr = ip2d
	dstaddr.Port = dst.(*net.UDPAddr).Port
	err = syscall.Connect(fd, &dstaddr)
	if err != nil {
		return nil, newError("failed to connect to source address").Base(err).AtWarning()
	}
	fdf := os.NewFile(uintptr(fd), "/dev/udp/")
	c, err := net.FileConn(fdf)
	if err != nil {
		return nil, newError("failed to create file conn").Base(err).AtWarning()
	}
	return c, nil
}
