package strmatcher_test

import (
	"testing"

	. "v2ray.com/core/common/strmatcher"
)

func TestFullMatcherGroup(t *testing.T) {
	g := new(FullMatcherGroup)
	g.Add("v2ray.com", 1)
	g.Add("google.com", 2)
	g.Add("x.a.com", 3)

	testCases := []struct {
		Domain string
		Result uint32
	}{
		{
			Domain: "v2ray.com",
			Result: 1,
		},
		{
			Domain: "y.com",
			Result: 0,
		},
	}

	for _, testCase := range testCases {
		r := g.Match(testCase.Domain)
		if r != testCase.Result {
			t.Error("Failed to match domain: ", testCase.Domain, ", expect ", testCase.Result, ", but got ", r)
		}
	}
}

func TestEmptyFullMatcherGroup(t *testing.T) {
	g := new(FullMatcherGroup)
	r := g.Match("v2ray.com")
	if r != 0 {
		t.Error("Expect 0, but ", r)
	}
}
