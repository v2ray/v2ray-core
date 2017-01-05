package scenarios

import (
	"net"
	"testing"

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
)

func TestDokodemoTCP(t *testing.T) {
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
			},
		},
		Outbound: []*core.OutboundConnectionConfig{
			{
				Settings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	clientPort := uint32(pickPort())
	clientPortRange := uint32(5)
	clientConfig := &core.Config{
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange: &v2net.PortRange{From: clientPort, To: clientPort + clientPortRange},
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
			},
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()
	assert.Error(InitializeServerConfig(clientConfig)).IsNil()

	for port := clientPort; port <= clientPort+clientPortRange; port++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(port),
		})

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
