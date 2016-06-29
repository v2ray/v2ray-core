package kcp

type ReceivingWindow struct {
	start uint32
	size  uint32
	list  []*DataSegment
}

func NewReceivingWindow(size uint32) *ReceivingWindow {
	return &ReceivingWindow{
		start: 0,
		size:  size,
		list:  make([]*DataSegment, size),
	}
}

func (this *ReceivingWindow) Size() uint32 {
	return this.size
}

func (this *ReceivingWindow) Position(idx uint32) uint32 {
	return (idx + this.start) % this.size
}

func (this *ReceivingWindow) Set(idx uint32, value *DataSegment) bool {
	pos := this.Position(idx)
	if this.list[pos] != nil {
		return false
	}
	this.list[pos] = value
	return true
}

func (this *ReceivingWindow) Remove(idx uint32) *DataSegment {
	pos := this.Position(idx)
	e := this.list[pos]
	this.list[pos] = nil
	return e
}

func (this *ReceivingWindow) RemoveFirst() *DataSegment {
	return this.Remove(0)
}

func (this *ReceivingWindow) Advance() {
	this.start++
	if this.start == this.size {
		this.start = 0
	}
}

type ACKList struct {
	timestamps []uint32
	numbers    []uint32
}

func (this *ACKList) Add(number uint32, timestamp uint32) {
	this.timestamps = append(this.timestamps, timestamp)
	this.numbers = append(this.numbers, number)
}

func (this *ACKList) Clear(una uint32) bool {
	count := 0
	for i := 0; i < len(this.numbers); i++ {
		if this.numbers[i] >= una {
			if i != count {
				this.numbers[count] = this.numbers[i]
				this.timestamps[count] = this.timestamps[i]
			}
			count++
		}
	}
	if count < len(this.numbers) {
		this.numbers = this.numbers[:count]
		this.timestamps = this.timestamps[:count]
		return true
	}
	return false
}

func (this *ACKList) AsSegment() *ACKSegment {
	count := len(this.numbers)
	if count == 0 {
		return nil
	}

	if count > 128 {
		count = 128
	}
	seg := &ACKSegment{
		Count:         byte(count),
		NumberList:    this.numbers[:count],
		TimestampList: this.timestamps[:count],
	}
	//this.numbers = nil
	//this.timestamps = nil
	return seg
}
