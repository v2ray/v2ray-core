// +build json

package transport

import (
	"encoding/json"

	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/ws"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		TCPConfig *tcp.Config `json:"tcpSettings"`
		KCPConfig *kcp.Config `json:"kcpSettings"`
		WSConfig  *ws.Config  `json:"wsSettings"`
	}
	jsonConfig := &JsonConfig{}
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}

	if jsonConfig.TCPConfig != nil {
		this.NetworkSettings = append(this.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_TCP,
			Settings: loader.NewTypedSettings(jsonConfig.TCPConfig),
		})
	}

	if jsonConfig.KCPConfig != nil {
		this.NetworkSettings = append(this.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_KCP,
			Settings: loader.NewTypedSettings(jsonConfig.KCPConfig),
		})
	}

	if jsonConfig.WSConfig != nil {
		this.NetworkSettings = append(this.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_WebSocket,
			Settings: loader.NewTypedSettings(jsonConfig.WSConfig),
		})
	}
	return nil
}
