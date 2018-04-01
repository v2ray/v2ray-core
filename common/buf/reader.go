package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

// BytesToBufferReader is a Reader that adjusts its reading speed automatically.
type BytesToBufferReader struct {
	io.Reader
	buffer []byte
}

// NewBytesToBufferReader returns a new BytesToBufferReader.
func NewBytesToBufferReader(reader io.Reader) Reader {
	return &BytesToBufferReader{
		Reader: reader,
	}
}

func (r *BytesToBufferReader) readSmall() (MultiBuffer, error) {
	b := New()
	for i := 0; i < 64; i++ {
		err := b.Reset(ReadFrom(r.Reader))
		if b.IsFull() {
			r.buffer = newBytes(Size + 1)
		}
		if !b.IsEmpty() {
			return NewMultiBufferValue(b), nil
		}
		if err != nil {
			b.Release()
			return nil, err
		}
	}

	return nil, newError("Reader returns too many empty payloads.")
}

func (r *BytesToBufferReader) freeBuffer() {
	freeBytes(r.buffer)
	r.buffer = nil
}

// ReadMultiBuffer implements Reader.
func (r *BytesToBufferReader) ReadMultiBuffer() (MultiBuffer, error) {
	if r.buffer == nil {
		return r.readSmall()
	}

	nBytes, err := r.Reader.Read(r.buffer)
	if nBytes > 0 {
		mb := NewMultiBufferCap(nBytes/Size + 1)
		mb.Write(r.buffer[:nBytes])
		if nBytes == len(r.buffer) && nBytes < int(largeSize) {
			freeBytes(r.buffer)
			r.buffer = newBytes(uint32(nBytes) + 1)
		} else if nBytes < Size {
			r.freeBuffer()
		}
		return mb, nil
	}

	r.freeBuffer()

	if err != nil {
		return nil, err
	}

	// Read() returns empty payload and nil err. We don't expect this to happen, but just in case.
	return r.readSmall()
}

// BufferedReader is a Reader that keeps its internal buffer.
type BufferedReader struct {
	stream   Reader
	leftOver MultiBuffer
	buffered bool
}

// NewBufferedReader returns a new BufferedReader.
func NewBufferedReader(reader Reader) *BufferedReader {
	return &BufferedReader{
		stream:   reader,
		buffered: true,
	}
}

// SetBuffered sets whether to keep the interal buffer.
func (r *BufferedReader) SetBuffered(f bool) {
	r.buffered = f
}

// IsBuffered returns true if internal buffer is used.
func (r *BufferedReader) IsBuffered() bool {
	return r.buffered
}

// BufferedBytes returns the number of bytes that is cached in this reader.
func (r *BufferedReader) BufferedBytes() int32 {
	return int32(r.leftOver.Len())
}

// ReadByte implements io.ByteReader.
func (r *BufferedReader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	return b[0], err
}

// Read implements io.Reader. It reads from internal buffer first (if available) and then reads from the underlying reader.
func (r *BufferedReader) Read(b []byte) (int, error) {
	if r.leftOver != nil {
		nBytes, _ := r.leftOver.Read(b)
		if r.leftOver.IsEmpty() {
			r.leftOver.Release()
			r.leftOver = nil
		}
		return nBytes, nil
	}

	if !r.buffered {
		if reader, ok := r.stream.(io.Reader); ok {
			return reader.Read(b)
		}
	}

	mb, err := r.stream.ReadMultiBuffer()
	if mb != nil {
		nBytes, _ := mb.Read(b)
		if !mb.IsEmpty() {
			r.leftOver = mb
		}
		return nBytes, err
	}
	return 0, err
}

// ReadMultiBuffer implements Reader.
func (r *BufferedReader) ReadMultiBuffer() (MultiBuffer, error) {
	if r.leftOver != nil {
		mb := r.leftOver
		r.leftOver = nil
		return mb, nil
	}

	return r.stream.ReadMultiBuffer()
}

// ReadAtMost returns a MultiBuffer with at most size.
func (r *BufferedReader) ReadAtMost(size int) (MultiBuffer, error) {
	if r.leftOver == nil {
		mb, err := r.stream.ReadMultiBuffer()
		if mb.IsEmpty() && err != nil {
			return nil, err
		}
		r.leftOver = mb
	}

	mb := r.leftOver.SliceBySize(int32(size))
	if r.leftOver.IsEmpty() {
		r.leftOver = nil
	}
	return mb, nil
}

func (r *BufferedReader) writeToInternal(writer io.Writer) (int64, error) {
	mbWriter := NewWriter(writer)
	totalBytes := int64(0)
	if r.leftOver != nil {
		totalBytes += int64(r.leftOver.Len())
		if err := mbWriter.WriteMultiBuffer(r.leftOver); err != nil {
			return 0, err
		}
		r.leftOver = nil
	}

	for {
		mb, err := r.stream.ReadMultiBuffer()
		if mb != nil {
			totalBytes += int64(mb.Len())
			if werr := mbWriter.WriteMultiBuffer(mb); werr != nil {
				return totalBytes, err
			}
		}
		if err != nil {
			return totalBytes, err
		}
	}
}

// WriteTo implements io.WriterTo.
func (r *BufferedReader) WriteTo(writer io.Writer) (int64, error) {
	nBytes, err := r.writeToInternal(writer)
	if errors.Cause(err) == io.EOF {
		return nBytes, nil
	}
	return nBytes, err
}
