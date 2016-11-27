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
	reader   io.Reader
	allocate func() *alloc.Buffer
}

// NewAdaptiveReader creates a new AdaptiveReader.
// The AdaptiveReader instance doesn't take the ownership of reader.
func NewAdaptiveReader(reader io.Reader) *AdaptiveReader {
	return &AdaptiveReader{
		reader:   reader,
		allocate: alloc.NewBuffer,
	}
}

// Read implements Reader.Read().
func (v *AdaptiveReader) Read() (*alloc.Buffer, error) {
	buffer := v.allocate().Clear()
	_, err := buffer.FillFrom(v.reader)
	if err != nil {
		buffer.Release()
		return nil, err
	}

	return buffer, nil
}

func (v *AdaptiveReader) Release() {
	v.reader = nil
}
