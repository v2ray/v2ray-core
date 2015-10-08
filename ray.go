package core

import (
	"github.com/v2ray/v2ray-core/common/alloc"
)

const (
	bufferSize = 16
)

// Ray is an internal tranport channel bewteen inbound and outbound connection.
type Ray struct {
	Input  chan *alloc.Buffer
	Output chan *alloc.Buffer
}

func NewRay() *Ray {
	return &Ray{
		Input:  make(chan *alloc.Buffer, bufferSize),
		Output: make(chan *alloc.Buffer, bufferSize),
	}
}

type OutboundRay interface {
	OutboundInput() <-chan *alloc.Buffer
	OutboundOutput() chan<- *alloc.Buffer
}

type InboundRay interface {
	InboundInput() chan<- *alloc.Buffer
	InboundOutput() <-chan *alloc.Buffer
}

func (ray *Ray) OutboundInput() <-chan *alloc.Buffer {
	return ray.Input
}

func (ray *Ray) OutboundOutput() chan<- *alloc.Buffer {
	return ray.Output
}

func (ray *Ray) InboundInput() chan<- *alloc.Buffer {
	return ray.Input
}

func (ray *Ray) InboundOutput() <-chan *alloc.Buffer {
	return ray.Output
}
