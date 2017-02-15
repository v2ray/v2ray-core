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

// NewReader creates a new BufferedReader based on an io.Reader.
func NewBufferedReader(rawReader io.Reader) *BufferedReader {
	return &BufferedReader{
		reader:   rawReader,
		buffer:   NewLocal(1024),
		buffered: true,
	}
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
		err := v.buffer.AppendSupplier(ReadFrom(v.reader))
		if err != nil {
			return 0, err
		}
	}

	if v.buffer.IsEmpty() {
		return 0, nil
	}

	return v.buffer.Read(b)
}
