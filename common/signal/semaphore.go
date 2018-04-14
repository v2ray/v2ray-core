package signal

// Semaphore is an implementation of semaphore.
type Semaphore struct {
	token chan struct{}
}

// NewSemaphore create a new Semaphore with n permits.
func NewSemaphore(n int) *Semaphore {
	s := &Semaphore{
		token: make(chan struct{}, n),
	}
	for i := 0; i < n; i++ {
		s.token <- struct{}{}
	}
	return s
}

// Wait returns a channel for acquiring a permit.
func (s *Semaphore) Wait() <-chan struct{} {
	return s.token
}

// Signal releases a permit into the Semaphore.
func (s *Semaphore) Signal() {
	s.token <- struct{}{}
}
