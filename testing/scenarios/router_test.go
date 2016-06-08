package scenarios

import (
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
)

func TestRouter(t *testing.T) {
	assert := assert.On(t)

	tcpServer := &tcp.Server{
		Port: v2net.Port(50024),
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	tcpServer2Accessed := false
	tcpServer2 := &tcp.Server{
		Port: v2net.Port(50025),
		MsgProcessor: func(data []byte) []byte {
			tcpServer2Accessed = true
			return data
		},
	}
	_, err = tcpServer2.Start()
	assert.Error(err).IsNil()
	defer tcpServer2.Close()

	assert.Error(InitializeServerSetOnce("test_3")).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(50020),
	})
	assert.Error(err).IsNil()

	payload := "direct dokodemo request."
	nBytes, err := conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	conn.CloseWrite()

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert.Error(err).IsNil()
	assert.String("Processed: " + payload).Equals(string(response[:nBytes]))
	conn.Close()

	conn, err = net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(50022),
	})
	assert.Error(err).IsNil()

	payload = "blocked dokodemo request."
	nBytes, err = conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	conn.CloseWrite()

	response = make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert.Error(err).IsNotNil()
	assert.Int(nBytes).Equals(0)
	assert.Bool(tcpServer2Accessed).IsFalse()
	conn.Close()

	CloseAllServers()
}
