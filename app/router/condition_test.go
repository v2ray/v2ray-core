package router_test

import (
	"context"
	"testing"

	. "v2ray.com/core/app/router"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/testing/assert"
)

func TestSubDomainMatcher(t *testing.T) {
	assert := assert.On(t)

	cases := []struct {
		pattern string
		input   context.Context
		output  bool
	}{
		{
			pattern: "v2ray.com",
			input:   proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.v2ray.com"), 80)),
			output:  true,
		},
		{
			pattern: "v2ray.com",
			input:   proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("v2ray.com"), 80)),
			output:  true,
		},
		{
			pattern: "v2ray.com",
			input:   proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.v3ray.com"), 80)),
			output:  false,
		},
		{
			pattern: "v2ray.com",
			input:   proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("2ray.com"), 80)),
			output:  false,
		},
		{
			pattern: "v2ray.com",
			input:   proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("xv2ray.com"), 80)),
			output:  false,
		},
	}
	for _, test := range cases {
		matcher := NewSubDomainMatcher(test.pattern)
		assert.Bool(matcher.Apply(test.input) == test.output).IsTrue()
	}
}
