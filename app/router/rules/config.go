package rules

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Rule struct {
	Tag       string
	Condition Condition
}

func (this *Rule) Apply(dest v2net.Destination) bool {
	return this.Condition.Apply(dest)
}

type RouterRuleConfig struct {
	rules []*Rule
}

func NewRouterRuleConfig() *RouterRuleConfig {
	return &RouterRuleConfig{
		rules: make([]*Rule, 0, 16),
	}
}

func (this *RouterRuleConfig) Add(rule *Rule) *RouterRuleConfig {
	this.rules = append(this.rules, rule)
	return this
}

func (this *RouterRuleConfig) Rules() []*Rule {
	return this.rules
}
