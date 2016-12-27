package crypto

import (
	"crypto/cipher"
	"io"

	"v2ray.com/core/common"
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

func (v *CryptionReader) Read(data []byte) (int, error) {
	if v.reader == nil {
		return 0, common.ErrObjectReleased
	}
	nBytes, err := v.reader.Read(data)
	if nBytes > 0 {
		v.stream.XORKeyStream(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

func (v *CryptionReader) Release() {
	common.Release(v.reader)
	common.Release(v.stream)
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

func (v *CryptionWriter) Write(data []byte) (int, error) {
	if v.writer == nil {
		return 0, common.ErrObjectReleased
	}
	v.stream.XORKeyStream(data, data)
	return v.writer.Write(data)
}

func (v *CryptionWriter) Release() {
	common.Release(v.writer)
	common.Release(v.stream)
}
