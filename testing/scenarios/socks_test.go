package scenarios

import (
	"testing"

	xproxy "golang.org/x/net/proxy"
	socks4 "h12.me/socks"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/proxy/socks"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/testing/servers/udp"
)

func TestSocksBridgeTCP(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
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
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
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
		Outbound: []*proxyman.OutboundHandlerConfig{
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
	assert.Error(err).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(clientPort),
	})
	assert.Error(err).IsNil()

	payload := "test payload"
	nBytes, err := conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert.Error(err).IsNil()
	assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
	assert.Error(conn.Close()).IsNil()

	CloseAllServers(servers)
}

func TestSocksBridageUDP(t *testing.T) {
	assert := assert.On(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert.Error(err).IsNil()
	defer udpServer.Close()

	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
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
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
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
		Outbound: []*proxyman.OutboundHandlerConfig{
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
	assert.Error(err).IsNil()

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(clientPort),
	})
	assert.Error(err).IsNil()

	payload := "dokodemo request."
	nBytes, err := conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert.Error(err).IsNil()
	assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
	assert.Error(conn.Close()).IsNil()

	CloseAllServers(servers)
}

func TestSocksConformance(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	authPort := pickPort()
	noAuthPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
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
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	servers, err := InitializeServerConfigs(serverConfig)
	assert.Error(err).IsNil()

	{
		noAuthDialer, err := xproxy.SOCKS5("tcp", net.TCPDestination(net.LocalHostIP, noAuthPort).NetAddr(), nil, xproxy.Direct)
		assert.Error(err).IsNil()
		conn, err := noAuthDialer.Dial("tcp", dest.NetAddr())
		assert.Error(err).IsNil()

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert.Error(err).IsNil()
		assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
		assert.Error(conn.Close()).IsNil()
	}

	{
		authDialer, err := xproxy.SOCKS5("tcp", net.TCPDestination(net.LocalHostIP, authPort).NetAddr(), &xproxy.Auth{User: "Test Account", Password: "Test Password"}, xproxy.Direct)
		assert.Error(err).IsNil()
		conn, err := authDialer.Dial("tcp", dest.NetAddr())
		assert.Error(err).IsNil()

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert.Error(err).IsNil()
		assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
		assert.Error(conn.Close()).IsNil()
	}

	{
		dialer := socks4.DialSocksProxy(socks4.SOCKS4, net.TCPDestination(net.LocalHostIP, noAuthPort).NetAddr())
		conn, err := dialer("tcp", dest.NetAddr())
		assert.Error(err).IsNil()

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert.Error(err).IsNil()
		assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
		assert.Error(conn.Close()).IsNil()
	}

	{
		dialer := socks4.DialSocksProxy(socks4.SOCKS4A, net.TCPDestination(net.LocalHostIP, noAuthPort).NetAddr())
		conn, err := dialer("tcp", net.TCPDestination(net.LocalHostDomain, tcpServer.Port).NetAddr())
		assert.Error(err).IsNil()

		payload := "test payload"
		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert.Error(err).IsNil()
		assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
		assert.Error(conn.Close()).IsNil()
	}

	CloseAllServers(servers)
}
