package kcp

type ReceivingWindow struct {
	start uint32
	size  uint32
	list  []*Segment
}

func NewReceivingWindow(size uint32) *ReceivingWindow {
	return &ReceivingWindow{
		start: 0,
		size:  size,
		list:  make([]*Segment, size),
	}
}

func (this *ReceivingWindow) Size() uint32 {
	return this.size
}

func (this *ReceivingWindow) Position(idx uint32) uint32 {
	return (idx + this.start) % this.size
}

func (this *ReceivingWindow) Set(idx uint32, value *Segment) bool {
	pos := this.Position(idx)
	if this.list[pos] != nil {
		return false
	}
	this.list[pos] = value
	return true
}

func (this *ReceivingWindow) Remove(idx uint32) *Segment {
	pos := this.Position(idx)
	e := this.list[pos]
	this.list[pos] = nil
	return e
}

func (this *ReceivingWindow) RemoveFirst() *Segment {
	return this.Remove(0)
}

func (this *ReceivingWindow) Advance() {
	this.start++
	if this.start == this.size {
		this.start = 0
	}
}
