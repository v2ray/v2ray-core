package kcp

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
)

type SegmentWriter interface {
	Write(seg ISegment)
}

type BufferedSegmentWriter struct {
	sync.Mutex
	mtu    uint32
	buffer *alloc.Buffer
	writer v2io.Writer
}

func NewSegmentWriter(writer *AuthenticationWriter) *BufferedSegmentWriter {
	return &BufferedSegmentWriter{
		mtu:    writer.Mtu(),
		writer: writer,
	}
}

func (this *BufferedSegmentWriter) Write(seg ISegment) {
	this.Lock()
	defer this.Unlock()

	nBytes := seg.ByteSize()
	if uint32(this.buffer.Len()+nBytes) > this.mtu {
		this.FlushWithoutLock()
	}

	if this.buffer == nil {
		this.buffer = alloc.NewSmallBuffer().Clear()
	}

	this.buffer.Append(seg.Bytes(nil))
}

func (this *BufferedSegmentWriter) FlushWithoutLock() {
	this.writer.Write(this.buffer)
	this.buffer = nil
}

func (this *BufferedSegmentWriter) Flush() {
	this.Lock()
	defer this.Unlock()

	if this.buffer.Len() == 0 {
		return
	}

	this.FlushWithoutLock()
}

type AuthenticationWriter struct {
	Authenticator Authenticator
	Writer        io.Writer
}

func (this *AuthenticationWriter) Write(payload *alloc.Buffer) error {
	defer payload.Release()

	this.Authenticator.Seal(payload)
	_, err := this.Writer.Write(payload.Value)
	return err
}

func (this *AuthenticationWriter) Release() {}

func (this *AuthenticationWriter) Mtu() uint32 {
	return effectiveConfig.Mtu - uint32(this.Authenticator.HeaderSize())
}
