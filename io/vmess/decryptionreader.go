package vmess

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
)

const (
	blockSize = 16
)

type DecryptionReader struct {
	cipher cipher.Block
	reader io.Reader
	buffer *bytes.Buffer
}

func NewDecryptionReader(reader io.Reader, key []byte) (*DecryptionReader, error) {
	decryptionReader := new(DecryptionReader)
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decryptionReader.cipher = cipher
	decryptionReader.reader = reader
	decryptionReader.buffer = bytes.NewBuffer(make([]byte, 0, 2*blockSize))
	return decryptionReader, nil
}

func (reader *DecryptionReader) readBlock() error {
	buffer := make([]byte, blockSize)
	nBytes, err := reader.reader.Read(buffer)
	if err != nil {
		return err
	}
	if nBytes < blockSize {
		return fmt.Errorf("Expected to read %d bytes, but got %d bytes", blockSize, nBytes)
	}
	reader.cipher.Decrypt(buffer, buffer)
	reader.buffer.Write(buffer)
	return nil
}

func (reader *DecryptionReader) Read(p []byte) (int, error) {
	if reader.buffer.Len() == 0 {
		err := reader.readBlock()
		if err != nil {
			return 0, err
		}
	}
	nBytes, err := reader.buffer.Read(p)
	if err != nil {
		return nBytes, err
	}
	if nBytes < len(p) {
		err = reader.readBlock()
		if err != nil {
			return nBytes, err
		}
		moreBytes, err := reader.buffer.Read(p[nBytes:])
		if err != nil {
			return nBytes, err
		}
		nBytes += moreBytes
		if nBytes != len(p) {
			return nBytes, fmt.Errorf("Unable to read %d bytes", len(p))
		}
	}
	return nBytes, err
}
