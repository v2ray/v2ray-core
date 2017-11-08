package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

// BufferToBytesWriter is a Writer that writes alloc.Buffer into underlying writer.
type BufferToBytesWriter struct {
	writer io.Writer
}

// Write implements Writer.Write(). Write() takes ownership of the given buffer.
func (w *BufferToBytesWriter) Write(mb MultiBuffer) error {
	defer mb.Release()

	bs := mb.ToNetBuffers()
	_, err := bs.WriteTo(w.writer)
	return err
}

type writerAdapter struct {
	writer MultiBufferWriter
}

// Write implements buf.MultiBufferWriter.
func (w *writerAdapter) Write(mb MultiBuffer) error {
	return w.writer.WriteMultiBuffer(mb)
}

type mergingWriter struct {
	writer io.Writer
	buffer []byte
}

func (w *mergingWriter) Write(mb MultiBuffer) error {
	defer mb.Release()

	for !mb.IsEmpty() {
		nBytes, _ := mb.Read(w.buffer)
		if _, err := w.writer.Write(w.buffer[:nBytes]); err != nil {
			return err
		}
	}
	return nil
}

type seqWriter struct {
	writer io.Writer
}

func (w *seqWriter) Write(mb MultiBuffer) error {
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

var (
	_ MultiBufferWriter = (*bytesToBufferWriter)(nil)
)

type bytesToBufferWriter struct {
	writer Writer
}

// Write implements io.Writer.
func (w *bytesToBufferWriter) Write(payload []byte) (int, error) {
	mb := NewMultiBufferCap(len(payload)/Size + 1)
	mb.Write(payload)
	if err := w.writer.Write(mb); err != nil {
		return 0, err
	}
	return len(payload), nil
}

func (w *bytesToBufferWriter) WriteMultiBuffer(mb MultiBuffer) error {
	return w.writer.Write(mb)
}

func (w *bytesToBufferWriter) ReadFrom(reader io.Reader) (int64, error) {
	mbReader := NewReader(reader)
	totalBytes := int64(0)
	for {
		mb, err := mbReader.Read()
		if errors.Cause(err) == io.EOF {
			break
		} else if err != nil {
			return totalBytes, err
		}
		totalBytes += int64(mb.Len())
		if err := w.writer.Write(mb); err != nil {
			return totalBytes, err
		}
	}
	return totalBytes, nil
}

type noOpWriter struct{}

func (noOpWriter) Write(b MultiBuffer) error {
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
