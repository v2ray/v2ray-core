package buf

import (
	"io"

	"v2ray.com/core/common"
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
		if b.IsFull() && largeSize > Size {
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
	if r.buffer == nil || largeSize == Size {
		return r.readSmall()
	}

	nBytes, err := r.Reader.Read(r.buffer)
	if nBytes > 0 {
		mb := NewMultiBufferCap(int32(nBytes/Size) + 1)
		common.Must2(mb.Write(r.buffer[:nBytes]))
		if nBytes == len(r.buffer) && nBytes < int(largeSize) {
			freeBytes(r.buffer)
			r.buffer = newBytes(int32(nBytes) + 1)
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
	// Reader is the underlying reader to be read from
	Reader Reader
	// Buffer is the internal buffer to be read from first
	Buffer MultiBuffer
	// Direct indicates whether or not to use the internal buffer
	Direct bool
}

// BufferedBytes returns the number of bytes that is cached in this reader.
func (r *BufferedReader) BufferedBytes() int32 {
	return r.Buffer.Len()
}

// ReadByte implements io.ByteReader.
func (r *BufferedReader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	return b[0], err
}

// Read implements io.Reader. It reads from internal buffer first (if available) and then reads from the underlying reader.
func (r *BufferedReader) Read(b []byte) (int, error) {
	if !r.Buffer.IsEmpty() {
		nBytes, err := r.Buffer.Read(b)
		common.Must(err)
		if r.Buffer.IsEmpty() {
			r.Buffer.Release()
			r.Buffer = nil
		}
		return nBytes, nil
	}

	if r.Direct {
		if reader, ok := r.Reader.(io.Reader); ok {
			return reader.Read(b)
		}
	}

	mb, err := r.Reader.ReadMultiBuffer()
	if err != nil {
		return 0, err
	}

	nBytes, err := mb.Read(b)
	common.Must(err)
	if !mb.IsEmpty() {
		r.Buffer = mb
	}
	return nBytes, err
}

// ReadMultiBuffer implements Reader.
func (r *BufferedReader) ReadMultiBuffer() (MultiBuffer, error) {
	if !r.Buffer.IsEmpty() {
		mb := r.Buffer
		r.Buffer = nil
		return mb, nil
	}

	return r.Reader.ReadMultiBuffer()
}

// ReadAtMost returns a MultiBuffer with at most size.
func (r *BufferedReader) ReadAtMost(size int32) (MultiBuffer, error) {
	if r.Buffer.IsEmpty() {
		mb, err := r.Reader.ReadMultiBuffer()
		if mb.IsEmpty() && err != nil {
			return nil, err
		}
		r.Buffer = mb
	}

	mb := r.Buffer.SliceBySize(size)
	if r.Buffer.IsEmpty() {
		r.Buffer = nil
	}
	return mb, nil
}

func (r *BufferedReader) writeToInternal(writer io.Writer) (int64, error) {
	mbWriter := NewWriter(writer)
	totalBytes := int64(0)
	if r.Buffer != nil {
		totalBytes += int64(r.Buffer.Len())
		if err := mbWriter.WriteMultiBuffer(r.Buffer); err != nil {
			return 0, err
		}
		r.Buffer = nil
	}

	for {
		mb, err := r.Reader.ReadMultiBuffer()
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

// Close implements io.Closer.
func (r *BufferedReader) Close() error {
	if !r.Buffer.IsEmpty() {
		r.Buffer.Release()
	}
	return common.Close(r.Reader)
}
