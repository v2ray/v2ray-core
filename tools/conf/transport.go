package conf

import (
	"errors"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

type TransportConfig struct {
	TCPConfig *TCPConfig       `json:"tcpSettings"`
	KCPConfig *KCPConfig       `json:"kcpSettings"`
	WSConfig  *WebSocketConfig `json:"wsSettings"`
}

func (this *TransportConfig) Build() (*transport.Config, error) {
	config := new(transport.Config)

	if this.TCPConfig != nil {
		ts, err := this.TCPConfig.Build()
		if err != nil {
			return nil, errors.New("Failed to build TCP config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_TCP,
			Settings: ts,
		})
	}

	if this.KCPConfig != nil {
		ts, err := this.KCPConfig.Build()
		if err != nil {
			return nil, errors.New("Failed to build KCP config: " + err.Error())
		}
		config.NetworkSettings = append(config.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_KCP,
			Settings: ts,
		})
	}

	if this.WSConfig != nil {
		ts, err := this.WSConfig.Build()
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
