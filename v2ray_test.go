package core_test

import (
	"testing"

	. "v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/uuid"
	_ "v2ray.com/core/main/distro/all"
	"v2ray.com/core/proxy/dokodemo"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/outbound"
	"v2ray.com/core/testing/assert"
)

func TestV2RayClose(t *testing.T) {
	assert := assert.On(t)

	port := net.Port(dice.RollUint16())
	config := &Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: net.SinglePortRange(port),
					Listen:    net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: net.NewIPOrDomain(net.LocalHostIP),
					Port:    uint32(0),
					NetworkList: &net.NetworkList{
						Network: []net.Network{net.Network_TCP, net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(0),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: uuid.New().String(),
									}),
								},
							},
						},
					},
				}),
			},
		},
	}

	server, err := New(config)
	assert.Error(err).IsNil()

	server.Close()
}
