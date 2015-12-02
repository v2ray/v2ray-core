package scenarios

import (
	"net"

	routerconfig "github.com/v2ray/v2ray-core/app/router/config/testing"
	_ "github.com/v2ray/v2ray-core/app/router/rules"
	rulesconfig "github.com/v2ray/v2ray-core/app/router/rules/config/testing"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	socksjson "github.com/v2ray/v2ray-core/proxy/socks/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
	vmessjson "github.com/v2ray/v2ray-core/proxy/vmess/config/json"
	"github.com/v2ray/v2ray-core/shell/point"
	"github.com/v2ray/v2ray-core/shell/point/config/testing/mocks"
)

const (
	socks5Version = byte(0x05)
)

func socks5AuthMethodRequest(methods ...byte) []byte {
	request := []byte{socks5Version, byte(len(methods))}
	request = append(request, methods...)
	return request
}

func appendAddress(request []byte, address v2net.Address) []byte {
	switch {
	case address.IsIPv4():
		request = append(request, byte(0x01))
		request = append(request, address.IP()...)

	case address.IsIPv6():
		request = append(request, byte(0x04))
		request = append(request, address.IP()...)

	case address.IsDomain():
		request = append(request, byte(0x03), byte(len(address.Domain())))
		request = append(request, []byte(address.Domain())...)

	}
	request = append(request, address.Port().Bytes()...)
	return request
}

func socks5Request(command byte, address v2net.Address) []byte {
	request := []byte{socks5Version, command, 0}
	request = appendAddress(request, address)
	return request
}

func socks5UDPRequest(address v2net.Address, payload []byte) []byte {
	request := make([]byte, 0, 1024)
	request = append(request, 0, 0, 0)
	request = appendAddress(request, address)
	request = append(request, payload...)
	return request
}

func setUpV2Ray(routing func(v2net.Destination) bool) (v2net.Port, v2net.Port, error) {
	id1, err := config.NewID("ad937d9d-6e23-4a5a-ba23-bce5092a7c51")
	if err != nil {
		return 0, 0, err
	}
	id2, err := config.NewID("93ccfc71-b136-4015-ac85-e037bd1ead9e")
	if err != nil {
		return 0, 0, err
	}
	users := []*vmessjson.ConfigUser{
		&vmessjson.ConfigUser{Id: id1},
		&vmessjson.ConfigUser{Id: id2},
	}

	portB := v2nettesting.PickPort()
	configB := mocks.Config{
		PortValue: portB,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &vmessjson.Inbound{
				AllowedClients: users,
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "freedom",
			SettingsValue: nil,
		},
	}
	pointB, err := point.NewPoint(&configB)
	if err != nil {
		return 0, 0, err
	}
	err = pointB.Start()
	if err != nil {
		return 0, 0, err
	}

	portA := v2nettesting.PickPort()
	portA2 := v2nettesting.PickPort()
	configA := mocks.Config{
		PortValue: portA,
		InboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "socks",
			SettingsValue: &socksjson.SocksConfig{
				AuthMethod: "noauth",
				UDPEnabled: true,
				HostIP:     socksjson.IPAddress(net.IPv4(127, 0, 0, 1)),
			},
		},
		OutboundConfigValue: &mocks.ConnectionConfig{
			ProtocolValue: "vmess",
			SettingsValue: &vmessjson.Outbound{
				[]*vmessjson.ConfigTarget{
					&vmessjson.ConfigTarget{
						Address: v2net.IPAddress([]byte{127, 0, 0, 1}, portB),
						Users:   users,
					},
				},
			},
		},
		InboundDetoursValue: []*mocks.InboundDetourConfig{
			&mocks.InboundDetourConfig{
				PortRangeValue: &mocks.PortRange{
					FromValue: portA2,
					ToValue:   portA2,
				},
				ConnectionConfig: &mocks.ConnectionConfig{
					ProtocolValue: "socks",
					SettingsValue: &socksjson.SocksConfig{
						AuthMethod: "noauth",
						UDPEnabled: false,
						HostIP:     socksjson.IPAddress(net.IPv4(127, 0, 0, 1)),
					},
				},
			},
		},
		OutboundDetoursValue: []*mocks.OutboundDetourConfig{
			&mocks.OutboundDetourConfig{
				TagValue: "direct",
				ConnectionConfig: &mocks.ConnectionConfig{
					ProtocolValue: "freedom",
					SettingsValue: nil,
				},
			},
		},
		RouterConfigValue: &routerconfig.RouterConfig{
			StrategyValue: "rules",
			SettingsValue: &rulesconfig.RouterRuleConfig{
				RuleList: []*rulesconfig.TestRule{
					&rulesconfig.TestRule{
						TagValue: "direct",
						Function: routing,
					},
				},
			},
		},
	}

	pointA, err := point.NewPoint(&configA)
	if err != nil {
		return 0, 0, err
	}
	err = pointA.Start()
	if err != nil {
		return 0, 0, err
	}

	return portA, portA2, nil
}
