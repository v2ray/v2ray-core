package shadowsocks_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	"github.com/v2ray/v2ray-core/proxy"
	. "github.com/v2ray/v2ray-core/proxy/shadowsocks"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport"
)

func TestNormalRequestParsing(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(1, 127, 0, 0, 1, 0, 80)

	request, err := ReadRequest(buffer, nil, false)
	assert.Error(err).IsNil()
	netassert.Address(request.Address).Equals(v2net.IPAddress([]byte{127, 0, 0, 1}))
	netassert.Port(request.Port).Equals(v2net.Port(80))
	assert.Bool(request.OTA).IsFalse()
}

func TestEmptyPayload(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	_, err := ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestSingleBytePayload(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear().AppendBytes(1)
	_, err := ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestWrongAddressType(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear().AppendBytes(5)
	_, err := ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestInsufficientAddressRequest(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear().AppendBytes(1, 1)
	_, err := ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)

	buffer = alloc.NewSmallBuffer().Clear().AppendBytes(4, 1)
	_, err = ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)

	buffer = alloc.NewSmallBuffer().Clear().AppendBytes(3, 255, 1)
	_, err = ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestInsufficientPortRequest(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear().AppendBytes(1, 1, 2, 3, 4, 5)
	_, err := ReadRequest(buffer, nil, false)
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestOTARequest(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(0x13, 13, 119, 119, 119, 46, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 0, 239, 115, 52, 212, 178, 172, 26, 6, 168, 0)

	auth := NewAuthenticator(HeaderKeyGenerator(
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5},
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5}))
	request, err := ReadRequest(buffer, auth, false)
	assert.Error(err).IsNil()
	netassert.Address(request.Address).Equals(v2net.DomainAddress("www.v2ray.com"))
	assert.Bool(request.OTA).IsTrue()
}

func TestInvalidOTARequest(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(0x13, 13, 119, 119, 119, 46, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 0, 239, 115, 52, 212, 178, 172, 26, 6, 168, 1)

	auth := NewAuthenticator(HeaderKeyGenerator(
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5},
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5}))
	_, err := ReadRequest(buffer, auth, false)
	assert.Error(err).Equals(proxy.ErrorInvalidAuthentication)
}

func TestUDPRequestParsing(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(1, 127, 0, 0, 1, 0, 80, 1, 2, 3, 4, 5, 6)

	request, err := ReadRequest(buffer, nil, true)
	assert.Error(err).IsNil()
	netassert.Address(request.Address).Equals(v2net.IPAddress([]byte{127, 0, 0, 1}))
	netassert.Port(request.Port).Equals(v2net.Port(80))
	assert.Bool(request.OTA).IsFalse()
	assert.Bytes(request.UDPPayload.Value).Equals([]byte{1, 2, 3, 4, 5, 6})
}

func TestUDPRequestWithOTA(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewSmallBuffer().Clear()
	buffer.AppendBytes(
		0x13, 13, 119, 119, 119, 46, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 0,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
		58, 32, 223, 30, 57, 199, 50, 139, 143, 101)

	auth := NewAuthenticator(HeaderKeyGenerator(
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5},
		[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5}))
	request, err := ReadRequest(buffer, auth, true)
	assert.Error(err).IsNil()
	netassert.Address(request.Address).Equals(v2net.DomainAddress("www.v2ray.com"))
	assert.Bool(request.OTA).IsTrue()
	assert.Bytes(request.UDPPayload.Value).Equals([]byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
}
