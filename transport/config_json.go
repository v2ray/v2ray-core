// +build json

package transport

import (
	"encoding/json"

	"v2ray.com/core/transport/internet/kcp"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/ws"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		TCPConfig *tcp.Config `json:"tcpSettings"`
		KCPConfig kcp.Config  `json:"kcpSettings"`
		WSConfig  *ws.Config  `json:"wsSettings"`
	}
	jsonConfig := &JsonConfig{
		KCPConfig: kcp.DefaultConfig(),
	}
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.tcpConfig = jsonConfig.TCPConfig
	this.kcpConfig = jsonConfig.KCPConfig
	this.wsConfig = jsonConfig.WSConfig
	return nil
}
