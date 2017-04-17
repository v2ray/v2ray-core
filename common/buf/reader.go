package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

// BytesToBufferReader is a Reader that adjusts its reading speed automatically.
type BytesToBufferReader struct {
	reader io.Reader
	buffer *Buffer
}

// Read implements Reader.Read().
func (r *BytesToBufferReader) Read() (MultiBuffer, error) {
	if err := r.buffer.Reset(ReadFrom(r.reader)); err != nil {
		return nil, err
	}

	mb := NewMultiBuffer()
	for !r.buffer.IsEmpty() {
		b := New()
		b.AppendSupplier(ReadFrom(r.buffer))
		mb.Append(b)
	}
	return mb, nil
}

func (r *BytesToBufferReader) WriteTo(writer io.Writer) (int64, error) {
	totalBytes := int64(0)
	eof := false
	for !eof {
		if err := r.buffer.Reset(ReadFrom(r.reader)); err != nil {
			if errors.Cause(err) == io.EOF {
				eof = true
			} else {
				return totalBytes, err
			}
		}
		nBytes, err := writer.Write(r.buffer.Bytes())
		totalBytes += int64(nBytes)
		if err != nil {
			return totalBytes, err
		}
	}
	return totalBytes, nil
}

type readerAdpater struct {
	MultiBufferReader
}

func (r *readerAdpater) Read() (MultiBuffer, error) {
	return r.ReadMultiBuffer()
}

type bufferToBytesReader struct {
	stream  Reader
	current MultiBuffer
	err     error
}

// fill fills in the internal buffer.
func (r *bufferToBytesReader) fill() {
	b, err := r.stream.Read()
	if err != nil {
		r.err = err
		return
	}
	r.current = b
}

func (r *bufferToBytesReader) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	if r.current == nil {
		r.fill()
		if r.err != nil {
			return 0, r.err
		}
	}
	nBytes, err := r.current.Read(b)
	if r.current.IsEmpty() {
		r.current.Release()
		r.current = nil
	}
	return nBytes, err
}

func (r *bufferToBytesReader) ReadMultiBuffer() (MultiBuffer, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.current == nil {
		r.fill()
		if r.err != nil {
			return nil, r.err
		}
	}
	b := r.current
	r.current = nil
	return b, nil
}

func (r *bufferToBytesReader) writeToInternal(writer io.Writer) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	mbWriter := NewWriter(writer)
	totalBytes := int64(0)
	for {
		if r.current == nil {
			r.fill()
			if r.err != nil {
				return totalBytes, r.err
			}
		}
		totalBytes := int64(r.current.Len())
		if err := mbWriter.Write(r.current); err != nil {
			return totalBytes, err
		}
		r.current = nil
	}
}

func (r *bufferToBytesReader) WriteTo(writer io.Writer) (int64, error) {
	nBytes, err := r.writeToInternal(writer)
	if errors.Cause(err) == io.EOF {
		return nBytes, nil
	}
	return nBytes, err
}
