package json

import (
	"encoding/json"

	v2routerconfigjson "github.com/v2ray/v2ray-core/app/router/config/json"
	"github.com/v2ray/v2ray-core/app/router/rules/config"
	"github.com/v2ray/v2ray-core/common/log"
)

type RouterRuleConfig struct {
	RuleList []json.RawMessage `json:"rules"`
}

func parseRule(msg json.RawMessage) config.Rule {
	rule := new(Rule)
	err := json.Unmarshal(msg, rule)
	if err != nil {
		log.Error("Invalid router rule: %v", err)
		return nil
	}
	if rule.Type == "field" {
		fieldrule := new(FieldRule)
		err = json.Unmarshal(msg, fieldrule)
		if err != nil {
			log.Error("Invalid field rule: %v", err)
			return nil
		}
		return fieldrule
	}
	log.Error("Unknown router rule type: %s", rule.Type)
	return nil
}

func (this *RouterRuleConfig) Rules() []config.Rule {
	rules := make([]config.Rule, len(this.RuleList))
	for idx, rawRule := range this.RuleList {
		rules[idx] = parseRule(rawRule)
	}
	return rules
}

func init() {
	v2routerconfigjson.RegisterRouterConfig("rules", func() interface{} {
		return new(RouterRuleConfig)
	})
}
