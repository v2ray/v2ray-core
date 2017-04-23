package crypto

import (
	"crypto/cipher"
	"io"

	"v2ray.com/core/common/buf"
)

type BytesGenerator interface {
	Next() []byte
}

type NoOpBytesGenerator struct {
	buffer [1]byte
}

func (v NoOpBytesGenerator) Next() []byte {
	return v.buffer[:0]
}

type StaticBytesGenerator struct {
	Content []byte
}

func (v StaticBytesGenerator) Next() []byte {
	return v.Content
}

type Authenticator interface {
	NonceSize() int
	Overhead() int
	Open(dst, cipherText []byte) ([]byte, error)
	Seal(dst, plainText []byte) ([]byte, error)
}

type AEADAuthenticator struct {
	cipher.AEAD
	NonceGenerator          BytesGenerator
	AdditionalDataGenerator BytesGenerator
}

func (v *AEADAuthenticator) Open(dst, cipherText []byte) ([]byte, error) {
	iv := v.NonceGenerator.Next()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	additionalData := v.AdditionalDataGenerator.Next()
	return v.AEAD.Open(dst, iv, cipherText, additionalData)
}

func (v *AEADAuthenticator) Seal(dst, plainText []byte) ([]byte, error) {
	iv := v.NonceGenerator.Next()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	additionalData := v.AdditionalDataGenerator.Next()
	return v.AEAD.Seal(dst, iv, plainText, additionalData), nil
}

type AuthenticationReader struct {
	auth       Authenticator
	buffer     *buf.Buffer
	reader     io.Reader
	sizeParser ChunkSizeDecoder
}

const (
	readerBufferSize = 32 * 1024
)

func NewAuthenticationReader(auth Authenticator, sizeParser ChunkSizeDecoder, reader io.Reader) *AuthenticationReader {
	return &AuthenticationReader{
		auth:       auth,
		buffer:     buf.NewLocal(readerBufferSize),
		reader:     reader,
		sizeParser: sizeParser,
	}
}

func (r *AuthenticationReader) readChunk() error {
	if err := r.buffer.Reset(buf.ReadFullFrom(r.reader, r.sizeParser.SizeBytes())); err != nil {
		return err
	}
	size, err := r.sizeParser.Decode(r.buffer.Bytes())
	if err != nil {
		return err
	}
	if size > readerBufferSize {
		return newError("size too large ", size).AtWarning()
	}

	if int(size) == r.auth.Overhead() {
		return io.EOF
	}

	if err := r.buffer.Reset(buf.ReadFullFrom(r.reader, int(size))); err != nil {
		return err
	}

	b, err := r.auth.Open(r.buffer.BytesTo(0), r.buffer.Bytes())
	if err != nil {
		return err
	}
	r.buffer.Slice(0, len(b))
	return nil
}

func (r *AuthenticationReader) Read() (buf.MultiBuffer, error) {
	if r.buffer.IsEmpty() {
		if err := r.readChunk(); err != nil {
			return nil, err
		}
	}

	mb := buf.NewMultiBuffer()
	for !r.buffer.IsEmpty() {
		b := buf.New()
		b.AppendSupplier(buf.ReadFrom(r.buffer))
		mb.Append(b)
	}
	return mb, nil
}

type AuthenticationWriter struct {
	auth       Authenticator
	buffer     []byte
	writer     io.Writer
	sizeParser ChunkSizeEncoder
}

func NewAuthenticationWriter(auth Authenticator, sizeParser ChunkSizeEncoder, writer io.Writer) *AuthenticationWriter {
	return &AuthenticationWriter{
		auth:       auth,
		buffer:     make([]byte, 32*1024),
		writer:     writer,
		sizeParser: sizeParser,
	}
}

func (w *AuthenticationWriter) writeInternal(b []byte) error {
	sizeBytes := w.sizeParser.SizeBytes()
	cipherChunk, err := w.auth.Seal(w.buffer[sizeBytes:sizeBytes], b)
	if err != nil {
		return err
	}

	w.sizeParser.Encode(uint16(len(cipherChunk)), w.buffer[:0])
	_, err = w.writer.Write(w.buffer[:sizeBytes+len(cipherChunk)])
	return err
}

func (w *AuthenticationWriter) Write(mb buf.MultiBuffer) error {
	defer mb.Release()

	const StartIndex = 17 * 1024
	for {
		payloadLen, _ := mb.Read(w.buffer[StartIndex:])
		err := w.writeInternal(w.buffer[StartIndex : StartIndex+payloadLen])
		if err != nil {
			return err
		}
		if mb.IsEmpty() {
			break
		}
	}
	return nil
}
