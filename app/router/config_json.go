// +build json

package router

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Strategy string          `json:"strategy"`
		Settings json.RawMessage `json:"settings"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	settings, err := CreateRouterConfig(jsonConfig.Strategy, []byte(jsonConfig.Settings))
	if err != nil {
		log.Error("Router: Failed to load router settings: ", err)
		return err
	}
	this.Strategy = jsonConfig.Strategy
	this.Settings = settings
	return nil
}
