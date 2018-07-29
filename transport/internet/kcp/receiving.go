package kcp

import (
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

type ReceivingWindow struct {
	cache map[uint32]*DataSegment
}

func NewReceivingWindow() *ReceivingWindow {
	return &ReceivingWindow{
		cache: make(map[uint32]*DataSegment),
	}
}

func (w *ReceivingWindow) Set(id uint32, value *DataSegment) bool {
	_, f := w.cache[id]
	if f {
		return false
	}
	w.cache[id] = value
	return true
}

func (w *ReceivingWindow) Has(id uint32) bool {
	_, f := w.cache[id]
	return f
}

func (w *ReceivingWindow) Remove(id uint32) *DataSegment {
	v, f := w.cache[id]
	if !f {
		return nil
	}
	delete(w.cache, id)
	return v
}

type AckList struct {
	writer     SegmentWriter
	timestamps []uint32
	numbers    []uint32
	nextFlush  []uint32

	flushCandidates []uint32
	dirty           bool
}

func NewAckList(writer SegmentWriter) *AckList {
	return &AckList{
		writer:          writer,
		timestamps:      make([]uint32, 0, 128),
		numbers:         make([]uint32, 0, 128),
		nextFlush:       make([]uint32, 0, 128),
		flushCandidates: make([]uint32, 0, 128),
	}
}

func (l *AckList) Add(number uint32, timestamp uint32) {
	l.timestamps = append(l.timestamps, timestamp)
	l.numbers = append(l.numbers, number)
	l.nextFlush = append(l.nextFlush, 0)
	l.dirty = true
}

func (l *AckList) Clear(una uint32) {
	count := 0
	for i := 0; i < len(l.numbers); i++ {
		if l.numbers[i] < una {
			continue
		}
		if i != count {
			l.numbers[count] = l.numbers[i]
			l.timestamps[count] = l.timestamps[i]
			l.nextFlush[count] = l.nextFlush[i]
		}
		count++
	}
	if count < len(l.numbers) {
		l.numbers = l.numbers[:count]
		l.timestamps = l.timestamps[:count]
		l.nextFlush = l.nextFlush[:count]
		l.dirty = true
	}
}

func (l *AckList) Flush(current uint32, rto uint32) {
	l.flushCandidates = l.flushCandidates[:0]

	seg := NewAckSegment()
	for i := 0; i < len(l.numbers); i++ {
		if l.nextFlush[i] > current {
			if len(l.flushCandidates) < cap(l.flushCandidates) {
				l.flushCandidates = append(l.flushCandidates, l.numbers[i])
			}
			continue
		}
		seg.PutNumber(l.numbers[i])
		seg.PutTimestamp(l.timestamps[i])
		timeout := rto / 2
		if timeout < 20 {
			timeout = 20
		}
		l.nextFlush[i] = current + timeout

		if seg.IsFull() {
			l.writer.Write(seg)
			seg.Release()
			seg = NewAckSegment()
			l.dirty = false
		}
	}

	if l.dirty || !seg.IsEmpty() {
		for _, number := range l.flushCandidates {
			if seg.IsFull() {
				break
			}
			seg.PutNumber(number)
		}
		l.writer.Write(seg)
		l.dirty = false
	}

	seg.Release()
}

type ReceivingWorker struct {
	sync.RWMutex
	conn       *Connection
	leftOver   buf.MultiBuffer
	window     *ReceivingWindow
	acklist    *AckList
	nextNumber uint32
	windowSize uint32
}

func NewReceivingWorker(kcp *Connection) *ReceivingWorker {
	worker := &ReceivingWorker{
		conn:       kcp,
		window:     NewReceivingWindow(),
		windowSize: kcp.Config.GetReceivingInFlightSize(),
	}
	worker.acklist = NewAckList(worker)
	return worker
}

func (w *ReceivingWorker) Release() {
	w.Lock()
	w.leftOver.Release()
	w.Unlock()
}

func (w *ReceivingWorker) ProcessSendingNext(number uint32) {
	w.Lock()
	defer w.Unlock()

	w.acklist.Clear(number)
}

func (w *ReceivingWorker) ProcessSegment(seg *DataSegment) {
	w.Lock()
	defer w.Unlock()

	number := seg.Number
	idx := number - w.nextNumber
	if idx >= w.windowSize {
		return
	}
	w.acklist.Clear(seg.SendingNext)
	w.acklist.Add(number, seg.Timestamp)

	if !w.window.Set(seg.Number, seg) {
		seg.Release()
	}
}

func (w *ReceivingWorker) ReadMultiBuffer() buf.MultiBuffer {
	if w.leftOver != nil {
		mb := w.leftOver
		w.leftOver = nil
		return mb
	}

	mb := buf.NewMultiBufferCap(32)

	w.Lock()
	defer w.Unlock()
	for {
		seg := w.window.Remove(w.nextNumber)
		if seg == nil {
			break
		}
		w.nextNumber++
		mb.Append(seg.Detach())
		seg.Release()
	}

	return mb
}

func (w *ReceivingWorker) Read(b []byte) int {
	mb := w.ReadMultiBuffer()
	nBytes, err := mb.Read(b)
	common.Must(err)
	if !mb.IsEmpty() {
		w.leftOver = mb
	}
	return nBytes
}

func (w *ReceivingWorker) IsDataAvailable() bool {
	w.RLock()
	defer w.RUnlock()
	return w.window.Has(w.nextNumber)
}

func (w *ReceivingWorker) NextNumber() uint32 {
	w.RLock()
	defer w.RUnlock()

	return w.nextNumber
}

func (w *ReceivingWorker) Flush(current uint32) {
	w.Lock()
	defer w.Unlock()

	w.acklist.Flush(current, w.conn.roundTrip.Timeout())
}

func (w *ReceivingWorker) Write(seg Segment) error {
	ackSeg := seg.(*AckSegment)
	ackSeg.Conv = w.conn.meta.Conversation
	ackSeg.ReceivingNext = w.nextNumber
	ackSeg.ReceivingWindow = w.nextNumber + w.windowSize
	ackSeg.Option = 0
	if w.conn.State() == StateReadyToClose {
		ackSeg.Option = SegmentOptionClose
	}
	return w.conn.output.Write(ackSeg)
}

func (*ReceivingWorker) CloseRead() {
}

func (w *ReceivingWorker) UpdateNecessary() bool {
	w.RLock()
	defer w.RUnlock()

	return len(w.acklist.numbers) > 0
}
