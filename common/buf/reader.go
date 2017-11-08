package buf

import (
	"io"

	"v2ray.com/core/common/errors"
)

// BytesToBufferReader is a Reader that adjusts its reading speed automatically.
type BytesToBufferReader struct {
	reader io.Reader
	buffer []byte
}

func NewBytesToBufferReader(reader io.Reader) Reader {
	return &BytesToBufferReader{
		reader: reader,
	}
}

func (r *BytesToBufferReader) readSmall() (MultiBuffer, error) {
	b := New()
	if err := b.Reset(ReadFrom(r.reader)); err != nil {
		b.Release()
		return nil, err
	}
	if b.IsFull() {
		r.buffer = make([]byte, 32*1024)
	}
	return NewMultiBufferValue(b), nil
}

// Read implements Reader.Read().
func (r *BytesToBufferReader) Read() (MultiBuffer, error) {
	if r.buffer == nil {
		return r.readSmall()
	}

	nBytes, err := r.reader.Read(r.buffer)
	if err != nil {
		return nil, err
	}

	mb := NewMultiBuffer()
	mb.Write(r.buffer[:nBytes])
	return mb, nil
}

type readerAdpater struct {
	MultiBufferReader
}

func (r *readerAdpater) Read() (MultiBuffer, error) {
	return r.ReadMultiBuffer()
}

type bufferToBytesReader struct {
	stream   Reader
	leftOver MultiBuffer
}

func (r *bufferToBytesReader) Read(b []byte) (int, error) {
	if r.leftOver != nil {
		nBytes, _ := r.leftOver.Read(b)
		if r.leftOver.IsEmpty() {
			r.leftOver.Release()
			r.leftOver = nil
		}
		return nBytes, nil
	}

	mb, err := r.stream.Read()
	if err != nil {
		return 0, err
	}

	nBytes, _ := mb.Read(b)
	if !mb.IsEmpty() {
		r.leftOver = mb
	}
	return nBytes, nil
}

func (r *bufferToBytesReader) ReadMultiBuffer() (MultiBuffer, error) {
	if r.leftOver != nil {
		mb := r.leftOver
		r.leftOver = nil
		return mb, nil
	}

	return r.stream.Read()
}

func (r *bufferToBytesReader) writeToInternal(writer io.Writer) (int64, error) {
	mbWriter := NewWriter(writer)
	totalBytes := int64(0)
	if r.leftOver != nil {
		totalBytes += int64(r.leftOver.Len())
		if err := mbWriter.Write(r.leftOver); err != nil {
			return 0, err
		}
	}

	for {
		mb, err := r.stream.Read()
		if err != nil {
			return totalBytes, err
		}
		totalBytes += int64(mb.Len())
		if err := mbWriter.Write(mb); err != nil {
			return totalBytes, err
		}
	}
}

func (r *bufferToBytesReader) WriteTo(writer io.Writer) (int64, error) {
	nBytes, err := r.writeToInternal(writer)
	if errors.Cause(err) == io.EOF {
		return nBytes, nil
	}
	return nBytes, err
}
