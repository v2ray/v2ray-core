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

func (v *BufferedSegmentWriter) Write(seg Segment) {
	v.Lock()
	defer v.Unlock()

	nBytes := seg.ByteSize()
	if uint32(v.buffer.Len()+nBytes) > v.mtu {
		v.FlushWithoutLock()
	}

	if v.buffer == nil {
		v.buffer = alloc.NewSmallBuffer()
	}

	v.buffer.AppendFunc(seg.Bytes())
}

func (v *BufferedSegmentWriter) FlushWithoutLock() {
	v.writer.Write(v.buffer)
	v.buffer = nil
}

func (v *BufferedSegmentWriter) Flush() {
	v.Lock()
	defer v.Unlock()

	if v.buffer.Len() == 0 {
		return
	}

	v.FlushWithoutLock()
}

type AuthenticationWriter struct {
	Authenticator internet.Authenticator
	Writer        io.Writer
	Config        *Config
}

func (v *AuthenticationWriter) Write(payload *alloc.Buffer) error {
	defer payload.Release()

	v.Authenticator.Seal(payload)
	_, err := v.Writer.Write(payload.Bytes())
	return err
}

func (v *AuthenticationWriter) Release() {}

func (v *AuthenticationWriter) Mtu() uint32 {
	return v.Config.Mtu.GetValue() - uint32(v.Authenticator.Overhead())
}
