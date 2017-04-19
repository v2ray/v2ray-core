package buf

import "net"

type MultiBufferWriter interface {
	WriteMultiBuffer(MultiBuffer) (int, error)
}

type MultiBufferReader interface {
	ReadMultiBuffer() (MultiBuffer, error)
}

type MultiBuffer []*Buffer

func NewMultiBuffer() MultiBuffer {
	return MultiBuffer(make([]*Buffer, 0, 128))
}

func NewMultiBufferValue(b ...*Buffer) MultiBuffer {
	return MultiBuffer(b)
}

func (b *MultiBuffer) Append(buf *Buffer) {
	*b = append(*b, buf)
}

func (b *MultiBuffer) AppendMulti(mb MultiBuffer) {
	*b = append(*b, mb...)
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

func (mb MultiBuffer) Len() int {
	size := 0
	for _, b := range mb {
		size += b.Len()
	}
	return size
}

func (mb MultiBuffer) IsEmpty() bool {
	for _, b := range mb {
		if !b.IsEmpty() {
			return false
		}
	}
	return true
}

func (mb MultiBuffer) Release() {
	for i, b := range mb {
		b.Release()
		mb[i] = nil
	}
}

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
