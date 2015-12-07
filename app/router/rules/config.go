package rules

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Rule interface {
	Tag() string
	Apply(dest v2net.Destination) bool
}

type RouterRuleConfig interface {
	Rules() []Rule
}
