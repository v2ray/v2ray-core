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

func TestRoutingRule(t *testing.T) {
	assert := assert.On(t)

	type ruleTest struct {
		input  context.Context
		output bool
	}

	cases := []struct {
		rule *RoutingRule
		test []ruleTest
	}{
		{
			rule: &RoutingRule{
				Domain: []*Domain{
					{
						Value: "v2ray.com",
						Type:  Domain_Plain,
					},
					{
						Value: "google.com",
						Type:  Domain_Domain,
					},
				},
			},
			test: []ruleTest{
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("v2ray.com"), 80)),
					output: true,
				},
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.v2ray.com.www"), 80)),
					output: true,
				},
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("v2ray.co"), 80)),
					output: false,
				},
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.google.com"), 80)),
					output: true,
				},
			},
		},
		{
			rule: &RoutingRule{
				Cidr: []*CIDR{
					{
						Ip:     []byte{8, 8, 8, 8},
						Prefix: 32,
					},
					{
						Ip:     []byte{8, 8, 8, 8},
						Prefix: 32,
					},
					{
						Ip:     net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334").IP(),
						Prefix: 128,
					},
				},
			},
			test: []ruleTest{
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)),
					output: true,
				},
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("8.8.4.4"), 80)),
					output: false,
				},
				ruleTest{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 80)),
					output: true,
				},
			},
		},
	}

	for _, test := range cases {
		cond, err := test.rule.BuildCondition()
		assert.Error(err).IsNil()

		for _, t := range test.test {
			assert.Bool(cond.Apply(t.input)).Equals(t.output)
		}
	}
}
