package scenarios

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

const (
	socks5Version = byte(0x05)
)

func socks5AuthMethodRequest(methods ...byte) []byte {
	request := []byte{socks5Version, byte(len(methods))}
	request = append(request, methods...)
	return request
}

func socks5Request(command byte, address v2net.Address) []byte {
	request := []byte{socks5Version, command, 0}
	switch {
	case address.IsIPv4():
		request = append(request, byte(0x01))
		request = append(request, address.IP()...)

	case address.IsIPv6():
		request = append(request, byte(0x04))
		request = append(request, address.IP()...)

	case address.IsDomain():
		request = append(request, byte(0x03), byte(len(address.Domain())))
		request = append(request, []byte(address.Domain())...)

	}
	request = append(request, address.PortBytes()...)
	return request
}
