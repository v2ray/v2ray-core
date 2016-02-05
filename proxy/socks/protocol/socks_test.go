package protocol

import (
	"bytes"
	"io"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	"github.com/v2ray/v2ray-core/proxy"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport"
)

func TestHasAuthenticationMethod(t *testing.T) {
	v2testing.Current(t)

	request := Socks5AuthenticationRequest{
		version:     socksVersion,
		nMethods:    byte(0x02),
		authMethods: [256]byte{0x01, 0x02},
	}

	assert.Bool(request.HasAuthMethod(byte(0x01))).IsTrue()

	request.authMethods[0] = byte(0x03)
	assert.Bool(request.HasAuthMethod(byte(0x01))).IsFalse()
}

func TestAuthenticationRequestRead(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewBuffer().Clear().AppendBytes(
		0x05, // version
		0x01, // nMethods
		0x02, // methods
	)
	request, _, err := ReadAuthentication(buffer)
	assert.Error(err).IsNil()
	assert.Byte(request.version).Named("Version").Equals(0x05)
	assert.Byte(request.nMethods).Named("#Methods").Equals(0x01)
	assert.Byte(request.authMethods[0]).Named("Auth Method").Equals(0x02)
}

func TestAuthenticationResponseWrite(t *testing.T) {
	v2testing.Current(t)

	response := NewAuthenticationResponse(byte(0x05))

	buffer := bytes.NewBuffer(make([]byte, 0, 10))
	WriteAuthentication(buffer, response)
	assert.Bytes(buffer.Bytes()).Equals([]byte{socksVersion, byte(0x05)})
}

func TestRequestRead(t *testing.T) {
	v2testing.Current(t)

	rawRequest := []byte{
		0x05,                   // version
		0x01,                   // cmd connect
		0x00,                   // reserved
		0x01,                   // ipv4 type
		0x72, 0x72, 0x72, 0x72, // 114.114.114.114
		0x00, 0x35, // port 53
	}
	request, err := ReadRequest(bytes.NewReader(rawRequest))
	assert.Error(err).IsNil()
	assert.Byte(request.Version).Named("Version").Equals(0x05)
	assert.Byte(request.Command).Named("Command").Equals(0x01)
	assert.Byte(request.AddrType).Named("Address Type").Equals(0x01)
	assert.Bytes(request.IPv4[:]).Named("IPv4").Equals([]byte{0x72, 0x72, 0x72, 0x72})
	v2netassert.Port(request.Port).Named("Port").Equals(v2net.Port(53))
}

func TestResponseWrite(t *testing.T) {
	v2testing.Current(t)

	response := Socks5Response{
		socksVersion,
		ErrorSuccess,
		AddrTypeIPv4,
		[4]byte{0x72, 0x72, 0x72, 0x72},
		"",
		[16]byte{},
		v2net.Port(53),
	}
	buffer := alloc.NewSmallBuffer().Clear()
	defer buffer.Release()

	response.Write(buffer)
	expectedBytes := []byte{
		socksVersion,
		ErrorSuccess,
		byte(0x00),
		AddrTypeIPv4,
		0x72, 0x72, 0x72, 0x72,
		byte(0x00), byte(0x035),
	}
	assert.Bytes(buffer.Value).Named("raw response").Equals(expectedBytes)
}

func TestSetIPv6(t *testing.T) {
	v2testing.Current(t)

	response := NewSocks5Response()
	response.SetIPv6([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})

	buffer := alloc.NewSmallBuffer().Clear()
	defer buffer.Release()
	response.Write(buffer)
	assert.Bytes(buffer.Value).Equals([]byte{
		socksVersion, 0, 0, AddrTypeIPv6, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 0})
}

func TestSetDomain(t *testing.T) {
	v2testing.Current(t)

	response := NewSocks5Response()
	response.SetDomain("v2ray.com")

	buffer := alloc.NewSmallBuffer().Clear()
	defer buffer.Release()
	response.Write(buffer)
	assert.Bytes(buffer.Value).Equals([]byte{
		socksVersion, 0, 0, AddrTypeDomain, 9, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 0})
}

func TestEmptyAuthRequest(t *testing.T) {
	v2testing.Current(t)

	_, _, err := ReadAuthentication(alloc.NewBuffer().Clear())
	assert.Error(err).Equals(io.EOF)
}

func TestSingleByteAuthRequest(t *testing.T) {
	v2testing.Current(t)

	_, _, err := ReadAuthentication(bytes.NewReader(make([]byte, 1)))
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestZeroAuthenticationMethod(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewBuffer().Clear().AppendBytes(5, 0)
	_, _, err := ReadAuthentication(buffer)
	assert.Error(err).Equals(proxy.ErrorInvalidAuthentication)
}
func TestWrongProtocolVersion(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewBuffer().Clear().AppendBytes(6, 1, 0)
	_, _, err := ReadAuthentication(buffer)
	assert.Error(err).Equals(proxy.ErrorInvalidProtocolVersion)
}

func TestEmptyRequest(t *testing.T) {
	v2testing.Current(t)

	_, err := ReadRequest(alloc.NewBuffer().Clear())
	assert.Error(err).Equals(io.EOF)
}

func TestIPv6Request(t *testing.T) {
	v2testing.Current(t)

	request, err := ReadRequest(alloc.NewBuffer().Clear().AppendBytes(5, 1, 0, 4, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 0, 8))
	assert.Error(err).IsNil()
	assert.Byte(request.Command).Equals(1)
	assert.Bytes(request.IPv6[:]).Equals([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6})
	v2netassert.Port(request.Port).Equals(8)
}
