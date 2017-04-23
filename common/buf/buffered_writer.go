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
func (w *BufferedWriter) Write(b []byte) (int, error) {
	if !w.buffered || w.buffer == nil {
		return w.writer.Write(b)
	}
	nBytes, err := w.buffer.Write(b)
	if err != nil {
		return 0, err
	}
	if w.buffer.IsFull() {
		if err := w.Flush(); err != nil {
			return 0, err
		}
		if nBytes < len(b) {
			if _, err := w.writer.Write(b[nBytes:]); err != nil {
				return nBytes, err
			}
		}
	}
	return len(b), nil
}

// Flush writes all buffered content into underlying writer, if any.
func (w *BufferedWriter) Flush() error {
	defer w.buffer.Clear()
	for !w.buffer.IsEmpty() {
		nBytes, err := w.writer.Write(w.buffer.Bytes())
		if err != nil {
			return err
		}
		w.buffer.SliceFrom(nBytes)
	}
	return nil
}

// IsBuffered returns true if this BufferedWriter holds a buffer.
func (w *BufferedWriter) IsBuffered() bool {
	return w.buffered
}

// SetBuffered controls whether the BufferedWriter holds a buffer for writing. If not buffered, any write() calls into underlying writer directly.
func (w *BufferedWriter) SetBuffered(cached bool) error {
	w.buffered = cached
	if !cached && !w.buffer.IsEmpty() {
		return w.Flush()
	}
	return nil
}
