package io

import (
	"hash/fnv"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/transport"
)

type AuthChunkReader struct {
	reader io.Reader
}

func NewAuthChunkReader(reader io.Reader) *AuthChunkReader {
	return &AuthChunkReader{
		reader: reader,
	}
}

func (this *AuthChunkReader) Read() (*alloc.Buffer, error) {
	buffer := alloc.NewBuffer()
	if _, err := io.ReadFull(this.reader, buffer.Value[:2]); err != nil {
		buffer.Release()
		return nil, err
	}

	length := serial.BytesLiteral(buffer.Value[:2]).Uint16Value()
	if _, err := io.ReadFull(this.reader, buffer.Value[:length]); err != nil {
		buffer.Release()
		return nil, err
	}
	buffer.Slice(0, int(length))

	fnvHash := fnv.New32a()
	fnvHash.Write(buffer.Value[4:])
	expAuth := serial.BytesLiteral(fnvHash.Sum(nil))
	actualAuth := serial.BytesLiteral(buffer.Value[:4])
	if !actualAuth.Equals(expAuth) {
		buffer.Release()
		return nil, transport.ErrorCorruptedPacket
	}
	buffer.SliceFrom(4)
	return buffer, nil
}
