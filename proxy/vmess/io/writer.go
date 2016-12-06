package io

import (
	"hash/fnv"

	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/serial"
)

type AuthChunkWriter struct {
	writer v2io.Writer
}

func NewAuthChunkWriter(writer v2io.Writer) *AuthChunkWriter {
	return &AuthChunkWriter{
		writer: writer,
	}
}

func (v *AuthChunkWriter) Write(buffer *alloc.Buffer) error {
	Authenticate(buffer)
	return v.writer.Write(buffer)
}

func (v *AuthChunkWriter) Release() {
	v.writer.Release()
	v.writer = nil
}

func Authenticate(buffer *alloc.Buffer) {
	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Bytes())
	buffer.PrependFunc(4, serial.WriteHash(fnvHash))

	buffer.PrependFunc(2, serial.WriteUint16(uint16(buffer.Len())))
}
