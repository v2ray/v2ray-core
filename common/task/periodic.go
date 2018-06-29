package task

import (
	"sync"
	"time"
)

// Periodic is a task that runs periodically.
type Periodic struct {
	// Interval of the task being run
	Interval time.Duration
	// Execute is the task function
	Execute func() error
	// OnFailure will be called when Execute returns non-nil error
	OnError func(error)

	access sync.RWMutex
	timer  *time.Timer
	closed bool
}

func (t *Periodic) setClosed(f bool) {
	t.access.Lock()
	t.closed = f
	t.access.Unlock()
}

func (t *Periodic) hasClosed() bool {
	t.access.RLock()
	defer t.access.RUnlock()

	return t.closed
}

func (t *Periodic) checkedExecute() error {
	if t.hasClosed() {
		return nil
	}

	if err := t.Execute(); err != nil {
		return err
	}

	t.access.Lock()
	defer t.access.Unlock()

	if t.closed {
		return nil
	}

	t.timer = time.AfterFunc(t.Interval, func() {
		if err := t.checkedExecute(); err != nil && t.OnError != nil {
			t.OnError(err)
		}
	})

	return nil
}

// Start implements common.Runnable. Start must not be called multiple times without Close being called.
func (t *Periodic) Start() error {
	t.setClosed(false)

	if err := t.checkedExecute(); err != nil {
		t.setClosed(true)
		return err
	}

	return nil
}

// Close implements common.Closable.
func (t *Periodic) Close() error {
	t.access.Lock()
	defer t.access.Unlock()

	t.closed = true
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	return nil
}
