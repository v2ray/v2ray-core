package serial

import "hash"

func WriteHash(h hash.Hash) func([]byte) (int, error) {
	return func(b []byte) (int, error) {
		h.Sum(b[:0])
		return h.Size(), nil
	}
}
