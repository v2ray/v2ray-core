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

	IsIPv4() bool   // True if this Address is an IPv4 address
	IsIPv6() bool   // True if this Address is an IPv6 address
	IsDomain() bool // True if this Address is an domain address

	String() string // String representation of this Address
	Equals(Address) bool
}

func ParseAddress(addr string) Address {
	ip := net.ParseIP(addr)
	if ip != nil {
		return IPAddress(ip)
	}
	return DomainAddress(addr)
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
func IPAddress(ip []byte) Address {
	switch len(ip) {
	case net.IPv4len:
		var addr IPv4Address = [4]byte{ip[0], ip[1], ip[2], ip[3]}
		return &addr
	case net.IPv6len:
		if allZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff {
			return IPAddress(ip[12:16])
		}
		var addr IPv6Address = [16]byte{
			ip[0], ip[1], ip[2], ip[3],
			ip[4], ip[5], ip[6], ip[7],
			ip[8], ip[9], ip[10], ip[11],
			ip[12], ip[13], ip[14], ip[15],
		}
		return &addr
	default:
		log.Error("Invalid IP format: ", ip)
		return nil
	}
}

// DomainAddress creates an Address with given domain and port.
func DomainAddress(domain string) Address {
	var addr DomainAddressImpl = DomainAddressImpl(domain)
	return &addr
}

type IPv4Address [4]byte

func (addr *IPv4Address) IP() net.IP {
	return net.IP(addr[:])
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
	return this.IP().String()
}

func (this *IPv4Address) Equals(another Address) bool {
	anotherIPv4, ok := another.(*IPv4Address)
	if !ok {
		return false
	}
	return this[0] == anotherIPv4[0] &&
		this[1] == anotherIPv4[1] &&
		this[2] == anotherIPv4[2] &&
		this[3] == anotherIPv4[3]
}

type IPv6Address [16]byte

func (addr *IPv6Address) IP() net.IP {
	return net.IP(addr[:])
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
	return "[" + this.IP().String() + "]"
}

func (this *IPv6Address) Equals(another Address) bool {
	anotherIPv6, ok := another.(*IPv6Address)
	if !ok {
		return false
	}
	for idx, v := range *this {
		if anotherIPv6[idx] != v {
			return false
		}
	}
	return true
}

type DomainAddressImpl string

func (addr *DomainAddressImpl) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (addr *DomainAddressImpl) Domain() string {
	return string(*addr)
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
	return this.Domain()
}

func (this *DomainAddressImpl) Equals(another Address) bool {
	anotherDomain, ok := another.(*DomainAddressImpl)
	if !ok {
		return false
	}
	return this.Domain() == anotherDomain.Domain()
}
