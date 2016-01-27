package shadowsocks

import (
	"io"

	"github.com/v2ray/v2ray-core/common/crypto"
)

type Cipher interface {
	KeySize() int
	IVSize() int
	NewEncodingStream(key []byte, iv []byte, writer io.Writer) (io.Writer, error)
	NewDecodingStream(key []byte, iv []byte, reader io.Reader) (io.Reader, error)
}

type AesCfb struct {
	KeyBytes int
}

func (this *AesCfb) KeySize() int {
	return this.KeyBytes
}

func (this *AesCfb) IVSize() int {
	return 16
}

func (this *AesCfb) NewEncodingStream(key []byte, iv []byte, writer io.Writer) (io.Writer, error) {
	stream, err := crypto.NewAesEncryptionStream(key, iv)
	if err != nil {
		return nil, err
	}
	aesWriter := crypto.NewCryptionWriter(stream, writer)
	return aesWriter, nil
}

func (this *AesCfb) NewDecodingStream(key []byte, iv []byte, reader io.Reader) (io.Reader, error) {
	stream, err := crypto.NewAesDecryptionStream(key, iv)
	if err != nil {
		return nil, err
	}
	aesReader := crypto.NewCryptionReader(stream, reader)
	return aesReader, nil
}

type Config struct {
	Cipher   Cipher
	Password string
}
