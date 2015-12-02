package net

import (
	"testing"

	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestTCPDestination(t *testing.T) {
	v2testing.Current(t)

	dest := NewTCPDestination(IPAddress([]byte{1, 2, 3, 4}, 80))
	assert.Bool(dest.IsTCP()).IsTrue()
	assert.Bool(dest.IsUDP()).IsFalse()
	assert.String(dest.String()).Equals("tcp:1.2.3.4:80")
}

func TestUDPDestination(t *testing.T) {
	v2testing.Current(t)

	dest := NewUDPDestination(IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}, 53))
	assert.Bool(dest.IsTCP()).IsFalse()
	assert.Bool(dest.IsUDP()).IsTrue()
	assert.String(dest.String()).Equals("udp:[2001:4860:4860::8888]:53")
}
