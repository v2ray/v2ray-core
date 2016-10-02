// +build json

package transport

import (
	"encoding/json"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/ws"

	"github.com/golang/protobuf/ptypes"
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
		any, err := ptypes.MarshalAny(jsonConfig.TCPConfig)
		if err != nil {
			return err
		}
		this.NetworkSettings = append(this.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_TCP,
			Settings: any,
		})
	}

	if jsonConfig.KCPConfig != nil {
		any, err := ptypes.MarshalAny(jsonConfig.KCPConfig)
		if err != nil {
			return err
		}
		this.NetworkSettings = append(this.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_KCP,
			Settings: any,
		})
	}

	if jsonConfig.WSConfig != nil {
		any, err := ptypes.MarshalAny(jsonConfig.WSConfig)
		if err != nil {
			return err
		}
		this.NetworkSettings = append(this.NetworkSettings, &internet.NetworkSettings{
			Network:  v2net.Network_WebSocket,
			Settings: any,
		})
	}
	return nil
}
