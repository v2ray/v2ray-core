package shadowsocks_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	. "github.com/v2ray/v2ray-core/proxy/shadowsocks"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestNormalRequestParsing(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(1, 127, 0, 0, 1, 0, 80)

	request, err := ReadRequest(buffer)
	assert.Error(err).IsNil()
	netassert.Address(request.Address).Equals(v2net.IPAddress([]byte{127, 0, 0, 1}))
	netassert.Port(request.Port).Equals(v2net.Port(80))
	assert.Bool(request.OTA).IsFalse()
}

func TestOTARequest(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(0x13, 13, 119, 119, 119, 46, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 0)

	request, err := ReadRequest(buffer)
	assert.Error(err).IsNil()
	netassert.Address(request.Address).Equals(v2net.DomainAddress("www.v2ray.com"))
	assert.Bool(request.OTA).IsTrue()
}
