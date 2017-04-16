package ray

import (
	"context"
	"io"
	"sync"
	"time"

	"v2ray.com/core/common/buf"
)

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
	access sync.Mutex
	data   buf.MultiBuffer
	ctx    context.Context
	wakeup chan bool
	close  bool
	err    bool
}

func NewStream(ctx context.Context) *Stream {
	return &Stream{
		ctx:    ctx,
		wakeup: make(chan bool, 1),
	}
}

func (s *Stream) getData() (buf.MultiBuffer, error) {
	s.access.Lock()
	defer s.access.Unlock()

	if s.data != nil {
		mb := s.data
		s.data = nil
		return mb, nil
	}

	if s.close {
		return nil, io.EOF
	}

	if s.err {
		return nil, io.ErrClosedPipe
	}

	return nil, nil
}

func (s *Stream) Read() (buf.MultiBuffer, error) {
	for {
		mb, err := s.getData()
		if err != nil {
			return nil, err
		}

		if mb != nil {
			return mb, nil
		}

		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-s.wakeup:
		}
	}
}

func (s *Stream) ReadTimeout(timeout time.Duration) (buf.MultiBuffer, error) {
	for {
		mb, err := s.getData()
		if err != nil {
			return nil, err
		}

		if mb != nil {
			return mb, nil
		}

		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-time.After(timeout):
			return nil, buf.ErrReadTimeout
		case <-s.wakeup:
		}
	}
}

func (s *Stream) Write(data buf.MultiBuffer) (err error) {
	if data.IsEmpty() {
		return
	}

	s.access.Lock()
	defer s.access.Unlock()

	if s.err {
		data.Release()
		return io.ErrClosedPipe
	}
	if s.close {
		data.Release()
		return io.ErrClosedPipe
	}

	if s.data == nil {
		s.data = data
	} else {
		s.data.AppendMulti(data)
	}
	s.wakeUp()

	return nil
}

func (s *Stream) wakeUp() {
	select {
	case s.wakeup <- true:
	default:
	}
}

func (s *Stream) Close() {
	s.access.Lock()
	s.close = true
	s.wakeUp()
	s.access.Unlock()
}

func (s *Stream) CloseError() {
	s.access.Lock()
	s.err = true
	if s.data != nil {
		s.data.Release()
		s.data = nil
	}
	s.wakeUp()
	s.access.Unlock()
}
