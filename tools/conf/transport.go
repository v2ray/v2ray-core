package conf

import (
	"v2ray.com/core/common/errors"
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
			return nil, errors.Base(err).Message("Failed to build TCP config.")
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_TCP,
			Settings: ts,
		})
	}

	if v.KCPConfig != nil {
		ts, err := v.KCPConfig.Build()
		if err != nil {
			return nil, errors.Base(err).Message("Failed to build mKCP config.")
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_MKCP,
			Settings: ts,
		})
	}

	if v.WSConfig != nil {
		ts, err := v.WSConfig.Build()
		if err != nil {
			return nil, errors.Base(err).Message("Failed to build WebSocket config.")
		}
		config.TransportSettings = append(config.TransportSettings, &internet.TransportConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			Settings: ts,
		})
	}
	return config, nil
}
