package crypto

import (
	"crypto/cipher"
	"io"
)

type cryptionReader struct {
	stream cipher.Stream
	reader io.Reader
}

func NewCryptionReader(stream cipher.Stream, reader io.Reader) io.Reader {
	return &cryptionReader{
		stream: stream,
		reader: reader,
	}
}

func (this *cryptionReader) Read(data []byte) (int, error) {
	nBytes, err := this.reader.Read(data)
	if nBytes > 0 {
		this.stream.XORKeyStream(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

type cryptionWriter struct {
	stream cipher.Stream
	writer io.Writer
}

func NewCryptionWriter(stream cipher.Stream, writer io.Writer) io.Writer {
	return &cryptionWriter{
		stream: stream,
		writer: writer,
	}
}

func (this *cryptionWriter) Write(data []byte) (int, error) {
	this.stream.XORKeyStream(data, data)
	return this.writer.Write(data)
}
