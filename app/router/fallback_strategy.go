package router

import (
	"fmt"
	"v2ray.com/core/features/outbound"
)

// FallbackStrategy is a balancer strategy that chooses the next alive outbound in the list when an outbound is dead
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
	if recorder, ok := handler.(outbound.FailedAttemptsRecorder); !ok {
		newError(fmt.Sprintf("invalid tag %s for fallback balancer", handler.Tag())).AtWarning().WriteToLog()
	} else {
		attempts := recorder.GetFailedAttempts()
		if attempts >= s.maxAttempts {
			recorder.ResetFailedAttempts()
			s.curIndex = (s.curIndex + 1) % len(s.tags)
			newError("balancer: switched to fallback " + s.tags[s.curIndex]).AtInfo().WriteToLog()
		}
	}
	return s.tags[s.curIndex]
}
