package kcp

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/predicate"
)

var (
	ErrIOTimeout        = newError("Read/Write timeout")
	ErrClosedListener   = newError("Listener closed.")
	ErrClosedConnection = newError("Connection closed.")
)

type State int32

func (s State) Is(states ...State) bool {
	for _, state := range states {
		if s == state {
			return true
		}
	}
	return false
}

const (
	StateActive          State = 0
	StateReadyToClose    State = 1
	StatePeerClosed      State = 2
	StateTerminating     State = 3
	StatePeerTerminating State = 4
	StateTerminated      State = 5
)

func nowMillisec() int64 {
	now := time.Now()
	return now.Unix()*1000 + int64(now.Nanosecond()/1000000)
}

type RoundTripInfo struct {
	sync.RWMutex
	variation        uint32
	srtt             uint32
	rto              uint32
	minRtt           uint32
	updatedTimestamp uint32
}

func (v *RoundTripInfo) UpdatePeerRTO(rto uint32, current uint32) {
	v.Lock()
	defer v.Unlock()

	if current-v.updatedTimestamp < 3000 {
		return
	}

	v.updatedTimestamp = current
	v.rto = rto
}

func (v *RoundTripInfo) Update(rtt uint32, current uint32) {
	if rtt > 0x7FFFFFFF {
		return
	}
	v.Lock()
	defer v.Unlock()

	// https://tools.ietf.org/html/rfc6298
	if v.srtt == 0 {
		v.srtt = rtt
		v.variation = rtt / 2
	} else {
		delta := rtt - v.srtt
		if v.srtt > rtt {
			delta = v.srtt - rtt
		}
		v.variation = (3*v.variation + delta) / 4
		v.srtt = (7*v.srtt + rtt) / 8
		if v.srtt < v.minRtt {
			v.srtt = v.minRtt
		}
	}
	var rto uint32
	if v.minRtt < 4*v.variation {
		rto = v.srtt + 4*v.variation
	} else {
		rto = v.srtt + v.variation
	}

	if rto > 10000 {
		rto = 10000
	}
	v.rto = rto * 5 / 4
	v.updatedTimestamp = current
}

func (v *RoundTripInfo) Timeout() uint32 {
	v.RLock()
	defer v.RUnlock()

	return v.rto
}

func (v *RoundTripInfo) SmoothedTime() uint32 {
	v.RLock()
	defer v.RUnlock()

	return v.srtt
}

type Updater struct {
	interval        int64
	shouldContinue  predicate.Predicate
	shouldTerminate predicate.Predicate
	updateFunc      func()
	notifier        chan bool
}

func NewUpdater(interval uint32, shouldContinue predicate.Predicate, shouldTerminate predicate.Predicate, updateFunc func()) *Updater {
	u := &Updater{
		interval:        int64(time.Duration(interval) * time.Millisecond),
		shouldContinue:  shouldContinue,
		shouldTerminate: shouldTerminate,
		updateFunc:      updateFunc,
		notifier:        make(chan bool, 1),
	}
	go u.Run()
	return u
}

func (v *Updater) WakeUp() {
	select {
	case v.notifier <- true:
	default:
	}
}

func (v *Updater) Run() {
	for <-v.notifier {
		if v.shouldTerminate() {
			return
		}
		interval := v.Interval()
		for v.shouldContinue() {
			v.updateFunc()
			time.Sleep(interval)
		}
	}
}

func (u *Updater) Interval() time.Duration {
	return time.Duration(atomic.LoadInt64(&u.interval))
}

func (u *Updater) SetInterval(d time.Duration) {
	atomic.StoreInt64(&u.interval, int64(d))
}

type SystemConnection interface {
	net.Conn
	Reset(func([]Segment))
	Overhead() int
}

var (
	_ buf.Reader = (*Connection)(nil)
)

// Connection is a KCP connection over UDP.
type Connection struct {
	conn       SystemConnection
	rd         time.Time
	wd         time.Time // write deadline
	since      int64
	dataInput  chan bool
	dataOutput chan bool
	Config     *Config

	conv             uint16
	state            State
	stateBeginTime   uint32
	lastIncomingTime uint32
	lastPingTime     uint32

	mss       uint32
	roundTrip *RoundTripInfo

	receivingWorker *ReceivingWorker
	sendingWorker   *SendingWorker

	output SegmentWriter

	dataUpdater *Updater
	pingUpdater *Updater
}

// NewConnection create a new KCP connection between local and remote.
func NewConnection(conv uint16, sysConn SystemConnection, config *Config) *Connection {
	log.Trace(newError("creating connection ", conv))

	conn := &Connection{
		conv:       conv,
		conn:       sysConn,
		since:      nowMillisec(),
		dataInput:  make(chan bool, 1),
		dataOutput: make(chan bool, 1),
		Config:     config,
		output:     NewRetryableWriter(NewSegmentWriter(sysConn)),
		mss:        config.GetMTUValue() - uint32(sysConn.Overhead()) - DataSegmentOverhead,
		roundTrip: &RoundTripInfo{
			rto:    100,
			minRtt: config.GetTTIValue(),
		},
	}
	sysConn.Reset(conn.Input)

	conn.receivingWorker = NewReceivingWorker(conn)
	conn.sendingWorker = NewSendingWorker(conn)

	isTerminating := func() bool {
		return conn.State().Is(StateTerminating, StateTerminated)
	}
	isTerminated := func() bool {
		return conn.State() == StateTerminated
	}
	conn.dataUpdater = NewUpdater(
		config.GetTTIValue(),
		predicate.Not(isTerminating).And(predicate.Any(conn.sendingWorker.UpdateNecessary, conn.receivingWorker.UpdateNecessary)),
		isTerminating,
		conn.updateTask)
	conn.pingUpdater = NewUpdater(
		5000, // 5 seconds
		predicate.Not(isTerminated),
		isTerminated,
		conn.updateTask)
	conn.pingUpdater.WakeUp()

	return conn
}

func (v *Connection) Elapsed() uint32 {
	return uint32(nowMillisec() - v.since)
}

func (v *Connection) OnDataInput() {
	select {
	case v.dataInput <- true:
	default:
	}
}

func (v *Connection) OnDataOutput() {
	select {
	case v.dataOutput <- true:
	default:
	}
}

// ReadMultiBuffer implements buf.Reader.
func (v *Connection) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if v == nil {
		return nil, io.EOF
	}

	for {
		if v.State().Is(StateReadyToClose, StateTerminating, StateTerminated) {
			return nil, io.EOF
		}
		mb := v.receivingWorker.ReadMultiBuffer()
		if !mb.IsEmpty() {
			return mb, nil
		}

		if v.State() == StatePeerTerminating {
			return nil, io.EOF
		}

		duration := time.Minute
		if !v.rd.IsZero() {
			duration = time.Until(v.rd)
			if duration < 0 {
				return nil, ErrIOTimeout
			}
		}

		select {
		case <-v.dataInput:
		case <-time.After(duration):
			if !v.rd.IsZero() && v.rd.Before(time.Now()) {
				return nil, ErrIOTimeout
			}
		}
	}
}

// Read implements the Conn Read method.
func (v *Connection) Read(b []byte) (int, error) {
	if v == nil {
		return 0, io.EOF
	}

	for {
		if v.State().Is(StateReadyToClose, StateTerminating, StateTerminated) {
			return 0, io.EOF
		}
		nBytes := v.receivingWorker.Read(b)
		if nBytes > 0 {
			return nBytes, nil
		}

		if v.State() == StatePeerTerminating {
			return 0, io.EOF
		}

		duration := time.Minute
		if !v.rd.IsZero() {
			duration = time.Until(v.rd)
			if duration < 0 {
				return 0, ErrIOTimeout
			}
		}

		select {
		case <-v.dataInput:
		case <-time.After(duration):
			if !v.rd.IsZero() && v.rd.Before(time.Now()) {
				return 0, ErrIOTimeout
			}
		}
	}
}

// Write implements the Conn Write method.
func (v *Connection) Write(b []byte) (int, error) {
	totalWritten := 0

	for {
		if v == nil || v.State() != StateActive {
			return totalWritten, io.ErrClosedPipe
		}

		nBytes := v.sendingWorker.Push(b[totalWritten:])
		v.dataUpdater.WakeUp()
		if nBytes > 0 {
			totalWritten += nBytes
			if totalWritten == len(b) {
				return totalWritten, nil
			}
		}

		duration := time.Minute
		if !v.wd.IsZero() {
			duration = time.Until(v.wd)
			if duration < 0 {
				return totalWritten, ErrIOTimeout
			}
		}

		select {
		case <-v.dataOutput:
		case <-time.After(duration):
			if !v.wd.IsZero() && v.wd.Before(time.Now()) {
				return totalWritten, ErrIOTimeout
			}
		}
	}
}

func (v *Connection) SetState(state State) {
	current := v.Elapsed()
	atomic.StoreInt32((*int32)(&v.state), int32(state))
	atomic.StoreUint32(&v.stateBeginTime, current)
	log.Trace(newError("#", v.conv, " entering state ", state, " at ", current).AtDebug())

	switch state {
	case StateReadyToClose:
		v.receivingWorker.CloseRead()
	case StatePeerClosed:
		v.sendingWorker.CloseWrite()
	case StateTerminating:
		v.receivingWorker.CloseRead()
		v.sendingWorker.CloseWrite()
		v.pingUpdater.SetInterval(time.Second)
	case StatePeerTerminating:
		v.sendingWorker.CloseWrite()
		v.pingUpdater.SetInterval(time.Second)
	case StateTerminated:
		v.receivingWorker.CloseRead()
		v.sendingWorker.CloseWrite()
		v.pingUpdater.SetInterval(time.Second)
		v.dataUpdater.WakeUp()
		v.pingUpdater.WakeUp()
		go v.Terminate()
	}
}

// Close closes the connection.
func (v *Connection) Close() error {
	if v == nil {
		return ErrClosedConnection
	}

	v.OnDataInput()
	v.OnDataOutput()

	state := v.State()
	if state.Is(StateReadyToClose, StateTerminating, StateTerminated) {
		return ErrClosedConnection
	}
	log.Trace(newError("closing connection to ", v.conn.RemoteAddr()))

	if state == StateActive {
		v.SetState(StateReadyToClose)
	}
	if state == StatePeerClosed {
		v.SetState(StateTerminating)
	}
	if state == StatePeerTerminating {
		v.SetState(StateTerminated)
	}

	return nil
}

// LocalAddr returns the local network address. The Addr returned is shared by all invocations of LocalAddr, so do not modify it.
func (v *Connection) LocalAddr() net.Addr {
	if v == nil {
		return nil
	}
	return v.conn.LocalAddr()
}

// RemoteAddr returns the remote network address. The Addr returned is shared by all invocations of RemoteAddr, so do not modify it.
func (v *Connection) RemoteAddr() net.Addr {
	if v == nil {
		return nil
	}
	return v.conn.RemoteAddr()
}

// SetDeadline sets the deadline associated with the listener. A zero time value disables the deadline.
func (v *Connection) SetDeadline(t time.Time) error {
	if err := v.SetReadDeadline(t); err != nil {
		return err
	}
	if err := v.SetWriteDeadline(t); err != nil {
		return err
	}
	return nil
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (v *Connection) SetReadDeadline(t time.Time) error {
	if v == nil || v.State() != StateActive {
		return ErrClosedConnection
	}
	v.rd = t
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (v *Connection) SetWriteDeadline(t time.Time) error {
	if v == nil || v.State() != StateActive {
		return ErrClosedConnection
	}
	v.wd = t
	return nil
}

// kcp update, input loop
func (v *Connection) updateTask() {
	v.flush()
}

func (v *Connection) Terminate() {
	if v == nil {
		return
	}
	log.Trace(newError("terminating connection to ", v.RemoteAddr()))

	//v.SetState(StateTerminated)
	v.OnDataInput()
	v.OnDataOutput()

	v.conn.Close()
	v.sendingWorker.Release()
	v.receivingWorker.Release()
}

func (v *Connection) HandleOption(opt SegmentOption) {
	if (opt & SegmentOptionClose) == SegmentOptionClose {
		v.OnPeerClosed()
	}
}

func (v *Connection) OnPeerClosed() {
	state := v.State()
	if state == StateReadyToClose {
		v.SetState(StateTerminating)
	}
	if state == StateActive {
		v.SetState(StatePeerClosed)
	}
}

// Input when you received a low level packet (eg. UDP packet), call it
func (v *Connection) Input(segments []Segment) {
	current := v.Elapsed()
	atomic.StoreUint32(&v.lastIncomingTime, current)

	for _, seg := range segments {
		if seg.Conversation() != v.conv {
			break
		}

		switch seg := seg.(type) {
		case *DataSegment:
			v.HandleOption(seg.Option)
			v.receivingWorker.ProcessSegment(seg)
			if v.receivingWorker.IsDataAvailable() {
				v.OnDataInput()
			}
			v.dataUpdater.WakeUp()
		case *AckSegment:
			v.HandleOption(seg.Option)
			v.sendingWorker.ProcessSegment(current, seg, v.roundTrip.Timeout())
			v.OnDataOutput()
			v.dataUpdater.WakeUp()
		case *CmdOnlySegment:
			v.HandleOption(seg.Option)
			if seg.Command() == CommandTerminate {
				state := v.State()
				if state == StateActive ||
					state == StatePeerClosed {
					v.SetState(StatePeerTerminating)
				} else if state == StateReadyToClose {
					v.SetState(StateTerminating)
				} else if state == StateTerminating {
					v.SetState(StateTerminated)
				}
			}
			if seg.Option == SegmentOptionClose || seg.Command() == CommandTerminate {
				v.OnDataInput()
				v.OnDataOutput()
			}
			v.sendingWorker.ProcessReceivingNext(seg.ReceivinNext)
			v.receivingWorker.ProcessSendingNext(seg.SendingNext)
			v.roundTrip.UpdatePeerRTO(seg.PeerRTO, current)
			seg.Release()
		default:
		}
	}
}

func (v *Connection) flush() {
	current := v.Elapsed()

	if v.State() == StateTerminated {
		return
	}
	if v.State() == StateActive && current-atomic.LoadUint32(&v.lastIncomingTime) >= 30000 {
		v.Close()
	}
	if v.State() == StateReadyToClose && v.sendingWorker.IsEmpty() {
		v.SetState(StateTerminating)
	}

	if v.State() == StateTerminating {
		log.Trace(newError("#", v.conv, " sending terminating cmd.").AtDebug())
		v.Ping(current, CommandTerminate)

		if current-atomic.LoadUint32(&v.stateBeginTime) > 8000 {
			v.SetState(StateTerminated)
		}
		return
	}
	if v.State() == StatePeerTerminating && current-atomic.LoadUint32(&v.stateBeginTime) > 4000 {
		v.SetState(StateTerminating)
	}

	if v.State() == StateReadyToClose && current-atomic.LoadUint32(&v.stateBeginTime) > 15000 {
		v.SetState(StateTerminating)
	}

	// flush acknowledges
	v.receivingWorker.Flush(current)
	v.sendingWorker.Flush(current)

	if current-atomic.LoadUint32(&v.lastPingTime) >= 3000 {
		v.Ping(current, CommandPing)
	}
}

func (v *Connection) State() State {
	return State(atomic.LoadInt32((*int32)(&v.state)))
}

func (v *Connection) Ping(current uint32, cmd Command) {
	seg := NewCmdOnlySegment()
	seg.Conv = v.conv
	seg.Cmd = cmd
	seg.ReceivinNext = v.receivingWorker.NextNumber()
	seg.SendingNext = v.sendingWorker.FirstUnacknowledged()
	seg.PeerRTO = v.roundTrip.Timeout()
	if v.State() == StateReadyToClose {
		seg.Option = SegmentOptionClose
	}
	v.output.Write(seg)
	atomic.StoreUint32(&v.lastPingTime, current)
	seg.Release()
}
