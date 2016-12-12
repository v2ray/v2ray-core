package bufio

import (
	"io"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

// BufferedReader is a reader with internal cache.
type BufferedReader struct {
	sync.Mutex
	reader io.Reader
	buffer *buf.Buffer
	cached bool
}

// NewReader creates a new BufferedReader based on an io.Reader.
func NewReader(rawReader io.Reader) *BufferedReader {
	return &BufferedReader{
		reader: rawReader,
		buffer: buf.New(),
		cached: true,
	}
}

// Release implements Releasable.Release().
func (v *BufferedReader) Release() {
	v.Lock()
	defer v.Unlock()

	v.buffer.Release()
	v.buffer = nil
	if releasable, ok := v.reader.(common.Releasable); ok {
		releasable.Release()
	}
	v.reader = nil
}

// Cached returns true if the internal cache is effective.
func (v *BufferedReader) Cached() bool {
	return v.cached
}

// SetCached is to enable or disable internal cache. If cache is disabled,
// Read() and Write() calls will be delegated to the underlying io.Reader directly.
func (v *BufferedReader) SetCached(cached bool) {
	v.cached = cached
}

// Read implements io.Reader.Read().
func (v *BufferedReader) Read(b []byte) (int, error) {
	v.Lock()
	defer v.Unlock()

	if v.reader == nil {
		return 0, io.EOF
	}

	if !v.cached {
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
