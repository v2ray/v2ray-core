package vmess

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"

	v2io "github.com/v2ray/v2ray-core/io"
)

const (
	blockSize = 16 // Decryption block size, inherited from AES
)

// DecryptionReader is a byte stream reader to decrypt AES-128 CBC (for now)
// encrypted content.
type DecryptionReader struct {
	reader *v2io.CryptionReader
	buffer *bytes.Buffer
}

// NewDecryptionReader creates a new DescriptionReader by given byte Reader and
// AES key.
func NewDecryptionReader(reader io.Reader, key []byte, iv []byte) (*DecryptionReader, error) {
	decryptionReader := new(DecryptionReader)
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesStream := cipher.NewCFBDecrypter(aesCipher, iv)
	decryptionReader.reader = v2io.NewCryptionReader(aesStream, reader)
	decryptionReader.buffer = bytes.NewBuffer(make([]byte, 0, 2*blockSize))
	return decryptionReader, nil
}

func (reader *DecryptionReader) readBlock() error {
	buffer := make([]byte, blockSize)
	nBytes, err := reader.reader.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}
	if nBytes < blockSize {
		return fmt.Errorf("Expected to read %d bytes, but got %d bytes", blockSize, nBytes)
	}
	reader.buffer.Write(buffer)
	return err
}

// Read returns decrypted bytes of given length
func (reader *DecryptionReader) Read(p []byte) (int, error) {
	nBytes, err := reader.buffer.Read(p)
	if err != nil && err != io.EOF {
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
