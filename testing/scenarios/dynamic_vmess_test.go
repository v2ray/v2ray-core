package scenarios

import (
	"net"
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
)

func TestDynamicVMess(t *testing.T) {
	assert := assert.On(t)

	tcpServer := &tcp.Server{
		Port: v2net.Port(50032),
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

	assert.Error(InitializeServerSetOnce("test_4")).IsNil()

	for i := 0; i < 100; i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: 50030,
		})
		assert.Error(err).IsNil()

		payload := "dokodemo request."
		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))

		expectedResponse := "Processed: " + payload
		response := readFrom(conn, time.Second, len(expectedResponse))
		assert.String(string(response)).Equals(expectedResponse)

		conn.Close()
	}

	CloseAllServers()
}
