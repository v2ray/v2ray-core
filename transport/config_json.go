// +build json

package transport

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/transport/internet/kcp"
	"github.com/v2ray/v2ray-core/transport/internet/tcp"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		TCPConfig *tcp.Config `json:"tcpSettings"`
		KCPCOnfig *kcp.Config `json:"kcpSettings"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.tcpConfig = jsonConfig.TCPConfig
	this.kcpConfig = jsonConfig.KCPCOnfig

	return nil
}
