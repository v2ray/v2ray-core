package buf

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
)

// ReadAllToBytes reads all content from the reader into a byte array, until EOF.
func ReadAllToBytes(reader io.Reader) ([]byte, error) {
	mb, err := ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	if mb.Len() == 0 {
		return nil, nil
	}
	b := make([]byte, mb.Len())
	mb, _ = SplitBytes(mb, b)
	ReleaseMulti(mb)
	return b, nil
}

// MultiBuffer is a list of Buffers. The order of Buffer matters.
type MultiBuffer []*Buffer

// MergeMulti merges content from src to dest, and returns the new address of dest and src
func MergeMulti(dest MultiBuffer, src MultiBuffer) (MultiBuffer, MultiBuffer) {
	dest = append(dest, src...)
	for idx := range src {
		src[idx] = nil
	}
	return dest, src[:0]
}

func MergeBytes(dest MultiBuffer, src []byte) MultiBuffer {
	n := len(dest)
	if n > 0 && !(dest)[n-1].IsFull() {
		nBytes, _ := (dest)[n-1].Write(src)
		src = src[nBytes:]
	}

	for len(src) > 0 {
		b := New()
		nBytes, _ := b.Write(src)
		src = src[nBytes:]
		dest = append(dest, b)
	}

	return dest
}

// ReleaseMulti release all content of the MultiBuffer, and returns an empty MultiBuffer.
func ReleaseMulti(mb MultiBuffer) MultiBuffer {
	for i := range mb {
		mb[i].Release()
		mb[i] = nil
	}
	return mb[:0]
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

// ReadFrom reads all content from reader until EOF.
func ReadFrom(reader io.Reader) (MultiBuffer, error) {
	mb := make(MultiBuffer, 0, 16)
	for {
		b := New()
		_, err := b.ReadFullFrom(reader, Size)
		if b.IsEmpty() {
			b.Release()
		} else {
			mb = append(mb, b)
		}
		if err != nil {
			if errors.Cause(err) == io.EOF || errors.Cause(err) == io.ErrUnexpectedEOF {
				return mb, nil
			}
			return mb, err
		}
	}
}

func SplitBytes(mb MultiBuffer, b []byte) (MultiBuffer, int) {
	totalBytes := 0

	for len(mb) > 0 {
		bb := mb[0]
		nBytes, _ := bb.Read(b)
		totalBytes += nBytes
		b = b[nBytes:]
		if !bb.IsEmpty() {
			break
		}
		bb.Release()
		mb = mb[1:]
	}

	return mb, totalBytes
}

// Len returns the total number of bytes in the MultiBuffer.
func (mb MultiBuffer) Len() int32 {
	if mb == nil {
		return 0
	}

	size := int32(0)
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

// String returns the content of the MultiBuffer in string.
func (mb MultiBuffer) String() string {
	v := make([]interface{}, len(mb))
	for i, b := range mb {
		v[i] = b
	}
	return serial.Concat(v...)
}

// SliceBySize splits the beginning of this MultiBuffer into another one, for at most size bytes.
func (mb *MultiBuffer) SliceBySize(size int32) MultiBuffer {
	slice := make(MultiBuffer, 0, 10)
	sliceSize := int32(0)
	endIndex := len(*mb)
	for i, b := range *mb {
		if b.Len()+sliceSize > size {
			endIndex = i
			break
		}
		sliceSize += b.Len()
		slice = append(slice, b)
		(*mb)[i] = nil
	}
	*mb = (*mb)[endIndex:]
	if endIndex == 0 && len(*mb) > 0 {
		b := New()
		common.Must2(b.ReadFullFrom((*mb)[0], size))
		return MultiBuffer{b}
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

type MultiBufferContainer struct {
	MultiBuffer
}

func (c *MultiBufferContainer) Read(b []byte) (int, error) {
	if c.MultiBuffer.IsEmpty() {
		return 0, io.EOF
	}

	mb, nBytes := SplitBytes(c.MultiBuffer, b)
	c.MultiBuffer = mb
	return nBytes, nil
}

func (c *MultiBufferContainer) ReadMultiBuffer() (MultiBuffer, error) {
	mb := c.MultiBuffer
	c.MultiBuffer = nil
	return mb, nil
}

func (c *MultiBufferContainer) Write(b []byte) (int, error) {
	c.MultiBuffer = MergeBytes(c.MultiBuffer, b)
	return len(b), nil
}

func (c *MultiBufferContainer) WriteMultiBuffer(b MultiBuffer) error {
	mb, _ := MergeMulti(c.MultiBuffer, b)
	c.MultiBuffer = mb
	return nil
}

func (c *MultiBufferContainer) Close() error {
	c.MultiBuffer = ReleaseMulti(c.MultiBuffer)
	return nil
}
