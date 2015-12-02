package net

import (
	"net"

	"github.com/v2ray/v2ray-core/common/log"
)

// Address represents a network address to be communicated with. It may be an IP address or domain
// address, not both. This interface doesn't resolve IP address for a given domain.
type Address interface {
	IP() net.IP     // IP of this Address
	Domain() string // Domain of this Address
	Port() Port     // Port of this Address

	IsIPv4() bool   // True if this Address is an IPv4 address
	IsIPv6() bool   // True if this Address is an IPv6 address
	IsDomain() bool // True if this Address is an domain address

	String() string // String representation of this Address
}

func allZeros(data []byte) bool {
	for _, v := range data {
		if v != 0 {
			return false
		}
	}
	return true
}

// IPAddress creates an Address with given IP and port.
func IPAddress(ip []byte, port uint16) Address {
	switch len(ip) {
	case net.IPv4len:
		return &IPv4Address{
			port: Port(port),
			ip:   [4]byte{ip[0], ip[1], ip[2], ip[3]},
		}
	case net.IPv6len:
		if allZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff {
			return IPAddress(ip[12:16], port)
		}
		return &IPv6Address{
			port: Port(port),
			ip: [16]byte{
				ip[0], ip[1], ip[2], ip[3],
				ip[4], ip[5], ip[6], ip[7],
				ip[8], ip[9], ip[10], ip[11],
				ip[12], ip[13], ip[14], ip[15],
			},
		}
	default:
		log.Error("Invalid IP format: %v", ip)
		return nil
	}
}

// DomainAddress creates an Address with given domain and port.
func DomainAddress(domain string, port uint16) Address {
	return &DomainAddressImpl{
		domain: domain,
		port:   Port(port),
	}
}

type IPv4Address struct {
	port Port
	ip   [4]byte
}

func (addr *IPv4Address) IP() net.IP {
	return net.IP(addr.ip[:])
}

func (this *IPv4Address) Port() Port {
	return this.port
}

func (addr *IPv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (addr *IPv4Address) IsIPv4() bool {
	return true
}

func (addr *IPv4Address) IsIPv6() bool {
	return false
}

func (addr *IPv4Address) IsDomain() bool {
	return false
}

func (this *IPv4Address) String() string {
	return this.IP().String() + ":" + this.port.String()
}

type IPv6Address struct {
	port Port
	ip   [16]byte
}

func (addr *IPv6Address) IP() net.IP {
	return net.IP(addr.ip[:])
}

func (this *IPv6Address) Port() Port {
	return this.port
}

func (addr *IPv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (addr *IPv6Address) IsIPv4() bool {
	return false
}

func (addr *IPv6Address) IsIPv6() bool {
	return true
}

func (addr *IPv6Address) IsDomain() bool {
	return false
}

func (this *IPv6Address) String() string {
	return "[" + this.IP().String() + "]:" + this.port.String()
}

type DomainAddressImpl struct {
	port   Port
	domain string
}

func (addr *DomainAddressImpl) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (this *DomainAddressImpl) Port() Port {
	return this.port
}

func (addr *DomainAddressImpl) Domain() string {
	return addr.domain
}

func (addr *DomainAddressImpl) IsIPv4() bool {
	return false
}

func (addr *DomainAddressImpl) IsIPv6() bool {
	return false
}

func (addr *DomainAddressImpl) IsDomain() bool {
	return true
}

func (this *DomainAddressImpl) String() string {
	return this.domain + ":" + this.port.String()
}
