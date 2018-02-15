package signal

type Semaphore struct {
	token chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	s := &Semaphore{
		token: make(chan struct{}, n),
	}
	for i := 0; i < n; i++ {
		s.token <- struct{}{}
	}
	return s
}

func (s *Semaphore) Wait() <-chan struct{} {
	return s.token
}

func (s *Semaphore) Signal() {
	s.token <- struct{}{}
}
