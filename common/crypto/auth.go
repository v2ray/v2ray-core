package crypto

import (
	"crypto/cipher"
	"errors"
	"io"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/serial"
)

var (
	ErrAuthenticationFailed = errors.New("Authentication failed.")
	errInsufficientBuffer   = errors.New("Insufficient buffer.")
)

type BytesGenerator func() []byte

type AuthenticationReader struct {
	aead     cipher.AEAD
	buffer   *alloc.Buffer
	reader   io.Reader
	ivGen    BytesGenerator
	extraGen BytesGenerator

	chunk []byte
}

func NewAuthenticationReader(aead cipher.AEAD, reader io.Reader, ivGen BytesGenerator, extraGen BytesGenerator) *AuthenticationReader {
	return &AuthenticationReader{
		aead:     aead,
		buffer:   alloc.NewLocalBuffer(32 * 1024),
		reader:   reader,
		ivGen:    ivGen,
		extraGen: extraGen,
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
	if size == v.aead.Overhead() {
		return io.EOF
	}
	cipherChunk := v.buffer.BytesRange(2, size+2)
	plainChunk, err := v.aead.Open(cipherChunk, v.ivGen(), cipherChunk, v.extraGen())
	if err != nil {
		return err
	}
	v.chunk = plainChunk
	v.buffer.SliceFrom(size + 2)
	return nil
}

func (v *AuthenticationReader) CopyChunk(b []byte) int {
	nBytes := copy(b, v.chunk)
	if nBytes == len(v.chunk) {
		v.chunk = nil
	} else {
		v.chunk = v.chunk[nBytes:]
	}
	return nBytes
}

func (v *AuthenticationReader) Read(b []byte) (int, error) {
	if len(v.chunk) > 0 {
		nBytes := v.CopyChunk(b)
		return nBytes, nil
	}

	err := v.NextChunk()
	if err == errInsufficientBuffer {
		_, err = v.buffer.FillFrom(v.reader)
	}

	if err != nil {
		return 0, err
	}

	totalBytes := 0
	for {
		totalBytes += v.CopyChunk(b)
		if len(b) == 0 {
			break
		}
		if err := v.NextChunk(); err != nil {
			break
		}
	}
	return totalBytes, nil
}

type AuthenticationWriter struct {
	aead     cipher.AEAD
	buffer   []byte
	writer   io.Writer
	ivGen    BytesGenerator
	extraGen BytesGenerator
}

func NewAuthenticationWriter(aead cipher.AEAD, writer io.Writer, ivGen BytesGenerator, extraGen BytesGenerator) *AuthenticationWriter {
	return &AuthenticationWriter{
		aead:     aead,
		buffer:   make([]byte, 32*1024),
		writer:   writer,
		ivGen:    ivGen,
		extraGen: extraGen,
	}
}

func (v *AuthenticationWriter) Write(b []byte) (int, error) {
	cipherChunk := v.aead.Seal(v.buffer[2:], v.ivGen(), b, v.extraGen())
	serial.Uint16ToBytes(uint16(len(cipherChunk)), b[:0])
	_, err := v.writer.Write(v.buffer[:2+len(cipherChunk)])
	return len(b), err
}
