package scenarios

import (
	"testing"

	"v2ray.com/core/app/router"

	xproxy "golang.org/x/net/proxy"
	socks4 "h12.io/socks"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/blackhole"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/proxy/socks"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/testing/servers/udp"
	. "v2ray.com/ext/assert"
)

func TestSocksBridgeTCP(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_PASSWORD,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    net.NewIPOrDomain(net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := tcp.PickPort()
	clientConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(clientPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&socks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&socks.Account{
										Username: "Test Account",
										Password: "Test Password",
									}),
								},
							},
						},
					},
				}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig, clientConfig)
	assert(err, IsNil)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(clientPort),
	})
	assert(err, IsNil)

	payload := "test payload"
	nBytes, err := conn.Write([]byte(payload))
	assert(err, IsNil)
	assert(nBytes, Equals, len(payload))

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert(err, IsNil)
	assert(response[:nBytes], Equals, xor([]byte(payload)))
	assert(conn.Close(), IsNil)

	CloseAllServers(servers)
}

func TestSocksBridageUDP(t *testing.T) {
	assert := With(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert(err, IsNil)
	defer udpServer.Close()

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_PASSWORD,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    net.NewIPOrDomain(net.LocalHostIP),
					UdpEnabled: true,
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := tcp.PickPort()
	clientConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(clientPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP, net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&socks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&socks.Account{
										Username: "Test Account",
										Password: "Test Password",
									}),
								},
							},
						},
					},
				}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig, clientConfig)
	assert(err, IsNil)

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(clientPort),
	})
	assert(err, IsNil)

	payload := "dokodemo request."
	nBytes, err := conn.Write([]byte(payload))
	assert(err, IsNil)
	assert(nBytes, Equals, len(payload))

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert(err, IsNil)
	assert(response[:nBytes], Equals, xor([]byte(payload)))
	assert(conn.Close(), IsNil)

	CloseAllServers(servers)
}

func TestSocksBridageUDPWithRouting(t *testing.T) {
	assert := With(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert(err, IsNil)
	defer udpServer.Close()

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&router.Config{
				Rule: []*router.RoutingRule{
					{
						Tag:        "out",
						InboundTag: []string{"socks"},
					},
				},
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				Tag: "socks",
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType:   socks.AuthType_NO_AUTH,
					Address:    net.NewIPOrDomain(net.LocalHostIP),
					UdpEnabled: true,
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&blackhole.Config{}),
			},
			{
				Tag:           "out",
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := tcp.PickPort()
	clientConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(clientPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP, net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&socks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
						},
					},
				}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig, clientConfig)
	assert(err, IsNil)

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(clientPort),
	})
	assert(err, IsNil)

	payload := "dokodemo request."
	nBytes, err := conn.Write([]byte(payload))
	assert(err, IsNil)
	assert(nBytes, Equals, len(payload))

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert(err, IsNil)
	assert(response[:nBytes], Equals, xor([]byte(payload)))
	assert(conn.Close(), IsNil)

	CloseAllServers(servers)
}

func TestSocksConformance(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	authPort := tcp.PickPort()
	noAuthPort := tcp.PickPort()
	serverConfig := &core.Config{
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(authPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_PASSWORD,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    net.NewIPOrDomain(net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(noAuthPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_NO_AUTH,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    net.NewIPOrDomain(net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig)
	assert(err, IsNil)

	{
		noAuthDialer, err := xproxy.SOCKS5("tcp", net.TCPDestination(net.LocalHostIP, noAuthPort).NetAddr(), nil, xproxy.Direct)
		assert(err, IsNil)
		conn, err := noAuthDialer.Dial("tcp", dest.NetAddr())
		assert(err, IsNil)

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert(err, IsNil)
		assert(nBytes, Equals, len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert(err, IsNil)
		assert(response[:nBytes], Equals, xor([]byte(payload)))
		assert(conn.Close(), IsNil)
	}

	{
		authDialer, err := xproxy.SOCKS5("tcp", net.TCPDestination(net.LocalHostIP, authPort).NetAddr(), &xproxy.Auth{User: "Test Account", Password: "Test Password"}, xproxy.Direct)
		assert(err, IsNil)
		conn, err := authDialer.Dial("tcp", dest.NetAddr())
		assert(err, IsNil)

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert(err, IsNil)
		assert(nBytes, Equals, len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert(err, IsNil)
		assert(response[:nBytes], Equals, xor([]byte(payload)))
		assert(conn.Close(), IsNil)
	}

	{
		dialer := socks4.DialSocksProxy(socks4.SOCKS4, net.TCPDestination(net.LocalHostIP, noAuthPort).NetAddr())
		conn, err := dialer("tcp", dest.NetAddr())
		assert(err, IsNil)

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert(err, IsNil)
		assert(nBytes, Equals, len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert(err, IsNil)
		assert(response[:nBytes], Equals, xor([]byte(payload)))
		assert(conn.Close(), IsNil)
	}

	{
		dialer := socks4.DialSocksProxy(socks4.SOCKS4A, net.TCPDestination(net.LocalHostIP, noAuthPort).NetAddr())
		conn, err := dialer("tcp", net.TCPDestination(net.LocalHostDomain, tcpServer.Port).NetAddr())
		assert(err, IsNil)

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert(err, IsNil)
		assert(nBytes, Equals, len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert(err, IsNil)
		assert(response[:nBytes], Equals, xor([]byte(payload)))
		assert(conn.Close(), IsNil)
	}

	CloseAllServers(servers)
}
