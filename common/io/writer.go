package io

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
)

// Writer extends io.Writer with alloc.Buffer.
type Writer interface {
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

// Write implements Writer.Write().
func (this *AdaptiveWriter) Write(buffer *alloc.Buffer) error {
	nBytes, err := this.writer.Write(buffer.Value)
	if nBytes < buffer.Len() {
		_, err = this.writer.Write(buffer.Value[nBytes:])
	}
	return err
}
