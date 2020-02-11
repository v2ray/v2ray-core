package strmatcher_test

import (
	"strconv"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/strmatcher"
)

func BenchmarkDomainMatcherGroup(b *testing.B) {
	g := new(DomainMatcherGroup)

	for i := 1; i <= 1024; i++ {
		g.Add(strconv.Itoa(i)+".v2ray.com", uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Match("0.v2ray.com")
	}
}

func BenchmarkFullMatcherGroup(b *testing.B) {
	g := new(FullMatcherGroup)

	for i := 1; i <= 1024; i++ {
		g.Add(strconv.Itoa(i)+".v2ray.com", uint32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Match("0.v2ray.com")
	}
}

func BenchmarkMarchGroup(b *testing.B) {
	g := new(MatcherGroup)
	for i := 1; i <= 1024; i++ {
		m, err := Domain.New(strconv.Itoa(i) + ".v2ray.com")
		common.Must(err)
		g.Add(m)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Match("0.v2ray.com")
	}
}

func BenchmarkDomainGroupMatcher(b *testing.B) {
	g := new(DomainGroupMatcher)

	for i := 1; i <= 1024; i++ {
		g.Add(strconv.Itoa(i) + ".v2ray.com")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Match("0.v2ray.com")
	}
}

func BenchmarkFullGroupMatcher(b *testing.B) {
	g := new(FullGroupMatcher)
	g.New()

	for i := 1; i <= 1024; i++ {
		g.Add(strconv.Itoa(i) + ".v2ray.com")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Match("0.v2ray.com")
	}
}

func BenchmarkOrMatcher(b *testing.B) {
	g := new(OrMatcher)
	g.New()

	external := make(map[string][]string)
	for i := 1; i <= 1024; i++ {
		err := g.ParsePattern("d"+strconv.Itoa(i)+".v2ray.com", external)
		common.Must(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = g.Match("0.v2ray.com")
	}
}
