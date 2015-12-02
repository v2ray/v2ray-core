package net

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestIPv4Address(t *testing.T) {
	assert := unit.Assert(t)

	ip := []byte{byte(1), byte(2), byte(3), byte(4)}
	port := NewPort(80)
	addr := IPAddress(ip, port.Value())

	assert.Bool(addr.IsIPv4()).IsTrue()
	assert.Bool(addr.IsIPv6()).IsFalse()
	assert.Bool(addr.IsDomain()).IsFalse()
	assert.Bytes(addr.IP()).Equals(ip)
	assert.Uint16(addr.Port().Value()).Equals(port.Value())
	assert.String(addr.String()).Equals("1.2.3.4:80")
}

func TestIPv6Address(t *testing.T) {
	assert := unit.Assert(t)

	ip := []byte{
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
	}
	port := NewPort(443)
	addr := IPAddress(ip, port.Value())

	assert.Bool(addr.IsIPv6()).IsTrue()
	assert.Bool(addr.IsIPv4()).IsFalse()
	assert.Bool(addr.IsDomain()).IsFalse()
	assert.Bytes(addr.IP()).Equals(ip)
	assert.Uint16(addr.Port().Value()).Equals(port.Value())
	assert.String(addr.String()).Equals("[102:304:102:304:102:304:102:304]:443")
}

func TestDomainAddress(t *testing.T) {
	assert := unit.Assert(t)

	domain := "v2ray.com"
	port := NewPort(443)
	addr := DomainAddress(domain, port.Value())

	assert.Bool(addr.IsDomain()).IsTrue()
	assert.Bool(addr.IsIPv4()).IsFalse()
	assert.Bool(addr.IsIPv6()).IsFalse()
	assert.String(addr.Domain()).Equals(domain)
	assert.Uint16(addr.Port().Value()).Equals(port.Value())
	assert.String(addr.String()).Equals("v2ray.com:443")
}

func TestNetIPv4Address(t *testing.T) {
	assert := unit.Assert(t)

	ip := net.IPv4(1, 2, 3, 4)
	port := NewPort(80)
	addr := IPAddress(ip, port.Value())
	assert.Bool(addr.IsIPv4()).IsTrue()
	assert.String(addr.String()).Equals("1.2.3.4:80")
}
