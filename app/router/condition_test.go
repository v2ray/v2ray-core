package router_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	proto "github.com/golang/protobuf/proto"

	. "v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/platform/filesystem"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/protocol/http"
	"v2ray.com/core/common/session"
)

func init() {
	wd, err := os.Getwd()
	common.Must(err)

	common.Must(filesystem.CopyFile(platform.GetAssetLocation("geoip.dat"), filepath.Join(wd, "..", "..", "release", "config", "geoip.dat")))
	common.Must(filesystem.CopyFile(platform.GetAssetLocation("geosite.dat"), filepath.Join(wd, "..", "..", "release", "config", "geosite.dat")))
}

func withOutbound(outbound *session.Outbound) *Context {
	return &Context{Outbound: outbound}
}

func withInbound(inbound *session.Inbound) *Context {
	return &Context{Inbound: inbound}
}

func TestRoutingRule(t *testing.T) {
	type ruleTest struct {
		input  *Context
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
					input:  &Context{},
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
					input:  &Context{},
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Geoip: []*GeoIP{
					{
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
					input:  &Context{},
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
					input:  &Context{},
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
					input:  &Context{Content: &session.Content{Protocol: (&http.SniffHeader{}).Protocol()}},
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
		{
			rule: &RoutingRule{
				PortList: &net.PortList{
					Range: []*net.PortRange{
						{From: 443, To: 443},
						{From: 1000, To: 1100},
					},
				},
			},
			test: []ruleTest{
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 443)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 1100)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 1005)}),
					output: true,
				},
				{
					input:  withOutbound(&session.Outbound{Target: net.TCPDestination(net.LocalHostIP, 53)}),
					output: false,
				},
			},
		},
		{
			rule: &RoutingRule{
				Protocol:   []string{"http"},
				Attributes: "attrs[':path'].startswith('/test')",
			},
			test: []ruleTest{
				{
					input:  &Context{Content: &session.Content{Protocol: "http/1.1", Attributes: map[string]interface{}{":path": "/test/1"}}},
					output: true,
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
	geositeBytes, err := filesystem.ReadAsset("geosite.dat")
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
	domains, err := loadGeoSite("CN")
	common.Must(err)

	matcher, err := NewDomainMatcher(domains)
	common.Must(err)

	type TestCase struct {
		Domain string
		Output bool
	}
	testCases := []TestCase{
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "163.com",
			Output: true,
		},
		{
			Domain: "164.com",
			Output: false,
		},
		{
			Domain: "164.com",
			Output: false,
		},
	}

	for i := 0; i < 1024; i++ {
		testCases = append(testCases, TestCase{Domain: strconv.Itoa(i) + ".not-exists.com", Output: false})
	}

	for _, testCase := range testCases {
		r := matcher.ApplyDomain(testCase.Domain)
		if r != testCase.Output {
			t.Error("expected output ", testCase.Output, " for domain ", testCase.Domain, " but got ", r)
		}
	}
}

func BenchmarkMultiGeoIPMatcher(b *testing.B) {
	var geoips []*GeoIP

	{
		ips, err := loadGeoIP("CN")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "CN",
			Cidr:        ips,
		})
	}

	{
		ips, err := loadGeoIP("JP")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "JP",
			Cidr:        ips,
		})
	}

	{
		ips, err := loadGeoIP("CA")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "CA",
			Cidr:        ips,
		})
	}

	{
		ips, err := loadGeoIP("US")
		common.Must(err)
		geoips = append(geoips, &GeoIP{
			CountryCode: "US",
			Cidr:        ips,
		})
	}

	matcher, err := NewMultiGeoIPMatcher(geoips, false)
	common.Must(err)

	ctx := withOutbound(&session.Outbound{Target: net.TCPDestination(net.ParseAddress("8.8.8.8"), 80)})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = matcher.Apply(ctx)
	}
}
