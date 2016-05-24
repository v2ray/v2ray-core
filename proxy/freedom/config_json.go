// +build json

package freedom

import (
	"encoding/json"
	"strings"

	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		DomainStrategy string `json:"domainStrategy"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.DomainStrategy = DomainStrategyAsIs
	domainStrategy := strings.ToLower(jsonConfig.DomainStrategy)
	if domainStrategy == "useip" {
		this.DomainStrategy = DomainStrategyUseIP
	}
	return nil
}

func init() {
	config.RegisterOutboundConfig("freedom",
		func(data []byte) (interface{}, error) {
			c := new(Config)
			if err := json.Unmarshal(data, c); err != nil {
				return nil, err
			}
			return c, nil
		})
}
