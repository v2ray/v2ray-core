package freedom_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	. "github.com/v2ray/v2ray-core/proxy/freedom"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	"github.com/v2ray/v2ray-core/transport/ray"
)

func TestSinglePacket(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	tcpServer := &tcp.Server{
		Port: port,
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	freedom := &FreedomConnection{}
	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := alloc.NewSmallBuffer().Clear().Append([]byte(data2Send))
	packet := v2net.NewPacket(v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), port), payload, false)

	err = freedom.Dispatch(packet, traffic)
	assert.Error(err).IsNil()
	close(traffic.InboundInput())

	respPayload := <-traffic.InboundOutput()
	defer respPayload.Release()
	assert.Bytes(respPayload.Value).Equals([]byte("Processed: Data to be sent to remote"))

	_, open := <-traffic.InboundOutput()
	assert.Bool(open).IsFalse()

	tcpServer.Close()
}

func TestUnreachableDestination(t *testing.T) {
	v2testing.Current(t)

	freedom := &FreedomConnection{}
	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := alloc.NewSmallBuffer().Clear().Append([]byte(data2Send))
	packet := v2net.NewPacket(v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 2}), 80), payload, false)

	err := freedom.Dispatch(packet, traffic)
	assert.Error(err).IsNotNil()

	_, open := <-traffic.InboundOutput()
	assert.Bool(open).IsFalse()
}
