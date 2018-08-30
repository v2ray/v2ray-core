package router_test

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	proto "github.com/golang/protobuf/proto"
	"v2ray.com/core/app/dispatcher"
	. "v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/protocol/http"
	"v2ray.com/core/proxy"
	. "v2ray.com/ext/assert"
	"v2ray.com/ext/sysio"
)

func TestRoutingRule(t *testing.T) {
	assert := With(t)

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
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("v2ray.com"), 80)),
					output: true,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.v2ray.com.www"), 80)),
					output: true,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("v2ray.co"), 80)),
					output: false,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.google.com"), 80)),
					output: true,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("facebook.com"), 80)),
					output: true,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("www.facebook.com"), 80)),
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
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)),
					output: true,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("8.8.4.4"), 80)),
					output: false,
				},
				{
					input:  proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.ParseAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 80)),
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
				UserEmail: []string{
					"admin@v2ray.com",
				},
			},
			test: []ruleTest{
				{
					input:  protocol.ContextWithUser(context.Background(), &protocol.MemoryUser{Email: "admin@v2ray.com"}),
					output: true,
				},
				{
					input:  protocol.ContextWithUser(context.Background(), &protocol.MemoryUser{Email: "love@v2ray.com"}),
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
	}

	for _, test := range cases {
		cond, err := test.rule.BuildCondition()
		assert(err, IsNil)

		for _, t := range test.test {
			assert(cond.Apply(t.input), Equals, t.output)
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
	time.Sleep(time.Second * 10)
	for i := 0; i < 1024; i++ {
		assert(matcher.ApplyDomain(strconv.Itoa(i)+".not-exists2.com"), IsFalse)
	}
}
