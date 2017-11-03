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
	mb, err := reader.Read()
	if err != nil {
		for _, handler := range h.onReadError {
			err = handler(err)
		}
	}
	return mb, err
}

func (h *copyHandler) writeTo(writer Writer, mb MultiBuffer) error {
	err := writer.Write(mb)
	if err != nil {
		for _, handler := range h.onWriteError {
			err = handler(err)
		}
	}
	return err
}

type CopyOption func(*copyHandler)

func IgnoreReaderError() CopyOption {
	return func(handler *copyHandler) {
		handler.onReadError = append(handler.onReadError, func(err error) error {
			return nil
		})
	}
}

func IgnoreWriterError() CopyOption {
	return func(handler *copyHandler) {
		handler.onWriteError = append(handler.onWriteError, func(err error) error {
			return nil
		})
	}
}

func UpdateActivity(timer signal.ActivityUpdater) CopyOption {
	return func(handler *copyHandler) {
		handler.onData = append(handler.onData, func(MultiBuffer) {
			timer.Update()
		})
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

		for _, handler := range handler.onData {
			handler(buffer)
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
