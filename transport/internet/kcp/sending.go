package kcp

import (
	"sync"
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

func (v *SendingWindow) Release() {
	if v == nil {
		return
	}
	v.len = 0
	for _, seg := range v.data {
		seg.Release()
	}
}

func (v *SendingWindow) Len() int {
	return int(v.len)
}

func (v *SendingWindow) IsEmpty() bool {
	return v.len == 0
}

func (v *SendingWindow) Size() uint32 {
	return v.cap
}

func (v *SendingWindow) IsFull() bool {
	return v.len == v.cap
}

func (v *SendingWindow) Push(number uint32, data []byte) {
	pos := (v.start + v.len) % v.cap
	v.data[pos].SetData(data)
	v.data[pos].Number = number
	v.data[pos].timeout = 0
	v.data[pos].transmit = 0
	v.inuse[pos] = true
	if v.len > 0 {
		v.next[v.last] = pos
		v.prev[pos] = v.last
	}
	v.last = pos
	v.len++
}

func (v *SendingWindow) FirstNumber() uint32 {
	return v.data[v.start].Number
}

func (v *SendingWindow) Clear(una uint32) {
	for !v.IsEmpty() && v.data[v.start].Number < una {
		v.Remove(0)
	}
}

func (v *SendingWindow) Remove(idx uint32) bool {
	if v.IsEmpty() {
		return false
	}

	pos := (v.start + idx) % v.cap
	if !v.inuse[pos] {
		return false
	}
	v.inuse[pos] = false
	v.totalInFlightSize--
	if pos == v.start && pos == v.last {
		v.len = 0
		v.start = 0
		v.last = 0
	} else if pos == v.start {
		delta := v.next[pos] - v.start
		if v.next[pos] < v.start {
			delta = v.next[pos] + v.cap - v.start
		}
		v.start = v.next[pos]
		v.len -= delta
	} else if pos == v.last {
		v.last = v.prev[pos]
	} else {
		v.next[v.prev[pos]] = v.next[pos]
		v.prev[v.next[pos]] = v.prev[pos]
	}
	return true
}

func (v *SendingWindow) HandleFastAck(number uint32, rto uint32) {
	if v.IsEmpty() {
		return
	}

	v.Visit(func(seg *DataSegment) bool {
		if number == seg.Number || number-seg.Number > 0x7FFFFFFF {
			return false
		}

		if seg.transmit > 0 && seg.timeout > rto/3 {
			seg.timeout -= rto / 3
		}
		return true
	})
}

func (v *SendingWindow) Visit(visitor func(seg *DataSegment) bool) {
	if v.IsEmpty() {
		return
	}

	for i := v.start; ; i = v.next[i] {
		if !visitor(&v.data[i]) || i == v.last {
			break
		}
	}
}

func (v *SendingWindow) Flush(current uint32, rto uint32, maxInFlightSize uint32) {
	if v.IsEmpty() {
		return
	}

	var lost uint32
	var inFlightSize uint32

	v.Visit(func(segment *DataSegment) bool {
		if current-segment.timeout >= 0x7FFFFFFF {
			return true
		}
		if segment.transmit == 0 {
			// First time
			v.totalInFlightSize++
		} else {
			lost++
		}
		segment.timeout = current + rto

		segment.Timestamp = current
		segment.transmit++
		v.writer.Write(segment)
		inFlightSize++
		if inFlightSize >= maxInFlightSize {
			return false
		}
		return true
	})

	if v.onPacketLoss != nil && inFlightSize > 0 && v.totalInFlightSize != 0 {
		rate := lost * 100 / v.totalInFlightSize
		v.onPacketLoss(rate)
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
	v.window.Release()
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

// Private: Visible for testing.
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

// Private: Visible for testing.
func (v *SendingWorker) ProcessAck(number uint32) bool {
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
		removed := v.ProcessAck(number)
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

func (v *SendingWorker) Push(b []byte) int {
	nBytes := 0
	v.Lock()
	defer v.Unlock()

	for len(b) > 0 && !v.window.IsFull() {
		var size int
		if len(b) > int(v.conn.mss) {
			size = int(v.conn.mss)
		} else {
			size = len(b)
		}
		v.window.Push(v.nextNumber, b[:size])
		v.nextNumber++
		b = b[size:]
		nBytes += size
	}
	return nBytes
}

// Private: Visible for testing.
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
	defer v.Unlock()

	cwnd := v.firstUnacknowledged + v.conn.Config.GetSendingInFlightSize()
	if cwnd > v.remoteNextNumber {
		cwnd = v.remoteNextNumber
	}
	if v.conn.Config.Congestion && cwnd > v.firstUnacknowledged+v.controlWindow {
		cwnd = v.firstUnacknowledged + v.controlWindow
	}

	if !v.window.IsEmpty() {
		v.window.Flush(current, v.conn.roundTrip.Timeout(), cwnd)
	} else if v.firstUnacknowledgedUpdated {
		v.conn.Ping(current, CommandPing)
	}

	v.firstUnacknowledgedUpdated = false
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
