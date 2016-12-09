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
			v.largeBuffer = buf.NewLocalBuffer(32 * 1024)
		}
		nBytes, err := v.largeBuffer.FillFrom(v.reader)
		if err != nil {
			return nil, err
		}
		if nBytes < buf.BufferSize {
			v.highVolumn = false
		}
	}

	buffer := buf.NewBuffer()
	if !v.largeBuffer.IsEmpty() {
		buffer.FillFrom(v.largeBuffer)
		return buffer, nil
	}

	_, err := buffer.FillFrom(v.reader)
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
