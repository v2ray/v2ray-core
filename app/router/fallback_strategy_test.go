package router

import (
	"context"
	"testing"
	"v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/app/stats"
	"v2ray.com/core/common"
	stats2 "v2ray.com/core/features/stats"
	"v2ray.com/core/transport"
)

type handler struct {
	tag            string
	failedAttempts stats.Counter
}

func (h *handler) Tag() string {
	return h.tag
}

func (h *handler) Start() error {
	return nil
}

func (h *handler) Close() error {
	return nil
}

func (h *handler) FailedAttempts() stats2.Counter {
	return &h.failedAttempts
}

func (h *handler) Dispatch(context.Context, *transport.Link) {
	if h.tag == "dead" {
		h.failedAttempts.Add(1)
	}
}

func TestFallbackStrategy(t *testing.T) {
	ctx := context.Background()

	ohm, _ := outbound.New(ctx, nil)
	common.Must(ohm.AddHandler(ctx, &handler{tag: "dead"}))
	common.Must(ohm.AddHandler(ctx, &handler{tag: "alive"}))

	strategy := NewFallbackStrategy(2)

	expect := []string{"dead", "dead", "alive", "alive"}
	for i, exp := range expect {
		tag := strategy.PickOutbound(ohm, []string{"dead", "alive"})
		if expect[i] != tag {
			t.Errorf("For run %d, expected %s, got %s", i, exp, tag)
		}
		ohm.GetHandler(tag).Dispatch(ctx, nil)
	}
}
