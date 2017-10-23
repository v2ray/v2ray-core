package ray

import (
	"context"
	"io"
	"sync"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform"
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

var streamSizeLimit uint64 = 10 * 1024 * 1024

func init() {
	const raySizeEnvKey = "v2ray.ray.buffer.size"
	size := platform.EnvFlag{
		Name:    raySizeEnvKey,
		AltName: platform.NormalizeEnvName(raySizeEnvKey),
	}.GetValueAsInt(10)
	streamSizeLimit = uint64(size) * 1024 * 1024
}

type Stream struct {
	access      sync.RWMutex
	data        buf.MultiBuffer
	size        uint64
	ctx         context.Context
	readSignal  chan bool
	writeSignal chan bool
	close       bool
	err         bool
}

func NewStream(ctx context.Context) *Stream {
	return &Stream{
		ctx:         ctx,
		readSignal:  make(chan bool, 1),
		writeSignal: make(chan bool, 1),
		size:        0,
	}
}

func (s *Stream) getData() (buf.MultiBuffer, error) {
	s.access.Lock()
	defer s.access.Unlock()

	if s.data != nil {
		mb := s.data
		s.data = nil
		s.size = 0
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

func (s *Stream) Peek(b *buf.Buffer) {
	s.access.RLock()
	defer s.access.RUnlock()

	b.Reset(func(data []byte) (int, error) {
		return s.data.Copy(data), nil
	})
}

func (s *Stream) Read() (buf.MultiBuffer, error) {
	for {
		mb, err := s.getData()
		if err != nil {
			return nil, err
		}

		if mb != nil {
			s.notifyRead()
			return mb, nil
		}

		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-s.writeSignal:
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
			s.notifyRead()
			return mb, nil
		}

		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-time.After(timeout):
			return nil, buf.ErrReadTimeout
		case <-s.writeSignal:
		}
	}
}

// Size returns the number of bytes hold in the Stream.
func (s *Stream) Size() uint64 {
	s.access.RLock()
	defer s.access.RUnlock()

	return s.size
}

func (s *Stream) waitForStreamSize() error {
	if streamSizeLimit == 0 {
		return nil
	}

	for s.Size() >= streamSizeLimit {
		select {
		case <-s.ctx.Done():
			return io.ErrClosedPipe
		case <-s.readSignal:
			if s.err || s.close {
				return io.ErrClosedPipe
			}
		}
	}

	return nil
}

func (s *Stream) Write(data buf.MultiBuffer) error {
	if data.IsEmpty() {
		return nil
	}

	if err := s.waitForStreamSize(); err != nil {
		data.Release()
		return err
	}

	s.access.Lock()
	defer s.access.Unlock()

	if s.err || s.close {
		data.Release()
		return io.ErrClosedPipe
	}

	if s.data == nil {
		s.data = data
	} else {
		s.data.AppendMulti(data)
	}
	s.size += uint64(data.Len())
	s.notifyWrite()

	return nil
}

func (s *Stream) notifyRead() {
	select {
	case s.readSignal <- true:
	default:
	}
}

func (s *Stream) notifyWrite() {
	select {
	case s.writeSignal <- true:
	default:
	}
}

func (s *Stream) Close() {
	s.access.Lock()
	s.close = true
	s.notifyRead()
	s.notifyWrite()
	s.access.Unlock()
}

func (s *Stream) CloseError() {
	s.access.Lock()
	s.err = true
	if s.data != nil {
		s.data.Release()
		s.data = nil
		s.size = 0
	}
	s.notifyRead()
	s.notifyWrite()
	s.access.Unlock()
}
