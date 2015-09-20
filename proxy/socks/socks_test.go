package socks

import (
	"bytes"
	"io/ioutil"
	"net"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestSocksTcpConnect(t *testing.T) {
	assert := unit.Assert(t)
	port := uint16(12385)

	och := &mocks.OutboundConnectionHandler{
		Data2Send:   bytes.NewBuffer(make([]byte, 0, 1024)),
		Data2Return: []byte("The data to be returned to socks server."),
	}

	core.RegisterOutboundConnectionHandlerFactory("mock_och", och)

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			ContentValue:  []byte("{\"auth\": \"noauth\"}"),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_och",
			ContentValue:  nil,
		},
	}

	point, err := core.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", "127.0.0.1:12385", nil, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "google.com:80"
	conn, err := socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNil()

	data2Send := "The data to be sent to remote server."
	conn.Write([]byte(data2Send))
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	dataReturned, err := ioutil.ReadAll(conn)
	assert.Error(err).IsNil()
	conn.Close()

	assert.Bytes([]byte(data2Send)).Equals(och.Data2Send.Bytes())
	assert.Bytes(dataReturned).Equals(och.Data2Return)
	assert.String(targetServer).Equals(och.Destination.Address().String())
}

func TestSocksTcpConnectWithUserPass(t *testing.T) {
	assert := unit.Assert(t)
	port := uint16(12386)

	och := &mocks.OutboundConnectionHandler{
		Data2Send:   bytes.NewBuffer(make([]byte, 0, 1024)),
		Data2Return: []byte("The data to be returned to socks server."),
	}

	core.RegisterOutboundConnectionHandlerFactory("mock_och", och)

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			ContentValue:  []byte("{\"auth\": \"password\",\"user\": \"userx\",\"pass\": \"passy\"}"),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_och",
			ContentValue:  nil,
		},
	}

	point, err := core.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", "127.0.0.1:12386", &proxy.Auth{"userx", "passy"}, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "1.2.3.4:443"
	conn, err := socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNil()

	data2Send := "The data to be sent to remote server."
	conn.Write([]byte(data2Send))
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	dataReturned, err := ioutil.ReadAll(conn)
	assert.Error(err).IsNil()
	conn.Close()

	assert.Bytes([]byte(data2Send)).Equals(och.Data2Send.Bytes())
	assert.Bytes(dataReturned).Equals(och.Data2Return)
	assert.String(targetServer).Equals(och.Destination.Address().String())
}
