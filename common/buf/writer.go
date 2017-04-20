package buf

import "io"

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

func (w *writerAdapter) Write(mb MultiBuffer) error {
	_, err := w.writer.WriteMultiBuffer(mb)
	return err
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

type bytesToBufferWriter struct {
	writer Writer
}

func (w *bytesToBufferWriter) Write(payload []byte) (int, error) {
	mb := NewMultiBuffer()
	for p := payload; len(p) > 0; {
		b := New()
		nBytes, _ := b.Write(p)
		p = p[nBytes:]
		mb.Append(b)
	}
	if err := w.writer.Write(mb); err != nil {
		return 0, err
	}
	return len(payload), nil
}

func (w *bytesToBufferWriter) WriteMulteBuffer(mb MultiBuffer) (int, error) {
	return mb.Len(), w.writer.Write(mb)
}

func (w *bytesToBufferWriter) ReadFrom(reader io.Reader) (int64, error) {
	mbReader := NewReader(reader)
	totalBytes := int64(0)
	eof := false
	for !eof {
		mb, err := mbReader.Read()
		if err == io.EOF {
			eof = true
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
