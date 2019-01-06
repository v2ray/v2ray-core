package buf

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
)

func readOneUDP(r io.Reader) (*Buffer, error) {
	b := New()
	for i := 0; i < 64; i++ {
		_, err := b.ReadFrom(r)
		if !b.IsEmpty() {
			return b, nil
		}
		if err != nil {
			b.Release()
			return nil, err
		}
	}

	b.Release()
	return nil, newError("Reader returns too many empty payloads.")
}

// ReadBuffer reads a Buffer from the given reader, without allocating large buffer in advance.
func ReadBuffer(r io.Reader) (*Buffer, error) {
	// Use an one-byte buffer to wait for incoming payload.
	var firstByte [1]byte
	nBytes, err := r.Read(firstByte[:])
	if err != nil {
		return nil, err
	}

	b := New()
	if nBytes > 0 {
		common.Must(b.WriteByte(firstByte[0]))
	}
	for i := 0; i < 64; i++ {
		_, err := b.ReadFrom(r)
		if !b.IsEmpty() {
			return b, nil
		}
		if err != nil {
			b.Release()
			return nil, err
		}
	}

	b.Release()
	return nil, newError("Reader returns too many empty payloads.")
}

// BufferedReader is a Reader that keeps its internal buffer.
type BufferedReader struct {
	// Reader is the underlying reader to be read from
	Reader Reader
	// Buffer is the internal buffer to be read from first
	Buffer MultiBuffer
	// Spliter is a function to read bytes from MultiBuffer
	Spliter func(MultiBuffer, []byte) (MultiBuffer, int)
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
	spliter := r.Spliter
	if spliter == nil {
		spliter = SplitBytes
	}

	if !r.Buffer.IsEmpty() {
		buffer, nBytes := spliter(r.Buffer, b)
		r.Buffer = buffer
		if r.Buffer.IsEmpty() {
			r.Buffer = nil
		}
		return nBytes, nil
	}

	mb, err := r.Reader.ReadMultiBuffer()
	if err != nil {
		return 0, err
	}

	mb, nBytes := spliter(mb, b)
	if !mb.IsEmpty() {
		r.Buffer = mb
	}
	return nBytes, nil
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

	rb, mb := SplitSize(r.Buffer, size)
	r.Buffer = rb
	if r.Buffer.IsEmpty() {
		r.Buffer = nil
	}
	return mb, nil
}

func (r *BufferedReader) writeToInternal(writer io.Writer) (int64, error) {
	mbWriter := NewWriter(writer)
	var sc SizeCounter
	if r.Buffer != nil {
		sc.Size = int64(r.Buffer.Len())
		if err := mbWriter.WriteMultiBuffer(r.Buffer); err != nil {
			return 0, err
		}
		r.Buffer = nil
	}

	err := Copy(r.Reader, mbWriter, CountSize(&sc))
	return sc.Size, err
}

// WriteTo implements io.WriterTo.
func (r *BufferedReader) WriteTo(writer io.Writer) (int64, error) {
	nBytes, err := r.writeToInternal(writer)
	if errors.Cause(err) == io.EOF {
		return nBytes, nil
	}
	return nBytes, err
}

// Interrupt implements common.Interruptible.
func (r *BufferedReader) Interrupt() {
	common.Interrupt(r.Reader)
}

// Close implements io.Closer.
func (r *BufferedReader) Close() error {
	return common.Close(r.Reader)
}

// SingleReader is a Reader that read one Buffer every time.
type SingleReader struct {
	io.Reader
}

// ReadMultiBuffer implements Reader.
func (r *SingleReader) ReadMultiBuffer() (MultiBuffer, error) {
	b, err := ReadBuffer(r.Reader)
	if err != nil {
		return nil, err
	}
	return MultiBuffer{b}, nil
}

// PacketReader is a Reader that read one Buffer every time.
type PacketReader struct {
	io.Reader
}

// ReadMultiBuffer implements Reader.
func (r *PacketReader) ReadMultiBuffer() (MultiBuffer, error) {
	b, err := readOneUDP(r.Reader)
	if err != nil {
		return nil, err
	}
	return MultiBuffer{b}, nil
}
