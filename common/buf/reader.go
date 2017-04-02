package buf

import "io"

// BytesToBufferReader is a Reader that adjusts its reading speed automatically.
type BytesToBufferReader struct {
	reader      io.Reader
	largeBuffer *Buffer
	highVolumn  bool
}

// Read implements Reader.Read().
func (v *BytesToBufferReader) Read() (*Buffer, error) {
	if v.highVolumn && v.largeBuffer.IsEmpty() {
		if v.largeBuffer == nil {
			v.largeBuffer = NewLocal(32 * 1024)
		}
		err := v.largeBuffer.AppendSupplier(ReadFrom(v.reader))
		if err != nil {
			return nil, err
		}
		if v.largeBuffer.Len() < Size {
			v.highVolumn = false
		}
	}

	buffer := New()
	if !v.largeBuffer.IsEmpty() {
		err := buffer.AppendSupplier(ReadFrom(v.largeBuffer))
		return buffer, err
	}

	err := buffer.AppendSupplier(ReadFrom(v.reader))
	if err != nil {
		buffer.Release()
		return nil, err
	}

	if buffer.IsFull() {
		v.highVolumn = true
	}

	return buffer, nil
}

type bufferToBytesReader struct {
	stream  Reader
	current *Buffer
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
