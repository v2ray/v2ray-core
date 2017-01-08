package scenarios

import (
	"net"
	"testing"

	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
)

func TestPassiveConnection(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
		SendFirst:    []byte("send first"),
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	serverPort := pickPort()
	serverConfig := &core.Config{
		Inbound: []*core.InboundConnectionConfig{
			{
				PortRange:              v2net.SinglePortRange(serverPort),
				ListenOn:               v2net.NewIPOrDomain(v2net.LocalHostIP),
				AllowPassiveConnection: true,
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
				Settings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	assert.Error(InitializeServerConfig(serverConfig)).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(serverPort),
	})
	assert.Error(err).IsNil()

	{
		response := make([]byte, 1024)
		nBytes, err := conn.Read(response)
		assert.Error(err).IsNil()
		assert.String(string(response[:nBytes])).Equals("send first")
	}

	payload := "dokodemo request."
	{

		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))
	}

	{
		response := make([]byte, 1024)
		nBytes, err := conn.Read(response)
		assert.Error(err).IsNil()
		assert.Bytes(response[:nBytes]).Equals(xor([]byte(payload)))
	}

	assert.Error(conn.Close()).IsNil()

	CloseAllServers()
}
