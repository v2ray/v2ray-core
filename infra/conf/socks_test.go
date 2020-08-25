package conf_test

import (
	"testing"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	. "v2ray.com/core/infra/conf"
	"v2ray.com/core/proxy/socks"
)

func TestSocksInboundConfig(t *testing.T) {
	creator := func() Buildable {
		return new(SocksServerConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"auth": "password",
				"accounts": [
					{
						"user": "my-username",
						"pass": "my-password"
					}
				],
				"udp": false,
				"ip": "127.0.0.1",
				"timeout": 5,
				"userLevel": 1
			}`,
			Parser: loadJSON(creator),
			Output: &socks.ServerConfig{
				AuthType: socks.AuthType_PASSWORD,
				Accounts: map[string]string{
					"my-username": "my-password",
				},
				UdpEnabled: false,
				Address: &net.IPOrDomain{
					Address: &net.IPOrDomain_Ip{
						Ip: []byte{127, 0, 0, 1},
					},
				},
				Timeout:   5,
				UserLevel: 1,
			},
		},
	})
}

func TestSocksOutboundConfig(t *testing.T) {
	creator := func() Buildable {
		return new(SocksClientConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"servers": [{
					"address": "127.0.0.1",
					"port": 1234,
					"users": [
						{"user": "test user", "pass": "test pass", "email": "test@email.com"}
					]
				}]
			}`,
			Parser: loadJSON(creator),
			Output: &socks.ClientConfig{
				Server: []*protocol.ServerEndpoint{
					{
						Address: &net.IPOrDomain{
							Address: &net.IPOrDomain_Ip{
								Ip: []byte{127, 0, 0, 1},
							},
						},
						Port: 1234,
						User: []*protocol.User{
							{
								Email: "test@email.com",
								Account: serial.ToTypedMessage(&socks.Account{
									Username: "test user",
									Password: "test pass",
								}),
							},
						},
					},
				},
			},
		},
	})
}
