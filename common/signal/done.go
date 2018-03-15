package signal

import (
	"sync"
)

// Done is an utility for notifications of something being done.
type Done struct {
	access sync.Mutex
	c      chan struct{}
	closed bool
}

// NewDone returns a new Done.
func NewDone() *Done {
	return &Done{
		c: make(chan struct{}),
	}
}

// Done returns true if Close() is called.
func (d *Done) Done() bool {
	select {
	case <-d.c:
		return true
	default:
		return false
	}
}

// C returns a channel for waiting for done.
func (d *Done) C() chan struct{} {
	return d.c
}

// Wait blocks until Close() is called.
func (d *Done) Wait() {
	<-d.c
}

// Close marks this Done 'done'. This method may be called multiple times. All calls after first call will have no effect on its status.
func (d *Done) Close() error {
	d.access.Lock()
	defer d.access.Unlock()

	if d.closed {
		return nil
	}

	d.closed = true
	close(d.c)

	return nil
}
