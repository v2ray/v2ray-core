package kcp

import (
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
)

type SendingWindow struct {
	start uint32
	cap   uint32
	len   uint32
	last  uint32

	data []*DataSegment
	prev []uint32
	next []uint32

	inFlightSize uint32
	writer       SegmentWriter
	onPacketLoss func(bool)
}

func NewSendingWindow(size uint32, inFlightSize uint32, writer SegmentWriter, onPacketLoss func(bool)) *SendingWindow {
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
		inFlightSize: inFlightSize,
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
	if this.len == 0 {
		return
	}

	pos := (this.start + idx) % this.cap
	seg := this.data[pos]
	if seg == nil {
		return
	}
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

func (this *SendingWindow) Flush(current uint32, resend uint32, rto uint32) {
	if this.Len() == 0 {
		return
	}

	lost := false
	var inFlightSize uint32

	for i := this.start; ; i = this.next[i] {
		segment := this.data[i]
		needsend := false
		if segment.transmit == 0 {
			needsend = true
			segment.transmit++
			segment.timeout = current + rto
		} else if _itimediff(current, segment.timeout) >= 0 {
			needsend = true
			segment.transmit++
			segment.timeout = current + rto
			lost = true
		} else if segment.ackSkipped >= resend {
			needsend = true
			segment.transmit++
			segment.ackSkipped = 0
			segment.timeout = current + rto
			lost = true
		}

		if needsend {
			this.writer.Write(segment)
			inFlightSize++
			if inFlightSize >= this.inFlightSize {
				break
			}
		}
		if i == this.last {
			break
		}
	}

	this.onPacketLoss(lost)
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

type SendingWorker struct {
	sync.Mutex
	kcp                 *KCP
	window              *SendingWindow
	queue               *SendingQueue
	windowSize          uint32
	firstUnacknowledged uint32
	nextNumber          uint32
	remoteNextNumber    uint32
	controlWindow       uint32
	fastResend          uint32
	updated             bool
}

func NewSendingWorker(kcp *KCP) *SendingWorker {
	worker := &SendingWorker{
		kcp:              kcp,
		queue:            NewSendingQueue(effectiveConfig.GetSendingQueueSize()),
		fastResend:       2,
		remoteNextNumber: 32,
		windowSize:       effectiveConfig.GetSendingWindowSize(),
		controlWindow:    effectiveConfig.GetSendingWindowSize(),
	}
	worker.window = NewSendingWindow(effectiveConfig.GetSendingWindowSize(), effectiveConfig.GetSendingInFlightSize(), worker, worker.OnPacketLoss)
	return worker
}

func (this *SendingWorker) ProcessReceivingNext(nextNumber uint32) {
	this.Lock()
	defer this.Unlock()

	this.window.Clear(nextNumber)
	this.FindFirstUnacknowledged()
}

// @Private
func (this *SendingWorker) FindFirstUnacknowledged() {
	prevUna := this.firstUnacknowledged
	if this.window.Len() > 0 {
		this.firstUnacknowledged = this.window.First().Number
	} else {
		this.firstUnacknowledged = this.nextNumber
	}
	if this.firstUnacknowledged != prevUna {
		this.updated = true
	}
}

func (this *SendingWorker) ProcessAck(number uint32) {
	if number-this.firstUnacknowledged > this.windowSize {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.window.Remove(number - this.firstUnacknowledged)
	this.FindFirstUnacknowledged()
}

func (this *SendingWorker) ProcessAckSegment(seg *AckSegment) {
	if this.remoteNextNumber < seg.ReceivingWindow {
		this.remoteNextNumber = seg.ReceivingWindow
	}
	this.ProcessReceivingNext(seg.ReceivingNext)
	var maxack uint32
	for i := 0; i < int(seg.Count); i++ {
		timestamp := seg.TimestampList[i]
		number := seg.NumberList[i]
		if this.kcp.current-timestamp > 10000 {
			this.kcp.update_ack(int32(this.kcp.current - timestamp))
		}
		this.ProcessAck(number)
		if maxack < number {
			maxack = number
		}
	}
	this.Lock()
	this.window.HandleFastAck(maxack)
	this.Unlock()
}

func (this *SendingWorker) Push(b []byte) int {
	nBytes := 0
	for len(b) > 0 && !this.queue.IsFull() {
		var size int
		if len(b) > int(this.kcp.mss) {
			size = int(this.kcp.mss)
		} else {
			size = len(b)
		}
		seg := &DataSegment{
			Data: alloc.NewSmallBuffer().Clear().Append(b[:size]),
		}
		this.Lock()
		this.queue.Push(seg)
		this.Unlock()
		b = b[size:]
		nBytes += size
	}
	return nBytes
}

func (this *SendingWorker) Write(seg ISegment) {
	dataSeg := seg.(*DataSegment)

	dataSeg.Conv = this.kcp.conv
	dataSeg.Timestamp = this.kcp.current
	dataSeg.SendingNext = this.firstUnacknowledged
	dataSeg.Opt = 0
	if this.kcp.state == StateReadyToClose {
		dataSeg.Opt = SegmentOptionClose
	}

	this.kcp.output.Write(dataSeg)
	this.updated = false
}

func (this *SendingWorker) PingNecessary() bool {
	return this.updated
}

func (this *SendingWorker) OnPacketLoss(lost bool) {
	if !effectiveConfig.Congestion {
		return
	}

	if lost {
		this.controlWindow = 3 * this.controlWindow / 4
	} else {
		this.controlWindow += this.controlWindow / 4
	}
	if this.controlWindow < 4 {
		this.controlWindow = 4
	}
	if this.controlWindow > 2*this.windowSize {
		this.controlWindow = 2 * this.windowSize
	}
}

func (this *SendingWorker) Flush() {
	this.Lock()
	defer this.Unlock()

	cwnd := this.firstUnacknowledged + this.windowSize
	if cwnd > this.remoteNextNumber {
		cwnd = this.remoteNextNumber
	}
	if effectiveConfig.Congestion && cwnd > this.firstUnacknowledged+this.controlWindow {
		cwnd = this.firstUnacknowledged + this.controlWindow
	}

	for !this.queue.IsEmpty() && _itimediff(this.nextNumber, cwnd) < 0 {
		seg := this.queue.Pop()
		seg.Number = this.nextNumber
		seg.timeout = this.kcp.current
		seg.ackSkipped = 0
		seg.transmit = 0
		this.window.Push(seg)
		this.nextNumber++
	}

	this.window.Flush(this.kcp.current, this.kcp.fastresend, this.kcp.rx_rto)
}

func (this *SendingWorker) CloseWrite() {
	this.Lock()
	defer this.Unlock()

	this.window.Clear(0xFFFFFFFF)
	this.queue.Clear()
}
