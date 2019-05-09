package net_test

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "v2ray.com/core/common/net"
)

func TestAddressProperty(t *testing.T) {
	type addrProprty struct {
		IP     []byte
		Domain string
		Family AddressFamily
		String string
	}

	testCases := []struct {
		Input  Address
		Output addrProprty
	}{
		{
			Input: IPAddress([]byte{byte(1), byte(2), byte(3), byte(4)}),
			Output: addrProprty{
				IP:     []byte{byte(1), byte(2), byte(3), byte(4)},
				Family: AddressFamilyIPv4,
				String: "1.2.3.4",
			},
		},
		{
			Input: IPAddress([]byte{
				byte(1), byte(2), byte(3), byte(4),
				byte(1), byte(2), byte(3), byte(4),
				byte(1), byte(2), byte(3), byte(4),
				byte(1), byte(2), byte(3), byte(4),
			}),
			Output: addrProprty{
				IP: []byte{
					byte(1), byte(2), byte(3), byte(4),
					byte(1), byte(2), byte(3), byte(4),
					byte(1), byte(2), byte(3), byte(4),
					byte(1), byte(2), byte(3), byte(4),
				},
				Family: AddressFamilyIPv6,
				String: "[102:304:102:304:102:304:102:304]",
			},
		},
		{
			Input: IPAddress([]byte{
				byte(0), byte(0), byte(0), byte(0),
				byte(0), byte(0), byte(0), byte(0),
				byte(0), byte(0), byte(255), byte(255),
				byte(1), byte(2), byte(3), byte(4),
			}),
			Output: addrProprty{
				IP:     []byte{byte(1), byte(2), byte(3), byte(4)},
				Family: AddressFamilyIPv4,
				String: "1.2.3.4",
			},
		},
		{
			Input: DomainAddress("v2ray.com"),
			Output: addrProprty{
				Domain: "v2ray.com",
				Family: AddressFamilyDomain,
				String: "v2ray.com",
			},
		},
		{
			Input: IPAddress(net.IPv4(1, 2, 3, 4)),
			Output: addrProprty{
				IP:     []byte{byte(1), byte(2), byte(3), byte(4)},
				Family: AddressFamilyIPv4,
				String: "1.2.3.4",
			},
		},
		{
			Input: ParseAddress("[2001:4860:0:2001::68]"),
			Output: addrProprty{
				IP:     []byte{0x20, 0x01, 0x48, 0x60, 0x00, 0x00, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68},
				Family: AddressFamilyIPv6,
				String: "[2001:4860:0:2001::68]",
			},
		},
		{
			Input: ParseAddress("::0"),
			Output: addrProprty{
				IP:     AnyIPv6.IP(),
				Family: AddressFamilyIPv6,
				String: "[::]",
			},
		},
		{
			Input: ParseAddress("[::ffff:123.151.71.143]"),
			Output: addrProprty{
				IP:     []byte{123, 151, 71, 143},
				Family: AddressFamilyIPv4,
				String: "123.151.71.143",
			},
		},
		{
			Input: NewIPOrDomain(ParseAddress("v2ray.com")).AsAddress(),
			Output: addrProprty{
				Domain: "v2ray.com",
				Family: AddressFamilyDomain,
				String: "v2ray.com",
			},
		},
		{
			Input: NewIPOrDomain(ParseAddress("8.8.8.8")).AsAddress(),
			Output: addrProprty{
				IP:     []byte{8, 8, 8, 8},
				Family: AddressFamilyIPv4,
				String: "8.8.8.8",
			},
		},
		{
			Input: NewIPOrDomain(ParseAddress("[2001:4860:0:2001::68]")).AsAddress(),
			Output: addrProprty{
				IP:     []byte{0x20, 0x01, 0x48, 0x60, 0x00, 0x00, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68},
				Family: AddressFamilyIPv6,
				String: "[2001:4860:0:2001::68]",
			},
		},
	}

	for _, testCase := range testCases {
		actual := addrProprty{
			Family: testCase.Input.Family(),
			String: testCase.Input.String(),
		}
		if testCase.Input.Family().IsIP() {
			actual.IP = testCase.Input.IP()
		} else {
			actual.Domain = testCase.Input.Domain()
		}

		if r := cmp.Diff(actual, testCase.Output); r != "" {
			t.Error("for input: ", testCase.Input, ":", r)
		}
	}
}

func TestInvalidAddressConvertion(t *testing.T) {
	panics := func(f func()) (ret bool) {
		defer func() {
			if r := recover(); r != nil {
				ret = true
			}
		}()
		f()
		return false
	}

	testCases := []func(){
		func() { ParseAddress("8.8.8.8").Domain() },
		func() { ParseAddress("2001:4860:0:2001::68").Domain() },
		func() { ParseAddress("v2ray.com").IP() },
	}
	for idx, testCase := range testCases {
		if !panics(testCase) {
			t.Error("case ", idx, " failed")
		}
	}
}

func BenchmarkParseAddressIPv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addr := ParseAddress("8.8.8.8")
		if addr.Family() != AddressFamilyIPv4 {
			panic("not ipv4")
		}
	}
}

func BenchmarkParseAddressIPv6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addr := ParseAddress("2001:4860:0:2001::68")
		if addr.Family() != AddressFamilyIPv6 {
			panic("not ipv6")
		}
	}
}

func BenchmarkParseAddressDomain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addr := ParseAddress("v2ray.com")
		if addr.Family() != AddressFamilyDomain {
			panic("not domain")
		}
	}
}
