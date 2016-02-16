package scenarios

import (
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"

	ssclient "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func TestShadowsocksTCP(t *testing.T) {
	v2testing.Current(t)

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

	assert.Error(InitializeServerServer("test_6")).IsNil()

	cipher, err := ssclient.NewCipher("aes-256-cfb", "v2ray-password")
	assert.Error(err).IsNil()

	rawAddr := []byte{1, 127, 0, 0, 1, 0xc3, 0x84} // 127.0.0.1:50052
	conn, err := ssclient.DialWithRawAddr(rawAddr, "127.0.0.1:50051", cipher)
	assert.Error(err).IsNil()

	payload := "shadowsocks request."
	nBytes, err := conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	conn.Conn.(*net.TCPConn).CloseWrite()

	response := make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert.Error(err).IsNil()
	assert.StringLiteral("Processed: " + payload).Equals(string(response[:nBytes]))
	conn.Close()

	cipher, err = ssclient.NewCipher("aes-128-cfb", "v2ray-another")
	assert.Error(err).IsNil()

	conn, err = ssclient.DialWithRawAddr(rawAddr, "127.0.0.1:50055", cipher)
	assert.Error(err).IsNil()

	payload = "shadowsocks request 2."
	nBytes, err = conn.Write([]byte(payload))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(len(payload))

	conn.Conn.(*net.TCPConn).CloseWrite()

	response = make([]byte, 1024)
	nBytes, err = conn.Read(response)
	assert.Error(err).IsNil()
	assert.StringLiteral("Processed: " + payload).Equals(string(response[:nBytes]))
	conn.Close()

	CloseAllServers()
}
