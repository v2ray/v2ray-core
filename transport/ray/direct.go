package ray

import (
	"context"
	"io"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/platform"
	"v2ray.com/core/common/signal"
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

// Stream is a sequential container for data in bytes.
type Stream struct {
	access      sync.RWMutex
	data        buf.MultiBuffer
	size        uint64
	ctx         context.Context
	readSignal  *signal.Notifier
	writeSignal *signal.Notifier
	close       bool
	err         bool
}

// NewStream creates a new Stream.
func NewStream(ctx context.Context) *Stream {
	return &Stream{
		ctx:         ctx,
		readSignal:  signal.NewNotifier(),
		writeSignal: signal.NewNotifier(),
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

// Peek fills in the given buffer with data from head of the Stream.
func (s *Stream) Peek(b *buf.Buffer) {
	s.access.RLock()
	defer s.access.RUnlock()

	common.Must(b.Reset(func(data []byte) (int, error) {
		return s.data.Copy(data), nil
	}))
}

// ReadMultiBuffer reads data from the Stream.
func (s *Stream) ReadMultiBuffer() (buf.MultiBuffer, error) {
	for {
		mb, err := s.getData()
		if err != nil {
			return nil, err
		}

		if mb != nil {
			s.readSignal.Signal()
			return mb, nil
		}

		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-s.writeSignal.Wait():
		}
	}
}

// ReadTimeout reads from the Stream with a specified timeout.
func (s *Stream) ReadTimeout(timeout time.Duration) (buf.MultiBuffer, error) {
	for {
		mb, err := s.getData()
		if err != nil {
			return nil, err
		}

		if mb != nil {
			s.readSignal.Signal()
			return mb, nil
		}

		select {
		case <-s.ctx.Done():
			return nil, io.EOF
		case <-time.After(timeout):
			return nil, buf.ErrReadTimeout
		case <-s.writeSignal.Wait():
		}
	}
}

// Size returns the number of bytes hold in the Stream.
func (s *Stream) Size() uint64 {
	s.access.RLock()
	defer s.access.RUnlock()

	return s.size
}

// waitForStreamSize waits until the Stream has room for more data, or any error happens.
func (s *Stream) waitForStreamSize() error {
	if streamSizeLimit == 0 {
		return nil
	}

	for s.Size() >= streamSizeLimit {
		select {
		case <-s.ctx.Done():
			return io.ErrClosedPipe
		case <-s.readSignal.Wait():
			if s.err || s.close {
				return io.ErrClosedPipe
			}
		}
	}

	return nil
}

// WriteMultiBuffer writes more data into the Stream.
func (s *Stream) WriteMultiBuffer(data buf.MultiBuffer) error {
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
		s.data = buf.NewMultiBufferCap(128)
	}

	s.data.AppendMulti(data)
	s.size += uint64(data.Len())
	s.writeSignal.Signal()

	return nil
}

// Close closes the stream for writing. Read() still works until EOF.
func (s *Stream) Close() error {
	s.access.Lock()
	s.close = true
	s.readSignal.Signal()
	s.writeSignal.Signal()
	s.access.Unlock()
	return nil
}

// CloseError closes the Stream with error. Read() will return an error afterwards.
func (s *Stream) CloseError() {
	s.access.Lock()
	s.err = true
	if s.data != nil {
		s.data.Release()
		s.data = nil
		s.size = 0
	}
	s.readSignal.Signal()
	s.writeSignal.Signal()
	s.access.Unlock()
}
