package socks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/v2ray/v2ray-core/app"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	v2proxy "github.com/v2ray/v2ray-core/proxy"
	_ "github.com/v2ray/v2ray-core/proxy/socks/json"
	proxytesting "github.com/v2ray/v2ray-core/proxy/testing"
	proxymocks "github.com/v2ray/v2ray-core/proxy/testing/mocks"
	"github.com/v2ray/v2ray-core/shell/point"
	"github.com/v2ray/v2ray-core/shell/point/testing/mocks"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSocksTcpConnect(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnOutput: connOutput,
		ConnInput:  bytes.NewReader(connInput),
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och", func(space app.Space, config interface{}) (v2proxy.OutboundConnectionHandler, error) {
		return och, nil
	})
	assert.Error(err).IsNil()

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: []byte(`
      {
        "auth": "noauth"
      }`),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: protocol,
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), nil, proxy.Direct)
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

	assert.Bytes([]byte(data2Send)).Equals(connOutput.Bytes())
	assert.Bytes(dataReturned).Equals(connInput)
	assert.StringLiteral(targetServer).Equals(och.Destination.NetAddr())
}

func TestSocksTcpConnectWithUserPass(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och", func(space app.Space, config interface{}) (v2proxy.OutboundConnectionHandler, error) {
		return och, nil
	})
	assert.Error(err).IsNil()

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: []byte(`
      {
        "auth": "password",
        "accounts": [
          {"user": "userx", "pass": "passy"}
        ]
      }`),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: protocol,
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), &proxy.Auth{"userx", "passy"}, proxy.Direct)
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

	assert.Bytes([]byte(data2Send)).Equals(connOutput.Bytes())
	assert.Bytes(dataReturned).Equals(connInput)
	assert.StringLiteral(targetServer).Equals(och.Destination.NetAddr())
}

func TestSocksTcpConnectWithWrongUserPass(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och", func(space app.Space, config interface{}) (v2proxy.OutboundConnectionHandler, error) {
		return och, nil
	})
	assert.Error(err).IsNil()

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: []byte(`
      {
        "auth": "password",
        "accounts": [
          {"user": "userx", "pass": "passy"}
        ]
      }`),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: protocol,
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), &proxy.Auth{"userx", "passz"}, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "1.2.3.4:443"
	_, err = socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNotNil()
}

func TestSocksTcpConnectWithWrongAuthMethod(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och", func(space app.Space, config interface{}) (v2proxy.OutboundConnectionHandler, error) {
		return och, nil
	})
	assert.Error(err).IsNil()

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: []byte(`
      {
        "auth": "password",
        "accounts": [
          {"user": "userx", "pass": "passy"}
        ]
      }`),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: protocol,
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), nil, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "1.2.3.4:443"
	_, err = socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNotNil()
}

func TestSocksUdpSend(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och",
		func(space app.Space, config interface{}) (v2proxy.OutboundConnectionHandler, error) {
			return och, nil
		})
	assert.Error(err).IsNil()

	config := mocks.Config{
		PortValue: port,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: []byte(`{"auth": "noauth", "udp": true}`),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: protocol,
			SettingsValue: nil,
		},
	}

	point, err := point.NewPoint(&config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(port),
		Zone: "",
	})

	assert.Error(err).IsNil()

	data2Send := []byte("Fake DNS request")

	buffer := make([]byte, 0, 1024)
	buffer = append(buffer, 0, 0, 0)
	buffer = append(buffer, 1, 8, 8, 4, 4, 0, 53)
	buffer = append(buffer, data2Send...)

	conn.Write(buffer)

	response := make([]byte, 1024)
	nBytes, err := conn.Read(response)

	assert.Error(err).IsNil()
	assert.Bytes(response[10:nBytes]).Equals(connInput)
	assert.Bytes(data2Send).Equals(connOutput.Bytes())
	assert.StringLiteral(och.Destination.String()).Equals("udp:8.8.4.4:53")
}
