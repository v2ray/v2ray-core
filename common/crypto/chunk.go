package crypto

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

// ChunkSizeDecoder is a utility class to decode size value from bytes.
type ChunkSizeDecoder interface {
	SizeBytes() int
	Decode([]byte) (uint16, error)
}

// ChunkSizeEncoder is a utility class to encode size value into bytes.
type ChunkSizeEncoder interface {
	SizeBytes() int
	Encode(uint16, []byte) []byte
}

type PlainChunkSizeParser struct{}

func (PlainChunkSizeParser) SizeBytes() int {
	return 2
}

func (PlainChunkSizeParser) Encode(size uint16, b []byte) []byte {
	return serial.Uint16ToBytes(size, b)
}

func (PlainChunkSizeParser) Decode(b []byte) (uint16, error) {
	return serial.BytesToUint16(b), nil
}

type AEADChunkSizeParser struct {
	Auth *AEADAuthenticator
}

func (p *AEADChunkSizeParser) SizeBytes() int {
	return 2 + p.Auth.Overhead()
}

func (p *AEADChunkSizeParser) Encode(size uint16, b []byte) []byte {
	b = serial.Uint16ToBytes(size-uint16(p.Auth.Overhead()), b)
	b, err := p.Auth.Seal(b[:0], b)
	common.Must(err)
	return b
}

func (p *AEADChunkSizeParser) Decode(b []byte) (uint16, error) {
	b, err := p.Auth.Open(b[:0], b)
	if err != nil {
		return 0, err
	}
	return serial.BytesToUint16(b) + uint16(p.Auth.Overhead()), nil
}

type ChunkStreamReader struct {
	sizeDecoder ChunkSizeDecoder
	reader      *buf.BufferedReader

	buffer       []byte
	leftOverSize int32
}

func NewChunkStreamReader(sizeDecoder ChunkSizeDecoder, reader io.Reader) *ChunkStreamReader {
	return &ChunkStreamReader{
		sizeDecoder: sizeDecoder,
		reader:      buf.NewBufferedReader(buf.NewReader(reader)),
		buffer:      make([]byte, sizeDecoder.SizeBytes()),
	}
}

func (r *ChunkStreamReader) readSize() (uint16, error) {
	if _, err := io.ReadFull(r.reader, r.buffer); err != nil {
		return 0, err
	}
	return r.sizeDecoder.Decode(r.buffer)
}

func (r *ChunkStreamReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	size := r.leftOverSize
	if size == 0 {
		nextSize, err := r.readSize()
		if err != nil {
			return nil, err
		}
		if nextSize == 0 {
			return nil, io.EOF
		}
		size = int32(nextSize)
	}
	r.leftOverSize = size

	mb, err := r.reader.ReadAtMost(size)
	if !mb.IsEmpty() {
		r.leftOverSize -= mb.Len()
		return mb, nil
	}
	return nil, err
}

type ChunkStreamWriter struct {
	sizeEncoder ChunkSizeEncoder
	writer      buf.Writer
}

func NewChunkStreamWriter(sizeEncoder ChunkSizeEncoder, writer io.Writer) *ChunkStreamWriter {
	return &ChunkStreamWriter{
		sizeEncoder: sizeEncoder,
		writer:      buf.NewWriter(writer),
	}
}

func (w *ChunkStreamWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	const sliceSize = 8192
	mbLen := mb.Len()
	mb2Write := buf.NewMultiBufferCap(mbLen/buf.Size + mbLen/sliceSize + 2)

	for {
		slice := mb.SliceBySize(sliceSize)

		b := buf.New()
		common.Must(b.Reset(func(buffer []byte) (int, error) {
			w.sizeEncoder.Encode(uint16(slice.Len()), buffer[:0])
			return w.sizeEncoder.SizeBytes(), nil
		}))
		mb2Write.Append(b)
		mb2Write.AppendMulti(slice)

		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.WriteMultiBuffer(mb2Write)
}
