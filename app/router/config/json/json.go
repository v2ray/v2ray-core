package json

import (
	"encoding/json"
)

type RouterConfig struct {
	StrategyValue string          `json:"strategy"`
	SettingsValue json.RawMessage `json:"settings"`
}

func (this *RouterConfig) Strategy() string {
	return this.StrategyValue
}

func (this *RouterConfig) Settings() interface{} {
	return CreateRouterConfig(this.Strategy())
}
