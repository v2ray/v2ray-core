package dns_test

import (
	"testing"

	. "v2ray.com/core/app/dns"
	"v2ray.com/core/common/net"
)

func TestFakeIPServer(t *testing.T) {
	external := map[string][]string{"geosite.dat:cn": []string{"dbaidu.com"}, "geosite.dat:us": []string{"dgoogle.com"}}
	err := InitFakeIPServer(
		&Config_Fake{
			FakeRules: []string{
				"dv2ray.com",
				"rv2ray.com",
				"egeosite.dat:cn",
				"egeosite.dat:cn",
				"egeosite.dat:us",
			},
			FakeNet:      "1.1.1.1/30",
			Regeneration: Config_Fake_LRU,
		}, external)
	if err != nil {
		t.Error("failed to initialize fake ip server")
	}
	cases := []struct {
		input  string
		output bool
		IP     net.IP
	}{
		{
			input:  "www.v2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 0},
		},
		{
			input:  "www.v2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 0},
		},
		{
			input:  "v2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 1},
		},
		{
			input:  "www.v3ray.com",
			output: false,
			IP:     net.IP{1, 1, 1, 1},
		},
		{
			input:  "2ray.com",
			output: false,
			IP:     net.IP{1, 1, 1, 1},
		},
		{
			input:  "xv2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 2},
		},
		{
			input:  "xv2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 2},
		},
		{
			input:  "v2rayxcom",
			output: true,
			IP:     net.IP{1, 1, 1, 3},
		},
		{
			input:  "www.v2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 0},
		},
		{
			input:  "v2ray.com",
			output: true,
			IP:     net.IP{1, 1, 1, 1},
		},
		{
			input:  "www.baidu.com",
			output: true,
			IP:     net.IP{1, 1, 1, 2},
		},
		{
			input:  "www.google.com",
			output: true,
			IP:     net.IP{1, 1, 1, 3},
		},
	}

	for _, test := range cases {
		res := GetFakeIPForDomain(test.input)
		if res == nil {
			if test.output {
				t.Error("excpet a fake IP, but got nil")
			}
			continue
		}
		if len(res) != 1 {
			t.Error("excpet 1 fake IP, but got ", len(res))
		}
		if !test.output {
			continue
		}
		for i, v := range res[0].IP() {
			if v != test.IP[i] {
				t.Error("excpet fake IP: ", test.IP, " but got: ", res[0])
			}
		}
		domain := GetDomainForFakeIP(res[0])
		if domain != test.input {
			t.Error("excpet origin domain name, but got ", domain)
		}
	}

	ResetFakeIPServer()

}
