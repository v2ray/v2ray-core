package io

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

// Reader extends io.Reader with alloc.Buffer.
type Reader interface {
	common.Releasable
	// Read reads content from underlying reader, and put it into an alloc.Buffer.
	Read() (*buf.Buffer, error)
}

// AdaptiveReader is a Reader that adjusts its reading speed automatically.
type AdaptiveReader struct {
	reader      io.Reader
	largeBuffer *buf.Buffer
	highVolumn  bool
}

// NewAdaptiveReader creates a new AdaptiveReader.
// The AdaptiveReader instance doesn't take the ownership of reader.
func NewAdaptiveReader(reader io.Reader) *AdaptiveReader {
	return &AdaptiveReader{
		reader: reader,
	}
}

// Read implements Reader.Read().
func (v *AdaptiveReader) Read() (*buf.Buffer, error) {
	if v.highVolumn && v.largeBuffer.IsEmpty() {
		if v.largeBuffer == nil {
			v.largeBuffer = buf.NewLocal(32 * 1024)
		}
		err := v.largeBuffer.AppendSupplier(buf.ReadFrom(v.reader))
		if err != nil {
			return nil, err
		}
		if v.largeBuffer.Len() < buf.Size {
			v.highVolumn = false
		}
	}

	buffer := buf.New()
	if !v.largeBuffer.IsEmpty() {
		buffer.AppendSupplier(buf.ReadFrom(v.largeBuffer))
		return buffer, nil
	}

	err := buffer.AppendSupplier(buf.ReadFrom(v.reader))
	if err != nil {
		buffer.Release()
		return nil, err
	}

	if buffer.IsFull() {
		v.highVolumn = true
	}

	return buffer, nil
}

// Release implements Releasable.Release().
func (v *AdaptiveReader) Release() {
	v.reader = nil
}
