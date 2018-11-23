package handshake

import (
	"crypto/cipher"
	"encoding/binary"

	"github.com/lucas-clemente/quic-go/internal/protocol"
)

type sealer struct {
	iv   []byte
	aead cipher.AEAD

	// use a single slice to avoid allocations
	nonceBuf []byte
}

var _ Sealer = &sealer{}

func newSealer(aead cipher.AEAD, iv []byte) Sealer {
	return &sealer{
		iv:       iv,
		aead:     aead,
		nonceBuf: make([]byte, aead.NonceSize()),
	}
}

func (s *sealer) Seal(dst, src []byte, pn protocol.PacketNumber, ad []byte) []byte {
	binary.BigEndian.PutUint64(s.nonceBuf[len(s.nonceBuf)-8:], uint64(pn))
	return s.aead.Seal(dst, s.nonceBuf, src, ad)
}

func (s *sealer) Overhead() int {
	return s.aead.Overhead()
}

type opener struct {
	iv   []byte
	aead cipher.AEAD

	// use a single slice to avoid allocations
	nonceBuf []byte
}

var _ Opener = &opener{}

func newOpener(aead cipher.AEAD, iv []byte) Opener {
	return &opener{
		iv:       iv,
		aead:     aead,
		nonceBuf: make([]byte, aead.NonceSize()),
	}
}

func (o *opener) Open(dst, src []byte, pn protocol.PacketNumber, ad []byte) ([]byte, error) {
	binary.BigEndian.PutUint64(o.nonceBuf[len(o.nonceBuf)-8:], uint64(pn))
	return o.aead.Open(dst, o.nonceBuf, src, ad)
}
