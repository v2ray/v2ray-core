package scenarios

import (
	"net"
	"testing"
	"time"

	"v2ray.com/core"
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
	"v2ray.com/core/transport/internet/headers/http"
	tcptransport "v2ray.com/core/transport/internet/tcp"
)

func TestNoOpConnectionHeader(t *testing.T) {
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
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: v2net.SinglePortRange(serverPort),
				ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
				Settings: serial.ToTypedMessage(&inbound.Config{
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&vmess.Account{
								Id: userID.String(),
							}),
						},
					},
				}),
				StreamSettings: &internet.StreamConfig{
					TransportSettings: []*internet.TransportConfig{
						{
							Protocol: internet.TransportProtocol_TCP,
							Settings: serial.ToTypedMessage(&tcptransport.Config{
								HeaderSettings: serial.ToTypedMessage(&http.Config{}),
							}),
						},
					},
				},
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
				Settings: serial.ToTypedMessage(&outbound.Config{
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
				StreamSettings: &internet.StreamConfig{
					TransportSettings: []*internet.TransportConfig{
						{
							Protocol: internet.TransportProtocol_TCP,
							Settings: serial.ToTypedMessage(&tcptransport.Config{
								HeaderSettings: serial.ToTypedMessage(&http.Config{}),
							}),
						},
					},
				},
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

	payload := "dokodemo request."
	nBytes, err := conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	response := readFrom(conn, time.Second*2, len(payload))
	assert.Bytes(response).Equals(xor([]byte(payload)))
	assert.Error(conn.Close()).IsNil()

	CloseAllServers()
}
