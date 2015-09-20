package net

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestIPv4Address(t *testing.T) {
	assert := unit.Assert(t)

	ip := []byte{byte(1), byte(2), byte(3), byte(4)}
	port := uint16(80)
	addr := IPAddress(ip, port)

	assert.Bool(addr.IsIPv4()).IsTrue()
	assert.Bytes(addr.IP()).Equals(ip)
	assert.Uint16(addr.Port()).Equals(port)
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
	port := uint16(443)
	addr := IPAddress(ip, port)

	assert.Bool(addr.IsIPv6()).IsTrue()
	assert.Bytes(addr.IP()).Equals(ip)
	assert.Uint16(addr.Port()).Equals(port)
	assert.String(addr.String()).Equals("[102:304:102:304:102:304:102:304]:443")
}

func TestDomainAddress(t *testing.T) {
	assert := unit.Assert(t)

	domain := "v2ray.com"
	port := uint16(443)
	addr := DomainAddress(domain, port)

	assert.Bool(addr.IsDomain()).IsTrue()
	assert.String(addr.Domain()).Equals(domain)
	assert.Uint16(addr.Port()).Equals(port)
	assert.String(addr.String()).Equals("v2ray.com:443")
}
