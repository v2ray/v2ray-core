package mux

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

type Reader struct {
	reader   io.Reader
	buffer   *buf.Buffer
	leftOver int
}

func NewReader(reader buf.Reader) *Reader {
	return &Reader{
		reader:   buf.ToBytesReader(reader),
		buffer:   buf.NewLocal(1024),
		leftOver: -1,
	}
}

func (r *Reader) ReadMetadata() (*FrameMetadata, error) {
	r.leftOver = -1

	b := r.buffer
	b.Clear()

	if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, 2)); err != nil {
		return nil, err
	}
	metaLen := serial.BytesToUint16(b.Bytes())
	if metaLen > 512 {
		return nil, newError("invalid metalen ", metaLen).AtWarning()
	}
	b.Clear()
	if err := b.AppendSupplier(buf.ReadFullFrom(r.reader, int(metaLen))); err != nil {
		return nil, err
	}
	return ReadFrameFrom(b.Bytes())
}

func (r *Reader) readSize() error {
	if err := r.buffer.Reset(buf.ReadFullFrom(r.reader, 2)); err != nil {
		return err
	}
	r.leftOver = int(serial.BytesToUint16(r.buffer.Bytes()))
	return nil
}

func (r *Reader) Read() (buf.MultiBuffer, error) {
	if r.leftOver == 0 {
		r.leftOver = -1
		return nil, io.EOF
	}
	if r.leftOver == -1 {
		if err := r.readSize(); err != nil {
			return nil, err
		}
	}

	mb := buf.NewMultiBuffer()
	for r.leftOver > 0 {
		readLen := buf.Size
		if r.leftOver < readLen {
			readLen = r.leftOver
		}
		b := buf.New()
		if err := b.AppendSupplier(func(bb []byte) (int, error) {
			return r.reader.Read(bb[:readLen])
		}); err != nil {
			mb.Release()
			return nil, err
		}
		r.leftOver -= b.Len()
		mb.Append(b)
		if b.Len() < readLen {
			break
		}
	}

	return mb, nil
}
