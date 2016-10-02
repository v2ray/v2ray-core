package kcp

import (
	"io"
	"sync"

	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/transport/internet"
)

type SegmentWriter interface {
	Write(seg Segment)
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

func (this *BufferedSegmentWriter) Write(seg Segment) {
	this.Lock()
	defer this.Unlock()

	nBytes := seg.ByteSize()
	if uint32(this.buffer.Len()+nBytes) > this.mtu {
		this.FlushWithoutLock()
	}

	if this.buffer == nil {
		this.buffer = alloc.NewLocalBuffer(2048).Clear()
	}

	this.buffer.Value = seg.Bytes(this.buffer.Value)
}

func (this *BufferedSegmentWriter) FlushWithoutLock() {
	go this.writer.Write(this.buffer)
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
	Authenticator internet.Authenticator
	Writer        io.Writer
	Config        *Config
}

func (this *AuthenticationWriter) Write(payload *alloc.Buffer) error {
	defer payload.Release()

	this.Authenticator.Seal(payload)
	_, err := this.Writer.Write(payload.Value)
	return err
}

func (this *AuthenticationWriter) Release() {}

func (this *AuthenticationWriter) Mtu() uint32 {
	return this.Config.Mtu.GetValue() - uint32(this.Authenticator.Overhead())
}
