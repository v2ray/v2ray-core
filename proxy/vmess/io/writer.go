package io

import (
	"hash/fnv"

	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
)

type AuthChunkWriter struct {
	writer v2io.Writer
}

func NewAuthChunkWriter(writer v2io.Writer) *AuthChunkWriter {
	return &AuthChunkWriter{
		writer: writer,
	}
}

func (this *AuthChunkWriter) Write(buffer *alloc.Buffer) error {
	Authenticate(buffer)
	return this.writer.Write(buffer)
}

func (this *AuthChunkWriter) Release() {
	this.writer.Release()
	this.writer = nil
}

func Authenticate(buffer *alloc.Buffer) {
	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Value)
	buffer.PrependHash(fnvHash)

	buffer.PrependUint16(uint16(buffer.Len()))
}
