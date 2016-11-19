package alloc

import (
	"os"
	"strconv"
	"sync"
)

type Pool interface {
	Allocate() *Buffer
	Free(*Buffer)
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
	rawBuffer := buffer.head
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
	mediumBufferByteSize = 8 * 1024
	BufferSize           = mediumBufferByteSize - defaultOffset

	largeBufferByteSize = 64 * 1024
	LargeBufferSize     = largeBufferByteSize - defaultOffset

	PoolSizeEnvKey = "v2ray.buffer.size"
)

var (
	mediumPool *BufferPool
	//largePool  *BufferPool
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
	totalByteSize := size * 1024 * 1024
	mediumPool = NewBufferPool(mediumBufferByteSize, totalByteSize/mediumBufferByteSize)
	//largePool = NewBufferPool(largeBufferByteSize, totalByteSize/4/largeBufferByteSize)
}
