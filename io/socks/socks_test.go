package socks

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestHasAuthenticationMethod(t *testing.T) {
	assert := unit.Assert(t)

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
	assert := unit.Assert(t)

	rawRequest := []byte{
		0x05, // version
		0x01, // nMethods
		0x02, // methods
	}
	request, err := ReadAuthentication(bytes.NewReader(rawRequest))
	assert.Error(err).IsNil()
	assert.Byte(request.version).Named("Version").Equals(0x05)
	assert.Byte(request.nMethods).Named("#Methods").Equals(0x01)
	assert.Byte(request.authMethods[0]).Named("Auth Method").Equals(0x02)
}

func TestAuthenticationResponseWrite(t *testing.T) {
	assert := unit.Assert(t)

	response := Socks5AuthenticationResponse{
		version:    socksVersion,
		authMethod: byte(0x05),
	}

	buffer := bytes.NewBuffer(make([]byte, 0, 10))
	WriteAuthentication(buffer, &response)
	assert.Bytes(buffer.Bytes()).Equals([]byte{socksVersion, byte(0x05)})
}

func TestRequestRead(t *testing.T) {
	assert := unit.Assert(t)

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
	assert.Uint16(request.Port).Named("Port").Equals(53)
}

func TestResponseToBytes(t *testing.T) {
	assert := unit.Assert(t)

	response := Socks5Response{
		socksVersion,
		ErrorSuccess,
		AddrTypeIPv4,
		[4]byte{0x72, 0x72, 0x72, 0x72},
		"",
		[16]byte{},
		uint16(53),
	}
	rawResponse := response.toBytes()
	expectedBytes := []byte{
		socksVersion,
		ErrorSuccess,
		byte(0x00),
		AddrTypeIPv4,
		0x72, 0x72, 0x72, 0x72,
		byte(0x00), byte(0x035),
	}
	assert.Bytes(rawResponse).Named("raw response").Equals(expectedBytes)
}
