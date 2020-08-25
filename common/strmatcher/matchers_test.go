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
