package kcp

import (
	"sync"
)

type SendingWindow struct {
	start uint32
	cap   uint32
	len   uint32
	last  uint32

	data []*DataSegment
	prev []uint32
	next []uint32

	totalInFlightSize uint32
	writer            SegmentWriter
	onPacketLoss      func(uint32)
}

func NewSendingWindow(size uint32, writer SegmentWriter, onPacketLoss func(uint32)) *SendingWindow {
	window := &SendingWindow{
		start:        0,
		cap:          size,
		len:          0,
		last:         0,
		data:         make([]*DataSegment, size),
		prev:         make([]uint32, size),
		next:         make([]uint32, size),
		writer:       writer,
		onPacketLoss: onPacketLoss,
	}
	return window
}

func (this *SendingWindow) Len() int {
	return int(this.len)
}

func (this *SendingWindow) IsEmpty() bool {
	return this.len == 0
}

func (this *SendingWindow) Size() uint32 {
	return this.cap
}

func (this *SendingWindow) IsFull() bool {
	return this.len == this.cap
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
	for !this.IsEmpty() && this.data[this.start].Number < una {
		this.Remove(0)
	}
}

func (this *SendingWindow) Remove(idx uint32) {
	if this.len == 0 {
		return
	}

	pos := (this.start + idx) % this.cap
	seg := this.data[pos]
	if seg == nil {
		return
	}
	this.totalInFlightSize--
	seg.Release()
	this.data[pos] = nil
	if pos == this.start && pos == this.last {
		this.len = 0
		this.start = 0
		this.last = 0
	} else if pos == this.start {
		delta := this.next[pos] - this.start
		if this.next[pos] < this.start {
			delta = this.next[pos] + this.cap - this.start
		}
		this.start = this.next[pos]
		this.len -= delta
	} else if pos == this.last {
		this.last = this.prev[pos]
	} else {
		this.next[this.prev[pos]] = this.next[pos]
		this.prev[this.next[pos]] = this.prev[pos]
	}
}

func (this *SendingWindow) HandleFastAck(number uint32) {
	if this.len == 0 {
		return
	}

	for i := this.start; ; i = this.next[i] {
		seg := this.data[i]
		if number-seg.Number > 0x7FFFFFFF {
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

func (this *SendingWindow) Flush(current uint32, resend uint32, rto uint32, maxInFlightSize uint32) {
	if this.IsEmpty() {
		return
	}

	var lost uint32
	var inFlightSize uint32

	for i := this.start; ; i = this.next[i] {
		segment := this.data[i]
		needsend := false
		if segment.transmit == 0 {
			needsend = true
			segment.transmit++
			segment.timeout = current + rto
			this.totalInFlightSize++
		} else if current-segment.timeout < 0x7FFFFFFF {
			needsend = true
			segment.transmit++
			segment.timeout = current + rto
			lost++
		} else if segment.ackSkipped >= resend {
			needsend = true
			segment.transmit++
			segment.ackSkipped = 0
			segment.timeout = current + rto
		}

		if needsend {
			segment.Timestamp = current
			this.writer.Write(segment)
			inFlightSize++
			if inFlightSize >= maxInFlightSize {
				break
			}
		}
		if i == this.last {
			break
		}
	}

	if this.onPacketLoss != nil && inFlightSize > 0 && this.totalInFlightSize != 0 {
		rate := lost * 100 / this.totalInFlightSize
		this.onPacketLoss(rate)
	}
}

type SendingWorker struct {
	sync.RWMutex
	conn                *Connection
	window              *SendingWindow
	firstUnacknowledged uint32
	nextNumber          uint32
	remoteNextNumber    uint32
	controlWindow       uint32
	fastResend          uint32
}

func NewSendingWorker(kcp *Connection) *SendingWorker {
	worker := &SendingWorker{
		conn:             kcp,
		fastResend:       2,
		remoteNextNumber: 32,
		controlWindow:    effectiveConfig.GetSendingInFlightSize(),
	}
	worker.window = NewSendingWindow(effectiveConfig.GetSendingBufferSize(), worker, worker.OnPacketLoss)
	return worker
}

func (this *SendingWorker) ProcessReceivingNext(nextNumber uint32) {
	this.Lock()
	defer this.Unlock()

	this.ProcessReceivingNextWithoutLock(nextNumber)
}

func (this *SendingWorker) ProcessReceivingNextWithoutLock(nextNumber uint32) {
	this.window.Clear(nextNumber)
	this.FindFirstUnacknowledged()
}

// Private: Visible for testing.
func (this *SendingWorker) FindFirstUnacknowledged() {
	if !this.window.IsEmpty() {
		this.firstUnacknowledged = this.window.First().Number
	} else {
		this.firstUnacknowledged = this.nextNumber
	}
}

// Private: Visible for testing.
func (this *SendingWorker) ProcessAck(number uint32) {
	// number < this.firstUnacknowledged || number >= this.nextNumber
	if number-this.firstUnacknowledged > 0x7FFFFFFF || number-this.nextNumber < 0x7FFFFFFF {
		return
	}

	this.window.Remove(number - this.firstUnacknowledged)
	this.FindFirstUnacknowledged()
}

func (this *SendingWorker) ProcessSegment(current uint32, seg *AckSegment) {
	defer seg.Release()

	this.Lock()
	defer this.Unlock()

	if this.remoteNextNumber < seg.ReceivingWindow {
		this.remoteNextNumber = seg.ReceivingWindow
	}
	this.ProcessReceivingNextWithoutLock(seg.ReceivingNext)
	if current-seg.Timestamp < 10000 {
		this.conn.roundTrip.Update(current-seg.Timestamp, current)
	}

	var maxack uint32
	for i := 0; i < int(seg.Count); i++ {
		number := seg.NumberList[i]

		this.ProcessAck(number)
		if maxack < number {
			maxack = number
		}
	}

	this.window.HandleFastAck(maxack)
}

func (this *SendingWorker) Push(b []byte) int {
	nBytes := 0
	this.Lock()
	defer this.Unlock()

	for len(b) > 0 && !this.window.IsFull() {
		var size int
		if len(b) > int(this.conn.mss) {
			size = int(this.conn.mss)
		} else {
			size = len(b)
		}
		seg := NewDataSegment()
		seg.Data = AllocateBuffer().Clear().Append(b[:size])
		seg.Number = this.nextNumber
		seg.timeout = 0
		seg.ackSkipped = 0
		seg.transmit = 0
		this.window.Push(seg)
		this.nextNumber++
		b = b[size:]
		nBytes += size
	}
	return nBytes
}

// Private: Visible for testing.
func (this *SendingWorker) Write(seg Segment) {
	dataSeg := seg.(*DataSegment)

	dataSeg.Conv = this.conn.conv
	dataSeg.SendingNext = this.firstUnacknowledged
	dataSeg.Option = 0
	if this.conn.State() == StateReadyToClose {
		dataSeg.Option = SegmentOptionClose
	}

	this.conn.output.Write(dataSeg)
}

func (this *SendingWorker) OnPacketLoss(lossRate uint32) {
	if !effectiveConfig.Congestion || this.conn.roundTrip.Timeout() == 0 {
		return
	}

	if lossRate >= 15 {
		this.controlWindow = 3 * this.controlWindow / 4
	} else if lossRate <= 5 {
		this.controlWindow += this.controlWindow / 4
	}
	if this.controlWindow < 16 {
		this.controlWindow = 16
	}
	if this.controlWindow > 2*effectiveConfig.GetSendingInFlightSize() {
		this.controlWindow = 2 * effectiveConfig.GetSendingInFlightSize()
	}
}

func (this *SendingWorker) Flush(current uint32) {
	this.Lock()
	defer this.Unlock()

	cwnd := this.firstUnacknowledged + effectiveConfig.GetSendingInFlightSize()
	if cwnd > this.remoteNextNumber {
		cwnd = this.remoteNextNumber
	}
	if effectiveConfig.Congestion && cwnd > this.firstUnacknowledged+this.controlWindow {
		cwnd = this.firstUnacknowledged + this.controlWindow
	}

	if !this.window.IsEmpty() {
		this.window.Flush(current, this.conn.fastresend, this.conn.roundTrip.Timeout(), cwnd)
	}
}

func (this *SendingWorker) CloseWrite() {
	this.Lock()
	defer this.Unlock()

	this.window.Clear(0xFFFFFFFF)
}

func (this *SendingWorker) IsEmpty() bool {
	this.RLock()
	defer this.RUnlock()

	return this.window.IsEmpty()
}
