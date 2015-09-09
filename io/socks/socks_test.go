package socks

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestAuthenticationRequestRead(t *testing.T) {
	assert := unit.Assert(t)

	rawRequest := []byte{
		0x05, // version
		0x01, // nMethods
		0x02, // methods
	}
	request, err := ReadAuthentication(bytes.NewReader(rawRequest))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	assert.Byte(request.version).Named("Version").Equals(0x05)
	assert.Byte(request.nMethods).Named("#Methods").Equals(0x01)
	assert.Byte(request.authMethods[0]).Named("Auth Method").Equals(0x02)
}

func TestAuthenticationResponseToBytes(t *testing.T) {
	assert := unit.Assert(t)

	socksVersion := uint8(5)
	authMethod := uint8(1)
	response := Socks5AuthenticationResponse{socksVersion, authMethod}
	bytes := response.ToBytes()

	assert.Byte(bytes[0]).Named("Version").Equals(socksVersion)
	assert.Byte(bytes[1]).Named("Auth Method").Equals(authMethod)
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
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
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
