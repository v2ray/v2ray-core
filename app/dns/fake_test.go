package dns_test

import (
	"testing"

	. "v2ray.com/core/app/dns"
)

func TestFakeIPServer(t *testing.T) {
	external := map[string][]string{"geosite.dat:cn": []string{"dbaidu.com"}, "geosite.dat:us": []string{"dgoogle.com"}}
	InitFakeIPServer([]string{
		"dv2ray.com",
		"rv2ray.com",
		"egeosite.dat:cn",
		"egeosite.dat:cn",
		"egeosite.dat:us",
	}, external)
	cases := []struct {
		input  string
		output bool
	}{
		{
			input:  "www.v2ray.com",
			output: true,
		},
		{
			input:  "v2ray.com",
			output: true,
		},
		{
			input:  "www.v3ray.com",
			output: false,
		},
		{
			input:  "2ray.com",
			output: false,
		},
		{
			input:  "xv2ray.com",
			output: false,
		},
		{
			input:  "v2ray.com",
			output: true,
		},
		{
			input:  "xv2ray.com",
			output: false,
		},
		{
			input:  "v2rayxcom",
			output: true,
		},
		{
			input:  "www.baidu.com",
			output: true,
		},
		{
			input:  "www.google.com",
			output: true,
		},
	}

	for _, test := range cases {
		res := GetFakeIPForDomain(test.input)
		if res == nil {
			if test.output {
				t.Error("excpet a fake IP, but got nil")
			}
			break
		}
		if len(res) != 1 {
			t.Error("excpet 1 fake IP, but got ", len(res))
		}
		domain := GetDomainForFakeIP(res[0])
		if domain != test.input {
			t.Error("excpet origin domain name, but got ", domain)
		}
	}

	ResetFakeIPServer()

}
