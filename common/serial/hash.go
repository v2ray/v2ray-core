package serial

import (
	"hash"

	"v2ray.com/core/common/buf"
)

func WriteHash(h hash.Hash) buf.BytesWriter {
	return func(b []byte) int {
		h.Sum(b[:0])
		return h.Size()
	}
}
