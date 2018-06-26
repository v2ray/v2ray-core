package signal

import "sync"

// Notifier is a utility for notifying changes. The change producer may notify changes multiple time, and the consumer may get notified asynchronously.
type Notifier struct {
	sync.Mutex
	waiters    []chan struct{}
	notCosumed bool
}

// NewNotifier creates a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Signal signals a change, usually by producer. This method never blocks.
func (n *Notifier) Signal() {
	n.Lock()
	defer n.Unlock()

	if len(n.waiters) == 0 {
		n.notCosumed = true
		return
	}

	for _, w := range n.waiters {
		close(w)
	}
	n.waiters = make([]chan struct{}, 0, 8)
}

// Wait returns a channel for waiting for changes.
func (n *Notifier) Wait() <-chan struct{} {
	n.Lock()
	defer n.Unlock()

	w := make(chan struct{})
	if n.notCosumed {
		n.notCosumed = false
		close(w)
		return w
	}

	n.waiters = append(n.waiters, w)
	return w
}
