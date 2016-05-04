package scenarios

import (
	"bytes"
	"io"
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
)

func TestDynamicVMess(t *testing.T) {
	v2testing.Current(t)

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

		conn.CloseWrite()

		response := bytes.NewBuffer(nil)
		_, err = io.Copy(response, conn)
		assert.Error(err).IsNil()
		assert.StringLiteral("Processed: " + payload).Equals(string(response.Bytes()))

		conn.Close()
	}

	CloseAllServers()
}
