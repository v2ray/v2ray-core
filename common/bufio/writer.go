package bufio

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
)

// BufferedWriter is an io.Writer with internal buffer. It writes to underlying writer when buffer is full or on demand.
// This type is not thread safe.
type BufferedWriter struct {
	writer   io.Writer
	buffer   *buf.Buffer
	buffered bool
}

// NewWriter creates a new BufferedWriter.
func NewWriter(rawWriter io.Writer) *BufferedWriter {
	return &BufferedWriter{
		writer:   rawWriter,
		buffer:   buf.NewSmall(),
		buffered: true,
	}
}

func (v *BufferedWriter) ReadFrom(reader io.Reader) (int64, error) {
	totalBytes := int64(0)
	for {
		oriSize := v.buffer.Len()
		err := v.buffer.AppendSupplier(buf.ReadFrom(reader))
		totalBytes += int64(v.buffer.Len() - oriSize)
		if err != nil {
			if errors.Cause(err) == io.EOF {
				return totalBytes, nil
			}
			return totalBytes, err
		}
		if err := v.Flush(); err != nil {
			return totalBytes, err
		}
	}
}

func (v *BufferedWriter) Write(b []byte) (int, error) {
	if !v.buffered || v.buffer == nil {
		return v.writer.Write(b)
	}
	nBytes, err := v.buffer.Write(b)
	if err != nil {
		return 0, err
	}
	if v.buffer.IsFull() {
		err := v.Flush()
		if err != nil {
			return 0, err
		}
		if nBytes < len(b) {
			if _, err := v.writer.Write(b[nBytes:]); err != nil {
				return nBytes, err
			}
		}
	}
	return len(b), nil
}

func (v *BufferedWriter) Flush() error {
	defer v.buffer.Clear()
	for !v.buffer.IsEmpty() {
		nBytes, err := v.writer.Write(v.buffer.Bytes())
		if err != nil {
			return err
		}
		v.buffer.SliceFrom(nBytes)
	}
	return nil
}

func (v *BufferedWriter) Buffered() bool {
	return v.buffered
}

func (v *BufferedWriter) SetBuffered(cached bool) {
	v.buffered = cached
	if !cached && !v.buffer.IsEmpty() {
		v.Flush()
	}
}

// Release implements common.Releasable.Release().
func (v *BufferedWriter) Release() {
	if !v.buffer.IsEmpty() {
		v.Flush()
	}

	if v.buffer != nil {
		v.buffer.Release()
		v.buffer = nil
	}
	common.Release(v.writer)
}
