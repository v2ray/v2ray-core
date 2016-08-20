package kcp

import (
	"hash/fnv"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type SimpleAuthenticator struct{}

func NewSimpleAuthenticator() internet.Authenticator {
	return &SimpleAuthenticator{}
}

func (this *SimpleAuthenticator) Overhead() int {
	return 6
}

func (this *SimpleAuthenticator) Seal(buffer *alloc.Buffer) {
	buffer.PrependUint16(uint16(buffer.Len()))
	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Value)
	buffer.PrependHash(fnvHash)

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
