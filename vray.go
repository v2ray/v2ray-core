package core

type VRay struct {
	Input  chan []byte
	Output chan []byte
}

func NewVRay() *VRay {
	ray := new(VRay)
	ray.Input = make(chan []byte, 128)
	ray.Output = make(chan []byte, 128)
	return ray
}

type OutboundVRay interface {
	OutboundInput() <-chan []byte
	OutboundOutput() chan<- []byte
}

type InboundVRay interface {
	InboundInput() chan<- []byte
	OutboundOutput() <-chan []byte
}

func (ray *VRay) OutboundInput() <-chan []byte {
	return ray.Input
}

func (ray *VRay) OutboundOutput() chan<- []byte {
	return ray.Output
}

func (ray *VRay) InboundInput() chan<- []byte {
	return ray.Input
}

func (ray *VRay) InboundOutput() <-chan []byte {
	return ray.Output
}
