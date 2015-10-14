package ray

import (
	"github.com/v2ray/v2ray-core/common/alloc"
)

const (
	bufferSize = 16
)

func NewRay() Ray {
	return &directRay{
		Input:  make(chan *alloc.Buffer, bufferSize),
		Output: make(chan *alloc.Buffer, bufferSize),
	}
}

type directRay struct {
	Input  chan *alloc.Buffer
	Output chan *alloc.Buffer
}

func (ray *directRay) OutboundInput() <-chan *alloc.Buffer {
	return ray.Input
}

func (ray *directRay) OutboundOutput() chan<- *alloc.Buffer {
	return ray.Output
}

func (ray *directRay) InboundInput() chan<- *alloc.Buffer {
	return ray.Input
}

func (ray *directRay) InboundOutput() <-chan *alloc.Buffer {
	return ray.Output
}
