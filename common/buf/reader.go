package buf

import "io"
import "v2ray.com/core/common"

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

// Release implements Releasable.Release().
func (v *BytesToBufferReader) Release() {
	common.Release(v.reader)
}

type BufferToBytesReader struct {
	stream  Reader
	current *Buffer
	eof     bool
}

// Fill fills in the internal buffer.
// Private: Visible for testing.
func (v *BufferToBytesReader) Fill() {
	b, err := v.stream.Read()
	v.current = b
	if err != nil {
		v.eof = true
		v.current = nil
	}
}

func (v *BufferToBytesReader) Read(b []byte) (int, error) {
	if v.eof {
		return 0, io.EOF
	}

	if v.current == nil {
		v.Fill()
		if v.eof {
			return 0, io.EOF
		}
	}
	nBytes, err := v.current.Read(b)
	if v.current.IsEmpty() {
		v.current.Release()
		v.current = nil
	}
	return nBytes, err
}

// Release implements Releasable.Release().
func (v *BufferToBytesReader) Release() {
	v.eof = true
	v.current.Release()
	v.current = nil
}
