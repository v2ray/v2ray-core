package buf

import (
	"io"
	"net"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
)

// BufferToBytesWriter is a Writer that writes alloc.Buffer into underlying writer.
type BufferToBytesWriter struct {
	io.Writer

	cache [][]byte
}

// WriteMultiBuffer implements Writer. This method takes ownership of the given buffer.
func (w *BufferToBytesWriter) WriteMultiBuffer(mb MultiBuffer) error {
	defer ReleaseMulti(mb)

	size := mb.Len()
	if size == 0 {
		return nil
	}

	if len(mb) == 1 {
		return WriteAllBytes(w.Writer, mb[0].Bytes())
	}

	bs := w.cache
	for _, b := range mb {
		bs = append(bs, b.Bytes())
	}
	w.cache = bs

	defer func() {
		for idx := range w.cache {
			w.cache[idx] = nil
		}
		w.cache = w.cache[:0]
	}()

	nb := net.Buffers(bs)

	for size > 0 {
		n, err := nb.WriteTo(w.Writer)
		if err != nil {
			return err
		}
		size -= int32(n)
	}

	return nil
}

// ReadFrom implements io.ReaderFrom.
func (w *BufferToBytesWriter) ReadFrom(reader io.Reader) (int64, error) {
	var sc SizeCounter
	err := Copy(NewReader(reader), w, CountSize(&sc))
	return sc.Size, err
}

// BufferedWriter is a Writer with internal buffer.
type BufferedWriter struct {
	sync.Mutex
	writer   Writer
	buffer   *Buffer
	buffered bool
}

// NewBufferedWriter creates a new BufferedWriter.
func NewBufferedWriter(writer Writer) *BufferedWriter {
	return &BufferedWriter{
		writer:   writer,
		buffer:   New(),
		buffered: true,
	}
}

// WriteByte implements io.ByteWriter.
func (w *BufferedWriter) WriteByte(c byte) error {
	return common.Error2(w.Write([]byte{c}))
}

// Write implements io.Writer.
func (w *BufferedWriter) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	w.Lock()
	defer w.Unlock()

	if !w.buffered {
		if writer, ok := w.writer.(io.Writer); ok {
			return writer.Write(b)
		}
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
			if err := w.flushInternal(); err != nil {
				return totalBytes, err
			}
		}
		b = b[nBytes:]
	}

	return totalBytes, nil
}

// WriteMultiBuffer implements Writer. It takes ownership of the given MultiBuffer.
func (w *BufferedWriter) WriteMultiBuffer(b MultiBuffer) error {
	if b.IsEmpty() {
		return nil
	}

	w.Lock()
	defer w.Unlock()

	if !w.buffered {
		return w.writer.WriteMultiBuffer(b)
	}

	defer ReleaseMulti(b)

	for !b.IsEmpty() {
		if w.buffer == nil {
			w.buffer = New()
		}
		if _, err := w.buffer.ReadFrom(&b); err != nil {
			return err
		}
		if w.buffer.IsFull() {
			if err := w.flushInternal(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Flush flushes buffered content into underlying writer.
func (w *BufferedWriter) Flush() error {
	w.Lock()
	defer w.Unlock()

	return w.flushInternal()
}

func (w *BufferedWriter) flushInternal() error {
	if w.buffer.IsEmpty() {
		return nil
	}

	b := w.buffer
	w.buffer = nil

	if writer, ok := w.writer.(io.Writer); ok {
		err := WriteAllBytes(writer, b.Bytes())
		b.Release()
		return err
	}

	return w.writer.WriteMultiBuffer(MultiBuffer{b})
}

// SetBuffered sets whether the internal buffer is used. If set to false, Flush() will be called to clear the buffer.
func (w *BufferedWriter) SetBuffered(f bool) error {
	w.Lock()
	defer w.Unlock()

	w.buffered = f
	if !f {
		return w.flushInternal()
	}
	return nil
}

// ReadFrom implements io.ReaderFrom.
func (w *BufferedWriter) ReadFrom(reader io.Reader) (int64, error) {
	if err := w.SetBuffered(false); err != nil {
		return 0, err
	}

	var sc SizeCounter
	err := Copy(NewReader(reader), w, CountSize(&sc))
	return sc.Size, err
}

// Close implements io.Closable.
func (w *BufferedWriter) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return common.Close(w.writer)
}

// SequentialWriter is a Writer that writes MultiBuffer sequentially into the underlying io.Writer.
type SequentialWriter struct {
	io.Writer
}

// WriteMultiBuffer implements Writer.
func (w *SequentialWriter) WriteMultiBuffer(mb MultiBuffer) error {
	defer ReleaseMulti(mb)

	for _, b := range mb {
		if b.IsEmpty() {
			continue
		}

		if err := WriteAllBytes(w.Writer, b.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

type noOpWriter byte

func (noOpWriter) WriteMultiBuffer(b MultiBuffer) error {
	ReleaseMulti(b)
	return nil
}

func (noOpWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (noOpWriter) ReadFrom(reader io.Reader) (int64, error) {
	b := New()
	defer b.Release()

	totalBytes := int64(0)
	for {
		b.Clear()
		_, err := b.ReadFrom(reader)
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
	Discard Writer = noOpWriter(0)

	// DiscardBytes is an io.Writer that swallows all contents written in.
	DiscardBytes io.Writer = noOpWriter(0)
)
