package handshake

import (
	"crypto/cipher"
	"encoding/binary"

	"github.com/lucas-clemente/quic-go/internal/protocol"
)

type sealer struct {
	iv          []byte
	aead        cipher.AEAD
	pnEncrypter cipher.Block

	// use a single slice to avoid allocations
	nonceBuf []byte
	pnMask   []byte

	// short headers protect 5 bits in the first byte, long headers only 4
	is1RTT bool
}

var _ Sealer = &sealer{}

func newSealer(aead cipher.AEAD, iv []byte, pnEncrypter cipher.Block, is1RTT bool) Sealer {
	return &sealer{
		iv:          iv,
		aead:        aead,
		nonceBuf:    make([]byte, aead.NonceSize()),
		is1RTT:      is1RTT,
		pnEncrypter: pnEncrypter,
		pnMask:      make([]byte, pnEncrypter.BlockSize()),
	}
}

func (s *sealer) Seal(dst, src []byte, pn protocol.PacketNumber, ad []byte) []byte {
	binary.BigEndian.PutUint64(s.nonceBuf[len(s.nonceBuf)-8:], uint64(pn))
	for i := 0; i < len(s.nonceBuf); i++ {
		s.nonceBuf[i] ^= s.iv[i]
	}
	sealed := s.aead.Seal(dst, s.nonceBuf, src, ad)
	for i := 0; i < len(s.nonceBuf); i++ {
		s.nonceBuf[i] = 0
	}
	return sealed
}

func (s *sealer) EncryptHeader(sample []byte, firstByte *byte, pnBytes []byte) {
	if len(sample) != s.pnEncrypter.BlockSize() {
		panic("invalid sample size")
	}
	s.pnEncrypter.Encrypt(s.pnMask, sample)
	if s.is1RTT {
		*firstByte ^= s.pnMask[0] & 0x1f
	} else {
		*firstByte ^= s.pnMask[0] & 0xf
	}
	for i := range pnBytes {
		pnBytes[i] ^= s.pnMask[i+1]
	}
}

func (s *sealer) Overhead() int {
	return s.aead.Overhead()
}

type opener struct {
	iv          []byte
	aead        cipher.AEAD
	pnDecrypter cipher.Block

	// use a single slice to avoid allocations
	nonceBuf []byte
	pnMask   []byte

	// short headers protect 5 bits in the first byte, long headers only 4
	is1RTT bool
}

var _ Opener = &opener{}

func newOpener(aead cipher.AEAD, iv []byte, pnDecrypter cipher.Block, is1RTT bool) Opener {
	return &opener{
		iv:          iv,
		aead:        aead,
		nonceBuf:    make([]byte, aead.NonceSize()),
		is1RTT:      is1RTT,
		pnDecrypter: pnDecrypter,
		pnMask:      make([]byte, pnDecrypter.BlockSize()),
	}
}

func (o *opener) Open(dst, src []byte, pn protocol.PacketNumber, ad []byte) ([]byte, error) {
	binary.BigEndian.PutUint64(o.nonceBuf[len(o.nonceBuf)-8:], uint64(pn))
	for i := 0; i < len(o.nonceBuf); i++ {
		o.nonceBuf[i] ^= o.iv[i]
	}
	opened, err := o.aead.Open(dst, o.nonceBuf, src, ad)
	for i := 0; i < len(o.nonceBuf); i++ {
		o.nonceBuf[i] = 0
	}
	return opened, err
}

func (o *opener) DecryptHeader(sample []byte, firstByte *byte, pnBytes []byte) {
	if len(sample) != o.pnDecrypter.BlockSize() {
		panic("invalid sample size")
	}
	o.pnDecrypter.Encrypt(o.pnMask, sample)
	if o.is1RTT {
		*firstByte ^= o.pnMask[0] & 0x1f
	} else {
		*firstByte ^= o.pnMask[0] & 0xf
	}
	for i := range pnBytes {
		pnBytes[i] ^= o.pnMask[i+1]
	}
}
