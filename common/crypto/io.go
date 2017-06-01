package crypto

import (
	"crypto/cipher"
	"io"

	"v2ray.com/core/common/buf"
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

func (r *CryptionReader) Read(data []byte) (int, error) {
	nBytes, err := r.reader.Read(data)
	if nBytes > 0 {
		r.stream.XORKeyStream(data[:nBytes], data[:nBytes])
	}
	return nBytes, err
}

var (
	_ buf.MultiBufferWriter = (*CryptionWriter)(nil)
)

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
func (w *CryptionWriter) Write(data []byte) (int, error) {
	w.stream.XORKeyStream(data, data)
	return w.writer.Write(data)
}

func (w *CryptionWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	bs := mb.ToNetBuffers()
	for _, b := range bs {
		w.stream.XORKeyStream(b, b)
	}
	_, err := bs.WriteTo(w.writer)
	return err
}
