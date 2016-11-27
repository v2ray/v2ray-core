package net

import (
	"net"

	"v2ray.com/core/common/log"
	"v2ray.com/core/common/predicate"
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

func (v AddressFamily) Either(fs ...AddressFamily) bool {
	for _, f := range fs {
		if v == f {
			return true
		}
	}
	return false
}

func (v AddressFamily) IsIPv4() bool {
	return v == AddressFamilyIPv4
}

func (v AddressFamily) IsIPv6() bool {
	return v == AddressFamilyIPv6
}

func (v AddressFamily) IsDomain() bool {
	return v == AddressFamilyDomain
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
		log.Error("Invalid IP format: ", ip)
		return nil
	}
}

// DomainAddress creates an Address with given domain.
func DomainAddress(domain string) Address {
	var addr domainAddress = domainAddress(domain)
	return addr
}

type ipv4Address [4]byte

func (addr ipv4Address) IP() net.IP {
	return net.IP(addr[:])
}

func (addr ipv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (addr ipv4Address) Family() AddressFamily {
	return AddressFamilyIPv4
}

func (v ipv4Address) String() string {
	return v.IP().String()
}

type ipv6Address [16]byte

func (addr ipv6Address) IP() net.IP {
	return net.IP(addr[:])
}

func (addr ipv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (v ipv6Address) Family() AddressFamily {
	return AddressFamilyIPv6
}

func (v ipv6Address) String() string {
	return "[" + v.IP().String() + "]"
}

type domainAddress string

func (addr domainAddress) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (addr domainAddress) Domain() string {
	return string(addr)
}

func (addr domainAddress) Family() AddressFamily {
	return AddressFamilyDomain
}

func (v domainAddress) String() string {
	return v.Domain()
}

func (v *IPOrDomain) AsAddress() Address {
	if v == nil {
		return nil
	}
	switch addr := v.Address.(type) {
	case *IPOrDomain_Ip:
		return IPAddress(addr.Ip)
	case *IPOrDomain_Domain:
		return DomainAddress(addr.Domain)
	}
	panic("Common|Net: Invalid address.")
}
