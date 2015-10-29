package vmess

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/app/point"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
	"github.com/v2ray/v2ray-core/proxy/vmess/config/json"
	"github.com/v2ray/v2ray-core/testing/mocks"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestVMessInAndOut(t *testing.T) {
	assert := unit.Assert(t)

	data2Send := "The data to be send to outbound server."
	testAccount, err := config.NewID("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	assert.Error(err).IsNil()

	portA := uint16(17392)
	ich := &mocks.InboundConnectionHandler{
		Data2Send:    []byte(data2Send),
		DataReturned: bytes.NewBuffer(make([]byte, 0, 1024)),
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
						Address:    v2net.IPAddress([]byte{127, 0, 0, 1}, 13829),
						TCPEnabled: true,
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

	portB := uint16(13829)

	och := &mocks.OutboundConnectionHandler{
		Data2Send:   bytes.NewBuffer(make([]byte, 0, 1024)),
		Data2Return: []byte("The data to be returned to inbound server."),
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
	assert.Bytes([]byte(data2Send)).Equals(och.Data2Send.Bytes())
	assert.Bytes(ich.DataReturned.Bytes()).Equals(och.Data2Return)
}

func TestVMessInAndOutUDP(t *testing.T) {
	assert := unit.Assert(t)

	data2Send := "The data to be send to outbound server."
	testAccount, err := config.NewID("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	assert.Error(err).IsNil()

	portA := uint16(17394)
	ich := &mocks.InboundConnectionHandler{
		Data2Send:    []byte(data2Send),
		DataReturned: bytes.NewBuffer(make([]byte, 0, 1024)),
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
						Address:    v2net.IPAddress([]byte{127, 0, 0, 1}, 13841),
						UDPEnabled: true,
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

	portB := uint16(13841)

	och := &mocks.OutboundConnectionHandler{
		Data2Send:   bytes.NewBuffer(make([]byte, 0, 1024)),
		Data2Return: []byte("The data to be returned to inbound server."),
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
				UDP: true,
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

	data2SendBuffer := alloc.NewBuffer()
	data2SendBuffer.Clear()
	data2SendBuffer.Append([]byte(data2Send))
	dest := v2net.NewUDPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}, 80))
	ich.Communicate(v2net.NewPacket(dest, data2SendBuffer, false))
	assert.Bytes([]byte(data2Send)).Equals(och.Data2Send.Bytes())
	assert.Bytes(ich.DataReturned.Bytes()).Equals(och.Data2Return)
}
