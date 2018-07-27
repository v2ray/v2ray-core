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

// WriteMultiBuffer implements buf.Writer.
func (w *CryptionWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer mb.Release()

	size := mb.Len()
	if size == 0 {
		return nil
	}

	bs := mb.ToNetBuffers()
	for _, b := range bs {
		w.stream.XORKeyStream(b, b)
	}

	for size > 0 {
		n, err := bs.WriteTo(w.writer)
		if err != nil {
			return err
		}
		size -= int32(n)
	}
	return nil
}
