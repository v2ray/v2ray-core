package kcp

import (
	"sync"

	"v2ray.com/core/common/alloc"
)

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

type AckList struct {
	writer     SegmentWriter
	timestamps []uint32
	numbers    []uint32
	nextFlush  []uint32
}

func NewAckList(writer SegmentWriter) *AckList {
	return &AckList{
		writer:     writer,
		timestamps: make([]uint32, 0, 32),
		numbers:    make([]uint32, 0, 32),
		nextFlush:  make([]uint32, 0, 32),
	}
}

func (this *AckList) Add(number uint32, timestamp uint32) {
	this.timestamps = append(this.timestamps, timestamp)
	this.numbers = append(this.numbers, number)
	this.nextFlush = append(this.nextFlush, 0)
}

func (this *AckList) Clear(una uint32) {
	count := 0
	for i := 0; i < len(this.numbers); i++ {
		if this.numbers[i] < una {
			continue
		}
		if i != count {
			this.numbers[count] = this.numbers[i]
			this.timestamps[count] = this.timestamps[i]
			this.nextFlush[count] = this.nextFlush[i]
		}
		count++
	}
	if count < len(this.numbers) {
		this.numbers = this.numbers[:count]
		this.timestamps = this.timestamps[:count]
		this.nextFlush = this.nextFlush[:count]
	}
}

func (this *AckList) Flush(current uint32, rto uint32) {
	seg := NewAckSegment()
	for i := 0; i < len(this.numbers) && !seg.IsFull(); i++ {
		if this.nextFlush[i] > current {
			continue
		}
		seg.PutNumber(this.numbers[i])
		seg.PutTimestamp(this.timestamps[i])
		timeout := rto / 4
		if timeout < 20 {
			timeout = 20
		}
		this.nextFlush[i] = current + timeout
	}
	if seg.Count > 0 {
		this.writer.Write(seg)
		seg.Release()
	}
}

type ReceivingWorker struct {
	sync.RWMutex
	conn       *Connection
	leftOver   *alloc.Buffer
	window     *ReceivingWindow
	acklist    *AckList
	nextNumber uint32
	windowSize uint32
}

func NewReceivingWorker(kcp *Connection) *ReceivingWorker {
	worker := &ReceivingWorker{
		conn:       kcp,
		window:     NewReceivingWindow(kcp.Config.GetReceivingBufferSize()),
		windowSize: kcp.Config.GetReceivingInFlightSize(),
	}
	worker.acklist = NewAckList(worker)
	return worker
}

func (this *ReceivingWorker) Release() {
	this.leftOver.Release()
}

func (this *ReceivingWorker) ProcessSendingNext(number uint32) {
	this.Lock()
	defer this.Unlock()

	this.acklist.Clear(number)
}

func (this *ReceivingWorker) ProcessSegment(seg *DataSegment) {
	this.Lock()
	defer this.Unlock()

	number := seg.Number
	idx := number - this.nextNumber
	if idx >= this.windowSize {
		return
	}
	this.acklist.Clear(seg.SendingNext)
	this.acklist.Add(number, seg.Timestamp)

	if !this.window.Set(idx, seg) {
		seg.Release()
	}
}

func (this *ReceivingWorker) Read(b []byte) int {
	this.Lock()
	defer this.Unlock()

	total := 0
	if this.leftOver != nil {
		nBytes := copy(b, this.leftOver.Value)
		if nBytes < this.leftOver.Len() {
			this.leftOver.SliceFrom(nBytes)
			return nBytes
		}
		this.leftOver.Release()
		this.leftOver = nil
		total += nBytes
	}

	for total < len(b) {
		seg := this.window.RemoveFirst()
		if seg == nil {
			break
		}
		this.window.Advance()
		this.nextNumber++

		nBytes := copy(b[total:], seg.Data.Value)
		total += nBytes
		if nBytes < seg.Data.Len() {
			seg.Data.SliceFrom(nBytes)
			this.leftOver = seg.Data
			seg.Data = nil
			seg.Release()
			break
		}
		seg.Release()
	}
	return total
}

func (this *ReceivingWorker) Flush(current uint32) {
	this.Lock()
	defer this.Unlock()

	this.acklist.Flush(current, this.conn.roundTrip.Timeout())
}

func (this *ReceivingWorker) Write(seg Segment) {
	ackSeg := seg.(*AckSegment)
	ackSeg.Conv = this.conn.conv
	ackSeg.ReceivingNext = this.nextNumber
	ackSeg.ReceivingWindow = this.nextNumber + this.windowSize
	if this.conn.state == StateReadyToClose {
		ackSeg.Option = SegmentOptionClose
	}
	this.conn.output.Write(ackSeg)
}

func (this *ReceivingWorker) CloseRead() {
}

func (this *ReceivingWorker) UpdateNecessary() bool {
	return len(this.acklist.numbers) > 0
}
