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

var pool2k = &sync.Pool{
	New: createAllocFunc(2 * 1024),
}

var pool8k = &sync.Pool{
	New: createAllocFunc(8 * 1024),
}

var pool64k = &sync.Pool{
	New: createAllocFunc(64 * 1024),
}

var pool128k = &sync.Pool{
	New: createAllocFunc(128 * 1024),
}
