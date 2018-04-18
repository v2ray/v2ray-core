package buf

import (
	"io"
	"net"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
)

// ReadAllToMultiBuffer reads all content from the reader into a MultiBuffer, until EOF.
func ReadAllToMultiBuffer(reader io.Reader) (MultiBuffer, error) {
	mb := NewMultiBufferCap(128)

	if _, err := mb.ReadFrom(reader); err != nil {
		mb.Release()
		return nil, err
	}

	return mb, nil
}

// ReadSizeToMultiBuffer reads specific number of bytes from reader into a MultiBuffer.
func ReadSizeToMultiBuffer(reader io.Reader, size int32) (MultiBuffer, error) {
	mb := NewMultiBufferCap(32)

	for size > 0 {
		bSize := size
		if bSize > Size {
			bSize = Size
		}
		b := NewSize(bSize)
		if err := b.Reset(ReadFullFrom(reader, bSize)); err != nil {
			mb.Release()
			return nil, err
		}
		size -= bSize
		mb.Append(b)
	}
	return mb, nil
}

// ReadAllToBytes reads all content from the reader into a byte array, until EOF.
func ReadAllToBytes(reader io.Reader) ([]byte, error) {
	mb, err := ReadAllToMultiBuffer(reader)
	if err != nil {
		return nil, err
	}
	b := make([]byte, mb.Len())
	common.Must2(mb.Read(b))
	mb.Release()
	return b, nil
}

// MultiBuffer is a list of Buffers. The order of Buffer matters.
type MultiBuffer []*Buffer

// NewMultiBufferCap creates a new MultiBuffer instance.
func NewMultiBufferCap(capacity int32) MultiBuffer {
	return MultiBuffer(make([]*Buffer, 0, capacity))
}

// NewMultiBufferValue wraps a list of Buffers into MultiBuffer.
func NewMultiBufferValue(b ...*Buffer) MultiBuffer {
	return MultiBuffer(b)
}

// Append appends buffer to the end of this MultiBuffer
func (mb *MultiBuffer) Append(buf *Buffer) {
	if buf != nil {
		*mb = append(*mb, buf)
	}
}

// AppendMulti appends a MultiBuffer to the end of this one.
func (mb *MultiBuffer) AppendMulti(buf MultiBuffer) {
	*mb = append(*mb, buf...)
}

// Copy copied the beginning part of the MultiBuffer into the given byte array.
func (mb MultiBuffer) Copy(b []byte) int {
	total := 0
	for _, bb := range mb {
		nBytes := copy(b[total:], bb.Bytes())
		total += nBytes
		if int32(nBytes) < bb.Len() {
			break
		}
	}
	return total
}

// ReadFrom implements io.ReaderFrom.
func (mb *MultiBuffer) ReadFrom(reader io.Reader) (int64, error) {
	totalBytes := int64(0)

	for {
		b := New()
		err := b.Reset(ReadFrom(reader))
		if b.IsEmpty() {
			b.Release()
		} else {
			mb.Append(b)
		}
		totalBytes += int64(b.Len())
		if err != nil {
			if errors.Cause(err) == io.EOF {
				return totalBytes, nil
			}
			return totalBytes, err
		}
	}
}

// Read implements io.Reader.
func (mb *MultiBuffer) Read(b []byte) (int, error) {
	if mb.Len() == 0 {
		return 0, io.EOF
	}
	endIndex := len(*mb)
	totalBytes := 0
	for i, bb := range *mb {
		nBytes, _ := bb.Read(b)
		totalBytes += nBytes
		b = b[nBytes:]
		if bb.IsEmpty() {
			bb.Release()
			(*mb)[i] = nil
		} else {
			endIndex = i
			break
		}
	}
	*mb = (*mb)[endIndex:]
	return totalBytes, nil
}

// WriteTo implements io.WriterTo.
func (mb *MultiBuffer) WriteTo(writer io.Writer) (int64, error) {
	defer mb.Release()

	totalBytes := int64(0)
	for _, b := range *mb {
		nBytes, err := writer.Write(b.Bytes())
		totalBytes += int64(nBytes)
		if err != nil {
			return totalBytes, err
		}
	}

	return totalBytes, nil
}

// Write implements io.Writer.
func (mb *MultiBuffer) Write(b []byte) (int, error) {
	totalBytes := len(b)

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

	return totalBytes, nil
}

// WriteMultiBuffer implements Writer.
func (mb *MultiBuffer) WriteMultiBuffer(b MultiBuffer) error {
	*mb = append(*mb, b...)
	return nil
}

// Len returns the total number of bytes in the MultiBuffer.
func (mb *MultiBuffer) Len() int32 {
	if mb == nil {
		return 0
	}

	size := int32(0)
	for _, b := range *mb {
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
func (mb *MultiBuffer) Release() {
	for i, b := range *mb {
		b.Release()
		(*mb)[i] = nil
	}
	*mb = nil
}

func (mb MultiBuffer) String() string {
	v := make([]interface{}, len(mb))
	for i, b := range mb {
		v[i] = b
	}
	return serial.Concat(v...)
}

// ToNetBuffers converts this MultiBuffer to net.Buffers. The return net.Buffers points to the same content of the MultiBuffer.
func (mb MultiBuffer) ToNetBuffers() net.Buffers {
	bs := make([][]byte, len(mb))
	for i, b := range mb {
		bs[i] = b.Bytes()
	}
	return bs
}

// SliceBySize splits the beginning of this MultiBuffer into another one, for at most size bytes.
func (mb *MultiBuffer) SliceBySize(size int32) MultiBuffer {
	slice := NewMultiBufferCap(10)
	sliceSize := int32(0)
	endIndex := len(*mb)
	for i, b := range *mb {
		if b.Len()+sliceSize > size {
			endIndex = i
			break
		}
		sliceSize += b.Len()
		slice.Append(b)
		(*mb)[i] = nil
	}
	*mb = (*mb)[endIndex:]
	if endIndex == 0 && len(*mb) > 0 {
		b := NewSize(size)
		common.Must(b.Reset(ReadFullFrom((*mb)[0], size)))
		return NewMultiBufferValue(b)
	}
	return slice
}

// SplitFirst splits out the first Buffer in this MultiBuffer.
func (mb *MultiBuffer) SplitFirst() *Buffer {
	if len(*mb) == 0 {
		return nil
	}
	b := (*mb)[0]
	(*mb)[0] = nil
	*mb = (*mb)[1:]
	return b
}
