package unixtest

import (
	"testing"

	"v2ray.com/core"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	clog "v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	. "v2ray.com/ext/assert"

	"v2ray.com/core/testing/servers/tcp"

	"v2ray.com/core/testing/scenarios"

	"v2ray.com/core/transport/internet/domainsocket"
)

func xor(b []byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		r[i] = v ^ 'c'
	}
	return r
}

func pickPort() net.Port {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	common.Must(err)
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return net.Port(addr.Port)
}

func TestUnix2tcpS(t *testing.T) {
	assert := With(t)
	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	Server := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.UnixReceiverConfig{
					DomainSockSettings: &domainsocket.DomainSocketSettings{Path: "\x00v2raytest"},
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
			{ProxySettings: serial.ToTypedMessage(&freedom.Config{})},
		},
	}
	serverPort := pickPort()
	Client := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{
				ErrorLogLevel: clog.Severity_Debug,
				ErrorLogType:  log.LogType_Console,
			}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				SenderSettings: serial.ToTypedMessage(&proxyman.UnixSenderConfig{
					DomainSockSettings: &domainsocket.DomainSocketSettings{Path: "\x00v2raytest"},
				}),
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
		Inbound: []*core.InboundHandlerConfig{
			{ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
				Address: net.NewIPOrDomain(dest.Address),
				Port:    uint32(dest.Port),
				NetworkList: &net.NetworkList{
					Network: []net.Network{net.Network_TCP},
				},
			}),
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(serverPort),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
			},
		},
	}

	servers, err := scenarios.InitializeServerConfigs(Server, Client)

	port := serverPort

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(port),
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

	scenarios.CloseAllServers(servers)
}
