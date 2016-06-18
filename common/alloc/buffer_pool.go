package alloc

import (
	"sync"
)

type BufferPool struct {
	chain     chan []byte
	allocator *sync.Pool
}

func NewBufferPool(bufferSize, poolSize int) *BufferPool {
	pool := &BufferPool{
		chain: make(chan []byte, poolSize),
		allocator: &sync.Pool{
			New: func() interface{} { return make([]byte, bufferSize) },
		},
	}
	for i := 0; i < poolSize/2; i++ {
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
	return &Buffer{
		head:   b,
		pool:   p,
		Value:  b[defaultOffset:],
		offset: defaultOffset,
	}
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
	SmallBufferSize = 1600 - defaultOffset
	BufferSize      = 8*1024 - defaultOffset
	LargeBufferSize = 64*1024 - defaultOffset
)

var smallPool = NewBufferPool(1600, 128)
var mediumPool = NewBufferPool(8*1024, 128)
var largePool = NewBufferPool(64*1024, 64)
