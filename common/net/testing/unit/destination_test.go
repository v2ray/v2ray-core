package unit

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestTCPDestination(t *testing.T) {
	v2testing.Current(t)

	dest := v2net.NewTCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}, 80))
	assert.Bool(dest.IsTCP()).IsTrue()
	assert.Bool(dest.IsUDP()).IsFalse()
	assert.StringLiteral(dest.String()).Equals("tcp:1.2.3.4:80")
}

func TestUDPDestination(t *testing.T) {
	v2testing.Current(t)

	dest := v2net.NewUDPDestination(v2net.IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}, 53))
	assert.Bool(dest.IsTCP()).IsFalse()
	assert.Bool(dest.IsUDP()).IsTrue()
	assert.StringLiteral(dest.String()).Equals("udp:[2001:4860:4860::8888]:53")
}
