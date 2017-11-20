package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

// BufferToBytesWriter is a Writer that writes alloc.Buffer into underlying writer.
type BufferToBytesWriter struct {
	io.Writer
}

func NewBufferToBytesWriter(writer io.Writer) *BufferToBytesWriter {
	return &BufferToBytesWriter{
		Writer: writer,
	}
}

// WriteMultiBuffer implements Writer. This method takes ownership of the given buffer.
func (w *BufferToBytesWriter) WriteMultiBuffer(mb MultiBuffer) error {
	defer mb.Release()

	bs := mb.ToNetBuffers()
	_, err := bs.WriteTo(w)
	return err
}

func (w *BufferToBytesWriter) ReadFrom(reader io.Reader) (int64, error) {
	if readerFrom, ok := w.Writer.(io.ReaderFrom); ok {
		return readerFrom.ReadFrom(reader)
	}

	var sc SizeCounter
	err := Copy(NewReader(reader), w, CountSize(&sc))
	return sc.Size, err
}

type BufferedWriter struct {
	writer       Writer
	legacyWriter io.Writer
	buffer       *Buffer
	buffered     bool
}

func NewBufferedWriter(writer Writer) *BufferedWriter {
	w := &BufferedWriter{
		writer:   writer,
		buffer:   New(),
		buffered: true,
	}
	if lw, ok := writer.(io.Writer); ok {
		w.legacyWriter = lw
	}
	return w
}

func (w *BufferedWriter) Write(b []byte) (int, error) {
	if !w.buffered && w.legacyWriter != nil {
		return w.legacyWriter.Write(b)
	}

	totalBytes := 0
	for len(b) > 0 {
		if w.buffer == nil {
			w.buffer = New()
		}

		nBytes, err := w.buffer.Write(b)
		totalBytes += nBytes
		if err != nil {
			return totalBytes, err
		}
		if !w.buffered || w.buffer.IsFull() {
			if err := w.Flush(); err != nil {
				return totalBytes, err
			}
		}
		b = b[nBytes:]
	}

	return totalBytes, nil
}

func (w *BufferedWriter) WriteMultiBuffer(b MultiBuffer) error {
	if !w.buffered {
		return w.writer.WriteMultiBuffer(b)
	}

	defer b.Release()

	for !b.IsEmpty() {
		if err := w.buffer.AppendSupplier(ReadFrom(&b)); err != nil {
			return err
		}
		if w.buffer.IsFull() {
			if err := w.Flush(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *BufferedWriter) Flush() error {
	if !w.buffer.IsEmpty() {
		if err := w.writer.WriteMultiBuffer(NewMultiBufferValue(w.buffer)); err != nil {
			return err
		}

		if w.buffered {
			w.buffer = New()
		} else {
			w.buffer = nil
		}
	}
	return nil
}

func (w *BufferedWriter) SetBuffered(f bool) error {
	w.buffered = f
	if !f {
		return w.Flush()
	}
	return nil
}

// ReadFrom implements io.ReaderFrom.
func (w *BufferedWriter) ReadFrom(reader io.Reader) (int64, error) {
	var sc SizeCounter
	if !w.buffer.IsEmpty() {
		sc.Size += int64(w.buffer.Len())
		if err := w.Flush(); err != nil {
			return sc.Size, err
		}
	}

	w.buffered = false
	err := Copy(NewReader(reader), w, CountSize(&sc))

	return sc.Size, err
}

type seqWriter struct {
	writer io.Writer
}

func (w *seqWriter) WriteMultiBuffer(mb MultiBuffer) error {
	defer mb.Release()

	for _, b := range mb {
		if b.IsEmpty() {
			continue
		}
		if _, err := w.writer.Write(b.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

type noOpWriter struct{}

func (noOpWriter) WriteMultiBuffer(b MultiBuffer) error {
	b.Release()
	return nil
}

type noOpBytesWriter struct{}

func (noOpBytesWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (noOpBytesWriter) ReadFrom(reader io.Reader) (int64, error) {
	b := New()
	defer b.Release()

	totalBytes := int64(0)
	for {
		err := b.Reset(ReadFrom(reader))
		totalBytes += int64(b.Len())
		if err != nil {
			if errors.Cause(err) == io.EOF {
				return totalBytes, nil
			}
			return totalBytes, err
		}
	}
}

var (
	// Discard is a Writer that swallows all contents written in.
	Discard Writer = noOpWriter{}

	// DiscardBytes is an io.Writer that swallows all contents written in.
	DiscardBytes io.Writer = noOpBytesWriter{}
)
