package vmess_test

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy"
	proxytesting "github.com/v2ray/v2ray-core/proxy/testing"
	proxymocks "github.com/v2ray/v2ray-core/proxy/testing/mocks"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/inbound"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
	"github.com/v2ray/v2ray-core/shell/point"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestVMessInAndOut(t *testing.T) {
	v2testing.Current(t)

	id, err := uuid.ParseString("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	assert.Error(err).IsNil()

	testAccount := proto.NewID(id)

	portA := v2nettesting.PickPort()
	portB := v2nettesting.PickPort()

	ichConnInput := []byte("The data to be send to outbound server.")
	ichConnOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	ich := &proxymocks.InboundConnectionHandler{
		ConnInput:  bytes.NewReader(ichConnInput),
		ConnOutput: ichConnOutput,
	}

	protocol, err := proxytesting.RegisterInboundConnectionHandlerCreator("mock_och", func(space app.Space, config interface{}) (proxy.InboundHandler, error) {
		ich.PacketDispatcher = space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
		return ich, nil
	})
	assert.Error(err).IsNil()

	configA := &point.Config{
		Port: portA,
		InboundConfig: &point.ConnectionConfig{
			Protocol: protocol,
			Settings: nil,
		},
		OutboundConfig: &point.ConnectionConfig{
			Protocol: "vmess",
			Settings: []byte(`{
        "vnext": [
          {
            "address": "127.0.0.1",
            "port": ` + portB.String() + `,
            "users": [
              {"id": "` + testAccount.String() + `"}
            ]
          }
        ]
      }`),
		},
	}

	pointA, err := point.NewPoint(configA)
	assert.Error(err).IsNil()

	err = pointA.Start()
	assert.Error(err).IsNil()

	ochConnInput := []byte("The data to be returned to inbound server.")
	ochConnOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(ochConnInput),
		ConnOutput: ochConnOutput,
	}

	protocol, err = proxytesting.RegisterOutboundConnectionHandlerCreator("mock_och", func(space app.Space, config interface{}) (proxy.OutboundHandler, error) {
		return och, nil
	})
	assert.Error(err).IsNil()

	configB := &point.Config{
		Port: portB,
		InboundConfig: &point.ConnectionConfig{
			Protocol: "vmess",
			Settings: []byte(`{
        "clients": [
          {"id": "` + testAccount.String() + `"}
        ]
      }`),
		},
		OutboundConfig: &point.ConnectionConfig{
			Protocol: protocol,
			Settings: nil,
		},
	}

	pointB, err := point.NewPoint(configB)
	assert.Error(err).IsNil()

	err = pointB.Start()
	assert.Error(err).IsNil()

	dest := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	ich.Communicate(v2net.NewPacket(dest, nil, true))
	assert.Bytes(ichConnInput).Equals(ochConnOutput.Bytes())
	assert.Bytes(ichConnOutput.Bytes()).Equals(ochConnInput)
}
