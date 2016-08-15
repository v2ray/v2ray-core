// +build json

package transport

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/transport/internet/kcp"
	"github.com/v2ray/v2ray-core/transport/internet/tcp"
	"github.com/v2ray/v2ray-core/transport/internet/ws"
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
