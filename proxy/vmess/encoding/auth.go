package encoding

import (
	"hash/fnv"
)

func Authenticate(b []byte) uint32 {
	fnv1hash := fnv.New32a()
	fnv1hash.Write(b)
	return fnv1hash.Sum32()
}
