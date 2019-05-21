package crypto

import (
	"crypto/cipher"
	"io"
	"math/rand"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bytespool"
	"v2ray.com/core/common/protocol"
)

type BytesGenerator func() []byte

func GenerateEmptyBytes() BytesGenerator {
	var b [1]byte
	return func() []byte {
		return b[:0]
	}
}

func GenerateStaticBytes(content []byte) BytesGenerator {
	return func() []byte {
		return content
	}
}

func GenerateIncreasingNonce(nonce []byte) BytesGenerator {
	c := append([]byte(nil), nonce...)
	return func() []byte {
		for i := range c {
			c[i]++
			if c[i] != 0 {
				break
			}
		}
		return c
	}
}

func GenerateInitialAEADNonce() BytesGenerator {
	return GenerateIncreasingNonce([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
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
	iv := v.NonceGenerator()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	var additionalData []byte
	if v.AdditionalDataGenerator != nil {
		additionalData = v.AdditionalDataGenerator()
	}
	return v.AEAD.Open(dst, iv, cipherText, additionalData)
}

func (v *AEADAuthenticator) Seal(dst, plainText []byte) ([]byte, error) {
	iv := v.NonceGenerator()
	if len(iv) != v.AEAD.NonceSize() {
		return nil, newError("invalid AEAD nonce size: ", len(iv))
	}

	var additionalData []byte
	if v.AdditionalDataGenerator != nil {
		additionalData = v.AdditionalDataGenerator()
	}
	return v.AEAD.Seal(dst, iv, plainText, additionalData), nil
}

type AuthenticationReader struct {
	auth         Authenticator
	reader       *buf.BufferedReader
	sizeParser   ChunkSizeDecoder
	sizeBytes    []byte
	transferType protocol.TransferType
	padding      PaddingLengthGenerator
	size         uint16
	paddingLen   uint16
	hasSize      bool
	done         bool
}

func NewAuthenticationReader(auth Authenticator, sizeParser ChunkSizeDecoder, reader io.Reader, transferType protocol.TransferType, paddingLen PaddingLengthGenerator) *AuthenticationReader {
	r := &AuthenticationReader{
		auth:         auth,
		sizeParser:   sizeParser,
		transferType: transferType,
		padding:      paddingLen,
		sizeBytes:    make([]byte, sizeParser.SizeBytes()),
	}
	if breader, ok := reader.(*buf.BufferedReader); ok {
		r.reader = breader
	} else {
		r.reader = &buf.BufferedReader{Reader: buf.NewReader(reader)}
	}
	return r
}

func (r *AuthenticationReader) readSize() (uint16, uint16, error) {
	if r.hasSize {
		r.hasSize = false
		return r.size, r.paddingLen, nil
	}
	if _, err := io.ReadFull(r.reader, r.sizeBytes); err != nil {
		return 0, 0, err
	}
	var padding uint16
	if r.padding != nil {
		padding = r.padding.NextPaddingLen()
	}
	size, err := r.sizeParser.Decode(r.sizeBytes)
	return size, padding, err
}

var errSoft = newError("waiting for more data")

func (r *AuthenticationReader) readBuffer(size int32, padding int32) (*buf.Buffer, error) {
	b := buf.New()
	if _, err := b.ReadFullFrom(r.reader, size); err != nil {
		b.Release()
		return nil, err
	}
	size -= padding
	rb, err := r.auth.Open(b.BytesTo(0), b.BytesTo(size))
	if err != nil {
		b.Release()
		return nil, err
	}
	b.Resize(0, int32(len(rb)))
	return b, nil
}

func (r *AuthenticationReader) readInternal(soft bool, mb *buf.MultiBuffer) error {
	if soft && r.reader.BufferedBytes() < r.sizeParser.SizeBytes() {
		return errSoft
	}

	if r.done {
		return io.EOF
	}

	size, padding, err := r.readSize()
	if err != nil {
		return err
	}

	if size == uint16(r.auth.Overhead())+padding {
		r.done = true
		return io.EOF
	}

	if soft && int32(size) > r.reader.BufferedBytes() {
		r.size = size
		r.paddingLen = padding
		r.hasSize = true
		return errSoft
	}

	if size <= buf.Size {
		b, err := r.readBuffer(int32(size), int32(padding))
		if err != nil {
			return nil
		}
		*mb = append(*mb, b)
		return nil
	}

	payload := bytespool.Alloc(int32(size))
	defer bytespool.Free(payload)

	if _, err := io.ReadFull(r.reader, payload[:size]); err != nil {
		return err
	}

	size -= padding

	rb, err := r.auth.Open(payload[:0], payload[:size])
	if err != nil {
		return err
	}

	*mb = buf.MergeBytes(*mb, rb)
	return nil
}

func (r *AuthenticationReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	const readSize = 16
	mb := make(buf.MultiBuffer, 0, readSize)
	if err := r.readInternal(false, &mb); err != nil {
		buf.ReleaseMulti(mb)
		return nil, err
	}

	for i := 1; i < readSize; i++ {
		err := r.readInternal(true, &mb)
		if err == errSoft || err == io.EOF {
			break
		}
		if err != nil {
			buf.ReleaseMulti(mb)
			return nil, err
		}
	}

	return mb, nil
}

type AuthenticationWriter struct {
	auth         Authenticator
	writer       buf.Writer
	sizeParser   ChunkSizeEncoder
	transferType protocol.TransferType
	padding      PaddingLengthGenerator
}

func NewAuthenticationWriter(auth Authenticator, sizeParser ChunkSizeEncoder, writer io.Writer, transferType protocol.TransferType, padding PaddingLengthGenerator) *AuthenticationWriter {
	w := &AuthenticationWriter{
		auth:         auth,
		writer:       buf.NewWriter(writer),
		sizeParser:   sizeParser,
		transferType: transferType,
	}
	if padding != nil {
		w.padding = padding
	}
	return w
}

func (w *AuthenticationWriter) seal(b []byte) (*buf.Buffer, error) {
	encryptedSize := int32(len(b) + w.auth.Overhead())
	var paddingSize int32
	if w.padding != nil {
		paddingSize = int32(w.padding.NextPaddingLen())
	}

	totalSize := encryptedSize + paddingSize
	if totalSize > buf.Size {
		return nil, newError("size too large: ", totalSize)
	}

	eb := buf.New()
	w.sizeParser.Encode(uint16(encryptedSize+paddingSize), eb.Extend(w.sizeParser.SizeBytes()))
	if _, err := w.auth.Seal(eb.Extend(encryptedSize)[:0], b); err != nil {
		eb.Release()
		return nil, err
	}
	if paddingSize > 0 {
		// With size of the chunk and padding length encrypted, the content of padding doesn't matter much.
		paddingBytes := eb.Extend(paddingSize)
		common.Must2(rand.Read(paddingBytes))
	}

	return eb, nil
}

func (w *AuthenticationWriter) writeStream(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)

	var maxPadding int32
	if w.padding != nil {
		maxPadding = int32(w.padding.MaxPaddingLen())
	}

	payloadSize := buf.Size - int32(w.auth.Overhead()) - w.sizeParser.SizeBytes() - maxPadding
	mb2Write := make(buf.MultiBuffer, 0, len(mb)+10)

	temp := buf.New()
	defer temp.Release()

	rawBytes := temp.Extend(payloadSize)

	for {
		nb, nBytes := buf.SplitBytes(mb, rawBytes)
		mb = nb

		eb, err := w.seal(rawBytes[:nBytes])

		if err != nil {
			buf.ReleaseMulti(mb2Write)
			return err
		}
		mb2Write = append(mb2Write, eb)
		if mb.IsEmpty() {
			break
		}
	}

	return w.writer.WriteMultiBuffer(mb2Write)
}

func (w *AuthenticationWriter) writePacket(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)

	mb2Write := make(buf.MultiBuffer, 0, len(mb)+1)

	for _, b := range mb {
		if b.IsEmpty() {
			continue
		}

		eb, err := w.seal(b.Bytes())
		if err != nil {
			continue
		}

		mb2Write = append(mb2Write, eb)
	}

	if mb2Write.IsEmpty() {
		return nil
	}

	return w.writer.WriteMultiBuffer(mb2Write)
}

// WriteMultiBuffer implements buf.Writer.
func (w *AuthenticationWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if mb.IsEmpty() {
		eb, err := w.seal([]byte{})
		common.Must(err)
		return w.writer.WriteMultiBuffer(buf.MultiBuffer{eb})
	}

	if w.transferType == protocol.TransferTypeStream {
		return w.writeStream(mb)
	}

	return w.writePacket(mb)
}
