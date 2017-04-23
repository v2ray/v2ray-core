package buf

import "net"

type MultiBufferWriter interface {
	WriteMultiBuffer(MultiBuffer) error
}

type MultiBufferReader interface {
	ReadMultiBuffer() (MultiBuffer, error)
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
