package ray

import (
	"context"
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
func NewRay(ctx context.Context) Ray {
	return &directRay{
		Input:  NewStream(ctx),
		Output: NewStream(ctx),
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
	ctx    context.Context
	close  chan bool
	err    chan bool
}

func NewStream(ctx context.Context) *Stream {
	return &Stream{
		ctx:    ctx,
		buffer: make(chan *buf.Buffer, bufferSize),
		close:  make(chan bool),
		err:    make(chan bool),
	}
}

func (v *Stream) Read() (*buf.Buffer, error) {
	select {
	case <-v.ctx.Done():
		return nil, io.ErrClosedPipe
	case <-v.err:
		return nil, io.ErrClosedPipe
	case b := <-v.buffer:
		return b, nil
	default:
		select {
		case <-v.ctx.Done():
			return nil, io.ErrClosedPipe
		case b := <-v.buffer:
			return b, nil
		case <-v.close:
			return nil, io.EOF
		case <-v.err:
			return nil, io.ErrClosedPipe
		}
	}
}

func (v *Stream) ReadTimeout(timeout time.Duration) (*buf.Buffer, error) {
	select {
	case <-v.ctx.Done():
		return nil, io.ErrClosedPipe
	case <-v.err:
		return nil, io.ErrClosedPipe
	case b := <-v.buffer:
		return b, nil
	default:
		select {
		case <-v.ctx.Done():
			return nil, io.ErrClosedPipe
		case b := <-v.buffer:
			return b, nil
		case <-v.close:
			return nil, io.EOF
		case <-v.err:
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
	case <-v.ctx.Done():
		return io.ErrClosedPipe
	case <-v.err:
		return io.ErrClosedPipe
	case <-v.close:
		return io.ErrClosedPipe
	default:
		select {
		case <-v.ctx.Done():
			return io.ErrClosedPipe
		case <-v.err:
			return io.ErrClosedPipe
		case <-v.close:
			return io.ErrClosedPipe
		case v.buffer <- data:
			return nil
		}
	}
}

func (v *Stream) Close() {
	defer swallowPanic()

	close(v.close)
}

func (v *Stream) CloseError() {
	defer swallowPanic()

	close(v.err)

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
