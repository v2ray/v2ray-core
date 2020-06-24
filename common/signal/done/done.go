package done

import (
	"sync"
)

// Instance is a utility for notifications of something being done.
type Instance struct {
	access sync.Mutex
	c      chan struct{}
	closed bool
}

// New returns a new Done.
func New() *Instance {
	return &Instance{
		c: make(chan struct{}),
	}
}

// Done returns true if Close() is called.
func (d *Instance) Done() bool {
	select {
	case <-d.Wait():
		return true
	default:
		return false
	}
}

// Wait returns a channel for waiting for done.
func (d *Instance) Wait() <-chan struct{} {
	return d.c
}

// Close marks this Done 'done'. This method may be called multiple times. All calls after first call will have no effect on its status.
func (d *Instance) Close() error {
	d.access.Lock()
	defer d.access.Unlock()

	if d.closed {
		return nil
	}

	d.closed = true
	close(d.c)

	return nil
}
