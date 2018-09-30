package strmatcher_test

import (
	"testing"

	. "v2ray.com/core/common/strmatcher"
)

func TestDomainMatcherGroup(t *testing.T) {
	g := new(DomainMatcherGroup)
	g.Add("v2ray.com", 1)
	g.Add("google.com", 2)
	g.Add("x.a.com", 3)
	g.Add("a.b.com", 4)
	g.Add("c.a.b.com", 5)

	testCases := []struct {
		Domain string
		Result uint32
	}{
		{
			Domain: "x.v2ray.com",
			Result: 1,
		},
		{
			Domain: "y.com",
			Result: 0,
		},
		{
			Domain: "a.b.com",
			Result: 4,
		},
		{
			Domain: "c.a.b.com",
			Result: 4,
		},
	}

	for _, testCase := range testCases {
		r := g.Match(testCase.Domain)
		if r != testCase.Result {
			t.Error("Failed to match domain: ", testCase.Domain, ", expect ", testCase.Result, ", but got ", r)
		}
	}
}

func TestEmptyDomainMatcherGroup(t *testing.T) {
	g := new(DomainMatcherGroup)
	r := g.Match("v2ray.com")
	if r != 0 {
		t.Error("Expect 0, but ", r)
	}
}
