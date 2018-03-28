package buf

import (
	"sync"
)

const (
	// Size of a regular buffer.
	Size = 2 * 1024
)

func createAllocFunc(size uint32) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

const (
	numPools  = 5
	sizeMulti = 4
)

var (
	pool     [numPools]sync.Pool
	poolSize [numPools]uint32
)

func init() {
	size := uint32(Size)
	for i := 0; i < numPools; i++ {
		pool[i] = sync.Pool{
			New: createAllocFunc(size),
		}
		poolSize[i] = size
		size *= sizeMulti
	}
}

func newBytes(size uint32) []byte {
	for idx, ps := range poolSize {
		if size <= ps {
			return pool[idx].Get().([]byte)
		}
	}
	return make([]byte, size)
}

func freeBytes(b []byte) {
	size := uint32(cap(b))
	b = b[0:cap(b)]
	for i := numPools - 1; i >= 0; i-- {
		if size >= poolSize[i] {
			pool[i].Put(b)
			return
		}
	}
}
