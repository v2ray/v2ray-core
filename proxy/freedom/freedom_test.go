package freedom

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/v2ray/v2ray-core/app/point"
	"github.com/v2ray/v2ray-core/app/point/config/testing/mocks"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	"github.com/v2ray/v2ray-core/proxy/socks/config/json"
	proxymocks "github.com/v2ray/v2ray-core/proxy/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	"github.com/v2ray/v2ray-core/testing/servers/udp"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestUDPSend(t *testing.T) {
	assert := unit.Assert(t)

	data2Send := "Data to be sent to remote"

	udpServer := &udp.Server{
		Port: 0,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}

	udpServerAddr, err := udpServer.Start()
	assert.Error(err).IsNil()

	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	ich := &proxymocks.InboundConnectionHandler{
		ConnInput:  bytes.NewReader([]byte("Not Used")),
		ConnOutput: connOutput,
	}

	connhandler.RegisterInboundConnectionHandlerFactory("mock_ich", ich)

	pointPort := v2nettesting.PickPort()
	config := mocks.Config{
		PortValue: pointPort,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_ich",
			SettingsValue: nil,
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "freedom",
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	data2SendBuffer := alloc.NewBuffer().Clear()
	data2SendBuffer.Append([]byte(data2Send))
	dest := v2net.NewUDPDestination(udpServerAddr)
	ich.Communicate(v2net.NewPacket(dest, data2SendBuffer, false))
	assert.Bytes(connOutput.Bytes()).Equals([]byte("Processed: Data to be sent to remote"))
}

func TestSocksTcpConnect(t *testing.T) {
	assert := unit.Assert(t)
	port := v2nettesting.PickPort()

	data2Send := "Data to be sent to remote"

	tcpServer := &tcp.Server{
		Port: port,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	pointPort := v2nettesting.PickPort()
	config := mocks.Config{
		PortValue: pointPort,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: &json.SocksConfig{
				AuthMethod: "auth",
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "freedom",
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", pointPort), nil, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNil()

	conn.Write([]byte(data2Send))
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	dataReturned, err := ioutil.ReadAll(conn)
	assert.Error(err).IsNil()
	conn.Close()

	assert.Bytes(dataReturned).Equals([]byte("Processed: Data to be sent to remote"))
}
