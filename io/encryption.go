package io

import (
	"crypto/cipher"
	"io"
)

// CryptionReader is a general purpose reader that applies
// block cipher on top of a regular reader.
type CryptionReader struct {
	mode   cipher.BlockMode
	reader io.Reader
}

func NewCryptionReader(mode cipher.BlockMode, reader io.Reader) *CryptionReader {
	this := new(CryptionReader)
	this.mode = mode
	this.reader = reader
	return this
}

// Read reads blocks from underlying reader, the length of blocks must be
// a multiply of BlockSize()
func (reader CryptionReader) Read(blocks []byte) (int, error) {
	nBytes, err := reader.reader.Read(blocks)
	if err != nil && err != io.EOF {
		return nBytes, err
	}
	if nBytes < len(blocks) {
		for i, _ := range blocks[nBytes:] {
			blocks[i] = 0
		}
	}
	reader.mode.CryptBlocks(blocks, blocks)
	return nBytes, err
}

func (reader CryptionReader) BlockSize() int {
	return reader.mode.BlockSize()
}

// Cryption writer is a general purpose of byte stream writer that applies
// block cipher on top of a regular writer.
type CryptionWriter struct {
	mode   cipher.BlockMode
	writer io.Writer
}

func NewCryptionWriter(mode cipher.BlockMode, writer io.Writer) *CryptionWriter {
	this := new(CryptionWriter)
	this.mode = mode
	this.writer = writer
	return this
}

// Write writes the give blocks to underlying writer. The length of the blocks
// must be a multiply of BlockSize()
func (writer CryptionWriter) Write(blocks []byte) (int, error) {
	writer.mode.CryptBlocks(blocks, blocks)
	return writer.writer.Write(blocks)
}

func (writer CryptionWriter) BlockSize() int {
	return writer.mode.BlockSize()
}
