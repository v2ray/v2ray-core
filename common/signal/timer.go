package signal

import (
	"context"
	"time"
)

type ActivityTimer struct {
	updated chan bool
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

func (t *ActivityTimer) Update() {
	select {
	case t.updated <- true:
	default:
	}
}

func (t *ActivityTimer) run() {
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

func CancelAfterInactivity(ctx context.Context, timeout time.Duration) (context.Context, *ActivityTimer) {
	ctx, cancel := context.WithCancel(ctx)
	timer := &ActivityTimer{
		ctx:     ctx,
		cancel:  cancel,
		timeout: timeout,
		updated: make(chan bool, 1),
	}
	go timer.run()
	return ctx, timer
}
