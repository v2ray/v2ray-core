package buf

import (
	"io"
	"time"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/signal"
)

// Reader extends io.Reader with alloc.Buffer.
type Reader interface {
	// Read reads content from underlying reader, and put it into an alloc.Buffer.
	Read() (MultiBuffer, error)
}

var ErrReadTimeout = newError("IO timeout")

type TimeoutReader interface {
	ReadTimeout(time.Duration) (MultiBuffer, error)
}

// Writer extends io.Writer with alloc.Buffer.
type Writer interface {
	// Write writes an alloc.Buffer into underlying writer.
	Write(MultiBuffer) error
}

// ReadFrom creates a Supplier to read from a given io.Reader.
func ReadFrom(reader io.Reader) Supplier {
	return func(b []byte) (int, error) {
		return reader.Read(b)
	}
}

// ReadFullFrom creates a Supplier to read full buffer from a given io.Reader.
func ReadFullFrom(reader io.Reader, size int) Supplier {
	return func(b []byte) (int, error) {
		return io.ReadFull(reader, b[:size])
	}
}

func copyInternal(timer signal.ActivityTimer, reader Reader, writer Writer) error {
	for {
		buffer, err := reader.Read()
		if err != nil {
			return err
		}

		timer.Update()

		if buffer.IsEmpty() {
			buffer.Release()
			continue
		}

		err = writer.Write(buffer)
		if err != nil {
			buffer.Release()
			return err
		}
	}
}

// Copy dumps all payload from reader to writer or stops when an error occurs.
// ActivityTimer gets updated as soon as there is a payload.
func Copy(timer signal.ActivityTimer, reader Reader, writer Writer) error {
	err := copyInternal(timer, reader, writer)
	if err != nil && errors.Cause(err) != io.EOF {
		return err
	}
	return nil
}

// NewReader creates a new Reader.
// The Reader instance doesn't take the ownership of reader.
func NewReader(reader io.Reader) Reader {
	if mr, ok := reader.(MultiBufferReader); ok {
		return &readerAdpater{
			MultiBufferReader: mr,
		}
	}

	return &BytesToBufferReader{
		reader: reader,
		buffer: NewLocal(32 * 1024),
	}
}

// ToBytesReader converts a Reaaer to io.Reader.
func ToBytesReader(stream Reader) io.Reader {
	return &bufferToBytesReader{
		stream: stream,
	}
}

// NewWriter creates a new Writer.
func NewWriter(writer io.Writer) Writer {
	return &BufferToBytesWriter{
		writer: writer,
	}
}

// ToBytesWriter converts a Writer to io.Writer
func ToBytesWriter(writer Writer) io.Writer {
	return &bytesToBufferWriter{
		writer: writer,
	}
}
