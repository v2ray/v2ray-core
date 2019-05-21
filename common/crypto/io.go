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
	_ buf.Writer = (*CryptionWriter)(nil)
)

type CryptionWriter struct {
	stream    cipher.Stream
	writer    io.Writer
	bufWriter buf.Writer
}

// NewCryptionWriter creates a new CryptionWriter.
func NewCryptionWriter(stream cipher.Stream, writer io.Writer) *CryptionWriter {
	return &CryptionWriter{
		stream:    stream,
		writer:    writer,
		bufWriter: buf.NewWriter(writer),
	}
}

// Write implements io.Writer.Write().
func (w *CryptionWriter) Write(data []byte) (int, error) {
	w.stream.XORKeyStream(data, data)

	if err := buf.WriteAllBytes(w.writer, data); err != nil {
		return 0, err
	}
	return len(data), nil
}

// WriteMultiBuffer implements buf.Writer.
func (w *CryptionWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	for _, b := range mb {
		w.stream.XORKeyStream(b.Bytes(), b.Bytes())
	}

	return w.bufWriter.WriteMultiBuffer(mb)
}
