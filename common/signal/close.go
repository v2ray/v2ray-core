package signal

// CancelSignal is a signal passed to goroutine, in order to cancel the goroutine on demand.
type CancelSignal struct {
	cancel chan struct{}
	done   chan struct{}
}

// NewCloseSignal creates a new CancelSignal.
func NewCloseSignal() *CancelSignal {
	return &CancelSignal{
		cancel: make(chan struct{}),
		done:   make(chan struct{}),
	}
}

// Cancel signals the goroutine to stop.
func (this *CancelSignal) Cancel() {
	close(this.cancel)
}

// WaitForCancel should be monitored by the goroutine for when to stop.
func (this *CancelSignal) WaitForCancel() <-chan struct{} {
	return this.cancel
}

// Done signals the caller that the goroutine has completely finished.
func (this *CancelSignal) Done() {
	close(this.done)
}

// WaitForDone is used by caller to wait for the goroutine finishes.
func (this *CancelSignal) WaitForDone() <-chan struct{} {
	return this.done
}
