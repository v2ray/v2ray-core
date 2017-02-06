package scenarios

import (
	"crypto/rand"
	"net"
	"testing"

	"time"

	"sync"

	"v2ray.com/core"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/inbound"
	"v2ray.com/core/proxy/vmess/outbound"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
)

func TestVMessDynamicPort(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(serverPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id: userID.String(),
							}),
						},
					},
					Detour: &inbound.DetourConfig{
						To: "detour",
					},
				}),
			},
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: &v2net.PortRange{
						From: uint32(serverPort + 1),
						To:   uint32(serverPort + 100),
					},
					Listen: v2net.NewIPOrDomain(v2net.LocalHostIP),
					AllocationStrategy: &proxyman.AllocationStrategy{
						Type: proxyman.AllocationStrategy_Random,
						Concurrency: &proxyman.AllocationStrategy_AllocationStrategyConcurrency{
							Value: 2,
						},
						Refresh: &proxyman.AllocationStrategy_AllocationStrategyRefresh{
							Value: 5,
						},
					},
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{}),
				Tag:           "detour",
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(clientPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: userID.String(),
									}),
								},
							},
						},
					},
				}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

	for i := 0; i < 10; i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
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
	}

	CloseAllServers()
}

func TestVMessGCM(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(serverPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id:      userID.String(),
								AlterId: 64,
							}),
						},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(clientPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id:      userID.String(),
										AlterId: 64,
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_AES128_GCM,
										},
									}),
								},
							},
						},
					},
				}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert.Error(err).IsNil()

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert.Error(err).IsNil()
			assert.Int(nBytes).Equals(len(payload))

			response := readFrom(conn, time.Second*10, 10240*1024)
			assert.Bytes(response).Equals(xor([]byte(payload)))
			assert.Error(conn.Close()).IsNil()
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers()
}

func TestVMessChacha20(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(serverPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id:      userID.String(),
								AlterId: 64,
							}),
						},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(clientPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id:      userID.String(),
										AlterId: 64,
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_CHACHA20_POLY1305,
										},
									}),
								},
							},
						},
					},
				}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert.Error(err).IsNil()

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert.Error(err).IsNil()
			assert.Int(nBytes).Equals(len(payload))

			response := readFrom(conn, time.Second*10, 10240*1024)
			assert.Bytes(response).Equals(xor([]byte(payload)))
			assert.Error(conn.Close()).IsNil()
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers()
}

func TestVMessNone(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(serverPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id:      userID.String(),
								AlterId: 64,
							}),
						},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(clientPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id:      userID.String(),
										AlterId: 64,
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_NONE,
										},
									}),
								},
							},
						},
					},
				}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert.Error(err).IsNil()

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert.Error(err).IsNil()
			assert.Int(nBytes).Equals(len(payload))

			response := readFrom(conn, time.Second*10, 10240*1024)
			assert.Bytes(response).Equals(xor([]byte(payload)))
			assert.Error(conn.Close()).IsNil()
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers()
}

func TestVMessKCP(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	userID := protocol.NewID(uuid.New())
	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(serverPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
					StreamSettings: &internet.StreamConfig{
						Protocol: internet.TransportProtocol_MKCP,
					},
				}),
				ProxySettings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id:      userID.String(),
								AlterId: 64,
							}),
						},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	clientPort := pickPort()
	clientConfig := &core.Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(clientPort),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(dest.Address),
					Port:    uint32(dest.Port),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id:      userID.String(),
										AlterId: 64,
										SecuritySettings: &protocol.SecurityConfig{
											Type: protocol.SecurityType_AES128_GCM,
										},
									}),
								},
							},
						},
					},
				}),
				SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
					StreamSettings: &internet.StreamConfig{
						Protocol: internet.TransportProtocol_MKCP,
					},
				}),
			},
		},
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: log.LogLevel_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert.Error(err).IsNil()

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert.Error(err).IsNil()
			assert.Int(nBytes).Equals(len(payload))

			response := readFrom(conn, time.Second*10, 10240*1024)
			assert.Bytes(response).Equals(xor([]byte(payload)))
			assert.Error(conn.Close()).IsNil()
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers()
}
