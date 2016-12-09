package serial

import (
	"hash"

	"v2ray.com/core/common/buf"
)

func WriteHash(h hash.Hash) buf.Supplier {
	return func(b []byte) (int, error) {
		h.Sum(b[:0])
		return h.Size(), nil
	}
}
