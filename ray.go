package core

const (
  bufferSize = 16
)

type Ray struct {
	Input  chan []byte
	Output chan []byte
}

func NewRay() Ray {
	return Ray{make(chan []byte, bufferSize), make(chan []byte, bufferSize)}
}

type OutboundRay interface {
	OutboundInput() <-chan []byte
	OutboundOutput() chan<- []byte
}

type InboundRay interface {
	InboundInput() chan<- []byte
	InboundOutput() <-chan []byte
}

func (ray Ray) OutboundInput() <-chan []byte {
	return ray.Input
}

func (ray Ray) OutboundOutput() chan<- []byte {
	return ray.Output
}

func (ray Ray) InboundInput() chan<- []byte {
	return ray.Input
}

func (ray Ray) InboundOutput() <-chan []byte {
	return ray.Output
}
