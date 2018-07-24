package buf

import (
	"io"
	"runtime"
	"syscall"
	"time"

	"v2ray.com/core/common/platform"
)

// Reader extends io.Reader with MultiBuffer.
type Reader interface {
	// ReadMultiBuffer reads content from underlying reader, and put it into a MultiBuffer.
	ReadMultiBuffer() (MultiBuffer, error)
}

// ErrReadTimeout is an error that happens with IO timeout.
var ErrReadTimeout = newError("IO timeout")

// TimeoutReader is a reader that returns error if Read() operation takes longer than the given timeout.
type TimeoutReader interface {
	ReadMultiBufferTimeout(time.Duration) (MultiBuffer, error)
}

// Writer extends io.Writer with MultiBuffer.
type Writer interface {
	// WriteMultiBuffer writes a MultiBuffer into underlying writer.
	WriteMultiBuffer(MultiBuffer) error
}

// ReadFrom creates a Supplier to read from a given io.Reader.
func ReadFrom(reader io.Reader) Supplier {
	return func(b []byte) (int, error) {
		return reader.Read(b)
	}
}

// ReadFullFrom creates a Supplier to read full buffer from a given io.Reader.
func ReadFullFrom(reader io.Reader, size int32) Supplier {
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

var useReadv = false

func init() {
	const defaultFlagValue = "NOT_DEFINED_AT_ALL"
	value := platform.NewEnvFlag("v2ray.buf.readv.disable").GetValue(func() string { return defaultFlagValue })
	if value != defaultFlagValue {
		useReadv = false
		return
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		newError("ReadV enabled").WriteToLog()
		useReadv = true
	}
}

// NewReader creates a new Reader.
// The Reader instance doesn't take the ownership of reader.
func NewReader(reader io.Reader) Reader {
	if mr, ok := reader.(Reader); ok {
		return mr
	}

	if useReadv {
		if sc, ok := reader.(syscall.Conn); ok {
			rawConn, err := sc.SyscallConn()
			if err != nil {
				newError("failed to get sysconn").Base(err).WriteToLog()
			} else {
				return NewReadVReader(reader, rawConn)
			}
		}
	}

	return NewBytesToBufferReader(reader)
}

// NewWriter creates a new Writer.
func NewWriter(writer io.Writer) Writer {
	if mw, ok := writer.(Writer); ok {
		return mw
	}

	return &BufferToBytesWriter{
		Writer: writer,
	}
}

// NewSequentialWriter returns a Writer that write Buffers in a MultiBuffer sequentially.
func NewSequentialWriter(writer io.Writer) Writer {
	return &seqWriter{
		writer: writer,
	}
}
