package conf

import (
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

type TransportConfig struct {
	TCPConfig *TCPConfig       `json:"tcpSettings"`
	KCPConfig *KCPConfig       `json:"kcpSettings"`
	WSConfig  *WebSocketConfig `json:"wsSettings"`
}

func (v *TransportConfig) Build() (*transport.Config, error) {
	config := new(transport.Config)

	if v.TCPConfig != nil {
		ts, err := v.TCPConfig.Build()
		if err != nil {
			return nil, errors.New("Failed to build TCP config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_TCP,
			Settings: ts,
		})
	}

	if v.KCPConfig != nil {
		ts, err := v.KCPConfig.Build()
		if err != nil {
			return nil, errors.New("Failed to build KCP config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_KCP,
			Settings: ts,
		})
	}

	if v.WSConfig != nil {
		ts, err := v.WSConfig.Build()
		if err != nil {
			return nil, errors.New("Failed to build WebSocket config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_WebSocket,
			Settings: ts,
		})
	}
	return config, nil
}
