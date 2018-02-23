package encoding

import (
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg encoding -path Proxy,VMess,Encoding

const (
	Version = byte(1)
)

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(0x02, net.AddressFamilyDomain),
	protocol.AddressFamilyByte(0x03, net.AddressFamilyIPv6),
	protocol.PortThenAddress(),
)
