package signal

// Notifier is an utility for notifying changes. The change producer may notify changes multiple time, and the consumer may get notified asychronously.
type Notifier struct {
	c chan struct{}
}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		c: make(chan struct{}, 1),
	}
}

// Signal signals a change, usually by producer. This method never blocks.
func (n *Notifier) Signal() {
	select {
	case n.c <- struct{}{}:
	default:
	}
}

// Wait returns a channel for waiting for changes. The returned channel never gets closed.
func (n *Notifier) Wait() <-chan struct{} {
	return n.c
}
