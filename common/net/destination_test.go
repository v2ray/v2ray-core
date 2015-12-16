package net_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestTCPDestination(t *testing.T) {
	v2testing.Current(t)

	dest := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	v2netassert.Destination(dest).IsTCP()
	v2netassert.Destination(dest).IsNotUDP()
	assert.String(dest).Equals("tcp:1.2.3.4:80")
}

func TestUDPDestination(t *testing.T) {
	v2testing.Current(t)

	dest := v2net.UDPDestination(v2net.IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}), 53)
	v2netassert.Destination(dest).IsNotTCP()
	v2netassert.Destination(dest).IsUDP()
	assert.String(dest).Equals("udp:[2001:4860:4860::8888]:53")
}
