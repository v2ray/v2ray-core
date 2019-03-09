package net_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "v2ray.com/core/common/net"
)

func TestDestinationProperty(t *testing.T) {
	testCases := []struct {
		Input     Destination
		Network   Network
		String    string
		NetString string
	}{
		{
			Input:     TCPDestination(IPAddress([]byte{1, 2, 3, 4}), 80),
			Network:   Network_TCP,
			String:    "tcp:1.2.3.4:80",
			NetString: "1.2.3.4:80",
		},
		{
			Input:     UDPDestination(IPAddress([]byte{0x20, 0x01, 0x48, 0x60, 0x48, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x88}), 53),
			Network:   Network_UDP,
			String:    "udp:[2001:4860:4860::8888]:53",
			NetString: "[2001:4860:4860::8888]:53",
		},
	}

	for _, testCase := range testCases {
		dest := testCase.Input
		if r := cmp.Diff(dest.Network, testCase.Network); r != "" {
			t.Error("unexpected Network in ", dest.String(), ": ", r)
		}
		if r := cmp.Diff(dest.String(), testCase.String); r != "" {
			t.Error(r)
		}
		if r := cmp.Diff(dest.NetAddr(), testCase.NetString); r != "" {
			t.Error(r)
		}
	}
}

func TestDestinationParse(t *testing.T) {
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
			if err != nil {
				t.Error("for test case: ", testcase.Input, " expected no error, but got ", err)
			}
			if d != testcase.Output {
				t.Error("for test case: ", testcase.Input, " expected output: ", testcase.Output.String(), " but got ", d.String())
			}
		} else {
			if err == nil {
				t.Error("for test case: ", testcase.Input, " expected error, but got nil")
			}
		}
	}
}
