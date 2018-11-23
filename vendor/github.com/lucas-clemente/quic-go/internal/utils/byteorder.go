package utils

import (
	"bytes"
	"io"
)

// A ByteOrder specifies how to convert byte sequences into 16-, 32-, or 64-bit unsigned integers.
type ByteOrder interface {
	ReadUintN(b io.ByteReader, length uint8) (uint64, error)
	ReadUint64(io.ByteReader) (uint64, error)
	ReadUint32(io.ByteReader) (uint32, error)
	ReadUint16(io.ByteReader) (uint16, error)

	WriteUint64(*bytes.Buffer, uint64)
	WriteUint32(*bytes.Buffer, uint32)
	WriteUint16(*bytes.Buffer, uint16)
}
