package conf_test

import (
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/proto"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
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
