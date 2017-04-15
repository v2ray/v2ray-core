package mux

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

type Reader struct {
	reader io.Reader
	buffer *buf.Buffer
}

func NewReader(reader buf.Reader) *Reader {
	return &Reader{
		reader: buf.ToBytesReader(reader),
		buffer: buf.NewLocal(1024),
	}
}

func (r *Reader) ReadMetadata() (*FrameMetadata, error) {
	b := r.buffer
	b.Clear()

	if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, 2)); err != nil {
		return nil, err
	}
	metaLen := serial.BytesToUint16(b.Bytes())
	if metaLen > 512 {
		return nil, newError("invalid metalen ", metaLen)
	}
	b.Clear()
	if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, int(metaLen))); err != nil {
		return nil, err
	}
	return ReadFrameFrom(b.Bytes())
}

func (r *Reader) Read() (buf.MultiBuffer, error) {
	r.buffer.Clear()
	if err := r.buffer.AppendSupplier(buf.ReadFullFrom(r.reader, 2)); err != nil {
		return nil, err
	}

	dataLen := int(serial.BytesToUint16(r.buffer.Bytes()))
	mb := buf.NewMultiBuffer()
	for dataLen > 0 {
		b := buf.New()
		readLen := buf.Size
		if dataLen < readLen {
			readLen = dataLen
		}
		if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, readLen)); err != nil {
			mb.Release()
			return nil, err
		}
		dataLen -= readLen
		mb.Append(b)
	}

	return mb, nil
}
