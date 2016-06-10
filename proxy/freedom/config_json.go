// +build json

package freedom

import (
	"encoding/json"
	"strings"

	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		DomainStrategy string `json:"domainStrategy"`
		Timeout        uint32 `json:"timeout"`
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
	this.Timeout = jsonConfig.Timeout
	return nil
}

func init() {
	internal.RegisterOutboundConfig("freedom", func() interface{} { return new(Config) })
}
