package buf

import (
	"sync"
)

// Pool provides functionality to generate and recycle buffers on demand.
type Pool struct {
	allocator *sync.Pool
}

// NewPool creates a SyncPool with given buffer size.
func NewPool(bufferSize uint32) *Pool {
	pool := &Pool{
		allocator: &sync.Pool{
			New: func() interface{} { return make([]byte, bufferSize) },
		},
	}
	return pool
}

// Allocate either returns a unused buffer from the pool, or generates a new one from system.
func (p *Pool) Allocate() *Buffer {
	return &Buffer{
		v:    p.allocator.Get().([]byte),
		pool: p,
	}
}

// // Free recycles the given buffer.
func (p *Pool) Free(buffer *Buffer) {
	if buffer.v != nil {
		p.allocator.Put(buffer.v)
	}
}

const (
	// Size of a regular buffer.
	Size = 2 * 1024
)

var mediumPool = NewPool(Size)
