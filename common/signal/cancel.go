package signal

import (
	"sync"
)

// CancelSignal is a signal passed to goroutine, in order to cancel the goroutine on demand.
type CancelSignal struct {
	cancel chan struct{}
	done   sync.WaitGroup
}

// NewCloseSignal creates a new CancelSignal.
func NewCloseSignal() *CancelSignal {
	return &CancelSignal{
		cancel: make(chan struct{}),
	}
}

func (this *CancelSignal) WaitThread() {
	this.done.Add(1)
}

// Cancel signals the goroutine to stop.
func (this *CancelSignal) Cancel() {
	close(this.cancel)
}

func (this *CancelSignal) Cancelled() bool {
	select {
	case <-this.cancel:
		return true
	default:
		return false
	}
}

// WaitForCancel should be monitored by the goroutine for when to stop.
func (this *CancelSignal) WaitForCancel() <-chan struct{} {
	return this.cancel
}

// FinishThread signals that current goroutine has finished.
func (this *CancelSignal) FinishThread() {
	this.done.Done()
}

// WaitForDone is used by caller to wait for the goroutine finishes.
func (this *CancelSignal) WaitForDone() {
	this.done.Wait()
}
