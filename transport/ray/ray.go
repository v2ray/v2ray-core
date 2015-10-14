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

// OutboundRay is a transport interface for outbound connections.
type OutboundRay interface {
	// OutboundInput provides a stream for the input of the outbound connection.
	// The outbound connection shall write all the input until it is closed.
	OutboundInput() <-chan *alloc.Buffer

	// OutboundOutput provides a stream to retrieve the response from the
	// outbound connection. The outbound connection shall close the channel
	// after all responses are receivced and put into the channel.
	OutboundOutput() chan<- *alloc.Buffer
}

// InboundRay is a transport interface for inbound connections.
type InboundRay interface {
	// InboundInput provides a stream to retrieve the request from client.
	// The inbound connection shall close the channel after the entire request
	// is received and put into the channel.
	InboundInput() chan<- *alloc.Buffer

	// InboudBound provides a stream of data for the inbound connection to write
	// as response. The inbound connection shall write all the data from the
	// channel until it is closed.
	InboundOutput() <-chan *alloc.Buffer
}

// Ray is an internal tranport channel bewteen inbound and outbound connection.
type Ray interface {
	InboundRay
	OutboundRay
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
