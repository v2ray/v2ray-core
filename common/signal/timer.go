package signal

import (
	"context"
	"time"
)

type ActivityTimer interface {
	Update()
}

type realActivityTimer struct {
	updated chan bool
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

func (t *realActivityTimer) Update() {
	select {
	case t.updated <- true:
	default:
	}
}

func (t *realActivityTimer) run() {
	for {
		select {
		case <-time.After(t.timeout):
		case <-t.ctx.Done():
			return
		}
		select {
		case <-t.updated:
		// Updated keep waiting.
		default:
			t.cancel()
			return
		}

	}
}

func CancelAfterInactivity(ctx context.Context, timeout time.Duration) (context.Context, ActivityTimer) {
	ctx, cancel := context.WithCancel(ctx)
	timer := &realActivityTimer{
		ctx:     ctx,
		cancel:  cancel,
		timeout: timeout,
		updated: make(chan bool, 1),
	}
	go timer.run()
	return ctx, timer
}

type noOpActivityTimer struct{}

func (noOpActivityTimer) Update() {}

func BackgroundTimer() ActivityTimer {
	return noOpActivityTimer{}
}
