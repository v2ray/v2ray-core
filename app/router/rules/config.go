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
	Rules         []*Rule
	ResolveDomain bool
}
