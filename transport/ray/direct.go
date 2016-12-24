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
	buffer    chan *buf.Buffer
	srcClose  chan bool
	destClose chan bool
}

func NewStream() *Stream {
	return &Stream{
		buffer:    make(chan *buf.Buffer, bufferSize),
		srcClose:  make(chan bool),
		destClose: make(chan bool),
	}
}

func (v *Stream) Read() (*buf.Buffer, error) {
	select {
	case <-v.destClose:
		return nil, io.ErrClosedPipe
	case b := <-v.buffer:
		return b, nil
	default:
		select {
		case b := <-v.buffer:
			return b, nil
		case <-v.srcClose:
			return nil, io.EOF
		}
	}
}

func (v *Stream) Write(data *buf.Buffer) (err error) {
	select {
	case <-v.destClose:
		return io.ErrClosedPipe
	case <-v.srcClose:
		return io.ErrClosedPipe
	default:
		select {
		case <-v.destClose:
			return io.ErrClosedPipe
		case <-v.srcClose:
			return io.ErrClosedPipe
		case v.buffer <- data:
			return nil
		}
	}
}

func (v *Stream) Close() {
	defer swallowPanic()

	close(v.srcClose)
}

func (v *Stream) Release() {
	defer swallowPanic()

	close(v.destClose)
	v.Close()

	n := len(v.buffer)
	for i := 0; i < n; i++ {
		select {
		case b := <-v.buffer:
			b.Release()
		default:
			return
		}
	}
}

func swallowPanic() {
	recover()
}
