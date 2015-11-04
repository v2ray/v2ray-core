package vmess

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/app/point"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	proxymocks "github.com/v2ray/v2ray-core/proxy/testing/mocks"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
	"github.com/v2ray/v2ray-core/proxy/vmess/config/json"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestVMessInAndOut(t *testing.T) {
	assert := unit.Assert(t)

	testAccount, err := config.NewID("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	assert.Error(err).IsNil()

	portA := v2nettesting.PickPort()
	portB := v2nettesting.PickPort()

	ichConnInput := []byte("The data to be send to outbound server.")
	ichConnOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	ich := &proxymocks.InboundConnectionHandler{
		ConnInput:  bytes.NewReader(ichConnInput),
		ConnOutput: ichConnOutput,
	}

	connhandler.RegisterInboundConnectionHandlerFactory("mock_ich", ich)

	configA := mocks.Config{
		PortValue: portA,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_ich",
			SettingsValue: nil,
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &json.Outbound{
				[]*json.ConfigTarget{
					&json.ConfigTarget{
						Address: v2net.IPAddress([]byte{127, 0, 0, 1}, portB),
						Users: []*json.ConfigUser{
							&json.ConfigUser{Id: testAccount},
						},
					},
				},
			},
		},
	}

	pointA, err := point.NewPoint(&configA)
	assert.Error(err).IsNil()

	err = pointA.Start()
	assert.Error(err).IsNil()

	ochConnInput := []byte("The data to be returned to inbound server.")
	ochConnOutput := bytes.NewBuffer(make([]byte, 0, 1024))
	och := &proxymocks.OutboundConnectionHandler{
		ConnInput:  bytes.NewReader(ochConnInput),
		ConnOutput: ochConnOutput,
	}

	connhandler.RegisterOutboundConnectionHandlerFactory("mock_och", och)

	configB := mocks.Config{
		PortValue: portB,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &json.Inbound{
				AllowedClients: []*json.ConfigUser{
					&json.ConfigUser{Id: testAccount},
				},
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_och",
			SettingsValue: nil,
		},
	}

	pointB, err := point.NewPoint(&configB)
	assert.Error(err).IsNil()

	err = pointB.Start()
	assert.Error(err).IsNil()

	dest := v2net.NewTCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}, 80))
	ich.Communicate(v2net.NewPacket(dest, nil, true))
	assert.Bytes(ichConnInput).Equals(ochConnOutput.Bytes())
	assert.Bytes(ichConnOutput.Bytes()).Equals(ochConnInput)
}
