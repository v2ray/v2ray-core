package rules

import (
	v2net "v2ray.com/core/common/net"
)

type Rule struct {
	Tag       string
	Condition Condition
}

func (this *Rule) Apply(dest v2net.Destination) bool {
	return this.Condition.Apply(dest)
}

type DomainStrategy int

var (
	DomainAsIs      = DomainStrategy(0)
	AlwaysUseIP     = DomainStrategy(1)
	UseIPIfNonMatch = DomainStrategy(2)
)

type RouterRuleConfig struct {
	Rules          []*Rule
	DomainStrategy DomainStrategy
}
