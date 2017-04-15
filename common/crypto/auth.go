package crypto

import (
	"crypto/cipher"
	"io"

	"golang.org/x/crypto/sha3"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

var (
	errInsufficientBuffer = newError("insufficient buffer")
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

type Uint16Generator interface {
	Next() uint16
}

type StaticUint16Generator uint16

func (g StaticUint16Generator) Next() uint16 {
	return uint16(g)
}

type ShakeUint16Generator struct {
	shake  sha3.ShakeHash
	buffer [2]byte
}

func NewShakeUint16Generator(nonce []byte) *ShakeUint16Generator {
	shake := sha3.NewShake128()
	shake.Write(nonce)
	return &ShakeUint16Generator{
		shake: shake,
	}
}

func (g *ShakeUint16Generator) Next() uint16 {
	g.shake.Read(g.buffer[:])
	return serial.BytesToUint16(g.buffer[:])
}

type AuthenticationReader struct {
	auth     Authenticator
	buffer   *buf.Buffer
	reader   io.Reader
	sizeMask Uint16Generator

	chunk []byte
}

const (
	readerBufferSize = 32 * 1024
)

func NewAuthenticationReader(auth Authenticator, reader io.Reader, sizeMask Uint16Generator) *AuthenticationReader {
	return &AuthenticationReader{
		auth:     auth,
		buffer:   buf.NewLocal(readerBufferSize),
		reader:   reader,
		sizeMask: sizeMask,
	}
}

func (v *AuthenticationReader) nextChunk(mask uint16) error {
	if v.buffer.Len() < 2 {
		return errInsufficientBuffer
	}
	size := int(serial.BytesToUint16(v.buffer.BytesTo(2)) ^ mask)
	if size > v.buffer.Len()-2 {
		return errInsufficientBuffer
	}
	if size > readerBufferSize-2 {
		return newError("size too large: ", size)
	}
	if size == v.auth.Overhead() {
		return io.EOF
	}
	if size < v.auth.Overhead() {
		return newError("invalid packet size: ", size)
	}
	cipherChunk := v.buffer.BytesRange(2, size+2)
	plainChunk, err := v.auth.Open(cipherChunk[:0], cipherChunk)
	if err != nil {
		return err
	}
	v.chunk = plainChunk
	v.buffer.SliceFrom(size + 2)
	return nil
}

func (v *AuthenticationReader) copyChunk(b []byte) int {
	if len(v.chunk) == 0 {
		return 0
	}
	nBytes := copy(b, v.chunk)
	if nBytes == len(v.chunk) {
		v.chunk = nil
	} else {
		v.chunk = v.chunk[nBytes:]
	}
	return nBytes
}

func (v *AuthenticationReader) ensureChunk() error {
	atHead := false
	if v.buffer.IsEmpty() {
		v.buffer.Clear()
		atHead = true
	}

	mask := v.sizeMask.Next()
	for {
		err := v.nextChunk(mask)
		if err != errInsufficientBuffer {
			return err
		}

		leftover := v.buffer.Bytes()
		if !atHead && len(leftover) > 0 {
			common.Must(v.buffer.Reset(func(b []byte) (int, error) {
				return copy(b, leftover), nil
			}))
		}

		if err := v.buffer.AppendSupplier(buf.ReadFrom(v.reader)); err != nil {
			return err
		}
	}
}

func (v *AuthenticationReader) Read(b []byte) (int, error) {
	if len(v.chunk) > 0 {
		nBytes := v.copyChunk(b)
		return nBytes, nil
	}

	err := v.ensureChunk()
	if err != nil {
		return 0, err
	}

	return v.copyChunk(b), nil
}

func (r *AuthenticationReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	err := r.ensureChunk()
	if err != nil {
		return nil, err
	}

	mb := buf.NewMultiBuffer()
	for len(r.chunk) > 0 {
		b := buf.New()
		nBytes, _ := b.Write(r.chunk)
		mb.Append(b)
		r.chunk = r.chunk[nBytes:]
	}
	r.chunk = nil
	return mb, nil
}

type AuthenticationWriter struct {
	auth     Authenticator
	buffer   []byte
	writer   io.Writer
	sizeMask Uint16Generator
}

func NewAuthenticationWriter(auth Authenticator, writer io.Writer, sizeMask Uint16Generator) *AuthenticationWriter {
	return &AuthenticationWriter{
		auth:     auth,
		buffer:   make([]byte, 32*1024),
		writer:   writer,
		sizeMask: sizeMask,
	}
}

func (w *AuthenticationWriter) Write(b []byte) (int, error) {
	cipherChunk, err := w.auth.Seal(w.buffer[2:2], b)
	if err != nil {
		return 0, err
	}

	size := uint16(len(cipherChunk)) ^ w.sizeMask.Next()
	serial.Uint16ToBytes(size, w.buffer[:0])
	_, err = w.writer.Write(w.buffer[:2+len(cipherChunk)])
	return len(b), err
}

func (w *AuthenticationWriter) WriteMultiBuffer(mb buf.MultiBuffer) (int, error) {
	const StartIndex = 17 * 1024
	var totalBytes int
	for {
		payloadLen, err := mb.Read(w.buffer[StartIndex:])
		if err != nil {
			return 0, err
		}
		nBytes, err := w.Write(w.buffer[StartIndex : StartIndex+payloadLen])
		totalBytes += nBytes
		if err != nil {
			return totalBytes, err
		}
		if mb.IsEmpty() {
			break
		}
	}
	return totalBytes, nil
}
