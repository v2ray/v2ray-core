package ray

import (
	"github.com/v2ray/v2ray-core/common/alloc"
)

const (
	bufferSize = 128
)

// NewRay creates a new Ray for direct traffic transport.
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

func (this *directRay) OutboundInput() <-chan *alloc.Buffer {
	return this.Input
}

func (this *directRay) OutboundOutput() chan<- *alloc.Buffer {
	return this.Output
}

func (this *directRay) InboundInput() chan<- *alloc.Buffer {
	return this.Input
}

func (this *directRay) InboundOutput() <-chan *alloc.Buffer {
	return this.Output
}
