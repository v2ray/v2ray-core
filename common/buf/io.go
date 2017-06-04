package buf

import (
	"io"
	"time"
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

// ReadAtLeastFrom create a Supplier to read at least size bytes from the given io.Reader.
func ReadAtLeastFrom(reader io.Reader, size int) Supplier {
	return func(b []byte) (int, error) {
		return io.ReadAtLeast(reader, b, size)
	}
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
		buffer: make([]byte, 32*1024),
	}
}

func NewMergingReader(reader io.Reader) Reader {
	return NewMergingReaderSize(reader, 32*1024)
}

func NewMergingReaderSize(reader io.Reader, size uint32) Reader {
	return &BytesToBufferReader{
		reader: reader,
		buffer: make([]byte, size),
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
	if mw, ok := writer.(MultiBufferWriter); ok {
		return &writerAdapter{
			writer: mw,
		}
	}

	return &BufferToBytesWriter{
		writer: writer,
	}
}

func NewMergingWriter(writer io.Writer) Writer {
	return NewMergingWriterSize(writer, 4096)
}

func NewMergingWriterSize(writer io.Writer, size uint32) Writer {
	return &mergingWriter{
		writer: writer,
		buffer: make([]byte, size),
	}
}

func NewSequentialWriter(writer io.Writer) Writer {
	return &seqWriter{
		writer: writer,
	}
}

// ToBytesWriter converts a Writer to io.Writer
func ToBytesWriter(writer Writer) io.Writer {
	return &bytesToBufferWriter{
		writer: writer,
	}
}
