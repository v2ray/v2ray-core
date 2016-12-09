package buf

import (
	"io"
	"sync"
)

// BufferToBytesWriter is a Writer that writes alloc.Buffer into underlying writer.
type BufferToBytesWriter struct {
	writer io.Writer
}

// Write implements Writer.Write(). Write() takes ownership of the given buffer.
func (v *BufferToBytesWriter) Write(buffer *Buffer) error {
	defer buffer.Release()
	for {
		nBytes, err := v.writer.Write(buffer.Bytes())
		if err != nil {
			return err
		}
		if nBytes == buffer.Len() {
			break
		}
		buffer.SliceFrom(nBytes)
	}
	return nil
}

// Release implements Releasable.Release().
func (v *BufferToBytesWriter) Release() {
	v.writer = nil
}

type BytesToBufferWriter struct {
	sync.Mutex
	writer Writer
}

func (v *BytesToBufferWriter) Write(payload []byte) (int, error) {
	v.Lock()
	defer v.Unlock()
	if v.writer == nil {
		return 0, io.ErrClosedPipe
	}

	bytesWritten := 0
	size := len(payload)
	for size > 0 {
		buffer := New()
		nBytes, _ := buffer.Write(payload)
		size -= nBytes
		payload = payload[nBytes:]
		bytesWritten += nBytes
		err := v.writer.Write(buffer)
		if err != nil {
			return bytesWritten, err
		}
	}

	return bytesWritten, nil
}

func (v *BytesToBufferWriter) Release() {
	v.Lock()
	v.writer.Release()
	v.writer = nil
	v.Unlock()
}
