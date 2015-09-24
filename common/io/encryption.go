package io

import (
	"crypto/cipher"
	"io"

	"github.com/v2ray/v2ray-core/common/log"
)

// CryptionReader is a general purpose reader that applies a stream cipher on top of a regular reader.
type CryptionReader struct {
	stream cipher.Stream
	reader io.Reader
}

// NewCryptionReader creates a new CryptionReader instance from given stream cipher and reader.
func NewCryptionReader(stream cipher.Stream, reader io.Reader) *CryptionReader {
	return &CryptionReader{
		stream: stream,
		reader: reader,
	}
}

// Read reads blocks from underlying reader, and crypt it. The content of blocks is modified in place.
func (reader CryptionReader) Read(blocks []byte) (int, error) {
	nBytes, err := reader.reader.Read(blocks)
	if nBytes > 0 {
		reader.stream.XORKeyStream(blocks[:nBytes], blocks[:nBytes])
	}
	return nBytes, err
}

// Cryption writer is a general purpose of byte stream writer that applies a stream cipher on top of a regular writer.
type CryptionWriter struct {
	stream cipher.Stream
	writer io.Writer
}

// NewCryptionWriter creates a new CryptionWriter from given stream cipher and writer.
func NewCryptionWriter(stream cipher.Stream, writer io.Writer) *CryptionWriter {
	return &CryptionWriter{
		stream: stream,
		writer: writer,
	}
}

// Crypt crypts the content of blocks without writing them into the underlying writer.
func (writer CryptionWriter) Crypt(blocks []byte) {
	writer.stream.XORKeyStream(blocks, blocks)
}

// Write crypts the content of blocks in place, and then writes the give blocks to underlying writer.
func (writer CryptionWriter) Write(blocks []byte) (int, error) {
	writer.Crypt(blocks)
	return writer.writer.Write(blocks)
}
