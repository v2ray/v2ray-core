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

func (b *Buffer) Clear() *Buffer {
	b.Value = b.head[:0]
	return b
}

func (b *Buffer) AppendBytes(bytes ...byte) *Buffer {
	b.Value = append(b.Value, bytes...)
	return b
}

func (b *Buffer) Append(data []byte) *Buffer {
	b.Value = append(b.Value, data...)
	return b
}

func (b *Buffer) Slice(from, to int) *Buffer {
	b.Value = b.Value[from:to]
	return b
}

func (b *Buffer) SliceFrom(from int) *Buffer {
	b.Value = b.Value[from:]
	return b
}

func (b *Buffer) Len() int {
	return len(b.Value)
}

func (b *Buffer) IsFull() bool {
	return len(b.Value) == cap(b.Value)
}

func (b *Buffer) Write(data []byte) (int, error) {
	b.Append(data)
	return len(data), nil
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
		for delta := p.buffers2Keep - pSize; delta > 0; delta-- {
			select {
			case p.chain <- make([]byte, p.bufferSize):
			default:
			}
		}
	}
}

var smallPool = newBufferPool(1024, 16, 64)
var mediumPool = newBufferPool(8*1024, 256, 2048)
var largePool = newBufferPool(64*1024, 128, 1024)

func NewSmallBuffer() *Buffer {
	return smallPool.allocate()
}

func NewBuffer() *Buffer {
	return mediumPool.allocate()
}

func NewLargeBuffer() *Buffer {
	return largePool.allocate()
}
