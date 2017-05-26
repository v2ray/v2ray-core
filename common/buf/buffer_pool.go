package buf

import (
	"runtime"
	"sync"

	"v2ray.com/core/common/platform"
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

// BufferPool is a Pool that utilizes an internal cache.
type BufferPool struct {
	chain chan []byte
	sub   Pool
}

// NewBufferPool creates a new BufferPool with given buffer size, and internal cache size.
func NewBufferPool(bufferSize, poolSize uint32) *BufferPool {
	pool := &BufferPool{
		chain: make(chan []byte, poolSize),
		sub:   NewSyncPool(bufferSize),
	}
	for i := uint32(0); i < poolSize; i++ {
		pool.chain <- make([]byte, bufferSize)
	}
	return pool
}

// Allocate implements Pool.Allocate().
func (p *BufferPool) Allocate() *Buffer {
	select {
	case b := <-p.chain:
		return &Buffer{
			v:    b,
			pool: p,
		}
	default:
		return p.sub.Allocate()
	}
}

// Free implements Pool.Free().
func (p *BufferPool) Free(buffer *Buffer) {
	if buffer.v == nil {
		return
	}
	select {
	case p.chain <- buffer.v:
	default:
		p.sub.Free(buffer)
	}
}

const (
	// Size of a regular buffer.
	Size = 2 * 1024

	poolSizeEnvKey = "v2ray.buffer.size"
)

var (
	mediumPool Pool
)

func getDefaultPoolSize() int {
	switch runtime.GOARCH {
	case "amd64", "386":
		return 20
	default:
		return 5
	}
}

func init() {
	f := platform.EnvFlag{
		Name:    poolSizeEnvKey,
		AltName: platform.NormalizeEnvName(poolSizeEnvKey),
	}
	size := f.GetValueAsInt(getDefaultPoolSize())
	if size > 0 {
		totalByteSize := uint32(size) * 1024 * 1024
		mediumPool = NewBufferPool(Size, totalByteSize/Size)
	} else {
		mediumPool = NewSyncPool(Size)
	}
}
