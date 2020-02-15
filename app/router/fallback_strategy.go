package router

import (
	"v2ray.com/core/features/outbound"
)

type FallbackStrategy struct {
	tags        []string
	curIndex    int
	maxAttempts int64
}

// NewFallbackStrategy returns a new instance of FallbackStrategy
func NewFallbackStrategy(maxAttempts int64) *FallbackStrategy {
	return &FallbackStrategy{
		tags:        nil,
		curIndex:    0,
		maxAttempts: maxAttempts,
	}
}

// PickOutbound picks an outbound with fallback strategy
func (s *FallbackStrategy) PickOutbound(ohm outbound.Manager, tags []string) string {
	if s.tags == nil {
		s.tags = tags
	}
	handler := ohm.GetHandler(s.tags[s.curIndex])
	attempts := handler.FailedAttempts()
	if attempts.Value() >= s.maxAttempts {
		attempts.Set(0)
		s.curIndex = (s.curIndex + 1) % len(s.tags)
		newError("balancer: switched to fallback " + s.tags[s.curIndex]).AtInfo().WriteToLog()
	}
	return s.tags[s.curIndex]
}
