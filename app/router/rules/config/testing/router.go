package testing

import (
	"github.com/v2ray/v2ray-core/app/router/rules/config"
)

type RouterRuleConfig struct {
	RuleList []*TestRule
}

func (this *RouterRuleConfig) Rules() []config.Rule {
	rules := make([]config.Rule, len(this.RuleList))
	for idx, rule := range this.RuleList {
		rules[idx] = rule
	}
	return rules
}
