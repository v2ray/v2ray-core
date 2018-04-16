package pipe

import (
	"time"

	"v2ray.com/core/common/buf"
)

type Reader struct {
	pipe *pipe
}

func (r *Reader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	return r.pipe.ReadMultiBuffer()
}

func (r *Reader) ReadMultiBufferWithTimeout(d time.Duration) (buf.MultiBuffer, error) {
	return r.pipe.ReadMultiBufferWithTimeout(d)
}

func (r *Reader) CloseError() {
	r.pipe.CloseError()
}
