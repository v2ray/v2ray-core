package kcp

import (
	"io"
	"sync"
	"time"

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
	sync.Mutex
	closed  bool
	cache   *alloc.Buffer
	queue   chan *alloc.Buffer
	timeout time.Time
}

func NewReceivingQueue() *ReceivingQueue {
	return &ReceivingQueue{
		queue: make(chan *alloc.Buffer, effectiveConfig.GetReceivingQueueSize()),
	}
}

func (this *ReceivingQueue) Read(buf []byte) (int, error) {
	if this.cache.Len() > 0 {
		nBytes, err := this.cache.Read(buf)
		if this.cache.IsEmpty() {
			this.cache.Release()
			this.cache = nil
		}
		return nBytes, err
	}

	var totalBytes int

L:
	for totalBytes < len(buf) {
		timeToSleep := time.Millisecond
		select {
		case payload, open := <-this.queue:
			if !open {
				return totalBytes, io.EOF
			}
			nBytes, err := payload.Read(buf)
			totalBytes += nBytes
			if err != nil {
				return totalBytes, err
			}
			if !payload.IsEmpty() {
				this.cache = payload
			}
			buf = buf[nBytes:]
		case <-time.After(timeToSleep):
			if totalBytes > 0 {
				break L
			}
			if !this.timeout.IsZero() && this.timeout.Before(time.Now()) {
				return totalBytes, errTimeout
			}
			timeToSleep += 500 * time.Millisecond
		}
	}

	return totalBytes, nil
}

func (this *ReceivingQueue) Put(payload *alloc.Buffer) bool {
	this.Lock()
	defer this.Unlock()

	if this.closed {
		payload.Release()
		return false
	}

	select {
	case this.queue <- payload:
		return true
	default:
		return false
	}
}

func (this *ReceivingQueue) SetReadDeadline(t time.Time) error {
	this.timeout = t
	return nil
}

func (this *ReceivingQueue) Close() {
	this.Lock()
	defer this.Unlock()

	if this.closed {
		return
	}
	this.closed = true
	close(this.queue)
}

type AckList struct {
	sync.Mutex
	writer     SegmentWriter
	timestamps []uint32
	numbers    []uint32
	nextFlush  []uint32
}

func NewACKList(writer SegmentWriter) *AckList {
	return &AckList{
		writer:     writer,
		timestamps: make([]uint32, 0, 32),
		numbers:    make([]uint32, 0, 32),
		nextFlush:  make([]uint32, 0, 32),
	}
}

func (this *AckList) Add(number uint32, timestamp uint32) {
	this.Lock()
	defer this.Unlock()

	this.timestamps = append(this.timestamps, timestamp)
	this.numbers = append(this.numbers, number)
	this.nextFlush = append(this.nextFlush, 0)
}

func (this *AckList) Clear(una uint32) {
	this.Lock()
	defer this.Unlock()

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
	seg := new(AckSegment)
	this.Lock()
	for i := 0; i < len(this.numbers); i++ {
		if this.nextFlush[i] <= current {
			seg.Count++
			seg.NumberList = append(seg.NumberList, this.numbers[i])
			seg.TimestampList = append(seg.TimestampList, this.timestamps[i])
			this.nextFlush[i] = current + rto/2
			if seg.Count == 128 {
				break
			}
		}
	}
	this.Unlock()
	if seg.Count > 0 {
		this.writer.Write(seg)
	}
}

type ReceivingWorker struct {
	kcp         *KCP
	queue       *ReceivingQueue
	window      *ReceivingWindow
	windowMutex sync.Mutex
	acklist     *AckList
	updated     bool
	nextNumber  uint32
	windowSize  uint32
}

func NewReceivingWorker(kcp *KCP) *ReceivingWorker {
	windowSize := effectiveConfig.GetReceivingWindowSize()
	worker := &ReceivingWorker{
		kcp:        kcp,
		queue:      NewReceivingQueue(),
		window:     NewReceivingWindow(windowSize),
		windowSize: windowSize,
	}
	worker.acklist = NewACKList(worker)
	return worker
}

func (this *ReceivingWorker) ProcessSendingNext(number uint32) {
	this.acklist.Clear(number)
}

func (this *ReceivingWorker) ProcessSegment(seg *DataSegment) {
	number := seg.Number
	if _itimediff(number, this.nextNumber+this.windowSize) >= 0 || _itimediff(number, this.nextNumber) < 0 {
		return
	}

	this.ProcessSendingNext(seg.SendingNext)

	this.acklist.Add(number, seg.Timestamp)
	this.windowMutex.Lock()
	idx := number - this.nextNumber

	if !this.window.Set(idx, seg) {
		seg.Release()
	}
	this.windowMutex.Unlock()

	this.DumpWindow()
}

// @Private
func (this *ReceivingWorker) DumpWindow() {
	this.windowMutex.Lock()
	defer this.windowMutex.Unlock()

	for {
		seg := this.window.RemoveFirst()
		if seg == nil {
			break
		}

		if !this.queue.Put(seg.Data) {
			this.window.Set(0, seg)
			break
		}

		seg.Data = nil
		this.window.Advance()
		this.nextNumber++
		this.updated = true
	}
}

func (this *ReceivingWorker) Read(b []byte) (int, error) {
	return this.queue.Read(b)
}

func (this *ReceivingWorker) SetReadDeadline(t time.Time) {
	this.queue.SetReadDeadline(t)
}

func (this *ReceivingWorker) Flush() {
	this.acklist.Flush(this.kcp.current, this.kcp.rx_rto)
}

func (this *ReceivingWorker) Write(seg ISegment) {
	ackSeg := seg.(*AckSegment)
	ackSeg.Conv = this.kcp.conv
	ackSeg.ReceivingNext = this.nextNumber
	ackSeg.ReceivingWindow = this.nextNumber + this.windowSize
	if this.kcp.state == StateReadyToClose {
		ackSeg.Opt = SegmentOptionClose
	}
	this.kcp.output.Write(ackSeg)
	this.updated = false
}

func (this *ReceivingWorker) CloseRead() {
	this.queue.Close()
}

func (this *ReceivingWorker) PingNecessary() bool {
	return this.updated
}
