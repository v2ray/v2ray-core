package io

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
)

// Reader extends io.Reader with alloc.Buffer.
type Reader interface {
	common.Releasable
	// Read reads content from underlying reader, and put it into an alloc.Buffer.
	Read() (*alloc.Buffer, error)
}

// AdaptiveReader is a Reader that adjusts its reading speed automatically.
type AdaptiveReader struct {
	reader      io.Reader
	largeBuffer *alloc.Buffer
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
func (v *AdaptiveReader) Read() (*alloc.Buffer, error) {
	if v.highVolumn && v.largeBuffer.IsEmpty() {
		if v.largeBuffer == nil {
			v.largeBuffer = alloc.NewLocalBuffer(256 * 1024).Clear()
		}
		nBytes, err := v.largeBuffer.FillFrom(v.reader)
		if err != nil {
			return nil, err
		}
		if nBytes < alloc.BufferSize {
			v.highVolumn = false
		}
	}

	buffer := alloc.NewBuffer().Clear()
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

func (v *AdaptiveReader) Release() {
	v.reader = nil
}
