package kcp

type SendingQueue struct {
	start uint32
	cap   uint32
	len   uint32
	list  []*DataSegment
}

func NewSendingQueue(size uint32) *SendingQueue {
	return &SendingQueue{
		start: 0,
		cap:   size,
		list:  make([]*DataSegment, size),
		len:   0,
	}
}

func (this *SendingQueue) IsFull() bool {
	return this.len == this.cap
}

func (this *SendingQueue) IsEmpty() bool {
	return this.len == 0
}

func (this *SendingQueue) Pop() *DataSegment {
	if this.IsEmpty() {
		return nil
	}
	seg := this.list[this.start]
	this.list[this.start] = nil
	this.len--
	this.start++
	if this.start == this.cap {
		this.start = 0
	}
	return seg
}

func (this *SendingQueue) Push(seg *DataSegment) {
	if this.IsFull() {
		return
	}
	this.list[(this.start+this.len)%this.cap] = seg
	this.len++
}

func (this *SendingQueue) Clear() {
	for i := uint32(0); i < this.len; i++ {
		this.list[(i+this.start)%this.cap].Release()
		this.list[(i+this.start)%this.cap] = nil
	}
	this.start = 0
	this.len = 0
}

func (this *SendingQueue) Len() uint32 {
	return this.len
}
