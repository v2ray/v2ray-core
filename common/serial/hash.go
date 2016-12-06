package serial

import (
	"hash"
	"v2ray.com/core/common/alloc"
)

func WriteHash(h hash.Hash) alloc.BytesWriter {
	return func(b []byte) int {
		h.Sum(b[:0])
		return h.Size()
	}
}
