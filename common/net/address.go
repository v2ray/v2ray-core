package net

import (
	"net"
	"strconv"

	"github.com/v2ray/v2ray-core/common/log"
)

type Address interface {
	IP() net.IP
	Domain() string
	Port() uint16
	PortBytes() []byte

	IsIPv4() bool
	IsIPv6() bool
	IsDomain() bool

	String() string
}

func IPAddress(ip []byte, port uint16) Address {
	switch len(ip) {
	case net.IPv4len:
		return IPv4Address{
			PortAddress: PortAddress{port: port},
			ip:          [4]byte{ip[0], ip[1], ip[2], ip[3]},
		}
	case net.IPv6len:
		return IPv6Address{
			PortAddress: PortAddress{port: port},
			ip:          [16]byte{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15]},
		}
	default:
		panic(log.Error("Unknown IP format: %v", ip))
	}
}

func DomainAddress(domain string, port uint16) Address {
	return DomainAddressImpl{
		domain:      domain,
		PortAddress: PortAddress{port: port},
	}
}

type PortAddress struct {
	port uint16
}

func (addr PortAddress) Port() uint16 {
	return addr.port
}

func (addr PortAddress) PortBytes() []byte {
	return []byte{byte(addr.port >> 8), byte(addr.port)}
}

type IPv4Address struct {
	PortAddress
	ip [4]byte
}

func (addr IPv4Address) IP() net.IP {
	return net.IP(addr.ip[:])
}

func (addr IPv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (addr IPv4Address) IsIPv4() bool {
	return true
}

func (addr IPv4Address) IsIPv6() bool {
	return false
}

func (addr IPv4Address) IsDomain() bool {
	return false
}

func (addr IPv4Address) String() string {
	return addr.IP().String() + ":" + strconv.Itoa(int(addr.PortAddress.port))
}

type IPv6Address struct {
	PortAddress
	ip [16]byte
}

func (addr IPv6Address) IP() net.IP {
	return net.IP(addr.ip[:])
}

func (addr IPv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (addr IPv6Address) IsIPv4() bool {
	return false
}

func (addr IPv6Address) IsIPv6() bool {
	return true
}

func (addr IPv6Address) IsDomain() bool {
	return false
}

func (addr IPv6Address) String() string {
	return "[" + addr.IP().String() + "]:" + strconv.Itoa(int(addr.PortAddress.port))
}

type DomainAddressImpl struct {
	PortAddress
	domain string
}

func (addr DomainAddressImpl) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (addr DomainAddressImpl) Domain() string {
	return addr.domain
}

func (addr DomainAddressImpl) IsIPv4() bool {
	return false
}

func (addr DomainAddressImpl) IsIPv6() bool {
	return false
}

func (addr DomainAddressImpl) IsDomain() bool {
	return true
}

func (addr DomainAddressImpl) String() string {
	return addr.domain + ":" + strconv.Itoa(int(addr.PortAddress.port))
}
