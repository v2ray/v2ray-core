package buf

import "io"

type MultiBufferWriter interface {
	WriteMultiBuffer(MultiBuffer) (int, error)
}

type MultiBufferReader interface {
	ReadMultiBuffer() (MultiBuffer, error)
}

type MultiBuffer []*Buffer

func NewMultiBuffer() MultiBuffer {
	return MultiBuffer(make([]*Buffer, 0, 32))
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
	if len(*mb) == 0 {
		return 0, io.EOF
	}
	endIndex := len(*mb)
	totalBytes := 0
	for i, bb := range *mb {
		nBytes, err := bb.Read(b)
		totalBytes += nBytes
		if err != nil {
			return totalBytes, err
		}
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
	for _, b := range mb {
		b.Release()
	}
}
