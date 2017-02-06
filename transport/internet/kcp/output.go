package kcp

import (
	"io"
	"sync"

	"v2ray.com/core/common/buf"
)

type SegmentWriter interface {
	Write(seg Segment) error
}

type SimpleSegmentWriter struct {
	sync.Mutex
	buffer *buf.Buffer
	writer io.Writer
}

func NewSegmentWriter(writer io.Writer) SegmentWriter {
	return &SimpleSegmentWriter{
		writer: writer,
		buffer: buf.NewSmall(),
	}
}

func (v *SimpleSegmentWriter) Write(seg Segment) error {
	v.Lock()
	defer v.Unlock()

	v.buffer.Reset(seg.Bytes())
	_, err := v.writer.Write(v.buffer.Bytes())
	return err
}
