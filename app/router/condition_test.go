package router_test

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	proto "github.com/golang/protobuf/proto"

	"v2ray.com/core/app/dispatcher"
	. "v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/protocol/http"
	"v2ray.com/core/common/session"
	. "v2ray.com/ext/assert"
	"v2ray.com/ext/sysio"
)

func withOutbound(outbound *session.Outbound) context.Context {
	return session.ContextWithOutbound(context.Background(), outbound)
}

func withInbound(inbound *session.Inbound) context.Context {
	return session.ContextWithInbound(context.Background(), inbound)
}

func TestRoutingRule(t *testing.T) {
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
					{
						Value: "^facebook\\.com$",
						Type:  Domain_Regex,
					},
				},
			},
			test: []ruleTest{
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("v2ray.com"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("www.v2ray.com.www"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("v2ray.co"), 80)}),
					output: false,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("www.google.com"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("facebook.com"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("www.facebook.com"), 80)}),
					output: false,
				},
				{
					input:  context.Background(),
					output: false,
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
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.4.4"), 80)}),
					output: false,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 80)}),
					output: true,
				},
				{
					input:  context.Background(),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				SourceCidr: []*CIDR{
					{
						Ip:     []byte{192, 168, 0, 0},
						Prefix: 16,
					},
				},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{Source: net.TCPDestination(net.ParseAddress("192.168.0.1"), 80)}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Source: net.TCPDestination(net.ParseAddress("10.0.0.1"), 80)}),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				UserEmail: []string{
					"admin@v2ray.com",
				},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{User: &protocol.MemoryUser{Email: "admin@v2ray.com"}}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{User: &protocol.MemoryUser{Email: "love@v2ray.com"}}),
					output: false,
				},
				{
					input:  context.Background(),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Protocol: []string{"http"},
			},
			test: []ruleTest{
				{
					input:  dispatcher.ContextWithSniffingResult(context.Background(), &http.SniffHeader{}),
					output: true,
				},
			},
		},
		{
			rule: &RoutingRule{
				InboundTag: []string{"test", "test1"},
			},
			test: []ruleTest{
				{
					input:  withInbound(&session.Inbound{Tag: "test"}),
					output: true,
				},
				{
					input:  withInbound(&session.Inbound{Tag: "test2"}),
					output: false,
				},
			},
		},
	}

	for _, test := range cases {
		cond, err := test.rule.BuildCondition()
		common.Must(err)

		for _, subtest := range test.test {
			actual := cond.Apply(subtest.input)
			if actual != subtest.output {
				t.Error("test case failed: ", subtest.input, " expected ", subtest.output, " but got ", actual)
			}
		}
	}
}

func loadGeoSite(country string) ([]*Domain, error) {
	geositeBytes, err := sysio.ReadAsset("geosite.dat")
	if err != nil {
		return nil, err
	}
	var geositeList GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return nil, err
	}

	for _, site := range geositeList.Entry {
		if site.CountryCode == country {
			return site.Domain, nil
		}
	}

	return nil, errors.New("country not found: " + country)
}

func TestChinaSites(t *testing.T) {
	assert := With(t)

	common.Must(sysio.CopyFile(platform.GetAssetLocation("geosite.dat"), filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "release", "config", "geosite.dat")))

	domains, err := loadGeoSite("CN")
	assert(err, IsNil)

	matcher, err := NewDomainMatcher(domains)
	common.Must(err)

	assert(matcher.ApplyDomain("163.com"), IsTrue)
	assert(matcher.ApplyDomain("163.com"), IsTrue)
	assert(matcher.ApplyDomain("164.com"), IsFalse)
	assert(matcher.ApplyDomain("164.com"), IsFalse)

	for i := 0; i < 1024; i++ {
		assert(matcher.ApplyDomain(strconv.Itoa(i)+".not-exists.com"), IsFalse)
	}
}
