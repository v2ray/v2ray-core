package bufio

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

// BufferedReader is a reader with internal cache.
type BufferedReader struct {
	reader   io.Reader
	buffer   *buf.Buffer
	buffered bool
}

// NewReader creates a new BufferedReader based on an io.Reader.
func NewReader(rawReader io.Reader) *BufferedReader {
	return &BufferedReader{
		reader:   rawReader,
		buffer:   buf.New(),
		buffered: true,
	}
}

// Release implements Releasable.Release().
func (v *BufferedReader) Release() {
	if v.buffer != nil {
		v.buffer.Release()
		v.buffer = nil
	}

	common.Release(v.reader)
}

// IsBuffered returns true if the internal cache is effective.
func (v *BufferedReader) IsBuffered() bool {
	return v.buffered
}

// SetBuffered is to enable or disable internal cache. If cache is disabled,
// Read() calls will be delegated to the underlying io.Reader directly.
func (v *BufferedReader) SetBuffered(cached bool) {
	v.buffered = cached
}

// Read implements io.Reader.Read().
func (v *BufferedReader) Read(b []byte) (int, error) {
	if !v.buffered || v.buffer == nil {
		if !v.buffer.IsEmpty() {
			return v.buffer.Read(b)
		}
		return v.reader.Read(b)
	}
	if v.buffer.IsEmpty() {
		err := v.buffer.AppendSupplier(buf.ReadFrom(v.reader))
		if err != nil {
			return 0, err
		}
	}

	if v.buffer.IsEmpty() {
		return 0, nil
	}

	return v.buffer.Read(b)
}
