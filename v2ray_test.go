package core_test

import (
	"testing"

	. "v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/dice"
	v2net "v2ray.com/core/common/net"
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

	port := v2net.Port(dice.RollUint16())
	config := &Config{
		Inbound: []*proxyman.InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortRange: v2net.SinglePortRange(port),
					Listen:    v2net.NewIPOrDomain(v2net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
					Port:    uint32(0),
					NetworkList: &v2net.NetworkList{
						Network: []v2net.Network{v2net.Network_TCP, v2net.Network_UDP},
					},
				}),
			},
		},
		Outbound: []*proxyman.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: v2net.NewIPOrDomain(v2net.LocalHostIP),
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
