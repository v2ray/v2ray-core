package io

import (
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
)

type Writer interface {
	Write(*alloc.Buffer) error
}

type AdaptiveWriter struct {
	writer io.Writer
}

func NewAdaptiveWriter(writer io.Writer) *AdaptiveWriter {
	return &AdaptiveWriter{
		writer: writer,
	}
}

func (this *AdaptiveWriter) Write(buffer *alloc.Buffer) error {
	nBytes, err := this.writer.Write(buffer.Value)
	if nBytes < buffer.Len() {
		_, err = this.writer.Write(buffer.Value[nBytes:])
	}
	return err
}
