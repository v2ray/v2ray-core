package bytespool

import "sync"

func createAllocFunc(size int32) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

// The following parameters controls the size of buffer pools.
// There are numPools pools. Starting from 2k size, the size of each pool is sizeMulti of the previous one.
// Package buf is guaranteed to not use buffers larger than the largest pool.
// Other packets may use larger buffers.
const (
	numPools  = 4
	sizeMulti = 4
)

var (
	pool     [numPools]sync.Pool
	poolSize [numPools]int32
)

func init() {
	size := int32(2048)
	for i := 0; i < numPools; i++ {
		pool[i] = sync.Pool{
			New: createAllocFunc(size),
		}
		poolSize[i] = size
		size *= sizeMulti
	}
}

// Alloc returns a byte slice with at least the given size. Minimum size of returned slice is 2048.
func Alloc(size int32) []byte {
	for idx, ps := range poolSize {
		if size <= ps {
			return pool[idx].Get().([]byte)
		}
	}
	return make([]byte, size)
}

// Free puts a byte slice into the internal pool.
func Free(b []byte) {
	size := int32(cap(b))
	b = b[0:cap(b)]
	for i := numPools - 1; i >= 0; i-- {
		if size >= poolSize[i] {
			pool[i].Put(b) // nolint: megacheck
			return
		}
	}
}
