package protocol_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol"
)

func TestAddressReading(t *testing.T) {
	data := []struct {
		Options []AddressOption
		Input   []byte
		Address net.Address
		Port    net.Port
		Error   bool
	}{
		{
			Options: []AddressOption{},
			Input:   []byte{},
			Error:   true,
		},
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
			Options: []AddressOption{AddressFamilyByte(0x01, net.AddressFamilyIPv4), PortThenAddress()},
			Input:   []byte{0, 53, 1, 0, 0, 0, 0},
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
		{
			Options: []AddressOption{AddressFamilyByte(0x03, net.AddressFamilyDomain)},
			Input:   []byte{3, 7, 56, 46, 56, 46, 56, 46, 56, 0, 80},
			Address: net.ParseAddress("8.8.8.8"),
			Port:    net.Port(80),
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x03, net.AddressFamilyDomain)},
			Input:   []byte{3, 7, 10, 46, 56, 46, 56, 46, 56, 0, 80},
			Error:   true,
		},
		{
			Options: []AddressOption{AddressFamilyByte(0x03, net.AddressFamilyDomain)},
			Input:   append(append([]byte{3, 24}, []byte("2a00:1450:4007:816::200e")...), 0, 80),
			Address: net.ParseAddress("2a00:1450:4007:816::200e"),
			Port:    net.Port(80),
		},
	}

	for _, tc := range data {
		b := buf.New()
		parser := NewAddressParser(tc.Options...)
		addr, port, err := parser.ReadAddressPort(b, bytes.NewReader(tc.Input))
		b.Release()
		if tc.Error {
			if err == nil {
				t.Errorf("Expect error but not: %v", tc)
			}
		} else {
			if err != nil {
				t.Errorf("Expect no error but: %s %v", err.Error(), tc)
			}

			if addr != tc.Address {
				t.Error("Got address ", addr.String(), " want ", tc.Address.String())
			}

			if tc.Port != port {
				t.Error("Got port ", port, " want ", tc.Port)
			}
		}
	}
}

func TestAddressWriting(t *testing.T) {
	data := []struct {
		Options []AddressOption
		Address net.Address
		Port    net.Port
		Bytes   []byte
		Error   bool
	}{
		{
			Options: []AddressOption{AddressFamilyByte(0x01, net.AddressFamilyIPv4)},
			Address: net.LocalHostIP,
			Port:    net.Port(80),
			Bytes:   []byte{1, 127, 0, 0, 1, 0, 80},
		},
	}

	for _, tc := range data {
		parser := NewAddressParser(tc.Options...)

		b := buf.New()
		err := parser.WriteAddressPort(b, tc.Address, tc.Port)
		if tc.Error {
			if err == nil {
				t.Error("Expect error but nil")
			}
		} else {
			common.Must(err)
			if diff := cmp.Diff(tc.Bytes, b.Bytes()); diff != "" {
				t.Error(err)
			}
		}
	}
}

func BenchmarkAddressReadingIPv4(b *testing.B) {
	parser := NewAddressParser(AddressFamilyByte(0x01, net.AddressFamilyIPv4))
	cache := buf.New()
	defer cache.Release()

	payload := buf.New()
	defer payload.Release()

	raw := []byte{1, 0, 0, 0, 0, 0, 53}
	payload.Write(raw)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.ReadAddressPort(cache, payload)
		common.Must(err)
		cache.Clear()
		payload.Clear()
		payload.Extend(int32(len(raw)))
	}
}

func BenchmarkAddressReadingIPv6(b *testing.B) {
	parser := NewAddressParser(AddressFamilyByte(0x04, net.AddressFamilyIPv6))
	cache := buf.New()
	defer cache.Release()

	payload := buf.New()
	defer payload.Release()

	raw := []byte{4, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 0, 80}
	payload.Write(raw)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.ReadAddressPort(cache, payload)
		common.Must(err)
		cache.Clear()
		payload.Clear()
		payload.Extend(int32(len(raw)))
	}
}

func BenchmarkAddressReadingDomain(b *testing.B) {
	parser := NewAddressParser(AddressFamilyByte(0x03, net.AddressFamilyDomain))
	cache := buf.New()
	defer cache.Release()

	payload := buf.New()
	defer payload.Release()

	raw := []byte{3, 9, 118, 50, 114, 97, 121, 46, 99, 111, 109, 0, 80}
	payload.Write(raw)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := parser.ReadAddressPort(cache, payload)
		common.Must(err)
		cache.Clear()
		payload.Clear()
		payload.Extend(int32(len(raw)))
	}
}

func BenchmarkAddressWritingIPv4(b *testing.B) {
	parser := NewAddressParser(AddressFamilyByte(0x01, net.AddressFamilyIPv4))
	writer := buf.New()
	defer writer.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		common.Must(parser.WriteAddressPort(writer, net.LocalHostIP, net.Port(80)))
		writer.Clear()
	}
}

func BenchmarkAddressWritingIPv6(b *testing.B) {
	parser := NewAddressParser(AddressFamilyByte(0x04, net.AddressFamilyIPv6))
	writer := buf.New()
	defer writer.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		common.Must(parser.WriteAddressPort(writer, net.LocalHostIPv6, net.Port(80)))
		writer.Clear()
	}
}

func BenchmarkAddressWritingDomain(b *testing.B) {
	parser := NewAddressParser(AddressFamilyByte(0x02, net.AddressFamilyDomain))
	writer := buf.New()
	defer writer.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		common.Must(parser.WriteAddressPort(writer, net.DomainAddress("www.v2ray.com"), net.Port(80)))
		writer.Clear()
	}
}
