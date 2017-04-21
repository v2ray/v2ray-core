package buf

import "io"

// BufferedWriter is an io.Writer with internal buffer. It writes to underlying writer when buffer is full or on demand.
// This type is not thread safe.
type BufferedWriter struct {
	writer   io.Writer
	buffer   *Buffer
	buffered bool
}

// NewBufferedWriter creates a new BufferedWriter.
func NewBufferedWriter(rawWriter io.Writer) *BufferedWriter {
	return &BufferedWriter{
		writer:   rawWriter,
		buffer:   NewLocal(1024),
		buffered: true,
	}
}

// Write implements io.Writer.
func (v *BufferedWriter) Write(b []byte) (int, error) {
	if !v.buffered || v.buffer == nil {
		return v.writer.Write(b)
	}
	nBytes, err := v.buffer.Write(b)
	if err != nil {
		return 0, err
	}
	if v.buffer.IsFull() {
		if err := v.Flush(); err != nil {
			return 0, err
		}
		if nBytes < len(b) {
			if _, err := v.writer.Write(b[nBytes:]); err != nil {
				return nBytes, err
			}
		}
	}
	return len(b), nil
}

// Flush writes all buffered content into underlying writer, if any.
func (v *BufferedWriter) Flush() error {
	defer v.buffer.Clear()
	for !v.buffer.IsEmpty() {
		nBytes, err := v.writer.Write(v.buffer.Bytes())
		if err != nil {
			return err
		}
		v.buffer.SliceFrom(nBytes)
	}
	return nil
}

// IsBuffered returns true if this BufferedWriter holds a buffer.
func (v *BufferedWriter) IsBuffered() bool {
	return v.buffered
}

// SetBuffered controls whether the BufferedWriter holds a buffer for writing. If not buffered, any write() calls into underlying writer directly.
func (v *BufferedWriter) SetBuffered(cached bool) error {
	v.buffered = cached
	if !cached && !v.buffer.IsEmpty() {
		return v.Flush()
	}
	return nil
}
