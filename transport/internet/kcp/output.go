package kcp

import (
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
)

type SegmentWriter struct {
	sync.Mutex
	mtu    uint32
	buffer *alloc.Buffer
	writer v2io.Writer
}

func NewSegmentWriter(mtu uint32, writer v2io.Writer) *SegmentWriter {
	return &SegmentWriter{
		mtu:    mtu,
		writer: writer,
	}
}

func (this *SegmentWriter) Write(seg ISegment) {
	this.Lock()
	defer this.Unlock()

	nBytes := seg.ByteSize()
	if uint32(this.buffer.Len()+nBytes) > this.mtu {
		this.FlushWithoutLock()
	}

	if this.buffer == nil {
		this.buffer = alloc.NewSmallBuffer().Clear()
	}

	this.buffer.Value = seg.Bytes(this.buffer.Value)
}

func (this *SegmentWriter) FlushWithoutLock() {
	this.writer.Write(this.buffer)
	this.buffer = nil
}

func (this *SegmentWriter) Flush() {
	this.Lock()
	defer this.Unlock()

	if this.buffer.Len() == 0 {
		return
	}

	this.FlushWithoutLock()
}
