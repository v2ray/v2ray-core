package signal

import (
	"sync"
)

type Done struct {
	access sync.Mutex
	c      chan struct{}
	closed bool
}

func NewDone() *Done {
	return &Done{
		c: make(chan struct{}),
	}
}

func (d *Done) Done() bool {
	select {
	case <-d.c:
		return true
	default:
		return false
	}
}

func (d *Done) C() chan struct{} {
	return d.c
}

func (d *Done) Wait() {
	<-d.c
}

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
