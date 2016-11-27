package net_test

import (
	"net"
	"testing"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestIPv4Address(t *testing.T) {
	assert := assert.On(t)

	ip := []byte{byte(1), byte(2), byte(3), byte(4)}
	addr := v2net.IPAddress(ip)

	assert.Address(addr).IsIPv4()
	assert.Address(addr).IsNotIPv6()
	assert.Address(addr).IsNotDomain()
	assert.Bytes(addr.IP()).Equals(ip)
	assert.Address(addr).EqualsString("1.2.3.4")
}

func TestIPv6Address(t *testing.T) {
	assert := assert.On(t)

	ip := []byte{
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
	}
	addr := v2net.IPAddress(ip)

	assert.Address(addr).IsIPv6()
	assert.Address(addr).IsNotIPv4()
	assert.Address(addr).IsNotDomain()
	assert.IP(addr.IP()).Equals(net.IP(ip))
	assert.Address(addr).EqualsString("[102:304:102:304:102:304:102:304]")
}

func TestIPv4Asv6(t *testing.T) {
	assert := assert.On(t)
	ip := []byte{
		byte(0), byte(0), byte(0), byte(0),
		byte(0), byte(0), byte(0), byte(0),
		byte(0), byte(0), byte(255), byte(255),
		byte(1), byte(2), byte(3), byte(4),
	}
	addr := v2net.IPAddress(ip)
	assert.Address(addr).EqualsString("1.2.3.4")
}

func TestDomainAddress(t *testing.T) {
	assert := assert.On(t)

	domain := "v2ray.com"
	addr := v2net.DomainAddress(domain)

	assert.Address(addr).IsDomain()
	assert.Address(addr).IsNotIPv6()
	assert.Address(addr).IsNotIPv4()
	assert.String(addr.Domain()).Equals(domain)
	assert.Address(addr).EqualsString("v2ray.com")
}

func TestNetIPv4Address(t *testing.T) {
	assert := assert.On(t)

	ip := net.IPv4(1, 2, 3, 4)
	addr := v2net.IPAddress(ip)
	assert.Address(addr).IsIPv4()
	assert.Address(addr).EqualsString("1.2.3.4")
}
