package buf

import (
	"io"

	"v2ray.com/core/common/bytespool"
)

const (
	// Size of a regular buffer.
	Size = 2048
)

// Buffer is a recyclable allocation of a byte array. Buffer.Release() recycles
// the buffer into an internal buffer pool, in order to recreate a buffer more
// quickly.
type Buffer struct {
	v     []byte
	start int32
	end   int32
}

// Release recycles the buffer into an internal buffer pool.
func (b *Buffer) Release() {
	if b == nil || b.v == nil {
		return
	}
	pool.Put(b.v)
	b.v = nil
	b.Clear()
}

// Clear clears the content of the buffer, results an empty buffer with
// Len() = 0.
func (b *Buffer) Clear() {
	b.start = 0
	b.end = 0
}

// Byte returns the bytes at index.
func (b *Buffer) Byte(index int32) byte {
	return b.v[b.start+index]
}

// SetByte sets the byte value at index.
func (b *Buffer) SetByte(index int32, value byte) {
	b.v[b.start+index] = value
}

// Bytes returns the content bytes of this Buffer.
func (b *Buffer) Bytes() []byte {
	return b.v[b.start:b.end]
}

// Extend increases the buffer size by n bytes, and returns the extended part.
// It panics if result size is larger than buf.Size.
func (b *Buffer) Extend(n int32) []byte {
	end := b.end + n
	if end > int32(len(b.v)) {
		panic(newError("out of bound: ", end))
	}
	ext := b.v[b.end:end]
	b.end = end
	return ext
}

// BytesRange returns a slice of this buffer with given from and to boundary.
func (b *Buffer) BytesRange(from, to int32) []byte {
	if from < 0 {
		from += b.Len()
	}
	if to < 0 {
		to += b.Len()
	}
	return b.v[b.start+from : b.start+to]
}

// BytesFrom returns a slice of this Buffer starting from the given position.
func (b *Buffer) BytesFrom(from int32) []byte {
	if from < 0 {
		from += b.Len()
	}
	return b.v[b.start+from : b.end]
}

// BytesTo returns a slice of this Buffer from start to the given position.
func (b *Buffer) BytesTo(to int32) []byte {
	if to < 0 {
		to += b.Len()
	}
	return b.v[b.start : b.start+to]
}

// Resize cuts the buffer at the given position.
func (b *Buffer) Resize(from, to int32) {
	if from < 0 {
		from += b.Len()
	}
	if to < 0 {
		to += b.Len()
	}
	if to < from {
		panic("Invalid slice")
	}
	b.end = b.start + to
	b.start += from
}

// Advance cuts the buffer at the given position.
func (b *Buffer) Advance(from int32) {
	if from < 0 {
		from += b.Len()
	}
	b.start += from
}

// Len returns the length of the buffer content.
func (b *Buffer) Len() int32 {
	if b == nil {
		return 0
	}
	return b.end - b.start
}

// IsEmpty returns true if the buffer is empty.
func (b *Buffer) IsEmpty() bool {
	return b.Len() == 0
}

// IsFull returns true if the buffer has no more room to grow.
func (b *Buffer) IsFull() bool {
	return b.end == int32(len(b.v))
}

// Write implements Write method in io.Writer.
func (b *Buffer) Write(data []byte) (int, error) {
	nBytes := copy(b.v[b.end:], data)
	b.end += int32(nBytes)
	return nBytes, nil
}

// WriteByte writes a single byte into the buffer.
func (b *Buffer) WriteByte(v byte) error {
	if b.IsFull() {
		return newError("buffer full")
	}
	b.v[b.end] = v
	b.end++
	return nil
}

// WriteString implements io.StringWriter.
func (b *Buffer) WriteString(s string) (int, error) {
	return b.Write([]byte(s))
}

// Read implements io.Reader.Read().
func (b *Buffer) Read(data []byte) (int, error) {
	if b.Len() == 0 {
		return 0, io.EOF
	}
	nBytes := copy(data, b.v[b.start:b.end])
	if int32(nBytes) == b.Len() {
		b.Clear()
	} else {
		b.start += int32(nBytes)
	}
	return nBytes, nil
}

// ReadFrom implements io.ReaderFrom.
func (b *Buffer) ReadFrom(reader io.Reader) (int64, error) {
	n, err := reader.Read(b.v[b.end:])
	b.end += int32(n)
	return int64(n), err
}

// ReadFullFrom reads exact size of bytes from given reader, or until error occurs.
func (b *Buffer) ReadFullFrom(reader io.Reader, size int32) (int64, error) {
	end := b.end + size
	if end > int32(len(b.v)) {
		return 0, newError("out of bound: ", end)
	}
	n, err := io.ReadFull(reader, b.v[b.end:end])
	b.end += int32(n)
	return int64(n), err
}

// String returns the string form of this Buffer.
func (b *Buffer) String() string {
	return string(b.Bytes())
}

var pool = bytespool.GetPool(Size)

// New creates a Buffer with 0 length and 2K capacity.
func New() *Buffer {
	return &Buffer{
		v: pool.Get().([]byte),
	}
}

// StackNew creates a new Buffer object on stack.
// This method is for buffers that is released in the same function.
func StackNew() Buffer {
	return Buffer{
		v: pool.Get().([]byte),
	}
}
