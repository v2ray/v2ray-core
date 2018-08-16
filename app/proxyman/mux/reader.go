package mux

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

// ReadMetadata reads FrameMetadata from the given reader.
func ReadMetadata(reader io.Reader) (*FrameMetadata, error) {
	metaLen, err := serial.ReadUint16(reader)
	if err != nil {
		return nil, err
	}
	if metaLen > 512 {
		return nil, newError("invalid metalen ", metaLen).AtError()
	}

	b := buf.New()
	defer b.Release()

	if err := b.Reset(buf.ReadFullFrom(reader, int32(metaLen))); err != nil {
		return nil, err
	}
	return ReadFrameFrom(b)
}

// PacketReader is an io.Reader that reads whole chunk of Mux frames every time.
type PacketReader struct {
	reader io.Reader
	eof    bool
}

// NewPacketReader creates a new PacketReader.
func NewPacketReader(reader io.Reader) *PacketReader {
	return &PacketReader{
		reader: reader,
		eof:    false,
	}
}

// ReadMultiBuffer implements buf.Reader.
func (r *PacketReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if r.eof {
		return nil, io.EOF
	}

	size, err := serial.ReadUint16(r.reader)
	if err != nil {
		return nil, err
	}

	if size > buf.Size {
		return nil, newError("packet size too large: ", size)
	}

	b := buf.New()
	if err := b.Reset(buf.ReadFullFrom(r.reader, int32(size))); err != nil {
		b.Release()
		return nil, err
	}
	r.eof = true
	return buf.NewMultiBufferValue(b), nil
}

// StreamReader reads Mux frame as a stream.
type StreamReader struct {
	reader   *buf.BufferedReader
	leftOver int32
}

// NewStreamReader creates a new StreamReader.
func NewStreamReader(reader *buf.BufferedReader) *StreamReader {
	return &StreamReader{
		reader:   reader,
		leftOver: -1,
	}
}

// ReadMultiBuffer implmenets buf.Reader.
func (r *StreamReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if r.leftOver == 0 {
		return nil, io.EOF
	}

	if r.leftOver == -1 {
		size, err := serial.ReadUint16(r.reader)
		if err != nil {
			return nil, err
		}
		r.leftOver = int32(size)
	}

	mb, err := r.reader.ReadAtMost(r.leftOver)
	r.leftOver -= mb.Len()
	return mb, err
}
