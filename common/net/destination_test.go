package net_test

import (
	"testing"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestTCPDestination(t *testing.T) {
	assert := assert.On(t)

	dest := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Destination(dest).IsTCP()
	assert.Destination(dest).IsNotUDP()
	assert.Destination(dest).EqualsString("tcp:1.2.3.4:80")
}

func TestUDPDestination(t *testing.T) {
	assert := assert.On(t)

	dest := v2net.UDPDestination(v2net.IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}), 53)
	assert.Destination(dest).IsNotTCP()
	assert.Destination(dest).IsUDP()
	assert.Destination(dest).EqualsString("udp:[2001:4860:4860::8888]:53")
}

func TestTCPDestinationEquals(t *testing.T) {
	assert := assert.On(t)

	dest := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Bool(dest.Equals(v2net.Destination{})).IsFalse()

	dest2 := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Bool(dest.Equals(dest2)).IsTrue()

	dest3 := v2net.UDPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Bool(dest.Equals(dest3)).IsFalse()

	dest4 := v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80)
	assert.Bool(dest.Equals(dest4)).IsFalse()
}

func TestUDPDestinationEquals(t *testing.T) {
	assert := assert.On(t)

	dest := v2net.UDPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Bool(dest.Equals(v2net.Destination{})).IsFalse()

	dest2 := v2net.UDPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Bool(dest.Equals(dest2)).IsTrue()

	dest3 := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	assert.Bool(dest.Equals(dest3)).IsFalse()

	dest4 := v2net.UDPDestination(v2net.DomainAddress("v2ray.com"), 80)
	assert.Bool(dest.Equals(dest4)).IsFalse()
}
