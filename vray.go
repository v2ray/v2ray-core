package core

type VRay struct {
	Input  chan []byte
	Output chan []byte
}

func NewVRay() VRay {
	return VRay{make(chan []byte, 128), make(chan []byte, 128)}
}

type OutboundVRay interface {
	OutboundInput() <-chan []byte
	OutboundOutput() chan<- []byte
}

type InboundVRay interface {
	InboundInput() chan<- []byte
	InboundOutput() <-chan []byte
}

func (ray VRay) OutboundInput() <-chan []byte {
	return ray.Input
}

func (ray VRay) OutboundOutput() chan<- []byte {
	return ray.Output
}

func (ray VRay) InboundInput() chan<- []byte {
	return ray.Input
}

func (ray VRay) InboundOutput() <-chan []byte {
	return ray.Output
}
