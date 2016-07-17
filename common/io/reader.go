package io // import "github.com/v2ray/v2ray-core/common/io"

import (
	"io"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
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
func (this *AdaptiveReader) Read() (*alloc.Buffer, error) {
	buffer := this.allocate().Clear()
	_, err := buffer.FillFrom(this.reader)
	if err != nil {
		buffer.Release()
		return nil, err
	}

	if buffer.Len() >= alloc.BufferSize {
		this.allocate = alloc.NewLargeBuffer
	} else {
		this.allocate = alloc.NewBuffer
	}

	return buffer, nil
}

func (this *AdaptiveReader) Release() {
	this.reader = nil
}
