package net_test

import (
	"testing"

	. "v2ray.com/core/common/net"
	. "v2ray.com/core/common/net/testing"
	. "v2ray.com/ext/assert"
)

func TestTCPDestination(t *testing.T) {
	assert := With(t)

	dest := TCPDestination(IPAddress([]byte{1, 2, 3, 4}), 80)
	assert(dest, IsTCP)
	assert(dest, Not(IsUDP))
	assert(dest.String(), Equals, "tcp:1.2.3.4:80")
}

func TestUDPDestination(t *testing.T) {
	assert := With(t)

	dest := UDPDestination(IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}), 53)
	assert(dest, Not(IsTCP))
	assert(dest, IsUDP)
	assert(dest.String(), Equals, "udp:[2001:4860:4860::8888]:53")
}

func TestDestinationParse(t *testing.T) {
	assert := With(t)

	cases := []struct {
		Input  string
		Output Destination
		Error  bool
	}{
		{
			Input:  "tcp:127.0.0.1:80",
			Output: TCPDestination(LocalHostIP, Port(80)),
		},
		{
			Input:  "udp:8.8.8.8:53",
			Output: UDPDestination(IPAddress([]byte{8, 8, 8, 8}), Port(53)),
		},
		{
			Input: "8.8.8.8:53",
			Output: Destination{
				Address: IPAddress([]byte{8, 8, 8, 8}),
				Port:    Port(53),
			},
		},
		{
			Input: ":53",
			Output: Destination{
				Address: AnyIP,
				Port:    Port(53),
			},
		},
		{
			Input: "8.8.8.8",
			Error: true,
		},
		{
			Input: "8.8.8.8:http",
			Error: true,
		},
	}

	for _, testcase := range cases {
		d, err := ParseDestination(testcase.Input)
		if !testcase.Error {
			assert(err, IsNil)
			assert(d, Equals, testcase.Output)
		} else {
			assert(err, IsNotNil)
		}
	}
}
