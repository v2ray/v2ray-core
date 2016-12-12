package crypto

import (
	"crypto/cipher"
	"errors"
	"io"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

var (
	ErrAuthenticationFailed = errors.New("Authentication failed.")

	errInsufficientBuffer = errors.New("Insufficient buffer.")
	errInvalidNonce       = errors.New("Invalid nonce.")
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
		return nil, errInvalidNonce
	}

	additionalData := v.AdditionalDataGenerator.Next()
	return v.AEAD.Open(dst, iv, cipherText, additionalData)
}

func (v *AEADAuthenticator) Seal(dst, plainText []byte) ([]byte, error) {
	iv := v.NonceGenerator.Next()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, errInvalidNonce
	}

	additionalData := v.AdditionalDataGenerator.Next()
	return v.AEAD.Seal(dst, iv, plainText, additionalData), nil
}

type AuthenticationReader struct {
	auth   Authenticator
	buffer *buf.Buffer
	reader io.Reader

	chunk      []byte
	aggressive bool
}

func NewAuthenticationReader(auth Authenticator, reader io.Reader, aggressive bool) *AuthenticationReader {
	return &AuthenticationReader{
		auth:       auth,
		buffer:     buf.NewLocal(32 * 1024),
		reader:     reader,
		aggressive: aggressive,
	}
}

func (v *AuthenticationReader) NextChunk() error {
	if v.buffer.Len() < 2 {
		return errInsufficientBuffer
	}
	size := int(serial.BytesToUint16(v.buffer.BytesTo(2)))
	if size > v.buffer.Len()-2 {
		return errInsufficientBuffer
	}
	if size == v.auth.Overhead() {
		return io.EOF
	}
	if size < v.auth.Overhead() {
		return errors.New("AuthenticationReader: invalid packet size.")
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

func (v *AuthenticationReader) CopyChunk(b []byte) int {
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

func (v *AuthenticationReader) EnsureChunk() error {
	for {
		err := v.NextChunk()
		if err == nil {
			return nil
		}
		if err == errInsufficientBuffer {
			if v.buffer.IsEmpty() {
				v.buffer.Clear()
			} else {
				leftover := v.buffer.Bytes()
				v.buffer.Reset(func(b []byte) (int, error) {
					return copy(b, leftover), nil
				})
			}
			if err := v.buffer.AppendSupplier(buf.ReadFrom(v.reader)); err == nil {
				continue
			}
		}
		return err
	}
}

func (v *AuthenticationReader) Read(b []byte) (int, error) {
	if len(v.chunk) > 0 {
		nBytes := v.CopyChunk(b)
		return nBytes, nil
	}

	err := v.EnsureChunk()
	if err != nil {
		return 0, err
	}

	nBytes := v.CopyChunk(b)
	return nBytes, nil
}

type AuthenticationWriter struct {
	auth     Authenticator
	buffer   []byte
	writer   io.Writer
	ivGen    BytesGenerator
	extraGen BytesGenerator
}

func NewAuthenticationWriter(auth Authenticator, writer io.Writer) *AuthenticationWriter {
	return &AuthenticationWriter{
		auth:   auth,
		buffer: make([]byte, 32*1024),
		writer: writer,
	}
}

func (v *AuthenticationWriter) Write(b []byte) (int, error) {
	cipherChunk, err := v.auth.Seal(v.buffer[2:2], b)
	if err != nil {
		return 0, err
	}

	serial.Uint16ToBytes(uint16(len(cipherChunk)), v.buffer[:0])
	_, err = v.writer.Write(v.buffer[:2+len(cipherChunk)])
	return len(b), err
}
