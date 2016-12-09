package scenarios

import (
	"fmt"
	"net"
	"testing"
	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
)

func TestShadowsocksTCP(t *testing.T) {
	assert := assert.On(t)

	tcpServer := &tcp.Server{
		Port: v2net.Port(50052),
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

	assert.Error(InitializeServerSetOnce("test_6")).IsNil()

	for i := 0; i < 1; i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: 50050,
		})
		assert.Error(err).IsNil()

		payload := "dokodemo request."
		nBytes, err := conn.Write([]byte(payload))
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(payload))

		//conn.CloseWrite()

		response := buf.NewBuffer()
		finished := false
		expectedResponse := "Processed: " + payload
		for {
			_, err := response.FillFrom(conn)
			assert.Error(err).IsNil()
			if err != nil {
				break
			}
			if response.String() == expectedResponse {
				finished = true
				break
			}
			if response.Len() > len(expectedResponse) {
				fmt.Printf("Unexpected response: %v\n", response.Bytes())
				break
			}
		}
		assert.Bool(finished).IsTrue()

		conn.Close()
	}

	CloseAllServers()
}
