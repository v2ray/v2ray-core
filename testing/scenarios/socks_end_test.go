package scenarios

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/app/point"
	"github.com/v2ray/v2ray-core/app/point/config/testing/mocks"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	socksjson "github.com/v2ray/v2ray-core/proxy/socks/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
	vmessjson "github.com/v2ray/v2ray-core/proxy/vmess/config/json"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func setUpV2Ray() (uint16, error) {
	id1, err := config.NewID("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	if err != nil {
		return 0, err
	}
	id2, err := config.NewID("93ccfc71-b136-4015-ac85-e037bd1ead9e")
	if err != nil {
		return 0, err
	}
	users := []*vmessjson.ConfigUser{
		&vmessjson.ConfigUser{Id: id1},
		&vmessjson.ConfigUser{Id: id2},
	}

	portB := v2nettesting.PickPort()
	configB := mocks.Config{
		PortValue: portB,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &vmessjson.Inbound{
				AllowedClients: users,
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "freedom",
			SettingsValue: nil,
		},
	}
	pointB, err := point.NewPoint(&configB)
	if err != nil {
		return 0, err
	}
	err = pointB.Start()
	if err != nil {
		return 0, err
	}

	portA := v2nettesting.PickPort()
	configA := mocks.Config{
		PortValue: portA,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: &socksjson.SocksConfig{
				AuthMethod: "noauth",
				UDPEnabled: true,
				HostIP:     socksjson.IPAddress(net.IPv4(127, 0, 0, 1)),
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &vmessjson.Outbound{
				[]*vmessjson.ConfigTarget{
					&vmessjson.ConfigTarget{
						Address: v2net.IPAddress([]byte{127, 0, 0, 1}, portB),
						Users:   users,
					},
				},
			},
		},
	}

	pointA, err := point.NewPoint(&configA)
	if err != nil {
		return 0, err
	}
	err = pointA.Start()
	if err != nil {
		return 0, err
	}

	return portA, nil
}

func TestTCPConnection(t *testing.T) {
	assert := unit.Assert(t)

	targetPort := v2nettesting.PickPort()
	tcpServer := &tcp.Server{
		Port: targetPort,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	v2rayPort, err := setUpV2Ray()
	assert.Error(err).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(v2rayPort),
	})

	authRequest := socks5AuthMethodRequest(byte(0))
	nBytes, err := conn.Write(authRequest)
	assert.Int(nBytes).Equals(len(authRequest))
	assert.Error(err).IsNil()

	authResponse := make([]byte, 1024)
	nBytes, err = conn.Read(authResponse)
	assert.Error(err).IsNil()
	assert.Bytes(authResponse[:nBytes]).Equals([]byte{socks5Version, 0})

	connectRequest := socks5Request(byte(1), v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort))
	nBytes, err = conn.Write(connectRequest)
	assert.Int(nBytes).Equals(len(connectRequest))
	assert.Error(err).IsNil()

	connectResponse := make([]byte, 1024)
	nBytes, err = conn.Read(connectResponse)
	assert.Error(err).IsNil()
	assert.Bytes(connectResponse[:nBytes]).Equals([]byte{socks5Version, 0, 0, 1, 0, 0, 0, 0, 6, 181})

	actualRequest := []byte("Request to target server.")
	nBytes, err = conn.Write(actualRequest)
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(actualRequest))

	actualRequest = []byte("Request to target server again.")
	nBytes, err = conn.Write(actualRequest)
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(actualRequest))

	conn.CloseWrite()

	actualResponse := make([]byte, 1024)
	nBytes, err = conn.Read(actualResponse)
	assert.Error(err).IsNil()
	assert.String(string(actualResponse[:nBytes])).Equals("Processed: Request to target server.Request to target server again.")

	conn.Close()
}

func TestTCPBind(t *testing.T) {
	assert := unit.Assert(t)

	targetPort := v2nettesting.PickPort()
	tcpServer := &tcp.Server{
		Port: targetPort,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	v2rayPort, err := setUpV2Ray()
	assert.Error(err).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(v2rayPort),
	})

	authRequest := socks5AuthMethodRequest(byte(0))
	nBytes, err := conn.Write(authRequest)
	assert.Int(nBytes).Equals(len(authRequest))
	assert.Error(err).IsNil()

	authResponse := make([]byte, 1024)
	nBytes, err = conn.Read(authResponse)
	assert.Error(err).IsNil()
	assert.Bytes(authResponse[:nBytes]).Equals([]byte{socks5Version, 0})

	connectRequest := socks5Request(byte(2), v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort))
	nBytes, err = conn.Write(connectRequest)
	assert.Int(nBytes).Equals(len(connectRequest))
	assert.Error(err).IsNil()

	connectResponse := make([]byte, 1024)
	nBytes, err = conn.Read(connectResponse)
	assert.Error(err).IsNil()
	assert.Bytes(connectResponse[:nBytes]).Equals([]byte{socks5Version, 7, 0, 1, 0, 0, 0, 0, 0, 0})

	conn.Close()
}
