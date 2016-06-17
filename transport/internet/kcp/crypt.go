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
	var length uint16 = uint16(buffer.Len())
	buffer.Prepend(serial.Uint16ToBytes(length))
	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Value)

	buffer.SliceBack(4)
	fnvHash.Sum(buffer.Value[:0])

	for i := 4; i < buffer.Len(); i++ {
		buffer.Value[i] ^= buffer.Value[i-4]
	}
}

func (this *SimpleAuthenticator) Open(buffer *alloc.Buffer) bool {
	for i := buffer.Len() - 1; i >= 4; i-- {
		buffer.Value[i] ^= buffer.Value[i-4]
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
