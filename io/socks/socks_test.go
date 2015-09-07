package socks

import (
	"testing"
)

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
