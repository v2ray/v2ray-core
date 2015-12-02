package unit

import (
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestIPv4Address(t *testing.T) {
	v2testing.Current(t)

	ip := []byte{byte(1), byte(2), byte(3), byte(4)}
	port := v2net.NewPort(80)
	addr := v2net.IPAddress(ip, port.Value())

	v2netassert.Address(addr).IsIPv4()
	v2netassert.Address(addr).IsNotIPv6()
	v2netassert.Address(addr).IsNotDomain()
	assert.Bytes(addr.IP()).Equals(ip)
	v2netassert.Port(addr.Port()).Equals(port)
	assert.String(addr).Equals("1.2.3.4:80")
}

func TestIPv6Address(t *testing.T) {
	v2testing.Current(t)

	ip := []byte{
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
		byte(1), byte(2), byte(3), byte(4),
	}
	port := v2net.NewPort(443)
	addr := v2net.IPAddress(ip, port.Value())

	v2netassert.Address(addr).IsIPv6()
	v2netassert.Address(addr).IsNotIPv4()
	v2netassert.Address(addr).IsNotDomain()
	assert.Bytes(addr.IP()).Equals(ip)
	assert.Uint16(addr.Port().Value()).Equals(port.Value())
	assert.String(addr).Equals("[102:304:102:304:102:304:102:304]:443")
}

func TestDomainAddress(t *testing.T) {
	v2testing.Current(t)

	domain := "v2ray.com"
	port := v2net.NewPort(443)
	addr := v2net.DomainAddress(domain, port.Value())

	v2netassert.Address(addr).IsDomain()
	v2netassert.Address(addr).IsNotIPv6()
	v2netassert.Address(addr).IsNotIPv4()
	assert.StringLiteral(addr.Domain()).Equals(domain)
	assert.Uint16(addr.Port().Value()).Equals(port.Value())
	assert.String(addr).Equals("v2ray.com:443")
}

func TestNetIPv4Address(t *testing.T) {
	v2testing.Current(t)

	ip := net.IPv4(1, 2, 3, 4)
	port := v2net.NewPort(80)
	addr := v2net.IPAddress(ip, port.Value())
	v2netassert.Address(addr).IsIPv4()
	assert.String(addr).Equals("1.2.3.4:80")
}
