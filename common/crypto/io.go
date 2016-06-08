package crypto

import (
	"crypto/cipher"
	"io"

	"github.com/v2ray/v2ray-core/common"
)

type CryptionReader struct {
	stream cipher.Stream
	reader io.Reader
}

func NewCryptionReader(stream cipher.Stream, reader io.Reader) *CryptionReader {
	return &CryptionReader{
		stream: stream,
		reader: reader,
	}
}

func (this *CryptionReader) Read(data []byte) (int, error) {
	if this.reader == nil {
		return 0, common.ErrorAlreadyReleased
	}
	nBytes, err := this.reader.Read(data)
	if nBytes > 0 {
		this.stream.XORKeyStream(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

func (this *CryptionReader) Release() {
	this.reader = nil
	this.stream = nil
}

type CryptionWriter struct {
	stream cipher.Stream
	writer io.Writer
}

func NewCryptionWriter(stream cipher.Stream, writer io.Writer) *CryptionWriter {
	return &CryptionWriter{
		stream: stream,
		writer: writer,
	}
}

func (this *CryptionWriter) Write(data []byte) (int, error) {
	if this.writer == nil {
		return 0, common.ErrorAlreadyReleased
	}
	this.stream.XORKeyStream(data, data)
	return this.writer.Write(data)
}

func (this *CryptionWriter) Release() {
	this.writer = nil
	this.stream = nil
}
