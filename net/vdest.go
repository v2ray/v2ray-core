package core

import (
	"net"
	"strconv"
)

const (
	AddrTypeIP     = byte(0x01)
	AddrTypeDomain = byte(0x03)
)

type VAddress struct {
	Type   byte
	IP     net.IP
	Domain string
	Port   uint16
}

func IPAddress(ip []byte, port uint16) VAddress {
	// TODO: check IP length
	return VAddress{
		AddrTypeIP,
		net.IP(ip),
		"",
		port}
}

func DomainAddress(domain string, port uint16) VAddress {
	return VAddress{
		AddrTypeDomain,
		nil,
		domain,
		port}
}

func (addr VAddress) IsIPv4() bool {
	return addr.Type == AddrTypeIP && len(addr.IP) == net.IPv4len
}

func (addr VAddress) IsIPv6() bool {
	return addr.Type == AddrTypeIP && len(addr.IP) == net.IPv6len
}

func (addr VAddress) IsDomain() bool {
	return addr.Type == AddrTypeDomain
}

func (addr VAddress) String() string {
	var host string
	switch addr.Type {
	case AddrTypeIP:
		host = addr.IP.String()
	case AddrTypeDomain:
		host = addr.Domain
	default:
		panic("Unknown Address Type " + strconv.Itoa(int(addr.Type)))
	}
	return host + ":" + strconv.Itoa(int(addr.Port))
}
