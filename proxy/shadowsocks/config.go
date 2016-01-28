package shadowsocks

import (
	"crypto/md5"
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
	Cipher Cipher
	Key    []byte
	UDP    bool
}

func PasswordToCipherKey(password string, keySize int) []byte {
	pwdBytes := []byte(password)
	key := make([]byte, 0, keySize)

	md5Sum := md5.Sum(pwdBytes)
	key = append(key, md5Sum[:]...)

	for len(key) < keySize {
		md5Hash := md5.New()
		md5Hash.Write(md5Sum[:])
		md5Hash.Write(pwdBytes)
		md5Hash.Sum(md5Sum[:0])

		key = append(key, md5Sum[:]...)
	}
	return key
}
