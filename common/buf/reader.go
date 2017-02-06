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

type BufferToBytesReader struct {
	stream  Reader
	current *Buffer
	err     error
}

// Fill fills in the internal buffer.
// Private: Visible for testing.
func (v *BufferToBytesReader) Fill() {
	b, err := v.stream.Read()
	v.current = b
	if err != nil {
		v.err = err
		v.current = nil
	}
}

func (v *BufferToBytesReader) Read(b []byte) (int, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.current == nil {
		v.Fill()
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
