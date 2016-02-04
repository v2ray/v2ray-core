package dokodemo_test

import (
	"net"
	"testing"

	testdispatcher "github.com/v2ray/v2ray-core/app/dispatcher/testing"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	. "github.com/v2ray/v2ray-core/proxy/dokodemo"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestDokodemoTCP(t *testing.T) {
	v2testing.Current(t)

	testPacketDispatcher := testdispatcher.NewTestPacketDispatcher(nil)

	data2Send := "Data to be sent to remote."

	dokodemo := NewDokodemoDoor(&Config{
		Address: v2net.IPAddress([]byte{1, 2, 3, 4}),
		Port:    128,
		Network: v2net.TCPNetwork.AsList(),
		Timeout: 600,
	}, testPacketDispatcher)
	defer dokodemo.Close()

	port := v2nettesting.PickPort()
	err := dokodemo.Listen(port)
	assert.Error(err).IsNil()
	netassert.Port(port).Equals(dokodemo.Port())

	tcpClient, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(port),
		Zone: "",
	})
	assert.Error(err).IsNil()

	tcpClient.Write([]byte(data2Send))
	tcpClient.CloseWrite()

	lastPacket := <-testPacketDispatcher.LastPacket

	response := make([]byte, 1024)
	nBytes, err := tcpClient.Read(response)
	assert.Error(err).IsNil()
	tcpClient.Close()

	assert.StringLiteral("Processed: " + data2Send).Equals(string(response[:nBytes]))
	assert.Bool(lastPacket.Destination().IsTCP()).IsTrue()
	netassert.Address(lastPacket.Destination().Address()).Equals(v2net.IPAddress([]byte{1, 2, 3, 4}))
	netassert.Port(lastPacket.Destination().Port()).Equals(128)
}

func TestDokodemoUDP(t *testing.T) {
	v2testing.Current(t)

	testPacketDispatcher := testdispatcher.NewTestPacketDispatcher(nil)

	data2Send := "Data to be sent to remote."

	dokodemo := NewDokodemoDoor(&Config{
		Address: v2net.IPAddress([]byte{5, 6, 7, 8}),
		Port:    256,
		Network: v2net.UDPNetwork.AsList(),
		Timeout: 600,
	}, testPacketDispatcher)
	defer dokodemo.Close()

	port := v2nettesting.PickPort()
	err := dokodemo.Listen(port)
	assert.Error(err).IsNil()
	netassert.Port(port).Equals(dokodemo.Port())

	udpClient, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(port),
		Zone: "",
	})
	assert.Error(err).IsNil()

	udpClient.Write([]byte(data2Send))
	udpClient.Close()

	lastPacket := <-testPacketDispatcher.LastPacket

	assert.StringLiteral(data2Send).Equals(string(lastPacket.Chunk().Value))
	assert.Bool(lastPacket.Destination().IsUDP()).IsTrue()
	netassert.Address(lastPacket.Destination().Address()).Equals(v2net.IPAddress([]byte{5, 6, 7, 8}))
	netassert.Port(lastPacket.Destination().Port()).Equals(256)
}
