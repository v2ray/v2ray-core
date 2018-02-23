package protocol_test

import (
	"bytes"
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol"
	. "v2ray.com/ext/assert"
)

func TestAddressParser(t *testing.T) {
	assert := With(t)

	data := []struct {
		Options []AddressOption
		Input   []byte
		Address net.Address
		Port    net.Port
		Error   bool
	}{
		{
			Options: []AddressOption{},
			Input:   []byte{0, 0, 0, 0, 0},
			Error:   true,
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x01, net.AddressFamilyIPv4)},
			Input:   []byte{1, 0, 0, 0, 0, 0, 53},
			Address: net.IPAddress([]byte{0, 0, 0, 0}),
			Port:    net.Port(53),
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x01, net.AddressFamilyIPv4)},
			Input:   []byte{1, 0, 0, 0, 0},
			Error:   true,
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x04, net.AddressFamilyIPv6)},
			Input:   []byte{4, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 0, 80},
			Address: net.IPAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}),
			Port:    net.Port(80),
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x03, net.AddressFamilyDomain)},
			Input:   []byte{3, 9, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 80},
			Address: net.DomainAddress("v2ray.com"),
			Port:    net.Port(80),
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x03, net.AddressFamilyDomain)},
			Input:   []byte{3, 9, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0},
			Error:   true,
		},
	}

	for _, tc := range data {
		b := buf.New()
		parser := NewAddressParser(tc.Options...)
		addr, port, err := parser.ReadAddressPort(b, bytes.NewReader(tc.Input))
		b.Release()
		if tc.Error {
			assert(err, IsNotNil)
		} else {
			assert(addr, Equals, tc.Address)
			assert(port, Equals, tc.Port)
		}
	}
}
