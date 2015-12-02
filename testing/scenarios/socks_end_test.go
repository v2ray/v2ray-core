package scenarios

import (
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	"github.com/v2ray/v2ray-core/testing/servers/udp"
)

var (
	EmptyRouting = func(v2net.Destination) bool {
		return false
	}
)

func TestTCPConnection(t *testing.T) {
	v2testing.Current(t)

	targetPort := v2nettesting.PickPort()
	tcpServer := &tcp.Server{
		Port: targetPort,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	v2rayPort, _, err := setUpV2Ray(EmptyRouting)
	assert.Error(err).IsNil()

	for i := 0; i < 100; i++ {
		conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
			IP:   []byte{127, 0, 0, 1},
			Port: int(v2rayPort),
		})

		authRequest := socks5AuthMethodRequest(byte(0))
		nBytes, err := conn.Write(authRequest)
		assert.Int(nBytes).Equals(len(authRequest))
		assert.Error(err).IsNil()

		authResponse := make([]byte, 1024)
		nBytes, err = conn.Read(authResponse)
		assert.Error(err).IsNil()
		assert.Bytes(authResponse[:nBytes]).Equals([]byte{socks5Version, 0})

		connectRequest := socks5Request(byte(1), v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort))
		nBytes, err = conn.Write(connectRequest)
		assert.Int(nBytes).Equals(len(connectRequest))
		assert.Error(err).IsNil()

		connectResponse := make([]byte, 1024)
		nBytes, err = conn.Read(connectResponse)
		assert.Error(err).IsNil()
		assert.Bytes(connectResponse[:nBytes]).Equals([]byte{socks5Version, 0, 0, 1, 0, 0, 0, 0, 6, 181})

		actualRequest := []byte("Request to target server.")
		nBytes, err = conn.Write(actualRequest)
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(actualRequest))

		actualRequest = []byte("Request to target server again.")
		nBytes, err = conn.Write(actualRequest)
		assert.Error(err).IsNil()
		assert.Int(nBytes).Equals(len(actualRequest))

		conn.CloseWrite()

		actualResponse := make([]byte, 1024)
		nBytes, err = conn.Read(actualResponse)
		assert.Error(err).IsNil()
		assert.StringLiteral(string(actualResponse[:nBytes])).Equals("Processed: Request to target server.Request to target server again.")

		conn.Close()
	}
}

func TestTCPBind(t *testing.T) {
	v2testing.Current(t)

	targetPort := v2nettesting.PickPort()
	tcpServer := &tcp.Server{
		Port: targetPort,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	v2rayPort, _, err := setUpV2Ray(EmptyRouting)
	assert.Error(err).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(v2rayPort),
	})

	authRequest := socks5AuthMethodRequest(byte(0))
	nBytes, err := conn.Write(authRequest)
	assert.Int(nBytes).Equals(len(authRequest))
	assert.Error(err).IsNil()

	authResponse := make([]byte, 1024)
	nBytes, err = conn.Read(authResponse)
	assert.Error(err).IsNil()
	assert.Bytes(authResponse[:nBytes]).Equals([]byte{socks5Version, 0})

	connectRequest := socks5Request(byte(2), v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort))
	nBytes, err = conn.Write(connectRequest)
	assert.Int(nBytes).Equals(len(connectRequest))
	assert.Error(err).IsNil()

	connectResponse := make([]byte, 1024)
	nBytes, err = conn.Read(connectResponse)
	assert.Error(err).IsNil()
	assert.Bytes(connectResponse[:nBytes]).Equals([]byte{socks5Version, 7, 0, 1, 0, 0, 0, 0, 0, 0})

	conn.Close()
}

func TestUDPAssociate(t *testing.T) {
	v2testing.Current(t)

	targetPort := v2nettesting.PickPort()
	udpServer := &udp.Server{
		Port: targetPort,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := udpServer.Start()
	assert.Error(err).IsNil()

	v2rayPort, _, err := setUpV2Ray(EmptyRouting)
	assert.Error(err).IsNil()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(v2rayPort),
	})

	authRequest := socks5AuthMethodRequest(byte(0))
	nBytes, err := conn.Write(authRequest)
	assert.Int(nBytes).Equals(len(authRequest))
	assert.Error(err).IsNil()

	authResponse := make([]byte, 1024)
	nBytes, err = conn.Read(authResponse)
	assert.Error(err).IsNil()
	assert.Bytes(authResponse[:nBytes]).Equals([]byte{socks5Version, 0})

	connectRequest := socks5Request(byte(3), v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort))
	nBytes, err = conn.Write(connectRequest)
	assert.Int(nBytes).Equals(len(connectRequest))
	assert.Error(err).IsNil()

	connectResponse := make([]byte, 1024)
	nBytes, err = conn.Read(connectResponse)
	assert.Error(err).IsNil()
	assert.Bytes(connectResponse[:nBytes]).Equals([]byte{socks5Version, 0, 0, 1, 127, 0, 0, 1, byte(v2rayPort >> 8), byte(v2rayPort)})

	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(v2rayPort),
	})
	assert.Error(err).IsNil()

	udpPayload := "UDP request to udp server."
	udpRequest := socks5UDPRequest(v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort), []byte(udpPayload))

	nBytes, err = udpConn.Write(udpRequest)
	assert.Int(nBytes).Equals(len(udpRequest))
	assert.Error(err).IsNil()

	udpResponse := make([]byte, 1024)
	nBytes, err = udpConn.Read(udpResponse)
	assert.Error(err).IsNil()
	assert.Bytes(udpResponse[:nBytes]).Equals(
		socks5UDPRequest(v2net.IPAddress([]byte{127, 0, 0, 1}, targetPort), []byte("Processed: UDP request to udp server.")))

	udpConn.Close()
	conn.Close()
}
