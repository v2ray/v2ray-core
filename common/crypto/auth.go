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
		if r.buffer.IsEmpty() {
			r.buffer.Clear()
		} else {
			common.Must(r.buffer.Reset(buf.ReadFrom(r.buffer)))
		}

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

func (r *AuthenticationReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b, err := r.readChunk(true)
	if err != nil {
		return nil, err
	}

	var mb buf.MultiBuffer
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

const (
	WriteSize = 1024
)

type AuthenticationWriter struct {
	auth         Authenticator
	writer       buf.Writer
	sizeParser   ChunkSizeEncoder
	transferType protocol.TransferType
}

func NewAuthenticationWriter(auth Authenticator, sizeParser ChunkSizeEncoder, writer io.Writer, transferType protocol.TransferType) *AuthenticationWriter {
	return &AuthenticationWriter{
		auth:         auth,
		writer:       buf.NewWriter(writer),
		sizeParser:   sizeParser,
		transferType: transferType,
	}
}

func (w *AuthenticationWriter) seal(b *buf.Buffer) (*buf.Buffer, error) {
	encryptedSize := b.Len() + w.auth.Overhead()

	eb := buf.New()
	common.Must(eb.Reset(func(bb []byte) (int, error) {
		w.sizeParser.Encode(uint16(encryptedSize), bb[:0])
		return w.sizeParser.SizeBytes(), nil
	}))
	if err := eb.AppendSupplier(func(bb []byte) (int, error) {
		_, err := w.auth.Seal(bb[:0], b.Bytes())
		return encryptedSize, err
	}); err != nil {
		eb.Release()
		return nil, err
	}

	return eb, nil
}

func (w *AuthenticationWriter) writeStream(mb buf.MultiBuffer) error {
	defer mb.Release()

	mb2Write := buf.NewMultiBufferCap(len(mb) + 10)

	for {
		b := buf.New()
		common.Must(b.Reset(func(bb []byte) (int, error) {
			return mb.Read(bb[:WriteSize])
		}))
		eb, err := w.seal(b)
		b.Release()

		if err != nil {
			mb2Write.Release()
			return err
		}
		mb2Write.Append(eb)
		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.WriteMultiBuffer(mb2Write)
}

func (w *AuthenticationWriter) writePacket(mb buf.MultiBuffer) error {
	defer mb.Release()

	mb2Write := buf.NewMultiBufferCap(len(mb) * 2)

	for {
		b := mb.SplitFirst()
		if b == nil {
			b = buf.New()
		}
		eb, err := w.seal(b)
		b.Release()
		if err != nil {
			mb2Write.Release()
			return err
		}
		mb2Write.Append(eb)
		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.WriteMultiBuffer(mb2Write)
}

func (w *AuthenticationWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if w.transferType == protocol.TransferTypeStream {
		return w.writeStream(mb)
	}

	return w.writePacket(mb)
}
