package socks

import (
	"bytes"
	"testing"
)

func TestAuthenticationRequestRead(t *testing.T) {
	rawRequest := []byte{
		0x05, // version
		0x01, // nMethods
		0x02, // methods
	}
	request, err := ReadAuthentication(bytes.NewReader(rawRequest))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if request.version != 0x05 {
		t.Errorf("Expected version 5, but got %d", request.version)
	}
	if request.nMethods != 0x01 {
		t.Errorf("Expected nMethod 1, but got %d", request.nMethods)
	}
	if request.authMethods[0] != 0x02 {
		t.Errorf("Expected method 2, but got %d", request.authMethods[0])
	}
}

func TestAuthenticationResponseToBytes(t *testing.T) {
	socksVersion := uint8(5)
	authMethod := uint8(1)
	response := Socks5AuthenticationResponse{socksVersion, authMethod}
	bytes := response.ToBytes()
	if bytes[0] != socksVersion {
		t.Errorf("Unexpected Socks version %d", bytes[0])
	}
	if bytes[1] != authMethod {
		t.Errorf("Unexpected Socks auth method %d", bytes[1])
	}
}

func TestRequestRead(t *testing.T) {
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
	if request.version != 0x05 {
		t.Errorf("Expected version 5, but got %d", request.version)
	}
	if request.command != 0x01 {
		t.Errorf("Expected command 1, but got %d", request.command)
	}
	if request.addrType != 0x01 {
		t.Errorf("Expected addresstype 1, but got %d", request.addrType)
	}
	if !bytes.Equal([]byte{0x72, 0x72, 0x72, 0x72}, request.ipv4[:]) {
		t.Errorf("Expected IPv4 address 114.114.114.114, but got %v", request.ipv4[:])
	}
	if request.port != 53 {
		t.Errorf("Expected port 53, but got %d", request.port)
	}
}

func TestResponseToBytes(t *testing.T) {
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
	if !bytes.Equal(rawResponse, expectedBytes) {
		t.Errorf("Expected response %v, but got %v", expectedBytes, rawResponse)
	}
}
