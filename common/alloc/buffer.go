package alloc

import (
	"time"
)

type Buffer struct {
	head  []byte
	pool  *bufferPool
	Value []byte
}

func (b *Buffer) Release() {
	b.pool.free(b)
	b.head = nil
	b.Value = nil
	b.pool = nil
}

func (b *Buffer) Clear() {
	b.Value = b.head[:0]
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
	chain        chan []byte
	bufferSize   int
	buffers2Keep int
}

func newBufferPool(bufferSize, buffers2Keep, poolSize int) *bufferPool {
	pool := &bufferPool{
		chain:        make(chan []byte, poolSize),
		bufferSize:   bufferSize,
		buffers2Keep: buffers2Keep,
	}
	for i := 0; i < buffers2Keep; i++ {
		pool.chain <- make([]byte, bufferSize)
	}
	go pool.cleanup(time.Tick(1 * time.Second))
	return pool
}

func (p *bufferPool) allocate() *Buffer {
	var b []byte
	select {
	case b = <-p.chain:
	default:
		b = make([]byte, p.bufferSize)
	}
	return &Buffer{
		head:  b,
		pool:  p,
		Value: b,
	}
}

func (p *bufferPool) free(buffer *Buffer) {
	select {
	case p.chain <- buffer.head:
	default:
	}
}

func (p *bufferPool) cleanup(tick <-chan time.Time) {
	for range tick {
		pSize := len(p.chain)
		if pSize > p.buffers2Keep {
			<-p.chain
			continue
		}
		for delta := pSize - p.buffers2Keep; delta > 0; delta-- {
			p.chain <- make([]byte, p.bufferSize)
		}
	}
}

var smallPool = newBufferPool(8*1024, 256, 2048)

func NewBuffer() *Buffer {
	return smallPool.allocate()
}
