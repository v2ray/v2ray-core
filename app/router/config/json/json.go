package json

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
)

type RouterConfig struct {
	StrategyValue string          `json:"strategy"`
	SettingsValue json.RawMessage `json:"settings"`
}

func (this *RouterConfig) Strategy() string {
	return this.StrategyValue
}

func (this *RouterConfig) Settings() interface{} {
	settings := CreateRouterConfig(this.Strategy())
	err := json.Unmarshal(this.SettingsValue, settings)
	if err != nil {
		log.Error("Failed to load router settings: %v", err)
		return nil
	}
	return settings
}
