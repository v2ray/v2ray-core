package kcp

import (
	"sync"

	"v2ray.com/core/common/buf"
)

type SendingWindow struct {
	start uint32
	cap   uint32
	len   uint32
	last  uint32

	data  []DataSegment
	inuse []bool
	prev  []uint32
	next  []uint32

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
		data:         make([]DataSegment, size),
		prev:         make([]uint32, size),
		next:         make([]uint32, size),
		inuse:        make([]bool, size),
		writer:       writer,
		onPacketLoss: onPacketLoss,
	}
	return window
}

func (sw *SendingWindow) Release() {
	if sw == nil {
		return
	}
	sw.len = 0
	for _, seg := range sw.data {
		seg.Release()
	}
}

func (sw *SendingWindow) Len() int {
	return int(sw.len)
}

func (sw *SendingWindow) IsEmpty() bool {
	return sw.len == 0
}

func (sw *SendingWindow) Size() uint32 {
	return sw.cap
}

func (sw *SendingWindow) IsFull() bool {
	return sw.len == sw.cap
}

func (sw *SendingWindow) Push(number uint32) *buf.Buffer {
	pos := (sw.start + sw.len) % sw.cap
	sw.data[pos].Number = number
	sw.data[pos].timeout = 0
	sw.data[pos].transmit = 0
	sw.inuse[pos] = true
	if sw.len > 0 {
		sw.next[sw.last] = pos
		sw.prev[pos] = sw.last
	}
	sw.last = pos
	sw.len++
	return sw.data[pos].Data()
}

func (sw *SendingWindow) FirstNumber() uint32 {
	return sw.data[sw.start].Number
}

func (sw *SendingWindow) Clear(una uint32) {
	for !sw.IsEmpty() && sw.data[sw.start].Number < una {
		sw.Remove(0)
	}
}

func (sw *SendingWindow) Remove(idx uint32) bool {
	if sw.IsEmpty() {
		return false
	}

	pos := (sw.start + idx) % sw.cap
	if !sw.inuse[pos] {
		return false
	}
	sw.inuse[pos] = false
	sw.totalInFlightSize--
	if pos == sw.start && pos == sw.last {
		sw.len = 0
		sw.start = 0
		sw.last = 0
	} else if pos == sw.start {
		delta := sw.next[pos] - sw.start
		if sw.next[pos] < sw.start {
			delta = sw.next[pos] + sw.cap - sw.start
		}
		sw.start = sw.next[pos]
		sw.len -= delta
	} else if pos == sw.last {
		sw.last = sw.prev[pos]
	} else {
		sw.next[sw.prev[pos]] = sw.next[pos]
		sw.prev[sw.next[pos]] = sw.prev[pos]
	}
	return true
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

	for i := sw.start; ; i = sw.next[i] {
		if !visitor(&sw.data[i]) || i == sw.last {
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

type SendingWorker struct {
	sync.RWMutex
	conn                       *Connection
	window                     *SendingWindow
	firstUnacknowledged        uint32
	firstUnacknowledgedUpdated bool
	nextNumber                 uint32
	remoteNextNumber           uint32
	controlWindow              uint32
	fastResend                 uint32
}

func NewSendingWorker(kcp *Connection) *SendingWorker {
	worker := &SendingWorker{
		conn:             kcp,
		fastResend:       2,
		remoteNextNumber: 32,
		controlWindow:    kcp.Config.GetSendingInFlightSize(),
	}
	worker.window = NewSendingWindow(kcp.Config.GetSendingBufferSize(), worker, worker.OnPacketLoss)
	return worker
}

func (v *SendingWorker) Release() {
	v.Lock()
	v.window.Release()
	v.Unlock()
}

func (v *SendingWorker) ProcessReceivingNext(nextNumber uint32) {
	v.Lock()
	defer v.Unlock()

	v.ProcessReceivingNextWithoutLock(nextNumber)
}

func (v *SendingWorker) ProcessReceivingNextWithoutLock(nextNumber uint32) {
	v.window.Clear(nextNumber)
	v.FindFirstUnacknowledged()
}

func (v *SendingWorker) FindFirstUnacknowledged() {
	first := v.firstUnacknowledged
	if !v.window.IsEmpty() {
		v.firstUnacknowledged = v.window.FirstNumber()
	} else {
		v.firstUnacknowledged = v.nextNumber
	}
	if first != v.firstUnacknowledged {
		v.firstUnacknowledgedUpdated = true
	}
}

func (v *SendingWorker) processAck(number uint32) bool {
	// number < v.firstUnacknowledged || number >= v.nextNumber
	if number-v.firstUnacknowledged > 0x7FFFFFFF || number-v.nextNumber < 0x7FFFFFFF {
		return false
	}

	removed := v.window.Remove(number - v.firstUnacknowledged)
	if removed {
		v.FindFirstUnacknowledged()
	}
	return removed
}

func (v *SendingWorker) ProcessSegment(current uint32, seg *AckSegment, rto uint32) {
	defer seg.Release()

	v.Lock()
	defer v.Unlock()

	if v.remoteNextNumber < seg.ReceivingWindow {
		v.remoteNextNumber = seg.ReceivingWindow
	}
	v.ProcessReceivingNextWithoutLock(seg.ReceivingNext)

	if seg.IsEmpty() {
		return
	}

	var maxack uint32
	var maxackRemoved bool
	for _, number := range seg.NumberList {
		removed := v.processAck(number)
		if maxack < number {
			maxack = number
			maxackRemoved = removed
		}
	}

	if maxackRemoved {
		v.window.HandleFastAck(maxack, rto)
		if current-seg.Timestamp < 10000 {
			v.conn.roundTrip.Update(current-seg.Timestamp, current)
		}
	}
}

func (v *SendingWorker) Push() *buf.Buffer {
	v.Lock()
	defer v.Unlock()

	if v.window.IsFull() {
		return nil
	}

	b := v.window.Push(v.nextNumber)
	v.nextNumber++
	return b
}

func (v *SendingWorker) Write(seg Segment) error {
	dataSeg := seg.(*DataSegment)

	dataSeg.Conv = v.conn.conv
	dataSeg.SendingNext = v.firstUnacknowledged
	dataSeg.Option = 0
	if v.conn.State() == StateReadyToClose {
		dataSeg.Option = SegmentOptionClose
	}

	return v.conn.output.Write(dataSeg)
}

func (v *SendingWorker) OnPacketLoss(lossRate uint32) {
	if !v.conn.Config.Congestion || v.conn.roundTrip.Timeout() == 0 {
		return
	}

	if lossRate >= 15 {
		v.controlWindow = 3 * v.controlWindow / 4
	} else if lossRate <= 5 {
		v.controlWindow += v.controlWindow / 4
	}
	if v.controlWindow < 16 {
		v.controlWindow = 16
	}
	if v.controlWindow > 2*v.conn.Config.GetSendingInFlightSize() {
		v.controlWindow = 2 * v.conn.Config.GetSendingInFlightSize()
	}
}

func (v *SendingWorker) Flush(current uint32) {
	v.Lock()

	cwnd := v.firstUnacknowledged + v.conn.Config.GetSendingInFlightSize()
	if cwnd > v.remoteNextNumber {
		cwnd = v.remoteNextNumber
	}
	if v.conn.Config.Congestion && cwnd > v.firstUnacknowledged+v.controlWindow {
		cwnd = v.firstUnacknowledged + v.controlWindow
	}

	if !v.window.IsEmpty() {
		v.window.Flush(current, v.conn.roundTrip.Timeout(), cwnd)
		v.firstUnacknowledgedUpdated = false
	}

	updated := v.firstUnacknowledgedUpdated
	v.firstUnacknowledgedUpdated = false

	v.Unlock()

	if updated {
		v.conn.Ping(current, CommandPing)
	}
}

func (v *SendingWorker) CloseWrite() {
	v.Lock()
	defer v.Unlock()

	v.window.Clear(0xFFFFFFFF)
}

func (v *SendingWorker) IsEmpty() bool {
	v.RLock()
	defer v.RUnlock()

	return v.window.IsEmpty()
}

func (v *SendingWorker) UpdateNecessary() bool {
	return !v.IsEmpty()
}

func (w *SendingWorker) FirstUnacknowledged() uint32 {
	w.RLock()
	defer w.RUnlock()

	return w.firstUnacknowledged
}
