package testing

import (
	"github.com/v2ray/v2ray-core/app/router/rules"
)

type RouterRuleConfig struct {
	RuleList []*TestRule
}

func (this *RouterRuleConfig) Rules() []rules.Rule {
	rules := make([]rules.Rule, len(this.RuleList))
	for idx, rule := range this.RuleList {
		rules[idx] = rule
	}
	return rules
}
