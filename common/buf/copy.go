package buf

import (
	"io"
	"time"
	"unsafe"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/signal"
)

type errorHandler func(error) error
type dataHandler func(MultiBuffer)

//go:notinheap
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

type readError struct {
	error
}

func (e readError) Error() string {
	return e.error.Error()
}

func (e readError) Inner() error {
	return e.error
}

func IsReadError(err error) bool {
	_, ok := err.(readError)
	return ok
}

type writeError struct {
	error
}

func (e writeError) Error() string {
	return e.error.Error()
}

func (e writeError) Inner() error {
	return e.error
}

func IsWriteError(err error) bool {
	_, ok := err.(writeError)
	return ok
}

func copyInternal(reader Reader, writer Writer, handler *copyHandler) error {
	for {
		buffer, err := handler.readFrom(reader)
		if !buffer.IsEmpty() {
			for _, handler := range handler.onData {
				handler(buffer)
			}

			if werr := handler.writeTo(writer, buffer); werr != nil {
				return writeError{werr}
			}
		}

		if err != nil {
			return readError{err}
		}
	}
}

// Copy dumps all payload from reader to writer or stops when an error occurs. It returns nil when EOF.
func Copy(reader Reader, writer Writer, options ...CopyOption) error {
	var handler copyHandler
	p := uintptr(unsafe.Pointer(&handler))
	h := (*copyHandler)(unsafe.Pointer(p))

	for _, option := range options {
		option(h)
	}
	err := copyInternal(reader, writer, h)
	if err != nil && errors.Cause(err) != io.EOF {
		return err
	}
	return nil
}

var ErrNotTimeoutReader = newError("not a TimeoutReader")

func CopyOnceTimeout(reader Reader, writer Writer, timeout time.Duration) error {
	timeoutReader, ok := reader.(TimeoutReader)
	if !ok {
		return ErrNotTimeoutReader
	}
	mb, err := timeoutReader.ReadMultiBufferTimeout(timeout)
	if err != nil {
		return err
	}
	return writer.WriteMultiBuffer(mb)
}
