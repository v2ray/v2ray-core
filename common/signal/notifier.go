package signal

type Notifier struct {
	c chan bool
}

func NewNotifier() *Notifier {
	return &Notifier{
		c: make(chan bool, 1),
	}
}

func (n *Notifier) Signal() {
	select {
	case n.c <- true:
	default:
	}
}

func (n *Notifier) Wait() <-chan bool {
	return n.c
}
