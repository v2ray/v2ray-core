package kcp

import (
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
)

const (
	NumDistro  = 5
	DistroSize = 1600
)

type Buffer struct {
	sync.Mutex
	buffer *alloc.Buffer

	next     int
	released int
	hold     bool
	distro   [NumDistro]*alloc.Buffer
}

func NewBuffer() *Buffer {
	b := &Buffer{
		next:     0,
		released: 0,
		hold:     true,
		buffer:   alloc.NewBuffer(),
	}
	for idx := range b.distro {
		content := b.buffer.Value[idx*DistroSize : (idx+1)*DistroSize]
		b.distro[idx] = alloc.CreateBuffer(content, b)
	}
	return b
}

func (this *Buffer) IsEmpty() bool {
	this.Lock()
	defer this.Unlock()

	return this.next == NumDistro
}

func (this *Buffer) Allocate() *alloc.Buffer {
	this.Lock()
	defer this.Unlock()
	if this.next == NumDistro {
		return nil
	}
	b := this.distro[this.next]
	this.next++
	return b
}

func (this *Buffer) Free(b *alloc.Buffer) {
	this.Lock()
	defer this.Unlock()

	this.released++
	if !this.hold && this.released == this.next {
		this.ReleaseBuffer()
	}
}

func (this *Buffer) Release() {
	this.Lock()
	defer this.Unlock()

	if this.next == this.released {
		this.ReleaseBuffer()
	}
	this.hold = false
}

func (this *Buffer) ReleaseBuffer() {
	this.buffer.Release()
	this.buffer = nil
	for idx := range this.distro {
		this.distro[idx] = nil
	}
}

var (
	globalBuffer       *Buffer
	globalBufferAccess sync.Mutex
)

func AllocateBuffer() *alloc.Buffer {
	globalBufferAccess.Lock()
	defer globalBufferAccess.Unlock()

	if globalBuffer == nil {
		globalBuffer = NewBuffer()
	}
	b := globalBuffer.Allocate()
	if globalBuffer.IsEmpty() {
		globalBuffer.Release()
		globalBuffer = nil
	}
	return b
}
