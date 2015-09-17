package net

import (
	"net"
	"strconv"
)

const (
	AddrTypeIP     = byte(0x01)
	AddrTypeDomain = byte(0x03)
)

type Address struct {
	Type   byte
	IP     net.IP
	Domain string
	Port   uint16
}

func IPAddress(ip []byte, port uint16) Address {
	ipCopy := make([]byte, 4)
	copy(ipCopy, ip)
	// TODO: check IP length
	return Address{
		Type:   AddrTypeIP,
		IP:     net.IP(ipCopy),
		Domain: "",
		Port:   port,
	}
}

func DomainAddress(domain string, port uint16) Address {
	return Address{
		Type:   AddrTypeDomain,
		IP:     nil,
		Domain: domain,
		Port:   port,
	}
}

func (addr Address) IsIPv4() bool {
	return addr.Type == AddrTypeIP && len(addr.IP) == net.IPv4len
}

func (addr Address) IsIPv6() bool {
	return addr.Type == AddrTypeIP && len(addr.IP) == net.IPv6len
}

func (addr Address) IsDomain() bool {
	return addr.Type == AddrTypeDomain
}

func (addr Address) String() string {
	var host string
	switch addr.Type {
	case AddrTypeIP:
		host = addr.IP.String()
		if len(addr.IP) == net.IPv6len {
			host = "[" + host + "]"
		}

	case AddrTypeDomain:
		host = addr.Domain
	default:
		panic("Unknown Address Type " + strconv.Itoa(int(addr.Type)))
	}
	return host + ":" + strconv.Itoa(int(addr.Port))
}
