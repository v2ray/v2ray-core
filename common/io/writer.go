package io

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
)

// Writer extends io.Writer with alloc.Buffer.
type Writer interface {
	common.Releasable
	// Write writes an alloc.Buffer into underlying writer.
	Write(*alloc.Buffer) error
}

// AdaptiveWriter is a Writer that writes alloc.Buffer into underlying writer.
type AdaptiveWriter struct {
	writer io.Writer
}

// NewAdaptiveWriter creates a new AdaptiveWriter.
func NewAdaptiveWriter(writer io.Writer) *AdaptiveWriter {
	return &AdaptiveWriter{
		writer: writer,
	}
}

// Write implements Writer.Write(). Write() takes ownership of the given buffer.
func (v *AdaptiveWriter) Write(buffer *alloc.Buffer) error {
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
func (v *AdaptiveWriter) Release() {
	v.writer = nil
}
