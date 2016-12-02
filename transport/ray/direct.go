package ray

import (
	"io"
	"sync"
	"time"

	"v2ray.com/core/common/alloc"
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
	access sync.RWMutex
	closed bool
	buffer chan *alloc.Buffer
}

func NewStream() *Stream {
	return &Stream{
		buffer: make(chan *alloc.Buffer, bufferSize),
	}
}

func (v *Stream) Read() (*alloc.Buffer, error) {
	if v.buffer == nil {
		return nil, io.EOF
	}
	v.access.RLock()
	if v.buffer == nil {
		v.access.RUnlock()
		return nil, io.EOF
	}
	channel := v.buffer
	v.access.RUnlock()
	result, open := <-channel
	if !open {
		return nil, io.EOF
	}
	return result, nil
}

func (v *Stream) Write(data *alloc.Buffer) error {
	for !v.closed {
		err := v.TryWriteOnce(data)
		if err != io.ErrNoProgress {
			return err
		}
	}
	return io.ErrClosedPipe
}

func (v *Stream) TryWriteOnce(data *alloc.Buffer) error {
	v.access.RLock()
	defer v.access.RUnlock()
	if v.closed {
		return io.ErrClosedPipe
	}
	select {
	case v.buffer <- data:
		return nil
	case <-time.After(2 * time.Second):
		return io.ErrNoProgress
	}
}

func (v *Stream) Close() {
	if v.closed {
		return
	}
	v.access.Lock()
	defer v.access.Unlock()
	if v.closed {
		return
	}
	v.closed = true
	close(v.buffer)
}

func (v *Stream) Release() {
	if v.buffer == nil {
		return
	}
	v.Close()
	v.access.Lock()
	defer v.access.Unlock()
	if v.buffer == nil {
		return
	}
	for data := range v.buffer {
		data.Release()
	}
	v.buffer = nil
}
