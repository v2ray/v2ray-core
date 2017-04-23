package crypto

import (
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

type ChunkSizeDecoder interface {
	SizeBytes() int
	Decode([]byte) (uint16, error)
}

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

type ChunkStreamReader struct {
	sizeDecoder ChunkSizeDecoder
	reader      buf.Reader

	buffer       []byte
	leftOver     buf.MultiBuffer
	leftOverSize uint16
}

func NewChunkStreamReader(sizeDecoder ChunkSizeDecoder, reader io.Reader) *ChunkStreamReader {
	return &ChunkStreamReader{
		sizeDecoder: sizeDecoder,
		reader:      buf.NewReader(reader),
		buffer:      make([]byte, sizeDecoder.SizeBytes()),
	}
}

func (r *ChunkStreamReader) readAtLeast(size int) (buf.MultiBuffer, error) {
	mb := r.leftOver
	for mb.Len() < size {
		extra, err := r.reader.Read()
		if err != nil {
			mb.Release()
			return nil, err
		}
		mb.AppendMulti(extra)
	}

	return mb, nil
}

func (r *ChunkStreamReader) readSize() (uint16, error) {
	if r.sizeDecoder.SizeBytes() > r.leftOver.Len() {
		mb, err := r.readAtLeast(r.sizeDecoder.SizeBytes() - r.leftOver.Len())
		if err != nil {
			return 0, err
		}
		r.leftOver.AppendMulti(mb)
	}
	r.leftOver.Read(r.buffer)
	return r.sizeDecoder.Decode(r.buffer)
}

func (r *ChunkStreamReader) Read() (buf.MultiBuffer, error) {
	size := int(r.leftOverSize)
	if size == 0 {
		nextSize, err := r.readSize()
		if err != nil {
			return nil, err
		}
		if nextSize == 0 {
			return nil, io.EOF
		}
		size = int(nextSize)
	}

	leftOver := r.leftOver
	if leftOver.IsEmpty() {
		mb, err := r.readAtLeast(1)
		if err != nil {
			return nil, err
		}
		leftOver = mb
	}

	if size >= leftOver.Len() {
		r.leftOverSize = uint16(size - leftOver.Len())
		r.leftOver = nil
		return leftOver, nil
	}

	mb := leftOver.SliceBySize(size)
	if mb.Len() != size {
		b := buf.New()
		b.AppendSupplier(buf.ReadFullFrom(&leftOver, size-mb.Len()))
		mb.Append(b)
	}

	r.leftOver = leftOver
	r.leftOverSize = 0
	return mb, nil
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

func (w *ChunkStreamWriter) Write(mb buf.MultiBuffer) error {
	mb2Write := buf.NewMultiBuffer()
	const sliceSize = 8192

	for {
		slice := mb.SliceBySize(sliceSize)

		b := buf.New()
		b.AppendSupplier(func(buffer []byte) (int, error) {
			w.sizeEncoder.Encode(uint16(slice.Len()), buffer[:0])
			return w.sizeEncoder.SizeBytes(), nil
		})
		mb2Write.Append(b)
		mb2Write.AppendMulti(slice)

		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.Write(mb2Write)
}
