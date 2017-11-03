package buf

import (
	"io"
)

// BufferedReader is a reader with internal cache.
type BufferedReader struct {
	reader   io.Reader
	buffer   *Buffer
	buffered bool
}

// NewBufferedReader creates a new BufferedReader based on an io.Reader.
func NewBufferedReader(rawReader io.Reader) *BufferedReader {
	return &BufferedReader{
		reader:   rawReader,
		buffer:   NewLocal(1024),
		buffered: true,
	}
}

// IsBuffered returns true if the internal cache is effective.
func (r *BufferedReader) IsBuffered() bool {
	return r.buffered
}

// SetBuffered is to enable or disable internal cache. If cache is disabled,
// Read() calls will be delegated to the underlying io.Reader directly.
func (r *BufferedReader) SetBuffered(cached bool) {
	r.buffered = cached
}

// Read implements io.Reader.Read().
func (r *BufferedReader) Read(b []byte) (int, error) {
	if !r.buffered || r.buffer == nil {
		if !r.buffer.IsEmpty() {
			return r.buffer.Read(b)
		}
		return r.reader.Read(b)
	}
	if r.buffer.IsEmpty() {
		if err := r.buffer.Reset(ReadFrom(r.reader)); err != nil {
			return 0, err
		}
	}

	if r.buffer.IsEmpty() {
		return 0, nil
	}

	return r.buffer.Read(b)
}
