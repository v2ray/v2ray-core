package net

import (
	"net"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common/predicate"
)

var (
	// LocalHostIP is a constant value for localhost IP in IPv4.
	LocalHostIP = IPAddress([]byte{127, 0, 0, 1})

	// AnyIP is a constant value for any IP in IPv4.
	AnyIP = IPAddress([]byte{0, 0, 0, 0})

	// LocalHostDomain is a constant value for localhost domain.
	LocalHostDomain = DomainAddress("localhost")

	// LocalHostIPv6 is a constant value for localhost IP in IPv6.
	LocalHostIPv6 = IPAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
)

// AddressFamily is the type of address.
type AddressFamily int

const (
	// AddressFamilyIPv4 represents address as IPv4
	AddressFamilyIPv4 = AddressFamily(0)

	// AddressFamilyIPv6 represents address as IPv6
	AddressFamilyIPv6 = AddressFamily(1)

	// AddressFamilyDomain represents address as Domain
	AddressFamilyDomain = AddressFamily(2)
)

// Either returns true if current AddressFamily matches any of the AddressFamilys provided.
func (af AddressFamily) Either(fs ...AddressFamily) bool {
	for _, f := range fs {
		if af == f {
			return true
		}
	}
	return false
}

// IsIPv4 returns true if current AddressFamily is IPv4.
func (af AddressFamily) IsIPv4() bool {
	return af == AddressFamilyIPv4
}

// IsIPv6 returns true if current AddressFamily is IPv6.
func (af AddressFamily) IsIPv6() bool {
	return af == AddressFamilyIPv6
}

// IsDomain returns true if current AddressFamily is Domain.
func (af AddressFamily) IsDomain() bool {
	return af == AddressFamilyDomain
}

// Address represents a network address to be communicated with. It may be an IP address or domain
// address, not both. This interface doesn't resolve IP address for a given domain.
type Address interface {
	IP() net.IP     // IP of this Address
	Domain() string // Domain of this Address
	Family() AddressFamily

	String() string // String representation of this Address
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
		return addr
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
		return addr
	default:
		log.Trace(newError("invalid IP format: ", ip).AtError())
		return nil
	}
}

// DomainAddress creates an Address with given domain.
func DomainAddress(domain string) Address {
	return domainAddress(domain)
}

type ipv4Address [4]byte

func (a ipv4Address) IP() net.IP {
	return net.IP(a[:])
}

func (ipv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (ipv4Address) Family() AddressFamily {
	return AddressFamilyIPv4
}

func (a ipv4Address) String() string {
	return a.IP().String()
}

type ipv6Address [16]byte

func (a ipv6Address) IP() net.IP {
	return net.IP(a[:])
}

func (ipv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (ipv6Address) Family() AddressFamily {
	return AddressFamilyIPv6
}

func (a ipv6Address) String() string {
	return "[" + a.IP().String() + "]"
}

type domainAddress string

func (domainAddress) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (a domainAddress) Domain() string {
	return string(a)
}

func (domainAddress) Family() AddressFamily {
	return AddressFamilyDomain
}

func (a domainAddress) String() string {
	return a.Domain()
}

// AsAddress translates IPOrDomain to Address.
func (d *IPOrDomain) AsAddress() Address {
	if d == nil {
		return nil
	}
	switch addr := d.Address.(type) {
	case *IPOrDomain_Ip:
		return IPAddress(addr.Ip)
	case *IPOrDomain_Domain:
		return DomainAddress(addr.Domain)
	}
	panic("Common|Net: Invalid address.")
}

// NewIPOrDomain translates Address to IPOrDomain
func NewIPOrDomain(addr Address) *IPOrDomain {
	switch addr.Family() {
	case AddressFamilyDomain:
		return &IPOrDomain{
			Address: &IPOrDomain_Domain{
				Domain: addr.Domain(),
			},
		}
	case AddressFamilyIPv4, AddressFamilyIPv6:
		return &IPOrDomain{
			Address: &IPOrDomain_Ip{
				Ip: addr.IP(),
			},
		}
	default:
		panic("Unknown Address type.")
	}
}
