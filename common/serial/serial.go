package serial

import (
	"encoding/binary"
	"io"
	"unsafe"

	"v2ray.com/core/common/stack"
)

// ReadUint16 reads first two bytes from the reader, and then coverts them to an uint16 value.
func ReadUint16(reader io.Reader) (uint16, error) {
	var b stack.TwoBytes
	s := b[:]
	p := uintptr(unsafe.Pointer(&s))
	v := (*[]byte)(unsafe.Pointer(p))

	if _, err := io.ReadFull(reader, *v); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(*v), nil
}

// WriteUint16 writes an uint16 value into writer.
func WriteUint16(writer io.Writer, value uint16) (int, error) {
	var b stack.TwoBytes
	s := b[:]
	p := uintptr(unsafe.Pointer(&s))
	v := (*[]byte)(unsafe.Pointer(p))

	binary.BigEndian.PutUint16(*v, value)
	return writer.Write(*v)
}

// WriteUint64 writes an uint64 value into writer.
func WriteUint64(writer io.Writer, value uint64) (int, error) {
	var b stack.EightBytes
	s := b[:]
	p := uintptr(unsafe.Pointer(&s))
	v := (*[]byte)(unsafe.Pointer(p))

	binary.BigEndian.PutUint64(*v, value)
	return writer.Write(*v)
}
