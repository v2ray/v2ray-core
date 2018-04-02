package scenarios

import (
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/buf"
	clog "v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/proxy/shadowsocks"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/testing/servers/udp"
	. "v2ray.com/ext/assert"

	ss "github.com/shadowsocks/go-shadowsocks2/core"
)

func TestShadowsocksAES256TCP(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "shadowsocks-password",
		CipherType: shadowsocks.CipherType_AES_256_CFB,
		Ota:        shadowsocks.Account_Enabled,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
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
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: account,
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert(err, IsNil)

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*20, 10240*1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksAES128UDP(t *testing.T) {
	assert := With(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert(err, IsNil)
	defer udpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "shadowsocks-password",
		CipherType: shadowsocks.CipherType_AES_128_CFB,
		Ota:        shadowsocks.Account_Enabled,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					UdpEnabled: true,
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
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
						Network: []net.Network{net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: account,
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert(err, IsNil)

			payload := make([]byte, 1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*5, 1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksChacha20TCP(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "shadowsocks-password",
		CipherType: shadowsocks.CipherType_CHACHA20_IETF,
		Ota:        shadowsocks.Account_Enabled,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
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
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: account,
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert(err, IsNil)

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*20, 10240*1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksAES256GCMTCP(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "shadowsocks-password",
		CipherType: shadowsocks.CipherType_AES_256_GCM,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
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
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: account,
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert(err, IsNil)

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*20, 10240*1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksAES128GCMUDP(t *testing.T) {
	assert := With(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert(err, IsNil)
	defer udpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "shadowsocks-password",
		CipherType: shadowsocks.CipherType_AES_128_GCM,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					UdpEnabled: true,
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
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
						Network: []net.Network{net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: account,
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert(err, IsNil)

			payload := make([]byte, 1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*5, 1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksAES128GCMUDPMux(t *testing.T) {
	assert := With(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert(err, IsNil)
	defer udpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "shadowsocks-password",
		CipherType: shadowsocks.CipherType_AES_128_GCM,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					UdpEnabled: false,
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
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
						Network: []net.Network{net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				SenderSettings: serial.ToTypedMessage(&proxyman.SenderConfig{
					MultiplexSettings: &proxyman.MultiplexingConfig{
						Enabled:     true,
						Concurrency: 8,
					},
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(serverPort),
							User: []*protocol.User{
								{
									Account: account,
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   []byte{127, 0, 0, 1},
				Port: int(clientPort),
			})
			assert(err, IsNil)

			payload := make([]byte, 1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*5, 1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksAES256GCMConformance(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "ss-password",
		CipherType: shadowsocks.CipherType_AES_256_GCM,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			cipher, err := ss.PickCipher("AES-256-GCM", nil, "ss-password")
			assert(err, IsNil)
			conn, err := ss.Dial("tcp", fmt.Sprintf(":%d", serverPort), cipher)
			assert(err, IsNil)
			_, err = conn.Write([]byte{1, 127, 0, 0, 1})
			assert(err, IsNil)
			_, err = conn.Write(serial.Uint16ToBytes(dest.Port.Value(), nil))
			assert(err, IsNil)

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*30, 10240*1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}

func TestShadowsocksChacha20Poly1305UDPConformance(t *testing.T) {
	assert := With(t)

	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	assert(err, IsNil)
	defer udpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "ss-password",
		CipherType: shadowsocks.CipherType_CHACHA20_POLY1305,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					UdpEnabled: true,
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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

	cipher, err := ss.PickCipher("CHACHA20-IETF-POLY1305", nil, "ss-password")
	assert(err, IsNil)
	conn, err := ss.ListenPacket("udp", ":0", cipher)
	assert(err, IsNil)

	for i := 0; i < 100; i++ {

		payload := buf.New()
		payload.AppendBytes(1, 127, 0, 0, 1)
		payload.AppendSupplier(serial.WriteUint16(dest.Port.Value()))

		payload.AppendSupplier(buf.ReadFullFrom(rand.Reader, 10))

		nBytes, err := conn.WriteTo(payload.Bytes(), &net.UDPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(serverPort),
		})
		assert(err, IsNil)
		assert(int32(nBytes), Equals, payload.Len())

		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		response := make([]byte, 10240)
		nBytes, _, err = conn.ReadFrom(response)
		assert(err, IsNil)
		assert(response[:7], Equals, payload.BytesTo(7))
		assert(response[7:nBytes], Equals, xor(payload.BytesFrom(7)))

	}

	assert(conn.Close(), IsNil)

	CloseAllServers(servers)
}

func TestShadowsocksChacha20Conformance(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	account := serial.ToTypedMessage(&shadowsocks.Account{
		Password:   "ss-password",
		CipherType: shadowsocks.CipherType_CHACHA20_IETF,
		Ota:        shadowsocks.Account_Disabled,
	})

	serverPort := tcp.PickPort()
	serverConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ServerConfig{
					User: &protocol.User{
						Account: account,
						Level:   1,
					},
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

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			cipher, err := ss.PickCipher("CHACHA20-IETF", nil, "ss-password")
			assert(err, IsNil)
			conn, err := ss.Dial("tcp", fmt.Sprintf(":%d", serverPort), cipher)
			assert(err, IsNil)
			_, err = conn.Write([]byte{1, 127, 0, 0, 1})
			assert(err, IsNil)
			_, err = conn.Write(serial.Uint16ToBytes(dest.Port.Value(), nil))
			assert(err, IsNil)

			payload := make([]byte, 10240*1024)
			rand.Read(payload)

			nBytes, err := conn.Write([]byte(payload))
			assert(err, IsNil)
			assert(nBytes, Equals, len(payload))

			response := readFrom(conn, time.Second*30, 10240*1024)
			assert(response, Equals, xor([]byte(payload)))
			assert(conn.Close(), IsNil)
			wg.Done()
		}()
	}
	wg.Wait()

	CloseAllServers(servers)
}
