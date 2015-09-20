package vmess

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestVMessInAndOut(t *testing.T) {
	assert := unit.Assert(t)

	data2Send := "The data to be send to outbound server."

	portA := uint16(17392)
	ich := &mocks.InboundConnectionHandler{
		Data2Send:    []byte(data2Send),
		DataReturned: bytes.NewBuffer(make([]byte, 0, 1024)),
	}

	core.RegisterInboundConnectionHandlerFactory("mock_ich", ich)

	configA := mocks.Config{
		PortValue: portA,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_ich",
			ContentValue:  nil,
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			ContentValue:  []byte("{\"vnext\":[{\"address\": \"127.0.0.1\", \"port\": 13829, \"users\":[{\"id\": \"ad937d9d-6e23-4a5a-ba23-bce5092a7c51\"}]}]}"),
		},
	}

	pointA, err := core.NewPoint(&configA)
	assert.Error(err).IsNil()

	err = pointA.Start()
	assert.Error(err).IsNil()

	portB := uint16(13829)

	och := &mocks.OutboundConnectionHandler{
		Data2Send:   bytes.NewBuffer(make([]byte, 0, 1024)),
		Data2Return: []byte("The data to be returned to inbound server."),
	}

	core.RegisterOutboundConnectionHandlerFactory("mock_och", och)

	configB := mocks.Config{
		PortValue: portB,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			ContentValue:  []byte("{\"clients\": [{\"id\": \"ad937d9d-6e23-4a5a-ba23-bce5092a7c51\"}]}"),
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_och",
			ContentValue:  nil,
		},
	}

	pointB, err := core.NewPoint(&configB)
	assert.Error(err).IsNil()

	err = pointB.Start()
	assert.Error(err).IsNil()

	dest := v2net.NewDestination(v2net.NetTCP, v2net.IPAddress([]byte{1, 2, 3, 4}, 80))
	ich.Communicate(dest)
	assert.Bytes([]byte(data2Send)).Equals(och.Data2Send.Bytes())
	assert.Bytes(ich.DataReturned.Bytes()).Equals(och.Data2Return)
}
