package mux

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

type Reader struct {
	reader          io.Reader
	remainingLength int
	buffer          *buf.Buffer
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

func (r *Reader) Read() (*buf.Buffer, bool, error) {
	b := buf.New()
	var dataLen int
	if r.remainingLength > 0 {
		dataLen = r.remainingLength
		r.remainingLength = 0
	} else {
		if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, 2)); err != nil {
			return nil, false, err
		}
		dataLen = int(serial.BytesToUint16(b.Bytes()))
		b.Clear()
	}

	if dataLen > buf.Size {
		r.remainingLength = dataLen - buf.Size
		dataLen = buf.Size
	}

	if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, dataLen)); err != nil {
		return nil, false, err
	}

	return b, (r.remainingLength > 0), nil
}
