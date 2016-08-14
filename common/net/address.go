package net

import (
	"net"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/predicate"
)

var (
	LocalHostIP = IPAddress([]byte{127, 0, 0, 1})
	AnyIP       = IPAddress([]byte{0, 0, 0, 0})
)

type AddressFamily int

const (
	AddressFamilyIPv4   = AddressFamily(0)
	AddressFamilyIPv6   = AddressFamily(1)
	AddressFamilyDomain = AddressFamily(2)
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

// ParseAddress parses a string into an Address. The return value will be an IPAddress when
// the string is in the form of IPv4 or IPv6 address, or a DomainAddress otherwise.
func ParseAddress(addr string) Address {
	ip := net.ParseIP(addr)
	if ip != nil {
		return IPAddress(ip)
	}
	return DomainAddress(addr)
}

// IPAddress creates an Address with given IP.
func IPAddress(ip []byte) Address {
	switch len(ip) {
	case net.IPv4len:
		var addr ipv4Address = [4]byte{ip[0], ip[1], ip[2], ip[3]}
		return &addr
	case net.IPv6len:
		if predicate.BytesAll(ip[0:10], 0) && predicate.BytesAll(ip[10:12], 0xff) {
			return IPAddress(ip[12:16])
		}
		var addr ipv6Address = [16]byte{
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

// DomainAddress creates an Address with given domain.
func DomainAddress(domain string) Address {
	var addr domainAddress = domainAddress(domain)
	return &addr
}

type ipv4Address [4]byte

func (addr *ipv4Address) IP() net.IP {
	return net.IP(addr[:])
}

func (addr *ipv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (addr *ipv4Address) IsIPv4() bool {
	return true
}

func (addr *ipv4Address) IsIPv6() bool {
	return false
}

func (addr *ipv4Address) IsDomain() bool {
	return false
}

func (this *ipv4Address) String() string {
	return this.IP().String()
}

func (this *ipv4Address) Equals(another Address) bool {
	anotherIPv4, ok := another.(*ipv4Address)
	if !ok {
		return false
	}
	return this[0] == anotherIPv4[0] &&
		this[1] == anotherIPv4[1] &&
		this[2] == anotherIPv4[2] &&
		this[3] == anotherIPv4[3]
}

type ipv6Address [16]byte

func (addr *ipv6Address) IP() net.IP {
	return net.IP(addr[:])
}

func (addr *ipv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (addr *ipv6Address) IsIPv4() bool {
	return false
}

func (addr *ipv6Address) IsIPv6() bool {
	return true
}

func (addr *ipv6Address) IsDomain() bool {
	return false
}

func (this *ipv6Address) String() string {
	return "[" + this.IP().String() + "]"
}

func (this *ipv6Address) Equals(another Address) bool {
	anotherIPv6, ok := another.(*ipv6Address)
	if !ok {
		return false
	}
	return this[0] == anotherIPv6[0] &&
		this[1] == anotherIPv6[1] &&
		this[2] == anotherIPv6[2] &&
		this[3] == anotherIPv6[3] &&
		this[4] == anotherIPv6[4] &&
		this[5] == anotherIPv6[5] &&
		this[6] == anotherIPv6[6] &&
		this[7] == anotherIPv6[7] &&
		this[8] == anotherIPv6[8] &&
		this[9] == anotherIPv6[9] &&
		this[10] == anotherIPv6[10] &&
		this[11] == anotherIPv6[11] &&
		this[12] == anotherIPv6[12] &&
		this[13] == anotherIPv6[13] &&
		this[14] == anotherIPv6[14] &&
		this[15] == anotherIPv6[15]
}

type domainAddress string

func (addr *domainAddress) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (addr *domainAddress) Domain() string {
	return string(*addr)
}

func (addr *domainAddress) IsIPv4() bool {
	return false
}

func (addr *domainAddress) IsIPv6() bool {
	return false
}

func (addr *domainAddress) IsDomain() bool {
	return true
}

func (this *domainAddress) String() string {
	return this.Domain()
}

func (this *domainAddress) Equals(another Address) bool {
	anotherDomain, ok := another.(*domainAddress)
	if !ok {
		return false
	}
	return this.Domain() == anotherDomain.Domain()
}
