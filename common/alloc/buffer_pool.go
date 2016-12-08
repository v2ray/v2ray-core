package alloc

import (
	"os"
	"strconv"
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
	return CreateBuffer(p.allocator.Get().([]byte), p)
}

// Free implements Pool.Free().
func (p *SyncPool) Free(buffer *Buffer) {
	rawBuffer := buffer.v
	if rawBuffer == nil {
		return
	}
	p.allocator.Put(rawBuffer)
}

type BufferPool struct {
	chain     chan []byte
	allocator *sync.Pool
}

func NewBufferPool(bufferSize, poolSize uint32) *BufferPool {
	pool := &BufferPool{
		chain: make(chan []byte, poolSize),
		allocator: &sync.Pool{
			New: func() interface{} { return make([]byte, bufferSize) },
		},
	}
	for i := uint32(0); i < poolSize; i++ {
		pool.chain <- make([]byte, bufferSize)
	}
	return pool
}

func (p *BufferPool) Allocate() *Buffer {
	var b []byte
	select {
	case b = <-p.chain:
	default:
		b = p.allocator.Get().([]byte)
	}
	return CreateBuffer(b, p)
}

func (p *BufferPool) Free(buffer *Buffer) {
	rawBuffer := buffer.v
	if rawBuffer == nil {
		return
	}
	select {
	case p.chain <- rawBuffer:
	default:
		p.allocator.Put(rawBuffer)
	}
}

const (
	BufferSize      = 8 * 1024
	SmallBufferSize = 2 * 1024

	PoolSizeEnvKey = "v2ray.buffer.size"
)

var (
	mediumPool Pool
	smallPool  = NewSyncPool(2048)
)

func init() {
	var size uint32 = 20
	sizeStr := os.Getenv(PoolSizeEnvKey)
	if len(sizeStr) > 0 {
		customSize, err := strconv.ParseUint(sizeStr, 10, 32)
		if err == nil {
			size = uint32(customSize)
		}
	}
	if size > 0 {
		totalByteSize := size * 1024 * 1024
		mediumPool = NewBufferPool(BufferSize, totalByteSize/BufferSize)
	} else {
		mediumPool = NewSyncPool(BufferSize)
	}
}
