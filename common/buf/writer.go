package buf

import (
	"io"
	"net"
)

// BufferToBytesWriter is a Writer that writes alloc.Buffer into underlying writer.
type BufferToBytesWriter struct {
	writer io.Writer
}

// Write implements Writer.Write(). Write() takes ownership of the given buffer.
func (v *BufferToBytesWriter) Write(mb MultiBuffer) error {
	if mw, ok := v.writer.(MultiBufferWriter); ok {
		_, err := mw.WriteMultiBuffer(mb)
		return err
	}

	defer mb.Release()
	bs := make([][]byte, len(mb))
	for i, b := range mb {
		bs[i] = b.Bytes()
	}
	nbs := net.Buffers(bs)
	_, err := nbs.WriteTo(v.writer)
	return err
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
