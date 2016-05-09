package ray

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/alloc"
)

const (
	bufferSize = 128
)

var (
	ErrorIOTimeout = errors.New("IO Timeout")
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

func (this *directRay) OutboundInput() InputStream {
	return this.Input
}

func (this *directRay) OutboundOutput() OutputStream {
	return this.Output
}

func (this *directRay) InboundInput() OutputStream {
	return this.Input
}

func (this *directRay) InboundOutput() InputStream {
	return this.Output
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

func (this *Stream) Read() (*alloc.Buffer, error) {
	if this.buffer == nil {
		return nil, io.EOF
	}
	this.access.RLock()
	if this.buffer == nil {
		this.access.RUnlock()
		return nil, io.EOF
	}
	channel := this.buffer
	this.access.RUnlock()
	result, open := <-channel
	if !open {
		return nil, io.EOF
	}
	return result, nil
}

func (this *Stream) Write(data *alloc.Buffer) error {
	if this.closed {
		return io.EOF
	}
	for {
		err := this.TryWriteOnce(data)
		if err != ErrorIOTimeout {
			return err
		}
	}
}

func (this *Stream) TryWriteOnce(data *alloc.Buffer) error {
	this.access.RLock()
	defer this.access.RUnlock()
	if this.closed {
		return io.EOF
	}
	select {
	case this.buffer <- data:
		return nil
	case <-time.Tick(time.Second):
		return ErrorIOTimeout
	}
}

func (this *Stream) Close() {
	if this.closed {
		return
	}
	this.access.Lock()
	defer this.access.Unlock()
	if this.closed {
		return
	}
	this.closed = true
	close(this.buffer)
}

func (this *Stream) Release() {
	if this.buffer == nil {
		return
	}
	this.Close()
	this.access.Lock()
	defer this.access.Unlock()
	if this.buffer == nil {
		return
	}
	for data := range this.buffer {
		data.Release()
	}
	this.buffer = nil
}
