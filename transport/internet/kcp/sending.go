// +build !confonly

package kcp

import (
	"container/list"
	"sync"

	"v2ray.com/core/common/buf"
)

type SendingWindow struct {
	cache             *list.List
	totalInFlightSize uint32
	writer            SegmentWriter
	onPacketLoss      func(uint32)
}

func NewSendingWindow(writer SegmentWriter, onPacketLoss func(uint32)) *SendingWindow {
	window := &SendingWindow{
		cache:        list.New(),
		writer:       writer,
		onPacketLoss: onPacketLoss,
	}
	return window
}

func (sw *SendingWindow) Release() {
	if sw == nil {
		return
	}
	for sw.cache.Len() > 0 {
		seg := sw.cache.Front().Value.(*DataSegment)
		seg.Release()
		sw.cache.Remove(sw.cache.Front())
	}
}

func (sw *SendingWindow) Len() uint32 {
	return uint32(sw.cache.Len())
}

func (sw *SendingWindow) IsEmpty() bool {
	return sw.cache.Len() == 0
}

func (sw *SendingWindow) Push(number uint32, b *buf.Buffer) {
	seg := NewDataSegment()
	seg.Number = number
	seg.payload = b

	sw.cache.PushBack(seg)
}

func (sw *SendingWindow) FirstNumber() uint32 {
	return sw.cache.Front().Value.(*DataSegment).Number
}

func (sw *SendingWindow) Clear(una uint32) {
	for !sw.IsEmpty() {
		seg := sw.cache.Front().Value.(*DataSegment)
		if seg.Number >= una {
			break
		}
		seg.Release()
		sw.cache.Remove(sw.cache.Front())
	}
}

func (sw *SendingWindow) HandleFastAck(number uint32, rto uint32) {
	if sw.IsEmpty() {
		return
	}

	sw.Visit(func(seg *DataSegment) bool {
		if number == seg.Number || number-seg.Number > 0x7FFFFFFF {
			return false
		}

		if seg.transmit > 0 && seg.timeout > rto/3 {
			seg.timeout -= rto / 3
		}
		return true
	})
}

func (sw *SendingWindow) Visit(visitor func(seg *DataSegment) bool) {
	if sw.IsEmpty() {
		return
	}

	for e := sw.cache.Front(); e != nil; e = e.Next() {
		seg := e.Value.(*DataSegment)
		if !visitor(seg) {
			break
		}
	}
}

func (sw *SendingWindow) Flush(current uint32, rto uint32, maxInFlightSize uint32) {
	if sw.IsEmpty() {
		return
	}

	var lost uint32
	var inFlightSize uint32

	sw.Visit(func(segment *DataSegment) bool {
		if current-segment.timeout >= 0x7FFFFFFF {
			return true
		}
		if segment.transmit == 0 {
			// First time
			sw.totalInFlightSize++
		} else {
			lost++
		}
		segment.timeout = current + rto

		segment.Timestamp = current
		segment.transmit++
		sw.writer.Write(segment)
		inFlightSize++
		if inFlightSize >= maxInFlightSize {
			return false
		}
		return true
	})

	if sw.onPacketLoss != nil && inFlightSize > 0 && sw.totalInFlightSize != 0 {
		rate := lost * 100 / sw.totalInFlightSize
		sw.onPacketLoss(rate)
	}
}

func (sw *SendingWindow) Remove(number uint32) bool {
	if sw.IsEmpty() {
		return false
	}

	for e := sw.cache.Front(); e != nil; e = e.Next() {
		seg := e.Value.(*DataSegment)
		if seg.Number > number {
			return false
		} else if seg.Number == number {
			if sw.totalInFlightSize > 0 {
				sw.totalInFlightSize--
			}
			seg.Release()
			sw.cache.Remove(e)
			return true
		}
	}

	return false
}

type SendingWorker struct {
	sync.RWMutex
	conn                       *Connection
	window                     *SendingWindow
	firstUnacknowledged        uint32
	nextNumber                 uint32
	remoteNextNumber           uint32
	controlWindow              uint32
	fastResend                 uint32
	windowSize                 uint32
	firstUnacknowledgedUpdated bool
	closed                     bool
}

func NewSendingWorker(kcp *Connection) *SendingWorker {
	worker := &SendingWorker{
		conn:             kcp,
		fastResend:       2,
		remoteNextNumber: 32,
		controlWindow:    kcp.Config.GetSendingInFlightSize(),
		windowSize:       kcp.Config.GetSendingBufferSize(),
	}
	worker.window = NewSendingWindow(worker, worker.OnPacketLoss)
	return worker
}

func (w *SendingWorker) Release() {
	w.Lock()
	w.window.Release()
	w.closed = true
	w.Unlock()
}

func (w *SendingWorker) ProcessReceivingNext(nextNumber uint32) {
	w.Lock()
	defer w.Unlock()

	w.ProcessReceivingNextWithoutLock(nextNumber)
}

func (w *SendingWorker) ProcessReceivingNextWithoutLock(nextNumber uint32) {
	w.window.Clear(nextNumber)
	w.FindFirstUnacknowledged()
}

func (w *SendingWorker) FindFirstUnacknowledged() {
	first := w.firstUnacknowledged
	if !w.window.IsEmpty() {
		w.firstUnacknowledged = w.window.FirstNumber()
	} else {
		w.firstUnacknowledged = w.nextNumber
	}
	if first != w.firstUnacknowledged {
		w.firstUnacknowledgedUpdated = true
	}
}

func (w *SendingWorker) processAck(number uint32) bool {
	// number < v.firstUnacknowledged || number >= v.nextNumber
	if number-w.firstUnacknowledged > 0x7FFFFFFF || number-w.nextNumber < 0x7FFFFFFF {
		return false
	}

	removed := w.window.Remove(number)
	if removed {
		w.FindFirstUnacknowledged()
	}
	return removed
}

func (w *SendingWorker) ProcessSegment(current uint32, seg *AckSegment, rto uint32) {
	defer seg.Release()

	w.Lock()
	defer w.Unlock()

	if w.closed {
		return
	}

	if w.remoteNextNumber < seg.ReceivingWindow {
		w.remoteNextNumber = seg.ReceivingWindow
	}
	w.ProcessReceivingNextWithoutLock(seg.ReceivingNext)

	if seg.IsEmpty() {
		return
	}

	var maxack uint32
	var maxackRemoved bool
	for _, number := range seg.NumberList {
		removed := w.processAck(number)
		if maxack < number {
			maxack = number
			maxackRemoved = removed
		}
	}

	if maxackRemoved {
		w.window.HandleFastAck(maxack, rto)
		if current-seg.Timestamp < 10000 {
			w.conn.roundTrip.Update(current-seg.Timestamp, current)
		}
	}
}

func (w *SendingWorker) Push(b *buf.Buffer) bool {
	w.Lock()
	defer w.Unlock()

	if w.closed {
		return false
	}

	if w.window.Len() > w.windowSize {
		return false
	}

	w.window.Push(w.nextNumber, b)
	w.nextNumber++
	return true
}

func (w *SendingWorker) Write(seg Segment) error {
	dataSeg := seg.(*DataSegment)

	dataSeg.Conv = w.conn.meta.Conversation
	dataSeg.SendingNext = w.firstUnacknowledged
	dataSeg.Option = 0
	if w.conn.State() == StateReadyToClose {
		dataSeg.Option = SegmentOptionClose
	}

	return w.conn.output.Write(dataSeg)
}

func (w *SendingWorker) OnPacketLoss(lossRate uint32) {
	if !w.conn.Config.Congestion || w.conn.roundTrip.Timeout() == 0 {
		return
	}

	if lossRate >= 15 {
		w.controlWindow = 3 * w.controlWindow / 4
	} else if lossRate <= 5 {
		w.controlWindow += w.controlWindow / 4
	}
	if w.controlWindow < 16 {
		w.controlWindow = 16
	}
	if w.controlWindow > 2*w.conn.Config.GetSendingInFlightSize() {
		w.controlWindow = 2 * w.conn.Config.GetSendingInFlightSize()
	}
}

func (w *SendingWorker) Flush(current uint32) {
	w.Lock()

	if w.closed {
		w.Unlock()
		return
	}

	cwnd := w.firstUnacknowledged + w.conn.Config.GetSendingInFlightSize()
	if cwnd > w.remoteNextNumber {
		cwnd = w.remoteNextNumber
	}
	if w.conn.Config.Congestion && cwnd > w.firstUnacknowledged+w.controlWindow {
		cwnd = w.firstUnacknowledged + w.controlWindow
	}

	if !w.window.IsEmpty() {
		w.window.Flush(current, w.conn.roundTrip.Timeout(), cwnd)
		w.firstUnacknowledgedUpdated = false
	}

	updated := w.firstUnacknowledgedUpdated
	w.firstUnacknowledgedUpdated = false

	w.Unlock()

	if updated {
		w.conn.Ping(current, CommandPing)
	}
}

func (w *SendingWorker) CloseWrite() {
	w.Lock()
	defer w.Unlock()

	w.window.Clear(0xFFFFFFFF)
}

func (w *SendingWorker) IsEmpty() bool {
	w.RLock()
	defer w.RUnlock()

	return w.window.IsEmpty()
}

func (w *SendingWorker) UpdateNecessary() bool {
	return !w.IsEmpty()
}

func (w *SendingWorker) FirstUnacknowledged() uint32 {
	w.RLock()
	defer w.RUnlock()

	return w.firstUnacknowledged
}
