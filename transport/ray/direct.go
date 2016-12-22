package ray

import (
	"io"

	"v2ray.com/core/common/buf"
)

const (
	bufferSize = 512
)

// NewRay creates a new Ray for direct traffic transport.
func NewRay() Ray {
	return &directRay{
		Input:  NewStream(),
		Output: NewStream(),
	}
}

type directRay struct {
	Input  *Stream
	Output *Stream
}

func (v *directRay) OutboundInput() InputStream {
	return v.Input
}

func (v *directRay) OutboundOutput() OutputStream {
	return v.Output
}

func (v *directRay) InboundInput() OutputStream {
	return v.Input
}

func (v *directRay) InboundOutput() InputStream {
	return v.Output
}

type Stream struct {
	buffer chan *buf.Buffer
}

func NewStream() *Stream {
	return &Stream{
		buffer: make(chan *buf.Buffer, bufferSize),
	}
}

func (v *Stream) Read() (*buf.Buffer, error) {
	buffer, open := <-v.buffer
	if !open {
		return nil, io.EOF
	}
	return buffer, nil
}

func (v *Stream) Write(data *buf.Buffer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = io.ErrClosedPipe
		}
	}()

	v.buffer <- data
	return nil
}

func (v *Stream) Close() {
	defer swallowPanic()

	close(v.buffer)
}

func (v *Stream) Release() {
	defer swallowPanic()

	close(v.buffer)

	for b := range v.buffer {
		b.Release()
	}
}

func swallowPanic() {
	recover()
}
