package kcp

import (
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
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

type ReceivingQueue struct {
	start uint32
	cap   uint32
	len   uint32
	data  []*alloc.Buffer
}

func NewReceivingQueue(size uint32) *ReceivingQueue {
	return &ReceivingQueue{
		cap:  size,
		data: make([]*alloc.Buffer, size),
	}
}

func (this *ReceivingQueue) IsEmpty() bool {
	return this.len == 0
}

func (this *ReceivingQueue) IsFull() bool {
	return this.len == this.cap
}

func (this *ReceivingQueue) Read(buf []byte) int {
	if this.IsEmpty() {
		return 0
	}

	totalBytes := 0
	lenBuf := len(buf)
	for !this.IsEmpty() && totalBytes < lenBuf {
		payload := this.data[this.start]
		nBytes, _ := payload.Read(buf)
		buf = buf[nBytes:]
		totalBytes += nBytes
		if payload.IsEmpty() {
			payload.Release()
			this.data[this.start] = nil
			this.start++
			if this.start == this.cap {
				this.start = 0
			}
			this.len--
			if this.len == 0 {
				this.start = 0
			}
		}
	}
	return totalBytes
}

func (this *ReceivingQueue) Put(payload *alloc.Buffer) {
	this.data[(this.start+this.len)%this.cap] = payload
	this.len++
}

func (this *ReceivingQueue) Close() {
	for i := uint32(0); i < this.len; i++ {
		this.data[(this.start+i)%this.cap].Release()
		this.data[(this.start+i)%this.cap] = nil
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
		if this.numbers[i] >= una {
			if i != count {
				this.numbers[count] = this.numbers[i]
				this.timestamps[count] = this.timestamps[i]
				this.nextFlush[count] = this.nextFlush[i]
			}
			count++
		}
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
		if this.nextFlush[i] <= current {
			seg.PutNumber(this.numbers[i], this.timestamps[i])
			this.nextFlush[i] = current + rto/2
		}
	}
	if seg.Count > 0 {
		this.writer.Write(seg)
		seg.Release()
	}
}

type ReceivingWorker struct {
	sync.RWMutex
	conn       *Connection
	queue      *ReceivingQueue
	window     *ReceivingWindow
	acklist    *AckList
	updated    bool
	nextNumber uint32
	windowSize uint32
}

func NewReceivingWorker(kcp *Connection) *ReceivingWorker {
	windowSize := effectiveConfig.GetReceivingWindowSize()
	worker := &ReceivingWorker{
		conn:       kcp,
		queue:      NewReceivingQueue(effectiveConfig.GetReceivingQueueSize()),
		window:     NewReceivingWindow(windowSize),
		windowSize: windowSize,
	}
	worker.acklist = NewAckList(worker)
	return worker
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

	for !this.queue.IsFull() {
		seg := this.window.RemoveFirst()
		if seg == nil {
			break
		}

		this.queue.Put(seg.Data)
		seg.Data = nil
		seg.Release()
		this.window.Advance()
		this.nextNumber++
		this.updated = true
	}
}

func (this *ReceivingWorker) Read(b []byte) int {
	this.Lock()
	defer this.Unlock()

	return this.queue.Read(b)
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
		ackSeg.Opt = SegmentOptionClose
	}
	this.conn.output.Write(ackSeg)
	this.updated = false
}

func (this *ReceivingWorker) CloseRead() {
	this.Lock()
	defer this.Unlock()

	this.queue.Close()
}

func (this *ReceivingWorker) PingNecessary() bool {
	this.RLock()
	defer this.RUnlock()
	return this.updated
}

func (this *ReceivingWorker) MarkPingNecessary(b bool) {
	this.Lock()
	defer this.Unlock()
	this.updated = b
}
