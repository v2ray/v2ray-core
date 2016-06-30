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
	sync.RWMutex
	closed  bool
	cache   *alloc.Buffer
	queue   chan *alloc.Buffer
	timeout time.Time
}

func NewReceivingQueue() *ReceivingQueue {
	return &ReceivingQueue{
		queue: make(chan *alloc.Buffer, effectiveConfig.ReadBuffer/effectiveConfig.Mtu),
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
			this.RLock()
			if !this.timeout.IsZero() && this.timeout.Before(time.Now()) {
				this.RUnlock()
				return totalBytes, errTimeout
			}
			this.RUnlock()
			timeToSleep += 500 * time.Millisecond
		}
	}

	return totalBytes, nil
}

func (this *ReceivingQueue) Put(payload *alloc.Buffer) {
	this.RLock()
	defer this.RUnlock()

	if this.closed {
		payload.Release()
		return
	}

	this.queue <- payload
}

func (this *ReceivingQueue) SetReadDeadline(t time.Time) error {
	this.Lock()
	defer this.Unlock()

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

type ACKList struct {
	kcp        *KCP
	timestamps []uint32
	numbers    []uint32
	nextFlush  []uint32
}

func NewACKList(kcp *KCP) *ACKList {
	return &ACKList{
		kcp:        kcp,
		timestamps: make([]uint32, 0, 32),
		numbers:    make([]uint32, 0, 32),
		nextFlush:  make([]uint32, 0, 32),
	}
}

func (this *ACKList) Add(number uint32, timestamp uint32) {
	this.timestamps = append(this.timestamps, timestamp)
	this.numbers = append(this.numbers, number)
	this.nextFlush = append(this.nextFlush, 0)
}

func (this *ACKList) Clear(una uint32) {
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

func (this *ACKList) Flush() bool {
	seg := &ACKSegment{
		Conv:            this.kcp.conv,
		ReceivingNext:   this.kcp.rcv_nxt,
		ReceivingWindow: this.kcp.rcv_nxt + this.kcp.rcv_wnd,
	}
	if this.kcp.state == StateReadyToClose {
		seg.Opt = SegmentOptionClose
	}
	current := this.kcp.current
	for i := 0; i < len(this.numbers); i++ {
		if this.nextFlush[i] <= current {
			seg.Count++
			seg.NumberList = append(seg.NumberList, this.numbers[i])
			seg.TimestampList = append(seg.TimestampList, this.timestamps[i])
			this.nextFlush[i] = current + 100
			if seg.Count == 128 {
				break
			}
		}
	}
	if seg.Count > 0 {
		this.kcp.output.Write(seg)
		return true
	}
	return false
}
