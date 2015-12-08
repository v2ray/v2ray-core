package json

import (
	"encoding/json"

	v2routerjson "github.com/v2ray/v2ray-core/app/router/json"
	"github.com/v2ray/v2ray-core/app/router/rules"
	"github.com/v2ray/v2ray-core/common/log"
)

type RouterRuleConfig struct {
	RuleList []json.RawMessage `json:"rules"`
}

func parseRule(msg json.RawMessage) rules.Rule {
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
	if rule.Type == "chinaip" {
		chinaiprule := new(ChinaIPRule)
		if err := json.Unmarshal(msg, chinaiprule); err != nil {
			log.Error("Invalid chinaip rule: %v", err)
			return nil
		}
		return chinaiprule
	}
	log.Error("Unknown router rule type: %s", rule.Type)
	return nil
}

func (this *RouterRuleConfig) Rules() []rules.Rule {
	rules := make([]rules.Rule, len(this.RuleList))
	for idx, rawRule := range this.RuleList {
		rules[idx] = parseRule(rawRule)
	}
	return rules
}

func init() {
	v2routerjson.RegisterRouterConfig("rules", func() interface{} {
		return new(RouterRuleConfig)
	})
}
