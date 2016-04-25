package signal

type CancelSignal struct {
	cancel chan struct{}
	done   chan struct{}
}

func NewCloseSignal() *CancelSignal {
	return &CancelSignal{
		cancel: make(chan struct{}),
		done:   make(chan struct{}),
	}
}

func (this *CancelSignal) Cancel() {
	close(this.cancel)
}

func (this *CancelSignal) WaitForCancel() <-chan struct{} {
	return this.cancel
}

func (this *CancelSignal) Done() {
	close(this.done)
}

func (this *CancelSignal) WaitForDone() <-chan struct{} {
	return this.done
}
