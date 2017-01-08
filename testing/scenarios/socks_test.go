package scenarios

import (
	"net"
	"testing"

	xproxy "golang.org/x/net/proxy"
	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
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
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: v2net.SinglePortRange(serverPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_PASSWORD,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    v2net.NewIPOrDomain(v2net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
		},
		Outbound: []*core.OutboundConnectionConfig{
			{
				Settings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: v2net.SinglePortRange(clientPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundConnectionConfig{
			{
				Settings: serial.ToTypedMessage(&socks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
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

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

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

	CloseAllServers()
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
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: v2net.SinglePortRange(serverPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_PASSWORD,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    v2net.NewIPOrDomain(v2net.LocalHostIP),
					UdpEnabled: true,
				}),
			},
		},
		Outbound: []*core.OutboundConnectionConfig{
			{
				Settings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: v2net.SinglePortRange(clientPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP, v2net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundConnectionConfig{
			{
				Settings: serial.ToTypedMessage(&socks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
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

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

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

	CloseAllServers()
}

func TestSocks5conformance(t *testing.T) {
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
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: v2net.SinglePortRange(authPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_PASSWORD,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    v2net.NewIPOrDomain(v2net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
			{
				PortRange: v2net.SinglePortRange(noAuthPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&socks.ServerConfig{
					AuthType: socks.AuthType_NO_AUTH,
					Accounts: map[string]string{
						"Test Account": "Test Password",
					},
					Address:    v2net.NewIPOrDomain(v2net.LocalHostIP),
					UdpEnabled: false,
				}),
			},
		},
		Outbound: []*core.OutboundConnectionConfig{
			{
				Settings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()

	{
		noAuthDialer, err := xproxy.SOCKS5("tcp", v2net.TCPDestination(v2net.LocalHostIP, noAuthPort).NetAddr(), nil, xproxy.Direct)
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
		authDialer, err := xproxy.SOCKS5("tcp", v2net.TCPDestination(v2net.LocalHostIP, noAuthPort).NetAddr(), &xproxy.Auth{User: "Test Account", Password: "Test Password"}, xproxy.Direct)
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

	CloseAllServers()
}
