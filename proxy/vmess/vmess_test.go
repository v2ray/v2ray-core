package vmess_test

import (
	"bytes"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	proxymocks "github.com/v2ray/v2ray-core/proxy/testing/mocks"
	vmess "github.com/v2ray/v2ray-core/proxy/vmess"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/inbound"
	inboundjson "github.com/v2ray/v2ray-core/proxy/vmess/inbound/json"
	vmessjson "github.com/v2ray/v2ray-core/proxy/vmess/json"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
	outboundjson "github.com/v2ray/v2ray-core/proxy/vmess/outbound/json"
	"github.com/v2ray/v2ray-core/shell/point"
	"github.com/v2ray/v2ray-core/shell/point/testing/mocks"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestVMessInAndOut(t *testing.T) {
	v2testing.Current(t)

	id, err := uuid.ParseString("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	assert.Error(err).IsNil()

	testAccount := vmess.NewID(id)

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
			SettingsValue: &outboundjson.Outbound{
				[]*outboundjson.ConfigTarget{
					&outboundjson.ConfigTarget{
						Destination: v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), portB),
						Users: []*vmessjson.ConfigUser{
							&vmessjson.ConfigUser{Id: testAccount},
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
			SettingsValue: &inboundjson.Inbound{
				AllowedClients: []*vmessjson.ConfigUser{
					&vmessjson.ConfigUser{Id: testAccount},
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

	dest := v2net.TCPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}), 80)
	ich.Communicate(v2net.NewPacket(dest, nil, true))
	assert.Bytes(ichConnInput).Equals(ochConnOutput.Bytes())
	assert.Bytes(ichConnOutput.Bytes()).Equals(ochConnInput)
}
