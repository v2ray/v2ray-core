package buf

import (
	"sync"
)

// Pool provides functionality to generate and recycle buffers on demand.
type Pool interface {
	// Allocate either returns a unused buffer from the pool, or generates a new one from system.
	Allocate() *Buffer
	// Free recycles the given buffer.
	Free(*Buffer)
}

// SyncPool is a buffer pool based on sync.Pool
type SyncPool struct {
	allocator *sync.Pool
}

// NewSyncPool creates a SyncPool with given buffer size.
func NewSyncPool(bufferSize uint32) *SyncPool {
	pool := &SyncPool{
		allocator: &sync.Pool{
			New: func() interface{} { return make([]byte, bufferSize) },
		},
	}
	return pool
}

// Allocate implements Pool.Allocate().
func (p *SyncPool) Allocate() *Buffer {
	return &Buffer{
		v:    p.allocator.Get().([]byte),
		pool: p,
	}
}

// Free implements Pool.Free().
func (p *SyncPool) Free(buffer *Buffer) {
	if buffer.v != nil {
		p.allocator.Put(buffer.v)
	}
}

const (
	// Size of a regular buffer.
	Size = 2 * 1024
)

var (
	mediumPool Pool = NewSyncPool(Size)
)
