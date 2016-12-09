package kcp

import (
	"io"
	"sync"

	"v2ray.com/core/common/buf"
)

type SegmentWriter interface {
	Write(seg Segment)
}

type BufferedSegmentWriter struct {
	sync.Mutex
	mtu    uint32
	buffer *buf.Buffer
	writer io.Writer
}

func NewSegmentWriter(writer io.Writer, mtu uint32) *BufferedSegmentWriter {
	return &BufferedSegmentWriter{
		mtu:    mtu,
		writer: writer,
		buffer: buf.NewSmallBuffer(),
	}
}

func (v *BufferedSegmentWriter) Write(seg Segment) {
	v.Lock()
	defer v.Unlock()

	nBytes := seg.ByteSize()
	if uint32(v.buffer.Len()+nBytes) > v.mtu {
		v.FlushWithoutLock()
	}

	v.buffer.AppendFunc(seg.Bytes())
}

func (v *BufferedSegmentWriter) FlushWithoutLock() {
	v.writer.Write(v.buffer.Bytes())
	v.buffer.Clear()
}

func (v *BufferedSegmentWriter) Flush() {
	v.Lock()
	defer v.Unlock()

	if v.buffer.IsEmpty() {
		return
	}

	v.FlushWithoutLock()
}
