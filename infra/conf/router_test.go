package conf_test

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core/app/router"
	. "v2ray.com/core/infra/conf"
)

func TestRouterConfig(t *testing.T) {
	createParser := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(RouterConfig)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			return config.Build()
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"strategy": "rules",
				"settings": {
					"domainStrategy": "AsIs",
					"rules": [
						{
							"type": "field",
							"domain": [
								"baidu.com",
								"qq.com"
							],
							"outboundTag": "direct"
						},
						{
							"type": "field",
							"ip": [
								"10.0.0.0/8",
								"::1/128"
							],
							"outboundTag": "test"
						}
					]
				},
				"balancers": [
					{
						"tag": "b1",
						"selector": ["test"]
					}
				]
			}`,
			Parser: createParser(),
			Output: &router.Config{
				DomainStrategy: router.Config_AsIs,
				BalancingRule: []*router.BalancingRule{
					{
						Tag:              "b1",
						OutboundSelector: []string{"test"},
					},
				},
				Rule: []*router.RoutingRule{
					{
						Domain: []*router.Domain{
							{
								Type:  router.Domain_Plain,
								Value: "baidu.com",
							},
							{
								Type:  router.Domain_Plain,
								Value: "qq.com",
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "direct",
						},
					},
					{
						Geoip: []*router.GeoIP{
							{
								Cidr: []*router.CIDR{
									{
										Ip:     []byte{10, 0, 0, 0},
										Prefix: 8,
									},
									{
										Ip:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
										Prefix: 128,
									},
								},
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "test",
						},
					},
				},
			},
		},
		{
			Input: `{
				"strategy": "rules",
				"settings": {
					"domainStrategy": "IPIfNonMatch",
					"rules": [
						{
							"type": "field",
							"domain": [
								"baidu.com",
								"qq.com"
							],
							"outboundTag": "direct"
						},
						{
							"type": "field",
							"ip": [
								"10.0.0.0/8",
								"::1/128"
							],
							"outboundTag": "test"
						}
					]
				}
			}`,
			Parser: createParser(),
			Output: &router.Config{
				DomainStrategy: router.Config_IpIfNonMatch,
				Rule: []*router.RoutingRule{
					{
						Domain: []*router.Domain{
							{
								Type:  router.Domain_Plain,
								Value: "baidu.com",
							},
							{
								Type:  router.Domain_Plain,
								Value: "qq.com",
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "direct",
						},
					},
					{
						Geoip: []*router.GeoIP{
							{
								Cidr: []*router.CIDR{
									{
										Ip:     []byte{10, 0, 0, 0},
										Prefix: 8,
									},
									{
										Ip:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
										Prefix: 128,
									},
								},
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "test",
						},
					},
				},
			},
		},
		{
			Input: `{
				"domainStrategy": "AsIs",
				"rules": [
					{
						"type": "field",
						"domain": [
							"baidu.com",
							"qq.com"
						],
						"outboundTag": "direct"
					},
					{
						"type": "field",
						"ip": [
							"10.0.0.0/8",
							"::1/128"
						],
						"outboundTag": "test"
					}
				]
			}`,
			Parser: createParser(),
			Output: &router.Config{
				DomainStrategy: router.Config_AsIs,
				Rule: []*router.RoutingRule{
					{
						Domain: []*router.Domain{
							{
								Type:  router.Domain_Plain,
								Value: "baidu.com",
							},
							{
								Type:  router.Domain_Plain,
								Value: "qq.com",
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "direct",
						},
					},
					{
						Geoip: []*router.GeoIP{
							{
								Cidr: []*router.CIDR{
									{
										Ip:     []byte{10, 0, 0, 0},
										Prefix: 8,
									},
									{
										Ip:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
										Prefix: 128,
									},
								},
							},
						},
						TargetTag: &router.RoutingRule_Tag{
							Tag: "test",
						},
					},
				},
			},
		},
	})
}
