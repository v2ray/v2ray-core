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

// ReadAtLeastFrom create a Supplier to read at least size bytes from the given io.Reader.
func ReadAtLeastFrom(reader io.Reader, size int) Supplier {
	return func(b []byte) (int, error) {
		return io.ReadAtLeast(reader, b, size)
	}
}

type copyHandler struct {
	onReadError  func(error) error
	onData       func()
	onWriteError func(error) error
}

func (h *copyHandler) readFrom(reader Reader) (MultiBuffer, error) {
	mb, err := reader.Read()
	if err != nil && h.onReadError != nil {
		err = h.onReadError(err)
	}
	return mb, err
}

func (h *copyHandler) writeTo(writer Writer, mb MultiBuffer) error {
	err := writer.Write(mb)
	if err != nil && h.onWriteError != nil {
		err = h.onWriteError(err)
	}
	return err
}

type CopyOption func(*copyHandler)

func IgnoreReaderError() CopyOption {
	return func(handler *copyHandler) {
		handler.onReadError = func(err error) error {
			return nil
		}
	}
}

func IgnoreWriterError() CopyOption {
	return func(handler *copyHandler) {
		handler.onWriteError = func(err error) error {
			return nil
		}
	}
}

func UpdateActivity(timer signal.ActivityTimer) CopyOption {
	return func(handler *copyHandler) {
		handler.onData = func() {
			timer.Update()
		}
	}
}

func copyInternal(reader Reader, writer Writer, handler *copyHandler) error {
	for {
		buffer, err := handler.readFrom(reader)
		if err != nil {
			return err
		}

		if buffer.IsEmpty() {
			buffer.Release()
			continue
		}

		if handler.onData != nil {
			handler.onData()
		}

		if err := handler.writeTo(writer, buffer); err != nil {
			buffer.Release()
			return err
		}
	}
}

// Copy dumps all payload from reader to writer or stops when an error occurs.
// ActivityTimer gets updated as soon as there is a payload.
func Copy(reader Reader, writer Writer, options ...CopyOption) error {
	handler := new(copyHandler)
	for _, option := range options {
		option(handler)
	}
	err := copyInternal(reader, writer, handler)
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
