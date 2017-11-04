package buf

import (
	"io"
	"net"

	"v2ray.com/core/common/errors"
)

type MultiBufferWriter interface {
	WriteMultiBuffer(MultiBuffer) error
}

type MultiBufferReader interface {
	ReadMultiBuffer() (MultiBuffer, error)
}

func ReadAllToMultiBuffer(reader io.Reader) (MultiBuffer, error) {
	mb := NewMultiBuffer()

	for {
		b := New()
		err := b.AppendSupplier(ReadFrom(reader))
		if !b.IsEmpty() {
			mb.Append(b)
		}
		if err != nil {
			if errors.Cause(err) == io.EOF {
				return mb, nil
			}
			mb.Release()
			return nil, err
		}
	}
}

func ReadAllToBytes(reader io.Reader) ([]byte, error) {
	mb, err := ReadAllToMultiBuffer(reader)
	if err != nil {
		return nil, err
	}
	b := make([]byte, mb.Len())
	mb.Read(b)
	mb.Release()
	return b, nil
}

// MultiBuffer is a list of Buffers. The order of Buffer matters.
type MultiBuffer []*Buffer

// NewMultiBuffer creates a new MultiBuffer instance.
func NewMultiBuffer() MultiBuffer {
	return MultiBuffer(make([]*Buffer, 0, 128))
}

// NewMultiBufferValue wraps a list of Buffers into MultiBuffer.
func NewMultiBufferValue(b ...*Buffer) MultiBuffer {
	return MultiBuffer(b)
}

func (mb *MultiBuffer) Append(buf *Buffer) {
	*mb = append(*mb, buf)
}

func (mb *MultiBuffer) AppendMulti(buf MultiBuffer) {
	*mb = append(*mb, buf...)
}

func (mb MultiBuffer) Copy(b []byte) int {
	total := 0
	for _, bb := range mb {
		nBytes := copy(b[total:], bb.Bytes())
		total += nBytes
		if nBytes < bb.Len() {
			break
		}
	}
	return total
}

func (mb *MultiBuffer) Read(b []byte) (int, error) {
	endIndex := len(*mb)
	totalBytes := 0
	for i, bb := range *mb {
		nBytes, _ := bb.Read(b)
		totalBytes += nBytes
		b = b[nBytes:]
		if bb.IsEmpty() {
			bb.Release()
		} else {
			endIndex = i
			break
		}
	}
	*mb = (*mb)[endIndex:]
	return totalBytes, nil
}

func (mb *MultiBuffer) Write(b []byte) {
	n := len(*mb)
	if n > 0 && !(*mb)[n-1].IsFull() {
		nBytes, _ := (*mb)[n-1].Write(b)
		b = b[nBytes:]
	}

	for len(b) > 0 {
		bb := New()
		nBytes, _ := bb.Write(b)
		b = b[nBytes:]
		mb.Append(bb)
	}
}

// Len returns the total number of bytes in the MultiBuffer.
func (mb MultiBuffer) Len() int {
	size := 0
	for _, b := range mb {
		size += b.Len()
	}
	return size
}

// IsEmpty return true if the MultiBuffer has no content.
func (mb MultiBuffer) IsEmpty() bool {
	for _, b := range mb {
		if !b.IsEmpty() {
			return false
		}
	}
	return true
}

// Release releases all Buffers in the MultiBuffer.
func (mb MultiBuffer) Release() {
	for i, b := range mb {
		b.Release()
		mb[i] = nil
	}
}

// ToNetBuffers converts this MultiBuffer to net.Buffers. The return net.Buffers points to the same content of the MultiBuffer.
func (mb MultiBuffer) ToNetBuffers() net.Buffers {
	bs := make([][]byte, len(mb))
	for i, b := range mb {
		bs[i] = b.Bytes()
	}
	return bs
}

func (mb *MultiBuffer) SliceBySize(size int) MultiBuffer {
	slice := NewMultiBuffer()
	sliceSize := 0
	endIndex := len(*mb)
	for i, b := range *mb {
		if b.Len()+sliceSize > size {
			endIndex = i
			break
		}
		sliceSize += b.Len()
		slice.Append(b)
	}
	*mb = (*mb)[endIndex:]
	return slice
}

func (mb *MultiBuffer) SplitFirst() *Buffer {
	if len(*mb) == 0 {
		return nil
	}
	b := (*mb)[0]
	*mb = (*mb)[1:]
	return b
}
