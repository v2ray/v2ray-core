package signal

import (
	"context"
	"time"
)

type ActivityUpdater interface {
	Update()
}

type ActivityTimer struct {
	updated chan bool
	timeout chan time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

func (t *ActivityTimer) Update() {
	select {
	case t.updated <- true:
	default:
	}
}

func (t *ActivityTimer) SetTimeout(timeout time.Duration) {
	t.timeout <- timeout
}

func (t *ActivityTimer) run() {
	ticker := time.NewTicker(<-t.timeout)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
		case <-t.ctx.Done():
			return
		case timeout := <-t.timeout:
			if timeout == 0 {
				t.cancel()
				return
			}

			ticker.Stop()
			ticker = time.NewTicker(timeout)
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

func CancelAfterInactivity(ctx context.Context, cancel context.CancelFunc, timeout time.Duration) *ActivityTimer {
	timer := &ActivityTimer{
		ctx:     ctx,
		cancel:  cancel,
		timeout: make(chan time.Duration, 1),
		updated: make(chan bool, 1),
	}
	timer.timeout <- timeout
	go timer.run()
	return timer
}
