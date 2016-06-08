package socks_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dns"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	v2proxy "github.com/v2ray/v2ray-core/proxy"
	proxytesting "github.com/v2ray/v2ray-core/proxy/testing"
	proxymocks "github.com/v2ray/v2ray-core/proxy/testing/mocks"
	"github.com/v2ray/v2ray-core/shell/point"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSocksTcpConnect(t *testing.T) {
	assert := assert.On(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnOutput: connOutput,
		ConnInput:  bytes.NewReader(connInput),
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och",
		func(space app.Space, config interface{}, meta *v2proxy.OutboundHandlerMeta) (v2proxy.OutboundHandler, error) {
			return och, nil
		})
	assert.Error(err).IsNil()

	config := &point.Config{
		Port: port,
		InboundConfig: &point.InboundConnectionConfig{
			Protocol: "socks",
			ListenOn: v2net.LocalHostIP,
			Settings: []byte(`
      {
        "auth": "noauth"
      }`),
		},
		DNSConfig: &dns.Config{
			NameServers: []v2net.Destination{
				v2net.UDPDestination(v2net.DomainAddress("localhost"), v2net.Port(53)),
			},
		},
		OutboundConfig: &point.OutboundConnectionConfig{
			Protocol: protocol,
			Settings: nil,
		},
	}

	point, err := point.NewPoint(config)
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
	assert.String(targetServer).Equals(och.Destination.NetAddr())
}

func TestSocksTcpConnectWithUserPass(t *testing.T) {
	assert := assert.On(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och",
		func(space app.Space, config interface{}, meta *v2proxy.OutboundHandlerMeta) (v2proxy.OutboundHandler, error) {
			return och, nil
		})
	assert.Error(err).IsNil()

	config := &point.Config{
		Port: port,
		InboundConfig: &point.InboundConnectionConfig{
			Protocol: "socks",
			ListenOn: v2net.LocalHostIP,
			Settings: []byte(`
      {
        "auth": "password",
        "accounts": [
          {"user": "userx", "pass": "passy"}
        ]
      }`),
		},
		DNSConfig: &dns.Config{
			NameServers: []v2net.Destination{
				v2net.UDPDestination(v2net.DomainAddress("localhost"), v2net.Port(53)),
			},
		},
		OutboundConfig: &point.OutboundConnectionConfig{
			Protocol: protocol,
			Settings: nil,
		},
	}

	point, err := point.NewPoint(config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), &proxy.Auth{User: "userx", Password: "passy"}, proxy.Direct)
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
	assert.String(targetServer).Equals(och.Destination.NetAddr())
}

func TestSocksTcpConnectWithWrongUserPass(t *testing.T) {
	assert := assert.On(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och",
		func(space app.Space, config interface{}, meta *v2proxy.OutboundHandlerMeta) (v2proxy.OutboundHandler, error) {
			return och, nil
		})
	assert.Error(err).IsNil()

	config := &point.Config{
		Port: port,
		InboundConfig: &point.InboundConnectionConfig{
			Protocol: "socks",
			ListenOn: v2net.LocalHostIP,
			Settings: []byte(`
      {
        "auth": "password",
        "accounts": [
          {"user": "userx", "pass": "passy"}
        ]
      }`),
		},
		DNSConfig: &dns.Config{
			NameServers: []v2net.Destination{
				v2net.UDPDestination(v2net.DomainAddress("localhost"), v2net.Port(53)),
			},
		},
		OutboundConfig: &point.OutboundConnectionConfig{
			Protocol: protocol,
			Settings: nil,
		},
	}

	point, err := point.NewPoint(config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), &proxy.Auth{User: "userx", Password: "passz"}, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "1.2.3.4:443"
	_, err = socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNotNil()
}

func TestSocksTcpConnectWithWrongAuthMethod(t *testing.T) {
	assert := assert.On(t)
	port := v2nettesting.PickPort()

	connInput := []byte("The data to be returned to socks server.")
	connOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(connInput),
		ConnOutput: connOutput,
	}

	protocol, err := proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och",
		func(space app.Space, config interface{}, meta *v2proxy.OutboundHandlerMeta) (v2proxy.OutboundHandler, error) {
			return och, nil
		})
	assert.Error(err).IsNil()

	config := &point.Config{
		Port: port,
		InboundConfig: &point.InboundConnectionConfig{
			ListenOn: v2net.LocalHostIP,
			Protocol: "socks",
			Settings: []byte(`
      {
        "auth": "password",
        "accounts": [
          {"user": "userx", "pass": "passy"}
        ]
      }`),
		},
		DNSConfig: &dns.Config{
			NameServers: []v2net.Destination{
				v2net.UDPDestination(v2net.DomainAddress("localhost"), v2net.Port(53)),
			},
		},
		OutboundConfig: &point.OutboundConnectionConfig{
			Protocol: protocol,
			Settings: nil,
		},
	}

	point, err := point.NewPoint(config)
	assert.Error(err).IsNil()

	err = point.Start()
	assert.Error(err).IsNil()

	socks5Client, err := proxy.SOCKS5("tcp", fmt.Sprintf("127.0.0.1:%d", port), nil, proxy.Direct)
	assert.Error(err).IsNil()

	targetServer := "1.2.3.4:443"
	_, err = socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNotNil()
}
