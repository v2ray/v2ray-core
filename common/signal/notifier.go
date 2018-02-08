package signal

type Notifier struct {
	c chan struct{}
}

func NewNotifier() *Notifier {
	return &Notifier{
		c: make(chan struct{}, 1),
	}
}

func (n *Notifier) Signal() {
	select {
	case n.c <- struct{}{}:
	default:
	}
}

func (n *Notifier) Wait() <-chan struct{} {
	return n.c
}
