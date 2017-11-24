package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

var (
	_ Reader    = (*BytesToBufferReader)(nil)
	_ io.Reader = (*BytesToBufferReader)(nil)
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
	if err := b.Reset(ReadFrom(r.Reader)); err != nil {
		b.Release()
		return nil, err
	}
	if b.IsFull() {
		r.buffer = make([]byte, 32*1024)
	}
	return NewMultiBufferValue(b), nil
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
		return mb, err
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
	stream       Reader
	legacyReader io.Reader
	leftOver     MultiBuffer
	buffered     bool
}

func NewBufferedReader(reader Reader) *BufferedReader {
	r := &BufferedReader{
		stream:   reader,
		buffered: true,
	}
	if lr, ok := reader.(io.Reader); ok {
		r.legacyReader = lr
	}
	return r
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

	if !r.buffered && r.legacyReader != nil {
		return r.legacyReader.Read(b)
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
