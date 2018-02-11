package signal

import (
	"sync"
	"time"
)

type PeriodicTask struct {
	Interval time.Duration
	Execute  func() error

	access sync.Mutex
	timer  *time.Timer
	closed bool
}

func (t *PeriodicTask) checkedExecute() error {
	t.access.Lock()
	defer t.access.Unlock()

	if t.closed {
		return nil
	}

	if err := t.Execute(); err != nil {
		return err
	}

	t.timer = time.AfterFunc(t.Interval, func() {
		t.checkedExecute()
	})

	return nil
}

func (t *PeriodicTask) Start() error {
	t.access.Lock()
	t.closed = false
	t.access.Unlock()

	if err := t.checkedExecute(); err != nil {
		t.closed = true
		return err
	}

	return nil
}

func (t *PeriodicTask) Close() error {
	t.access.Lock()
	defer t.access.Unlock()

	t.closed = true
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}

	return nil
}
