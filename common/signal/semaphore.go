package signal

type Semaphore struct {
	token chan bool
}

func NewSemaphore(n int) *Semaphore {
	s := &Semaphore{
		token: make(chan bool, n),
	}
	for i := 0; i < n; i++ {
		s.token <- true
	}
	return s
}

func (s *Semaphore) Wait() <-chan bool {
	return s.token
}

func (s *Semaphore) Signal() {
	s.token <- true
}
