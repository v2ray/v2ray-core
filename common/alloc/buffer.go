package alloc

import (
	//"fmt"
	"time"
)

type Buffer struct {
	head  []byte
	pool  *bufferPool
	Value []byte
}

func (b *Buffer) Release() {
	b.pool.free(b)
}

func (b *Buffer) Clear() {
	b.Value = b.Value[:0]
}

func (b *Buffer) Append(data []byte) {
	b.Value = append(b.Value, data...)
}

func (b *Buffer) Slice(from, to int) {
	b.Value = b.Value[from:to]
}

func (b *Buffer) SliceFrom(from int) {
	b.Value = b.Value[from:]
}

func (b *Buffer) Len() int {
	return len(b.Value)
}

type bufferPool struct {
	chain         chan *Buffer
	allocator     func(*bufferPool) *Buffer
	elements2Keep int
}

func newBufferPool(allocator func(*bufferPool) *Buffer, elements2Keep, size int) *bufferPool {
	pool := &bufferPool{
		chain:         make(chan *Buffer, size),
		allocator:     allocateSmall,
		elements2Keep: elements2Keep,
	}
	for i := 0; i < elements2Keep; i++ {
		pool.chain <- allocator(pool)
	}
	go pool.cleanup(time.Tick(1 * time.Second))
	return pool
}

func (p *bufferPool) allocate() *Buffer {
	//fmt.Printf("Pool size: %d\n", len(p.chain))
	var b *Buffer
	select {
	case b = <-p.chain:
	default:
		b = p.allocator(p)
	}
	b.Value = b.head
	return b
}

func (p *bufferPool) free(buffer *Buffer) {
	select {
	case p.chain <- buffer:
	default:
	}
	//fmt.Printf("Pool size: %d\n", len(p.chain))
}

func (p *bufferPool) cleanup(tick <-chan time.Time) {
	for range tick {
		pSize := len(p.chain)
		if pSize > p.elements2Keep {
			<-p.chain
			continue
		}
		for delta := pSize - p.elements2Keep; delta > 0; delta-- {
			p.chain <- p.allocator(p)
		}
	}
}

func allocateSmall(pool *bufferPool) *Buffer {
	b := &Buffer{
		head: make([]byte, 8*1024),
	}
	b.Value = b.head
	b.pool = pool
	return b
}

var smallPool = newBufferPool(allocateSmall, 256, 2048)

func NewBuffer() *Buffer {
	return smallPool.allocate()
}
