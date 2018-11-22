package net_test

import (
	"net"
	"testing"

	. "v2ray.com/core/common/net"
	. "v2ray.com/core/common/net/testing"
	. "v2ray.com/ext/assert"
)

func TestIPv4Address(t *testing.T) {
	assert := With(t)

	ip := []byte{byte(1), byte(2), byte(3), byte(4)}
	addr := IPAddress(ip)

	assert(addr, IsIPv4)
	assert(addr, Not(IsIPv6))
	assert(addr, Not(IsDomain))
	assert([]byte(addr.IP()), Equals, ip)
	assert(addr.String(), Equals, "1.2.3.4")
}

func TestIPv6Address(t *testing.T) {
	assert := With(t)

	ip := []byte{
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
	}
	addr := IPAddress(ip)

	assert(addr, IsIPv6)
	assert(addr, Not(IsIPv4))
	assert(addr, Not(IsDomain))
	assert(addr.IP(), Equals, net.IP(ip))
	assert(addr.String(), Equals, "[102:304:102:304:102:304:102:304]")
}

func TestIPv4Asv6(t *testing.T) {
	assert := With(t)
	ip := []byte{
		byte(0), byte(0), byte(0), byte(0),
		byte(0), byte(0), byte(0), byte(0),
		byte(0), byte(0), byte(255), byte(255),
		byte(1), byte(2), byte(3), byte(4),
	}
	addr := IPAddress(ip)
	assert(addr.String(), Equals, "1.2.3.4")
}

func TestDomainAddress(t *testing.T) {
	assert := With(t)

	domain := "v2ray.com"
	addr := DomainAddress(domain)

	assert(addr, IsDomain)
	assert(addr, Not(IsIPv6))
	assert(addr, Not(IsIPv4))
	assert(addr.Domain(), Equals, domain)
	assert(addr.String(), Equals, "v2ray.com")
}

func TestNetIPv4Address(t *testing.T) {
	assert := With(t)

	ip := net.IPv4(1, 2, 3, 4)
	addr := IPAddress(ip)
	assert(addr, IsIPv4)
	assert(addr.String(), Equals, "1.2.3.4")
}

func TestParseIPv6Address(t *testing.T) {
	assert := With(t)

	ip := ParseAddress("[2001:4860:0:2001::68]")
	assert(ip, IsIPv6)
	assert(ip.String(), Equals, "[2001:4860:0:2001::68]")

	ip = ParseAddress("[::ffff:123.151.71.143]")
	assert(ip, IsIPv4)
	assert(ip.String(), Equals, "123.151.71.143")
}

func TestInvalidAddressConvertion(t *testing.T) {
	assert := With(t)

	assert(func() { ParseAddress("8.8.8.8").Domain() }, Panics)
	assert(func() { ParseAddress("2001:4860:0:2001::68").Domain() }, Panics)
	assert(func() { ParseAddress("v2ray.com").IP() }, Panics)
}

func TestIPOrDomain(t *testing.T) {
	assert := With(t)

	assert(NewIPOrDomain(ParseAddress("v2ray.com")).AsAddress(), Equals, ParseAddress("v2ray.com"))
	assert(NewIPOrDomain(ParseAddress("8.8.8.8")).AsAddress(), Equals, ParseAddress("8.8.8.8"))
	assert(NewIPOrDomain(ParseAddress("2001:4860:0:2001::68")).AsAddress(), Equals, ParseAddress("2001:4860:0:2001::68"))
}

func BenchmarkParseAddressIPv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addr := ParseAddress("8.8.8.8")
		if addr.Family() != AddressFamilyIPv4 {
			panic("not ipv4")
		}
	}
}

func BenchmarkParseAddressIPv6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addr := ParseAddress("2001:4860:0:2001::68")
		if addr.Family() != AddressFamilyIPv6 {
			panic("not ipv6")
		}
	}
}

func BenchmarkParseAddressDomain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addr := ParseAddress("v2ray.com")
		if addr.Family() != AddressFamilyDomain {
			panic("not domain")
		}
	}
}
