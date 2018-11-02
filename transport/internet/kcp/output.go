package kcp

import (
	"io"
	"sync"

	"v2ray.com/core/common/retry"

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
		buffer: buf.New(),
	}
}

func (w *SimpleSegmentWriter) Write(seg Segment) error {
	w.Lock()
	defer w.Unlock()

	w.buffer.Clear()
	rawBytes := w.buffer.Extend(seg.ByteSize())
	seg.Serialize(rawBytes)
	_, err := w.writer.Write(w.buffer.Bytes())
	return err
}

type RetryableWriter struct {
	writer SegmentWriter
}

func NewRetryableWriter(writer SegmentWriter) SegmentWriter {
	return &RetryableWriter{
		writer: writer,
	}
}

func (w *RetryableWriter) Write(seg Segment) error {
	return retry.Timed(5, 100).On(func() error {
		return w.writer.Write(seg)
	})
}
