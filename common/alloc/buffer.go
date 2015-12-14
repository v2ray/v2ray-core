package alloc

import (
	"time"
)

// Buffer is a recyclable allocation of a byte array. Buffer.Release() recycles
// the buffer into an internal buffer pool, in order to recreate a buffer more
// quickly.
type Buffer struct {
	head  []byte
	pool  *bufferPool
	Value []byte
}

// Release recycles the buffer into an internal buffer pool.
func (b *Buffer) Release() {
	b.pool.free(b)
	b.head = nil
	b.Value = nil
	b.pool = nil
}

// Clear clears the content of the buffer, results an empty buffer with
// Len() = 0.
func (b *Buffer) Clear() *Buffer {
	b.Value = b.head[:0]
	return b
}

// AppendBytes appends one or more bytes to the end of the buffer.
func (b *Buffer) AppendBytes(bytes ...byte) *Buffer {
	b.Value = append(b.Value, bytes...)
	return b
}

// Append appends a byte array to the end of the buffer.
func (b *Buffer) Append(data []byte) *Buffer {
	b.Value = append(b.Value, data...)
	return b
}

func (b *Buffer) Bytes() []byte {
	return b.Value
}

// Slice cuts the buffer at the given position.
func (b *Buffer) Slice(from, to int) *Buffer {
	b.Value = b.Value[from:to]
	return b
}

// SliceFrom cuts the buffer at the given position.
func (b *Buffer) SliceFrom(from int) *Buffer {
	b.Value = b.Value[from:]
	return b
}

// Len returns the length of the buffer content.
func (b *Buffer) Len() int {
	return len(b.Value)
}

// IsFull returns true if the buffer has no more room to grow.
func (b *Buffer) IsFull() bool {
	return len(b.Value) == cap(b.Value)
}

// Write implements Write method in io.Writer.
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

var smallPool = newBufferPool(1024, 64, 512)
var mediumPool = newBufferPool(8*1024, 256, 2048)
var largePool = newBufferPool(64*1024, 128, 1024)

// NewSmallBuffer creates a Buffer with 1K bytes of arbitrary content.
func NewSmallBuffer() *Buffer {
	return smallPool.allocate()
}

// NewBuffer creates a Buffer with 8K bytes of arbitrary content.
func NewBuffer() *Buffer {
	return mediumPool.allocate()
}

// NewLargeBuffer creates a Buffer with 64K bytes of arbitrary content.
func NewLargeBuffer() *Buffer {
	return largePool.allocate()
}
