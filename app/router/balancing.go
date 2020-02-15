// +build !confonly

package router

import (
	"v2ray.com/core/common/dice"
	"v2ray.com/core/features/outbound"
)

type BalancingStrategy interface {
	PickOutbound(outbound.Manager, []string) string
}

type RandomStrategy struct {
}

func (s *RandomStrategy) PickOutbound(_ outbound.Manager, tags []string) string {
	n := len(tags)
	if n == 0 {
		panic("0 tags")
	}

	return tags[dice.Roll(n)]
}

type Balancer struct {
	selectors []string
	strategy  BalancingStrategy
	ohm       outbound.Manager
}

func (b *Balancer) PickOutbound() (string, error) {
	hs, ok := b.ohm.(outbound.HandlerSelector)
	if !ok {
		return "", newError("outbound.Manager is not a HandlerSelector")
	}
	tags := hs.Select(b.selectors)
	if len(tags) == 0 {
		return "", newError("no available outbounds selected")
	}
	tag := b.strategy.PickOutbound(b.ohm, tags)
	if tag == "" {
		return "", newError("balancing strategy returns empty tag")
	}
	return tag, nil
}
