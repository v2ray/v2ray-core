package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

// Reader extends io.Reader with alloc.Buffer.
type Reader interface {
	Release()
	// Read reads content from underlying reader, and put it into an alloc.Buffer.
	Read() (*Buffer, error)
}

// Writer extends io.Writer with alloc.Buffer.
type Writer interface {
	Release()
	// Write writes an alloc.Buffer into underlying writer.
	Write(*Buffer) error
}

func ReadFrom(reader io.Reader) Supplier {
	return func(b []byte) (int, error) {
		return reader.Read(b)
	}
}

func ReadFullFrom(reader io.Reader, size int) Supplier {
	return func(b []byte) (int, error) {
		return io.ReadFull(reader, b[:size])
	}
}

// Pipe dumps all content from reader to writer, until an error happens.
func Pipe(reader Reader, writer Writer) error {
	for {
		buffer, err := reader.Read()
		if err != nil {
			return err
		}

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

// PipeUntilEOF behaves the same as Pipe(). The only difference is PipeUntilEOF returns nil on EOF.
func PipeUntilEOF(reader Reader, writer Writer) error {
	err := Pipe(reader, writer)
	if err != nil && errors.Cause(err) != io.EOF {
		return err
	}
	return nil
}

// NewReader creates a new Reader.
// The Reader instance doesn't take the ownership of reader.
func NewReader(reader io.Reader) Reader {
	return &BytesToBufferReader{
		reader: reader,
	}
}

func NewBytesReader(stream Reader) *BufferToBytesReader {
	return &BufferToBytesReader{
		stream: stream,
	}
}

// NewWriter creates a new Writer.
func NewWriter(writer io.Writer) Writer {
	return &BufferToBytesWriter{
		writer: writer,
	}
}

func NewBytesWriter(writer Writer) *BytesToBufferWriter {
	return &BytesToBufferWriter{
		writer: writer,
	}
}
