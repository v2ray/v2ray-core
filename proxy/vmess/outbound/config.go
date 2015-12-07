package outbound

type Config interface {
	Receivers() []*Receiver
}
