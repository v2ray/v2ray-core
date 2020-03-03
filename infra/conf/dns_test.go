package conf_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/platform/filesystem"
	. "v2ray.com/core/infra/conf"
)

func init() {
	wd, err := os.Getwd()
	common.Must(err)

	common.Must(filesystem.CopyFile(platform.GetAssetLocation("geoip.dat"), filepath.Join(wd, "..", "..", "release", "config", "geoip.dat")))

	geositeFilePath := platform.GetAssetLocation("geosite.dat")
	geositeFile, err := os.OpenFile(geositeFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	common.Must(err)
	defer geositeFile.Close()

	list := &router.GeoSiteList{
		Entry: []*router.GeoSite{
			{
				CountryCode: "TEST",
				Domain: []*router.Domain{
					{Type: router.Domain_Full, Value: "example.com"},
				},
			},
		},
	}

	listBytes, err := proto.Marshal(list)
	common.Must(err)
	common.Must2(geositeFile.Write(listBytes))
}

type hostRulesWarpper struct {
	mappings []*dns.Config_HostMapping
}

func (hrw hostRulesWarpper) Len() int {
	return len(hrw.mappings)
}
func (hrw hostRulesWarpper) Swap(i, j int) {
	hrw.mappings[i], hrw.mappings[j] = hrw.mappings[j], hrw.mappings[i]
}
func (hrw hostRulesWarpper) Less(i, j int) bool {
	return hrw.mappings[i].Domain < hrw.mappings[j].Domain
}

func TestDnsConfigParsing(t *testing.T) {
	geositePath := platform.GetAssetLocation("geosite.dat")
	defer func() {
		os.Remove(geositePath)
	}()

	parserCreator := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(DnsConfig)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			con, err := config.Build()
			sort.Sort(hostRulesWarpper{con.StaticHosts})
			return con, err
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"servers": [{
					"address": "8.8.8.8",
					"port": 5353,
					"domains": ["domain:v2ray.com"]
				}],
				"hosts": {
					"v2ray.com": "127.0.0.1",
					"domain:example.com": "google.com",
					"geosite:test": "10.0.0.1",
					"keyword:google": "8.8.8.8",
					"regexp:.*\\.com": "8.8.4.4"
				},
				"fake": {
					"fakeRules": ["geosite:test"]
				},
				"clientIp": "10.0.0.1"
			}`,
			Parser: parserCreator(),
			Output: &dns.Config{
				NameServer: []*dns.NameServer{
					{
						Address: &net.Endpoint{
							Address: &net.IPOrDomain{
								Address: &net.IPOrDomain_Ip{
									Ip: []byte{8, 8, 8, 8},
								},
							},
							Network: net.Network_UDP,
							Port:    5353,
						},
						PrioritizedDomain: []*dns.NameServer_PriorityDomain{
							{
								Type:   dns.DomainMatchingType_New,
								Domain: "dv2ray.com",
							},
						},
					},
				},
				StaticHosts: []*dns.Config_HostMapping{
					{
						Type:          dns.DomainMatchingType_New,
						Domain:        "dexample.com",
						ProxiedDomain: "google.com",
					},
					{
						Type:   dns.DomainMatchingType_New,
						Domain: "egeosite.dat:test",
						Ip:     [][]byte{{10, 0, 0, 1}},
					},
					{
						Type:   dns.DomainMatchingType_New,
						Domain: "fv2ray.com",
						Ip:     [][]byte{{127, 0, 0, 1}},
					},
					{
						Type:   dns.DomainMatchingType_New,
						Domain: "kgoogle",
						Ip:     [][]byte{{8, 8, 8, 8}},
					},
					{
						Type:   dns.DomainMatchingType_New,
						Domain: "r.*\\.com",
						Ip:     [][]byte{{8, 8, 4, 4}},
					},
				},
				ExternalRules: map[string]*dns.ConfigPatterns{
					"geosite.dat:test": {
						Patterns: []string{
							"fexample.com",
						},
					},
				},
				Fake: &dns.Config_Fake{
					FakeNet:      "224.0.0.0/22",
					FakeRules:    []string{"egeosite.dat:test"},
					Regeneration: dns.Config_Fake_LRU,
					Path:         "",
				},
				ClientIp: []byte{10, 0, 0, 1},
			},
		},
	})
}
