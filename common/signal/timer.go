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

func (t *ActivityTimer) run(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(<-t.timeout)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		case timeout := <-t.timeout:
			if timeout == 0 {
				cancel()
				return
			}

			ticker.Stop()
			ticker = time.NewTicker(timeout)
			continue
		}

		select {
		case <-t.updated:
		// Updated keep waiting.
		default:
			cancel()
			return
		}
	}
}

func CancelAfterInactivity(ctx context.Context, cancel context.CancelFunc, timeout time.Duration) *ActivityTimer {
	timer := &ActivityTimer{
		timeout: make(chan time.Duration, 1),
		updated: make(chan bool, 1),
	}
	timer.timeout <- timeout
	go timer.run(ctx, cancel)
	return timer
}
