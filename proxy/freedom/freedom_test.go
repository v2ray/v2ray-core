package freedom

import (
	"bytes"
	"io/ioutil"
	"net"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/common/net"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	"github.com/v2ray/v2ray-core/testing/mocks"
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

	ich := &mocks.InboundConnectionHandler{
		Data2Send:    []byte("Not Used"),
		DataReturned: bytes.NewBuffer(make([]byte, 0, 1024)),
	}

	core.RegisterInboundConnectionHandlerFactory("mock_ich", ich)

	pointPort := uint16(38724)
	config := mocks.Config{
		PortValue: pointPort,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_ich",
			ContentValue:  nil,
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "freedom",
			ContentValue:  nil,
		},
	}

	point, err := core.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	dest := v2net.NewUDPDestination(udpServerAddr)
	ich.Communicate(v2net.NewPacket(dest, []byte(data2Send), false))
	assert.Bytes(ich.DataReturned.Bytes()).Equals([]byte("Processed: Data to be sent to remote"))
}

func TestSocksTcpConnect(t *testing.T) {
	assert := unit.Assert(t)
	port := uint16(38293)

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

	pointPort := uint16(38724)
	config := mocks.Config{
		PortValue: pointPort,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			ContentValue:  []byte("{\"auth\": \"noauth\"}"),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "freedom",
			ContentValue:  nil,
		},
	}

	point, err := core.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", "127.0.0.1:38724", nil, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "127.0.0.1:38293"
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
