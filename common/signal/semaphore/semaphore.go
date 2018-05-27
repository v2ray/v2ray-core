package semaphore

// Instance is an implementation of semaphore.
type Instance struct {
	token chan struct{}
}

// New create a new Semaphore with n permits.
func New(n int) *Instance {
	s := &Instance{
		token: make(chan struct{}, n),
	}
	for i := 0; i < n; i++ {
		s.token <- struct{}{}
	}
	return s
}

// Wait returns a channel for acquiring a permit.
func (s *Instance) Wait() <-chan struct{} {
	return s.token
}

// Signal releases a permit into the semaphore.
func (s *Instance) Signal() {
	s.token <- struct{}{}
}
