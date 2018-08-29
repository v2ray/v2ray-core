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

	access  sync.Mutex
	timer   *time.Timer
	running bool
}

func (t *Periodic) hasClosed() bool {
	t.access.Lock()
	defer t.access.Unlock()

	return !t.running
}

func (t *Periodic) checkedExecute() error {
	if t.hasClosed() {
		return nil
	}

	if err := t.Execute(); err != nil {
		t.access.Lock()
		t.running = false
		t.access.Unlock()
		return err
	}

	t.access.Lock()
	defer t.access.Unlock()

	if !t.running {
		return nil
	}

	t.timer = time.AfterFunc(t.Interval, func() {
		t.checkedExecute() // nolint: errcheck
	})

	return nil
}

// Start implements common.Runnable.
func (t *Periodic) Start() error {
	t.access.Lock()
	if t.running {
		return nil
	}
	t.running = true
	t.access.Unlock()

	if err := t.checkedExecute(); err != nil {
		t.access.Lock()
		t.running = false
		t.access.Unlock()
		return err
	}

	return nil
}

// Close implements common.Closable.
func (t *Periodic) Close() error {
	t.access.Lock()
	defer t.access.Unlock()

	t.running = false
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	return nil
}
