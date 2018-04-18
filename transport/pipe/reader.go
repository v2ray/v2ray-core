package pipe

import (
	"time"

	"v2ray.com/core/common/buf"
)

// Reader is a buf.Reader that reads content from a pipe.
type Reader struct {
	pipe *pipe
}

// ReadMultiBuffer implements buf.Reader.
func (r *Reader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	return r.pipe.ReadMultiBuffer()
}

// ReadMultiBufferWithTimeout reads content from a pipe within the given duration, or returns buf.ErrTimeout otherwise.
func (r *Reader) ReadMultiBufferWithTimeout(d time.Duration) (buf.MultiBuffer, error) {
	return r.pipe.ReadMultiBufferWithTimeout(d)
}

// CloseError sets the pipe to error state. Both reading and writing from/to the pipe will return io.ErrClosedPipe.
func (r *Reader) CloseError() {
	r.pipe.CloseError()
}
