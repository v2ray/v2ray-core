package buf

import "io"

// BytesToBufferReader is a Reader that adjusts its reading speed automatically.
type BytesToBufferReader struct {
	reader io.Reader
	buffer *Buffer
}

// Read implements Reader.Read().
func (v *BytesToBufferReader) Read() (MultiBuffer, error) {
	if err := v.buffer.Reset(ReadFrom(v.reader)); err != nil {
		return nil, err
	}

	mb := NewMultiBuffer()
	for !v.buffer.IsEmpty() {
		b := New()
		b.AppendSupplier(ReadFrom(v.buffer))
		mb.Append(b)
	}
	return mb, nil
}

type bufferToBytesReader struct {
	stream  Reader
	current MultiBuffer
	err     error
}

// fill fills in the internal buffer.
func (v *bufferToBytesReader) fill() {
	b, err := v.stream.Read()
	if err != nil {
		v.err = err
		return
	}
	v.current = b
}

func (v *bufferToBytesReader) Read(b []byte) (int, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.current == nil {
		v.fill()
		if v.err != nil {
			return 0, v.err
		}
	}
	nBytes, err := v.current.Read(b)
	if v.current.IsEmpty() {
		v.current.Release()
		v.current = nil
	}
	return nBytes, err
}
