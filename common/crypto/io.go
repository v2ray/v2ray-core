package crypto

import (
	"crypto/cipher"
	"io"
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
	nBytes, err := v.reader.Read(data)
	if nBytes > 0 {
		v.stream.XORKeyStream(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

type CryptionWriter struct {
	stream cipher.Stream
	writer io.Writer
}

// NewCryptionWriter creates a new CryptionWriter.
func NewCryptionWriter(stream cipher.Stream, writer io.Writer) *CryptionWriter {
	return &CryptionWriter{
		stream: stream,
		writer: writer,
	}
}

// Write implements io.Writer.Write().
func (v *CryptionWriter) Write(data []byte) (int, error) {
	v.stream.XORKeyStream(data, data)
	return v.writer.Write(data)
}
