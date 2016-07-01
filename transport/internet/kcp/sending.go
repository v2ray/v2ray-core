package kcp

type SendingWindow struct {
	start uint32
	cap   uint32
	len   uint32
	last  uint32

	data []*DataSegment
	prev []uint32
	next []uint32

	kcp *KCP
}

func NewSendingWindow(kcp *KCP, size uint32) *SendingWindow {
	window := &SendingWindow{
		start: 0,
		cap:   size,
		len:   0,
		last:  0,
		data:  make([]*DataSegment, size),
		prev:  make([]uint32, size),
		next:  make([]uint32, size),
		kcp:   kcp,
	}
	return window
}

func (this *SendingWindow) Len() int {
	return int(this.len)
}

func (this *SendingWindow) Push(seg *DataSegment) {
	pos := (this.start + this.len) % this.cap
	this.data[pos] = seg
	if this.len > 0 {
		this.next[this.last] = pos
		this.prev[pos] = this.last
	}
	this.last = pos
	this.len++
}

func (this *SendingWindow) First() *DataSegment {
	return this.data[this.start]
}

func (this *SendingWindow) Clear(una uint32) {
	for this.Len() > 0 && this.data[this.start].Number < una {
		this.Remove(0)
	}
}

func (this *SendingWindow) Remove(idx uint32) {
	pos := (this.start + idx) % this.cap
	seg := this.data[pos]
	if seg == nil {
		return
	}
	seg.Release()
	this.data[pos] = nil
	if pos == this.start {
		if this.start == this.last {
			this.len = 0
			this.start = 0
			this.last = 0
		} else {
			delta := this.next[pos] - this.start
			this.start = this.next[pos]
			this.len -= delta
		}
	} else if pos == this.last {
		this.last = this.prev[pos]
	} else {
		this.next[this.prev[pos]] = this.next[pos]
		this.prev[this.next[pos]] = this.prev[pos]
	}
}

func (this *SendingWindow) HandleFastAck(number uint32) {

	for i := this.start; ; i = this.next[i] {
		seg := this.data[i]
		if _itimediff(number, seg.Number) < 0 {
			break
		}
		if number != seg.Number {
			seg.ackSkipped++
		}
		if i == this.last {
			break
		}
	}
}

func (this *SendingWindow) Flush() bool {
	if this.Len() == 0 {
		return false
	}

	current := this.kcp.current
	resent := uint32(this.kcp.fastresend)
	if this.kcp.fastresend <= 0 {
		resent = 0xffffffff
	}
	lost := false
	segSent := false

	for i := this.start; ; i = this.next[i] {
		segment := this.data[i]
		needsend := false
		if segment.transmit == 0 {
			needsend = true
			segment.transmit++
			segment.timeout = current + this.kcp.rx_rto
		} else if _itimediff(current, segment.timeout) >= 0 {
			needsend = true
			segment.transmit++
			segment.timeout = current + this.kcp.rx_rto
			lost = true
		} else if segment.ackSkipped >= resent {
			needsend = true
			segment.transmit++
			segment.ackSkipped = 0
			segment.timeout = current + this.kcp.rx_rto
			lost = true
		}

		if needsend {
			segment.Timestamp = current
			segment.SendingNext = this.kcp.snd_una
			segment.Opt = 0
			if this.kcp.state == StateReadyToClose {
				segment.Opt = SegmentOptionClose
			}

			this.kcp.output.Write(segment)
			segSent = true
		}
		if i == this.last {
			break
		}
	}

	this.kcp.HandleLost(lost)

	return segSent
}

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
