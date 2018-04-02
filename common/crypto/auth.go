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

type IncreasingAEADNonceGenerator struct {
	nonce []byte
}

func NewIncreasingAEADNonceGenerator() *IncreasingAEADNonceGenerator {
	return &IncreasingAEADNonceGenerator{
		nonce: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	}
}

func (g *IncreasingAEADNonceGenerator) Next() []byte {
	for i := range g.nonce {
		g.nonce[i]++
		if g.nonce[i] != 0 {
			break
		}
	}
	return g.nonce
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

	var additionalData []byte
	if v.AdditionalDataGenerator != nil {
		additionalData = v.AdditionalDataGenerator.Next()
	}
	return v.AEAD.Open(dst, iv, cipherText, additionalData)
}

func (v *AEADAuthenticator) Seal(dst, plainText []byte) ([]byte, error) {
	iv := v.NonceGenerator.Next()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	var additionalData []byte
	if v.AdditionalDataGenerator != nil {
		additionalData = v.AdditionalDataGenerator.Next()
	}
	return v.AEAD.Seal(dst, iv, plainText, additionalData), nil
}

type AuthenticationReader struct {
	auth         Authenticator
	reader       *buf.BufferedReader
	sizeParser   ChunkSizeDecoder
	transferType protocol.TransferType
	size         int32
}

func NewAuthenticationReader(auth Authenticator, sizeParser ChunkSizeDecoder, reader io.Reader, transferType protocol.TransferType) *AuthenticationReader {
	return &AuthenticationReader{
		auth:         auth,
		reader:       buf.NewBufferedReader(buf.NewReader(reader)),
		sizeParser:   sizeParser,
		transferType: transferType,
		size:         -1,
	}
}

func (r *AuthenticationReader) readSize() (int32, error) {
	if r.size != -1 {
		s := r.size
		r.size = -1
		return s, nil
	}
	sizeBytes := make([]byte, r.sizeParser.SizeBytes())
	_, err := io.ReadFull(r.reader, sizeBytes)
	if err != nil {
		return 0, err
	}
	size, err := r.sizeParser.Decode(sizeBytes)
	return int32(size), err
}

var errSoft = newError("waiting for more data")

func (r *AuthenticationReader) readInternal(soft bool) (*buf.Buffer, error) {
	if soft && r.reader.BufferedBytes() < int32(r.sizeParser.SizeBytes()) {
		return nil, errSoft
	}

	size, err := r.readSize()
	if err != nil {
		return nil, err
	}

	if size == -2 || size == int32(r.auth.Overhead()) {
		r.size = -2
		return nil, io.EOF
	}

	if soft && size > r.reader.BufferedBytes() {
		r.size = size
		return nil, errSoft
	}

	b := buf.NewSize(size)
	if err := b.Reset(buf.ReadFullFrom(r.reader, size)); err != nil {
		b.Release()
		return nil, err
	}

	rb, err := r.auth.Open(b.BytesTo(0), b.BytesTo(size))
	if err != nil {
		b.Release()
		return nil, err
	}
	b.Slice(0, int32(len(rb)))

	return b, nil
}

func (r *AuthenticationReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b, err := r.readInternal(false)
	if err != nil {
		return nil, err
	}

	mb := buf.NewMultiBufferCap(32)
	mb.Append(b)

	for {
		b, err := r.readInternal(true)
		if err == errSoft || err == io.EOF {
			break
		}
		if err != nil {
			mb.Release()
			return nil, err
		}
		mb.Append(b)
	}

	return mb, nil
}

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
	encryptedSize := int(b.Len()) + w.auth.Overhead()

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

	payloadSize := buf.Size - w.auth.Overhead() - w.sizeParser.SizeBytes()
	mb2Write := buf.NewMultiBufferCap(int32(len(mb) + 10))

	for {
		b := buf.New()
		common.Must(b.Reset(func(bb []byte) (int, error) {
			return mb.Read(bb[:payloadSize])
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

	if mb.IsEmpty() {
		b := buf.New()
		defer b.Release()

		eb, _ := w.seal(b)
		return w.writer.WriteMultiBuffer(buf.NewMultiBufferValue(eb))
	}

	mb2Write := buf.NewMultiBufferCap(int32(len(mb)) + 1)

	for !mb.IsEmpty() {
		b := mb.SplitFirst()
		if b == nil {
			continue
		}
		eb, err := w.seal(b)
		b.Release()
		if err != nil {
			mb2Write.Release()
			return err
		}
		mb2Write.Append(eb)
	}

	return w.writer.WriteMultiBuffer(mb2Write)
}

func (w *AuthenticationWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if w.transferType == protocol.TransferTypeStream {
		return w.writeStream(mb)
	}

	return w.writePacket(mb)
}
