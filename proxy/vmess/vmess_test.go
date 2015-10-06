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
			SettingsValue: nil,
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &VMessOutboundConfig{
				[]VNextConfig{
					VNextConfig{
						Address: "127.0.0.1",
						Port:    13829,
						Network: "tcp",
						Users: []VMessUser{
							VMessUser{Id: "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"},
						},
					},
				},
			},
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
			SettingsValue: &VMessInboundConfig{
				AllowedClients: []VMessUser{
					VMessUser{Id: "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"},
				},
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_och",
			SettingsValue: nil,
		},
	}

	pointB, err := core.NewPoint(&configB)
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

	portA := uint16(17394)
	ich := &mocks.InboundConnectionHandler{
		Data2Send:    []byte(data2Send),
		DataReturned: bytes.NewBuffer(make([]byte, 0, 1024)),
	}

	core.RegisterInboundConnectionHandlerFactory("mock_ich", ich)

	configA := mocks.Config{
		PortValue: portA,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_ich",
			SettingsValue: nil,
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &VMessOutboundConfig{
				[]VNextConfig{
					VNextConfig{
						Address: "127.0.0.1",
						Port:    13841,
						Network: "udp",
						Users: []VMessUser{
							VMessUser{Id: "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"},
						},
					},
				},
			},
		},
	}

	pointA, err := core.NewPoint(&configA)
	assert.Error(err).IsNil()

	err = pointA.Start()
	assert.Error(err).IsNil()

	portB := uint16(13841)

	och := &mocks.OutboundConnectionHandler{
		Data2Send:   bytes.NewBuffer(make([]byte, 0, 1024)),
		Data2Return: []byte("The data to be returned to inbound server."),
	}

	core.RegisterOutboundConnectionHandlerFactory("mock_och", och)

	configB := mocks.Config{
		PortValue: portB,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &VMessInboundConfig{
				AllowedClients: []VMessUser{
					VMessUser{Id: "ad937d9d-6e23-4a5a-ba23-bce5092a7c51"},
				},
				UDPEnabled: true,
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "mock_och",
			SettingsValue: nil,
		},
	}

	pointB, err := core.NewPoint(&configB)
	assert.Error(err).IsNil()

	err = pointB.Start()
	assert.Error(err).IsNil()

	dest := v2net.NewUDPDestination(v2net.IPAddress([]byte{1, 2, 3, 4}, 80))
	ich.Communicate(v2net.NewPacket(dest, []byte(data2Send), false))
	assert.Bytes([]byte(data2Send)).Equals(och.Data2Send.Bytes())
	assert.Bytes(ich.DataReturned.Bytes()).Equals(och.Data2Return)
}
