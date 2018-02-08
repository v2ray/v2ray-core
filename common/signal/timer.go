package signal

import (
	"context"
	"time"
)

type ActivityUpdater interface {
	Update()
}

type ActivityTimer struct {
	updated chan struct{}
	timeout chan time.Duration
	closing chan struct{}
}

func (t *ActivityTimer) Update() {
	select {
	case t.updated <- struct{}{}:
	default:
	}
}

func (t *ActivityTimer) SetTimeout(timeout time.Duration) {
	select {
	case <-t.closing:
	case t.timeout <- timeout:
	}
}

func (t *ActivityTimer) run(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		close(t.closing)
	}()

	timeout := <-t.timeout
	if timeout == 0 {
		return
	}

	ticker := time.NewTicker(timeout)
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
			return
		}
	}
}

func CancelAfterInactivity(ctx context.Context, cancel context.CancelFunc, timeout time.Duration) *ActivityTimer {
	timer := &ActivityTimer{
		timeout: make(chan time.Duration, 1),
		updated: make(chan struct{}, 1),
		closing: make(chan struct{}),
	}
	timer.timeout <- timeout
	go timer.run(ctx, cancel)
	return timer
}
