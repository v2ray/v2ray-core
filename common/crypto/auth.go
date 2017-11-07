package crypto

import (
	"crypto/cipher"
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/protocol"
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
	auth         Authenticator
	buffer       *buf.Buffer
	reader       io.Reader
	sizeParser   ChunkSizeDecoder
	size         int
	transferType protocol.TransferType
}

const (
	readerBufferSize = 32 * 1024
)

func NewAuthenticationReader(auth Authenticator, sizeParser ChunkSizeDecoder, reader io.Reader, transferType protocol.TransferType) *AuthenticationReader {
	return &AuthenticationReader{
		auth:         auth,
		buffer:       buf.NewLocal(readerBufferSize),
		reader:       reader,
		sizeParser:   sizeParser,
		size:         -1,
		transferType: transferType,
	}
}

func (r *AuthenticationReader) readSize() error {
	if r.size >= 0 {
		return nil
	}

	sizeBytes := r.sizeParser.SizeBytes()
	if r.buffer.Len() < sizeBytes {
		common.Must(r.buffer.Reset(buf.ReadFrom(r.buffer)))
		delta := sizeBytes - r.buffer.Len()
		if err := r.buffer.AppendSupplier(buf.ReadAtLeastFrom(r.reader, delta)); err != nil {
			return err
		}
	}
	size, err := r.sizeParser.Decode(r.buffer.BytesTo(sizeBytes))
	if err != nil {
		return err
	}
	r.size = int(size)
	r.buffer.SliceFrom(sizeBytes)
	return nil
}

func (r *AuthenticationReader) readChunk(waitForData bool) ([]byte, error) {
	if err := r.readSize(); err != nil {
		return nil, err
	}
	if r.size > readerBufferSize-r.sizeParser.SizeBytes() {
		return nil, newError("size too large ", r.size).AtWarning()
	}

	if r.size == r.auth.Overhead() {
		return nil, io.EOF
	}

	if r.buffer.Len() < r.size {
		if !waitForData {
			return nil, io.ErrNoProgress
		}

		if r.buffer.IsEmpty() {
			r.buffer.Clear()
		} else {
			common.Must(r.buffer.Reset(buf.ReadFrom(r.buffer)))
		}

		delta := r.size - r.buffer.Len()
		if err := r.buffer.AppendSupplier(buf.ReadAtLeastFrom(r.reader, delta)); err != nil {
			return nil, err
		}
	}

	b, err := r.auth.Open(r.buffer.BytesTo(0), r.buffer.BytesTo(r.size))
	if err != nil {
		return nil, err
	}
	r.buffer.SliceFrom(r.size)
	r.size = -1
	return b, nil
}

func (r *AuthenticationReader) Read() (buf.MultiBuffer, error) {
	b, err := r.readChunk(true)
	if err != nil {
		return nil, err
	}

	mb := buf.NewMultiBuffer()
	if r.transferType == protocol.TransferTypeStream {
		mb.Write(b)
	} else {
		var bb *buf.Buffer
		if len(b) < buf.Size {
			bb = buf.New()
		} else {
			bb = buf.NewLocal(len(b))
		}
		bb.Append(b)
		mb.Append(bb)
	}

	for r.buffer.Len() >= r.sizeParser.SizeBytes() {
		b, err := r.readChunk(false)
		if err != nil {
			break
		}
		if r.transferType == protocol.TransferTypeStream {
			mb.Write(b)
		} else {
			var bb *buf.Buffer
			if len(b) < buf.Size {
				bb = buf.New()
			} else {
				bb = buf.NewLocal(len(b))
			}
			bb.Append(b)
			mb.Append(bb)
		}
	}

	return mb, nil
}

type AuthenticationWriter struct {
	auth         Authenticator
	buffer       []byte
	payload      []byte
	writer       *buf.BufferedWriter
	sizeParser   ChunkSizeEncoder
	transferType protocol.TransferType
}

func NewAuthenticationWriter(auth Authenticator, sizeParser ChunkSizeEncoder, writer io.Writer, transferType protocol.TransferType) *AuthenticationWriter {
	const payloadSize = 1024
	return &AuthenticationWriter{
		auth:         auth,
		buffer:       make([]byte, payloadSize+sizeParser.SizeBytes()+auth.Overhead()),
		payload:      make([]byte, payloadSize),
		writer:       buf.NewBufferedWriterSize(writer, readerBufferSize),
		sizeParser:   sizeParser,
		transferType: transferType,
	}
}

func (w *AuthenticationWriter) append(b []byte) error {
	encryptedSize := len(b) + w.auth.Overhead()
	buffer := w.sizeParser.Encode(uint16(encryptedSize), w.buffer[:0])

	buffer, err := w.auth.Seal(buffer, b)
	if err != nil {
		return err
	}

	if _, err := w.writer.Write(buffer); err != nil {
		return err
	}

	return nil
}

func (w *AuthenticationWriter) writeStream(mb buf.MultiBuffer) error {
	defer mb.Release()

	for {
		n, _ := mb.Read(w.payload)
		if err := w.append(w.payload[:n]); err != nil {
			return err
		}
		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.Flush()
}

func (w *AuthenticationWriter) writePacket(mb buf.MultiBuffer) error {
	defer mb.Release()

	for {
		b := mb.SplitFirst()
		if b == nil {
			b = buf.New()
		}
		if err := w.append(b.Bytes()); err != nil {
			b.Release()
			return err
		}
		b.Release()
		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.Flush()
}

func (w *AuthenticationWriter) Write(mb buf.MultiBuffer) error {
	if w.transferType == protocol.TransferTypeStream {
		return w.writeStream(mb)
	}

	return w.writePacket(mb)
}
