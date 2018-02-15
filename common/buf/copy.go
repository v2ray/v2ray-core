package buf

import (
	"io"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/signal"
)

type errorHandler func(error) error
type dataHandler func(MultiBuffer)

type copyHandler struct {
	onReadError  []errorHandler
	onData       []dataHandler
	onWriteError []errorHandler
}

func (h *copyHandler) readFrom(reader Reader) (MultiBuffer, error) {
	mb, err := reader.ReadMultiBuffer()
	if err != nil {
		for _, handler := range h.onReadError {
			err = handler(err)
		}
	}
	return mb, err
}

func (h *copyHandler) writeTo(writer Writer, mb MultiBuffer) error {
	err := writer.WriteMultiBuffer(mb)
	if err != nil {
		for _, handler := range h.onWriteError {
			err = handler(err)
		}
	}
	return err
}

// SizeCounter is for counting bytes copied by Copy().
type SizeCounter struct {
	Size int64
}

// CopyOption is an option for copying data.
type CopyOption func(*copyHandler)

// IgnoreReaderError is a CopyOption that ignores errors from reader. Copy will continue in such case.
func IgnoreReaderError() CopyOption {
	return func(handler *copyHandler) {
		handler.onReadError = append(handler.onReadError, func(err error) error {
			return nil
		})
	}
}

// IgnoreWriterError is a CopyOption that ignores errors from writer. Copy will continue in such case.
func IgnoreWriterError() CopyOption {
	return func(handler *copyHandler) {
		handler.onWriteError = append(handler.onWriteError, func(err error) error {
			return nil
		})
	}
}

// UpdateActivity is a CopyOption to update activity on each data copy operation.
func UpdateActivity(timer signal.ActivityUpdater) CopyOption {
	return func(handler *copyHandler) {
		handler.onData = append(handler.onData, func(MultiBuffer) {
			timer.Update()
		})
	}
}

// CountSize is a CopyOption that sums the total size of data copied into the given SizeCounter.
func CountSize(sc *SizeCounter) CopyOption {
	return func(handler *copyHandler) {
		handler.onData = append(handler.onData, func(b MultiBuffer) {
			sc.Size += int64(b.Len())
		})
	}
}

func copyInternal(reader Reader, writer Writer, handler *copyHandler) error {
	for {
		buffer, err := handler.readFrom(reader)
		if !buffer.IsEmpty() {
			for _, handler := range handler.onData {
				handler(buffer)
			}

			if werr := handler.writeTo(writer, buffer); werr != nil {
				buffer.Release()
				return werr
			}
		} else if err != nil {
			return err
		}
	}
}

// Copy dumps all payload from reader to writer or stops when an error occurs. It returns nil when EOF.
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
