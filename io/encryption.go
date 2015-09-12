package io

import (
	"crypto/cipher"
	"io"

	"github.com/v2ray/v2ray-core/log"
)

// CryptionReader is a general purpose reader that applies
// block cipher on top of a regular reader.
type CryptionReader struct {
	stream cipher.Stream
	reader io.Reader
}

func NewCryptionReader(stream cipher.Stream, reader io.Reader) *CryptionReader {
	this := new(CryptionReader)
	this.stream = stream
	this.reader = reader
	return this
}

// Read reads blocks from underlying reader, the length of blocks must be
// a multiply of BlockSize()
func (reader CryptionReader) Read(blocks []byte) (int, error) {
	nBytes, err := reader.reader.Read(blocks)
	log.Debug("CryptionReader: Read %d bytes", nBytes)
	if nBytes > 0 {
		reader.stream.XORKeyStream(blocks[:nBytes], blocks[:nBytes])
	}
	if err != nil {
		log.Error("Error reading blocks: %v", err)
	}
	return nBytes, err
}

// Cryption writer is a general purpose of byte stream writer that applies
// block cipher on top of a regular writer.
type CryptionWriter struct {
	stream cipher.Stream
	writer io.Writer
}

func NewCryptionWriter(stream cipher.Stream, writer io.Writer) *CryptionWriter {
	this := new(CryptionWriter)
	this.stream = stream
	this.writer = writer
	return this
}

// Write writes the give blocks to underlying writer. The length of the blocks
// must be a multiply of BlockSize()
func (writer CryptionWriter) Write(blocks []byte) (int, error) {
	log.Debug("CryptionWriter writing %d bytes", len(blocks))
	writer.stream.XORKeyStream(blocks, blocks)
	return writer.writer.Write(blocks)
}
