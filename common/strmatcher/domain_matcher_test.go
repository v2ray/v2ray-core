package strmatcher_test

import (
	"reflect"
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
	g.Add("x.y.com", 4)
	g.Add("x.y.com", 6)

	testCases := []struct {
		Domain string
		Result []uint32
	}{
		{
			Domain: "x.v2ray.com",
			Result: []uint32{1},
		},
		{
			Domain: "y.com",
			Result: nil,
		},
		{
			Domain: "a.b.com",
			Result: []uint32{4},
		},
		{ // Matches [c.a.b.com, a.b.com]
			Domain: "c.a.b.com",
			Result: []uint32{5, 4},
		},
		{
			Domain: "c.a..b.com",
			Result: nil,
		},
		{
			Domain: ".com",
			Result: nil,
		},
		{
			Domain: "com",
			Result: nil,
		},
		{
			Domain: "",
			Result: nil,
		},
		{
			Domain: "x.y.com",
			Result: []uint32{4, 6},
		},
	}

	for _, testCase := range testCases {
		r := g.Match(testCase.Domain)
		if !reflect.DeepEqual(r, testCase.Result) {
			t.Error("Failed to match domain: ", testCase.Domain, ", expect ", testCase.Result, ", but got ", r)
		}
	}
}

func TestEmptyDomainMatcherGroup(t *testing.T) {
	g := new(DomainMatcherGroup)
	r := g.Match("v2ray.com")
	if len(r) != 0 {
		t.Error("Expect [], but ", r)
	}
}
