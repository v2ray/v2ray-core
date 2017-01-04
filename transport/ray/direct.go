package ray

import (
	"errors"
	"io"

	"time"

	"v2ray.com/core/common/buf"
)

const (
	bufferSize = 512
)

var ErrReadTimeout = errors.New("Ray: timeout.")

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

func (v *directRay) AddInspector(inspector Inspector) {
	if inspector == nil {
		return
	}
	v.Input.inspector.AddInspector(inspector)
	v.Output.inspector.AddInspector(inspector)
}

type Stream struct {
	buffer    chan *buf.Buffer
	srcClose  chan bool
	destClose chan bool
	inspector *InspectorChain
}

func NewStream() *Stream {
	return &Stream{
		buffer:    make(chan *buf.Buffer, bufferSize),
		srcClose:  make(chan bool),
		destClose: make(chan bool),
		inspector: &InspectorChain{},
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
		case <-v.destClose:
			return nil, io.ErrClosedPipe
		}
	}
}

func (v *Stream) ReadTimeout(timeout time.Duration) (*buf.Buffer, error) {
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
		case <-v.destClose:
			return nil, io.ErrClosedPipe
		case <-time.After(timeout):
			return nil, ErrReadTimeout
		}
	}
}

func (v *Stream) Write(data *buf.Buffer) (err error) {
	if data.IsEmpty() {
		return
	}

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
			v.inspector.Input(data)
			return nil
		}
	}
}

func (v *Stream) Close() {
	defer swallowPanic()

	close(v.srcClose)
}

func (v *Stream) ForceClose() {
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

func (v *Stream) Release() {}

func swallowPanic() {
	recover()
}
