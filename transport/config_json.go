// +build json

package transport

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/transport/hub/kcpv"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		ConnectionReuse bool         `json:"connectionReuse"`
		EnableKcp       bool         `json:"EnableKCP,omitempty"`
		KcpConfig       *kcpv.Config `json:"KcpConfig,omitempty"`
	}
	jsonConfig := &JsonConfig{
		ConnectionReuse: true,
		EnableKcp:       false,
	}
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.ConnectionReuse = jsonConfig.ConnectionReuse
	this.enableKcp = jsonConfig.EnableKcp
	if jsonConfig.KcpConfig != nil {
		this.kcpConfig = jsonConfig.KcpConfig
		if jsonConfig.KcpConfig.AdvancedConfigs == nil {
			jsonConfig.KcpConfig.AdvancedConfigs = kcpv.DefaultAdvancedConfigs
		}
	} else {
		if jsonConfig.EnableKcp {
			log.Error("transport: You have enabled KCP but no configure is given")
		}
	}

	return nil
}
