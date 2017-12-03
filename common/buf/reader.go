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

func NewBytesToBufferReader(reader io.Reader) Reader {
	return &BytesToBufferReader{
		Reader: reader,
	}
}

func (r *BytesToBufferReader) readSmall() (MultiBuffer, error) {
	b := New()
	err := b.Reset(ReadFrom(r.Reader))
	if b.IsFull() {
		r.buffer = make([]byte, 32*1024)
	}
	if !b.IsEmpty() {
		return NewMultiBufferValue(b), nil
	}
	b.Release()
	return nil, err
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
		return mb, nil
	}
	return nil, err
}

var (
	_ Reader        = (*BufferedReader)(nil)
	_ io.Reader     = (*BufferedReader)(nil)
	_ io.ByteReader = (*BufferedReader)(nil)
	_ io.WriterTo   = (*BufferedReader)(nil)
)

type BufferedReader struct {
	stream   Reader
	leftOver MultiBuffer
	buffered bool
}

func NewBufferedReader(reader Reader) *BufferedReader {
	return &BufferedReader{
		stream:   reader,
		buffered: true,
	}
}

func (r *BufferedReader) SetBuffered(f bool) {
	r.buffered = f
}

func (r *BufferedReader) IsBuffered() bool {
	return r.buffered
}

func (r *BufferedReader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	return b[0], err
}

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

	mb := r.leftOver.SliceBySize(size)
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

func (r *BufferedReader) WriteTo(writer io.Writer) (int64, error) {
	nBytes, err := r.writeToInternal(writer)
	if errors.Cause(err) == io.EOF {
		return nBytes, nil
	}
	return nBytes, err
}
