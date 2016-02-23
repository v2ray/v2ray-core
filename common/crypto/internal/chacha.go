package internal

//go:generate go run chacha_core_gen.go

import (
	"encoding/binary"
)

const (
	wordSize  = 4                    // the size of ChaCha20's words
	stateSize = 16                   // the size of ChaCha20's state, in words
	blockSize = stateSize * wordSize // the size of ChaCha20's block, in bytes
)

type ChaCha20Stream struct {
	state  [stateSize]uint32 // the state as an array of 16 32-bit words
	block  [blockSize]byte   // the keystream as an array of 64 bytes
	offset int               // the offset of used bytes in block
	rounds int
}

func NewChaCha20Stream(key []byte, nonce []byte, rounds int) *ChaCha20Stream {
	s := new(ChaCha20Stream)
	// the magic constants for 256-bit keys
	s.state[0] = 0x61707865
	s.state[1] = 0x3320646e
	s.state[2] = 0x79622d32
	s.state[3] = 0x6b206574

	for i := 0; i < 8; i++ {
		s.state[i+4] = binary.LittleEndian.Uint32(key[i*4 : i*4+4])
	}

	switch len(nonce) {
	case 8:
		s.state[14] = binary.LittleEndian.Uint32(nonce[0:])
		s.state[15] = binary.LittleEndian.Uint32(nonce[4:])
	case 12:
		s.state[13] = binary.LittleEndian.Uint32(nonce[0:4])
		s.state[14] = binary.LittleEndian.Uint32(nonce[4:8])
		s.state[15] = binary.LittleEndian.Uint32(nonce[8:12])
	default:
		panic("bad nonce length")
	}

	s.rounds = rounds
	s.advance()
	return s
}

func (s *ChaCha20Stream) XORKeyStream(dst, src []byte) {
	// Stride over the input in 64-byte blocks, minus the amount of keystream
	// previously used. This will produce best results when processing blocks
	// of a size evenly divisible by 64.
	i := 0
	max := len(src)
	for i < max {
		gap := blockSize - s.offset

		limit := i + gap
		if limit > max {
			limit = max
		}

		o := s.offset
		for j := i; j < limit; j++ {
			dst[j] = src[j] ^ s.block[o]
			o++
		}

		i += gap
		s.offset = o

		if o == blockSize {
			s.advance()
		}
	}
}

func (s *ChaCha20Stream) advance() {
	ChaCha20Block(&s.state, s.block[:], s.rounds)

	s.offset = 0
	s.state[12]++
}
