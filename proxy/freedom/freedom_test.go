package freedom

import (
	"io/ioutil"
	"net"
	"sync"
	"testing"

	"golang.org/x/net/proxy"

	"github.com/v2ray/v2ray-core"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestSocksTcpConnect(t *testing.T) {
	assert := unit.Assert(t)
	port := 48274

	data2Send := "Data to be sent to remote"
	data2Return := "Data to be returned to local"

	var serverReady sync.Mutex
	serverReady.Lock()

	go func() {
		listener, err := net.ListenTCP("tcp", &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: port,
			Zone: "",
		})
		assert.Error(err).IsNil()

		serverReady.Unlock()
		conn, err := listener.Accept()
		assert.Error(err).IsNil()

		buffer := make([]byte, 1024)
		nBytes, err := conn.Read(buffer)
		assert.Error(err).IsNil()

		if string(buffer[:nBytes]) == data2Send {
			_, err = conn.Write([]byte(data2Return))
			assert.Error(err).IsNil()
		}

		conn.Close()
		listener.Close()
	}()

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

	serverReady.Lock()

	targetServer := "127.0.0.1:48274"
	conn, err := socks5Client.Dial("tcp", targetServer)
	assert.Error(err).IsNil()

	conn.Write([]byte(data2Send))
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	dataReturned, err := ioutil.ReadAll(conn)
	assert.Error(err).IsNil()
	conn.Close()

	assert.Bytes(dataReturned).Equals([]byte(data2Return))
}
