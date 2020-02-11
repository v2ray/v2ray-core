package strmatcher_test

import (
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/strmatcher"
)

func TestMatcher(t *testing.T) {
	cases := []struct {
		pattern string
		mType   Type
		input   string
		output  bool
	}{
		{
			pattern: "v2ray.com",
			mType:   Domain,
			input:   "www.v2ray.com",
			output:  true,
		},
		{
			pattern: "v2ray.com",
			mType:   Domain,
			input:   "v2ray.com",
			output:  true,
		},
		{
			pattern: "v2ray.com",
			mType:   Domain,
			input:   "www.v3ray.com",
			output:  false,
		},
		{
			pattern: "v2ray.com",
			mType:   Domain,
			input:   "2ray.com",
			output:  false,
		},
		{
			pattern: "v2ray.com",
			mType:   Domain,
			input:   "xv2ray.com",
			output:  false,
		},
		{
			pattern: "v2ray.com",
			mType:   Full,
			input:   "v2ray.com",
			output:  true,
		},
		{
			pattern: "v2ray.com",
			mType:   Full,
			input:   "xv2ray.com",
			output:  false,
		},
		{
			pattern: "v2ray.com",
			mType:   Regex,
			input:   "v2rayxcom",
			output:  true,
		},
	}
	for _, test := range cases {
		matcher, err := test.mType.New(test.pattern)
		common.Must(err)
		if m := matcher.Match(test.input); m != test.output {
			t.Error("unexpected output: ", m, " for test case ", test)
		}
	}
}

func TestOrMatcher(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		output  bool
	}{
		{
			pattern: "dv2ray.com",
			input:   "www.v2ray.com",
			output:  true,
		},
		{
			pattern: "dv2ray.com",
			input:   "v2ray.com",
			output:  true,
		},
		{
			pattern: "dv2ray.com",
			input:   "www.v3ray.com",
			output:  false,
		},
		{
			pattern: "dv2ray.com",
			input:   "2ray.com",
			output:  false,
		},
		{
			pattern: "dv2ray.com",
			input:   "xv2ray.com",
			output:  false,
		},
		{
			pattern: "fv2ray.com",
			input:   "v2ray.com",
			output:  true,
		},
		{
			pattern: "fv2ray.com",
			input:   "xv2ray.com",
			output:  false,
		},
		{
			pattern: "rv2ray.com",
			input:   "v2rayxcom",
			output:  true,
		},
		{
			pattern: "egeosite.dat:cn",
			input:   "www.baidu.com",
			output:  true,
		},
		{
			pattern: "egeosite.dat:cn",
			input:   "www.google.com",
			output:  false,
		},
		{
			pattern: "egeosite.dat:us",
			input:   "www.google.com",
			output:  true,
		},
	}
	external := map[string][]string{"geosite.dat:cn": []string{"dbaidu.com"}, "geosite.dat:us": []string{"dgoogle.com"}}
	for _, test := range cases {
		om := new(OrMatcher)
		om.New()
		err := om.ParsePattern(test.pattern, external)
		common.Must(err)
		if m := om.Match(test.input); m != test.output {
			t.Error("unexpected output: ", m, " for test case ", test)
		}
	}
}
