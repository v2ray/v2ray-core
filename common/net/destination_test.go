package net_test

import (
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestTCPDestination(t *testing.T) {
	assert := assert.On(t)

	dest := TCPDestination(IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Destination(dest).IsTCP()
	assert.Destination(dest).IsNotUDP()
	assert.Destination(dest).EqualsString("tcp:1.2.3.4:80")
}

func TestUDPDestination(t *testing.T) {
	assert := assert.On(t)

	dest := UDPDestination(IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}), 53)
	assert.Destination(dest).IsNotTCP()
	assert.Destination(dest).IsUDP()
	assert.Destination(dest).EqualsString("udp:[2001:4860:4860::8888]:53")
}
