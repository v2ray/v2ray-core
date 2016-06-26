package kcp

import (
	"hash/fnv"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/serial"
)

type Authenticator interface {
	HeaderSize() int
	// Encrypt encrypts the whole block in src into dst.
	// Dst and src may point at the same memory.
	Seal(buffer *alloc.Buffer)

	// Decrypt decrypts the whole block in src into dst.
	// Dst and src may point at the same memory.
	Open(buffer *alloc.Buffer) bool
}

type SimpleAuthenticator struct{}

func NewSimpleAuthenticator() Authenticator {
	return &SimpleAuthenticator{}
}

func (this *SimpleAuthenticator) HeaderSize() int {
	return 6
}

func (this *SimpleAuthenticator) Seal(buffer *alloc.Buffer) {
	buffer.PrependUint16(uint16(buffer.Len()))
	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Value)

	buffer.SliceBack(4)
	fnvHash.Sum(buffer.Value[:0])

	len := buffer.Len()
	xtra := 4 - len%4
	if xtra != 0 {
		buffer.Slice(0, len+xtra)
	}
	xorfwd(buffer.Value)
	if xtra != 0 {
		buffer.Slice(0, len)
	}
}

func (this *SimpleAuthenticator) Open(buffer *alloc.Buffer) bool {
	len := buffer.Len()
	xtra := 4 - len%4
	if xtra != 0 {
		buffer.Slice(0, len+xtra)
	}
	xorbkd(buffer.Value)
	if xtra != 0 {
		buffer.Slice(0, len)
	}

	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Value[4:])
	if serial.BytesToUint32(buffer.Value[:4]) != fnvHash.Sum32() {
		return false
	}

	length := serial.BytesToUint16(buffer.Value[4:6])
	if buffer.Len()-6 != int(length) {
		return false
	}

	buffer.SliceFrom(6)

	return true
}
