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

func (v *CancelSignal) WaitThread() {
	v.done.Add(1)
}

// Cancel signals the goroutine to stop.
func (v *CancelSignal) Cancel() {
	close(v.cancel)
}

func (v *CancelSignal) Cancelled() bool {
	select {
	case <-v.cancel:
		return true
	default:
		return false
	}
}

// WaitForCancel should be monitored by the goroutine for when to stop.
func (v *CancelSignal) WaitForCancel() <-chan struct{} {
	return v.cancel
}

// FinishThread signals that current goroutine has finished.
func (v *CancelSignal) FinishThread() {
	v.done.Done()
}

// WaitForDone is used by caller to wait for the goroutine finishes.
func (v *CancelSignal) WaitForDone() {
	v.done.Wait()
}
