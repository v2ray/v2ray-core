package scenarios

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/app/commander"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/servers/tcp"
	. "v2ray.com/ext/assert"
)

func TestCommanderRemoveHandler(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert(err, IsNil)
	defer tcpServer.Close()

	clientPort := pickPort()
	cmdPort := pickPort()
	clientConfig := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&commander.Config{Tag: "api"}),
			serial.ToTypedMessage(&router.Config{
				Rule: []*router.RoutingRule{
					{
						InboundTag: []string{"api"},
						Tag:        "api",
					},
				},
			}),
			serial.ToTypedMessage(&command.Config{}),
		},
		Inbound: []*core.InboundHandlerConfig{
			{
				Tag: "d",
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
			{
				Tag: "api",
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(cmdPort),
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
				Tag:           "default-outbound",
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	servers, err := InitializeServerConfigs(clientConfig)
	assert(err, IsNil)

	{
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(clientPort),
		})
		assert(err, IsNil)

		payload := "commander request."
		nBytes, err := conn.Write([]byte(payload))
		assert(err, IsNil)
		assert(nBytes, Equals, len(payload))

		response := make([]byte, 1024)
		nBytes, err = conn.Read(response)
		assert(err, IsNil)
		assert(response[:nBytes], Equals, xor([]byte(payload)))
		assert(conn.Close(), IsNil)
	}

	cmdConn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", cmdPort), grpc.WithInsecure())
	assert(err, IsNil)

	hsClient := command.NewHandlerServiceClient(cmdConn)
	resp, err := hsClient.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: "d",
	})
	assert(err, IsNil)
	assert(resp, IsNotNil)

	{
		_, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(clientPort),
		})
		assert(err, IsNotNil)
	}

	CloseAllServers(servers)
}
