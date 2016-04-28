package io // import "github.com/v2ray/v2ray-core/common/io"

import (
	"io"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
)

// ReadFrom reads from a reader and put all content to a buffer.
// If buffer is nil, ReadFrom creates a new normal buffer.
func ReadFrom(reader io.Reader, buffer *alloc.Buffer) (*alloc.Buffer, error) {
	if buffer == nil {
		buffer = alloc.NewBuffer()
	}
	nBytes, err := reader.Read(buffer.Value)
	buffer.Slice(0, nBytes)
	return buffer, err
}

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
	isLarge  bool
}

// NewAdaptiveReader creates a new AdaptiveReader.
// The AdaptiveReader instance doesn't take the ownership of reader.
func NewAdaptiveReader(reader io.Reader) *AdaptiveReader {
	return &AdaptiveReader{
		reader:   reader,
		allocate: alloc.NewBuffer,
		isLarge:  false,
	}
}

// Read implements Reader.Read().
func (this *AdaptiveReader) Read() (*alloc.Buffer, error) {
	buffer, err := ReadFrom(this.reader, this.allocate())

	if buffer.IsFull() && !this.isLarge {
		this.allocate = alloc.NewLargeBuffer
		this.isLarge = true
	} else if !buffer.IsFull() {
		this.allocate = alloc.NewBuffer
		this.isLarge = false
	}

	if err != nil {
		buffer.Release()
		return nil, err
	}
	return buffer, nil
}

func (this *AdaptiveReader) Release() {
	this.reader = nil
}
