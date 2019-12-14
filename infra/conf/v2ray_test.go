package conf_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	clog "v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	. "v2ray.com/core/infra/conf"
	"v2ray.com/core/proxy/blackhole"
	dns_proxy "v2ray.com/core/proxy/dns"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/http"
	"v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/websocket"
)

func TestV2RayConfig(t *testing.T) {
	createParser := func() func(string) (proto.Message, error) {
		return func(s string) (proto.Message, error) {
			config := new(Config)
			if err := json.Unmarshal([]byte(s), config); err != nil {
				return nil, err
			}
			return config.Build()
		}
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"outbound": {
					"protocol": "freedom",
					"settings": {}
				},
				"log": {
					"access": "/var/log/v2ray/access.log",
					"loglevel": "error",
					"error": "/var/log/v2ray/error.log"
				},
				"inbound": {
					"streamSettings": {
						"network": "ws",
						"wsSettings": {
							"headers": {
								"host": "example.domain"
							},
							"path": ""
						},
						"tlsSettings": {
							"alpn": "h2"
						},
						"security": "tls"
					},
					"protocol": "vmess",
					"port": 443,
					"settings": {
						"clients": [
							{
								"alterId": 100,
								"security": "aes-128-gcm",
								"id": "0cdf8a45-303d-4fed-9780-29aa7f54175e"
							}
						]
					}
				},
				"inbounds": [{
					"streamSettings": {
						"network": "ws",
						"wsSettings": {
							"headers": {
								"host": "example.domain"
							},
							"path": ""
						},
						"tlsSettings": {
							"alpn": "h2"
						},
						"security": "tls"
					},
					"protocol": "vmess",
					"port": "443-500",
					"allocate": {
						"strategy": "random",
						"concurrency": 3
					},
					"settings": {
						"clients": [
							{
								"alterId": 100,
								"security": "aes-128-gcm",
								"id": "0cdf8a45-303d-4fed-9780-29aa7f54175e"
							}
						]
					}
				}],
				"outboundDetour": [
					{
						"tag": "blocked",
						"protocol": "blackhole"
					},
					{
						"protocol": "dns"
					}
				],
				"routing": {
					"strategy": "rules",
					"settings": {
						"rules": [
							{
								"ip": [
									"10.0.0.0/8"
								],
								"type": "field",
								"outboundTag": "blocked"
							}
						]
					}
				},
				"transport": {
					"httpSettings": {
						"path": "/test"
					}
				}
			}`,
			Parser: createParser(),
			Output: &core.Config{
				App: []*serial.TypedMessage{
					serial.ToTypedMessage(&log.Config{
						ErrorLogType:  log.LogType_File,
						ErrorLogPath:  "/var/log/v2ray/error.log",
						ErrorLogLevel: clog.Severity_Error,
						AccessLogType: log.LogType_File,
						AccessLogPath: "/var/log/v2ray/access.log",
					}),
					serial.ToTypedMessage(&dispatcher.Config{}),
					serial.ToTypedMessage(&proxyman.InboundConfig{}),
					serial.ToTypedMessage(&proxyman.OutboundConfig{}),
					serial.ToTypedMessage(&router.Config{
						DomainStrategy: router.Config_AsIs,
						Rule: []*router.RoutingRule{
							{
								Geoip: []*router.GeoIP{
									{
										Cidr: []*router.CIDR{
											{
												Ip:     []byte{10, 0, 0, 0},
												Prefix: 8,
											},
										},
									},
								},
								TargetTag: &router.RoutingRule_Tag{
									Tag: "blocked",
								},
							},
						},
					}),
				},
				Outbound: []*core.OutboundHandlerConfig{
					{
						SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "tcp",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&freedom.Config{
							DomainStrategy: freedom.Config_AS_IS,
							UserLevel:      0,
						}),
					},
					{
						Tag: "blocked",
						SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "tcp",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&blackhole.Config{}),
					},
					{
						SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "tcp",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&dns_proxy.Config{
							Server: &net.Endpoint{},
						}),
					},
				},
				Inbound: []*core.InboundHandlerConfig{
					{
						ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
							PortRange: &net.PortRange{
								From: 443,
								To:   443,
							},
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "websocket",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "websocket",
										Settings: serial.ToTypedMessage(&websocket.Config{
											Header: []*websocket.Header{
												{
													Key:   "host",
													Value: "example.domain",
												},
											},
										}),
									},
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
								SecurityType: "v2ray.core.transport.internet.tls.Config",
								SecuritySettings: []*serial.TypedMessage{
									serial.ToTypedMessage(&tls.Config{
										NextProtocol: []string{"h2"},
									}),
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&inbound.Config{
							User: []*protocol.User{
								{
									Level: 0,
									Account: serial.ToTypedMessage(&vmess.Account{
										Id:      "0cdf8a45-303d-4fed-9780-29aa7f54175e",
										AlterId: 100,
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_AES128_GCM,
										},
									}),
								},
							},
						}),
					},
					{
						ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
							PortRange: &net.PortRange{
								From: 443,
								To:   500,
							},
							AllocationStrategy: &proxyman.AllocationStrategy{
								Type: proxyman.AllocationStrategy_Random,
								Concurrency: &proxyman.AllocationStrategy_AllocationStrategyConcurrency{
									Value: 3,
								},
							},
							StreamSettings: &internet.StreamConfig{
								ProtocolName: "websocket",
								TransportSettings: []*internet.TransportConfig{
									{
										ProtocolName: "websocket",
										Settings: serial.ToTypedMessage(&websocket.Config{
											Header: []*websocket.Header{
												{
													Key:   "host",
													Value: "example.domain",
												},
											},
										}),
									},
									{
										ProtocolName: "http",
										Settings: serial.ToTypedMessage(&http.Config{
											Path: "/test",
										}),
									},
								},
								SecurityType: "v2ray.core.transport.internet.tls.Config",
								SecuritySettings: []*serial.TypedMessage{
									serial.ToTypedMessage(&tls.Config{
										NextProtocol: []string{"h2"},
									}),
								},
							},
						}),
						ProxySettings: serial.ToTypedMessage(&inbound.Config{
							User: []*protocol.User{
								{
									Level: 0,
									Account: serial.ToTypedMessage(&vmess.Account{
										Id:      "0cdf8a45-303d-4fed-9780-29aa7f54175e",
										AlterId: 100,
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_AES128_GCM,
										},
									}),
								},
							},
						}),
					},
				},
			},
		},
	})
}

func TestMuxConfig_Build(t *testing.T) {
	tests := []struct {
		name   string
		fields string
		want   *proxyman.MultiplexingConfig
	}{
		{"default", `{"enabled": true, "concurrency": 16}`, &proxyman.MultiplexingConfig{
			Enabled:     true,
			Concurrency: 16,
		}},
		{"empty def", `{}`, &proxyman.MultiplexingConfig{
			Enabled:     false,
			Concurrency: 8,
		}},
		{"not enable", `{"enabled": false, "concurrency": 4}`, &proxyman.MultiplexingConfig{
			Enabled:     false,
			Concurrency: 4,
		}},
		{"forbidden", `{"enabled": false, "concurrency": -1}`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MuxConfig{}
			common.Must(json.Unmarshal([]byte(tt.fields), m))
			if got := m.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MuxConfig.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Override(t *testing.T) {
	tests := []struct {
		name string
		orig *Config
		over *Config
		fn   string
		want *Config
	}{
		{"combine/empty",
			&Config{},
			&Config{
				LogConfig:    &LogConfig{},
				RouterConfig: &RouterConfig{},
				DNSConfig:    &DnsConfig{},
				Transport:    &TransportConfig{},
				Policy:       &PolicyConfig{},
				Api:          &ApiConfig{},
				Stats:        &StatsConfig{},
				Reverse:      &ReverseConfig{},
			},
			"",
			&Config{
				LogConfig:    &LogConfig{},
				RouterConfig: &RouterConfig{},
				DNSConfig:    &DnsConfig{},
				Transport:    &TransportConfig{},
				Policy:       &PolicyConfig{},
				Api:          &ApiConfig{},
				Stats:        &StatsConfig{},
				Reverse:      &ReverseConfig{},
			},
		},
		{"combine/newattr",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "old"}}},
			&Config{LogConfig: &LogConfig{}}, "",
			&Config{LogConfig: &LogConfig{}, InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "old"}}}},
		{"replace/inbounds",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos0"}, InboundDetourConfig{Protocol: "vmess", Tag: "pos1"}}},
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}},
			"",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos0"}, InboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}}},
		{"replace/inbounds-all",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos0"}, InboundDetourConfig{Protocol: "vmess", Tag: "pos1"}}},
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos1", Protocol: "kcp"}, InboundDetourConfig{Tag: "pos2", Protocol: "kcp"}}},
			"",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos1", Protocol: "kcp"}, InboundDetourConfig{Tag: "pos2", Protocol: "kcp"}}}},
		{"replace/notag",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{}, InboundDetourConfig{Protocol: "vmess"}}},
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}},
			"",
			&Config{InboundConfigs: []InboundDetourConfig{InboundDetourConfig{}, InboundDetourConfig{Protocol: "vmess"}, InboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}}},
		{"replace/outbounds",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos0"}, OutboundDetourConfig{Protocol: "vmess", Tag: "pos1"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}},
			"",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos0"}, OutboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}}},
		{"replace/outbounds-append",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos0"}, OutboundDetourConfig{Protocol: "vmess", Tag: "pos1"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos1", Protocol: "kcp"}, OutboundDetourConfig{Tag: "pos2", Protocol: "kcp"}}},
			"",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos1", Protocol: "kcp"}, OutboundDetourConfig{Tag: "pos2", Protocol: "kcp"}}}},
		{"replace/outbounds-append",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos0"}, OutboundDetourConfig{Protocol: "vmess", Tag: "pos1"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos2", Protocol: "kcp"}}},
			"config_tail.json",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos0"}, OutboundDetourConfig{Protocol: "vmess", Tag: "pos1"}, OutboundDetourConfig{Tag: "pos2", Protocol: "kcp"}}}},
		{"replace/outbound-prepend",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{}, OutboundDetourConfig{Protocol: "vmess"}}},
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos1", Protocol: "kcp"}}},
			"config.json",
			&Config{OutboundConfigs: []OutboundDetourConfig{OutboundDetourConfig{Tag: "pos1", Protocol: "kcp"}, OutboundDetourConfig{}, OutboundDetourConfig{Protocol: "vmess"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.orig.Override(tt.over, tt.fn)
			if r := cmp.Diff(tt.orig, tt.want); r != "" {
				t.Error(r)
			}
		})
	}
}
