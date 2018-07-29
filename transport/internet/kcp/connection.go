package kcp

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/signal/semaphore"
)

var (
	ErrIOTimeout        = newError("Read/Write timeout")
	ErrClosedListener   = newError("Listener closed.")
	ErrClosedConnection = newError("Connection closed.")
)

// State of the connection
type State int32

// Is returns true if current State is one of the candidates.
func (s State) Is(states ...State) bool {
	for _, state := range states {
		if s == state {
			return true
		}
	}
	return false
}

const (
	StateActive          State = 0 // Connection is active
	StateReadyToClose    State = 1 // Connection is closed locally
	StatePeerClosed      State = 2 // Connection is closed on remote
	StateTerminating     State = 3 // Connection is ready to be destroyed locally
	StatePeerTerminating State = 4 // Connection is ready to be destroyed on remote
	StateTerminated      State = 5 // Connection is destroyed.
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

func (info *RoundTripInfo) UpdatePeerRTO(rto uint32, current uint32) {
	info.Lock()
	defer info.Unlock()

	if current-info.updatedTimestamp < 3000 {
		return
	}

	info.updatedTimestamp = current
	info.rto = rto
}

func (info *RoundTripInfo) Update(rtt uint32, current uint32) {
	if rtt > 0x7FFFFFFF {
		return
	}
	info.Lock()
	defer info.Unlock()

	// https://tools.ietf.org/html/rfc6298
	if info.srtt == 0 {
		info.srtt = rtt
		info.variation = rtt / 2
	} else {
		delta := rtt - info.srtt
		if info.srtt > rtt {
			delta = info.srtt - rtt
		}
		info.variation = (3*info.variation + delta) / 4
		info.srtt = (7*info.srtt + rtt) / 8
		if info.srtt < info.minRtt {
			info.srtt = info.minRtt
		}
	}
	var rto uint32
	if info.minRtt < 4*info.variation {
		rto = info.srtt + 4*info.variation
	} else {
		rto = info.srtt + info.variation
	}

	if rto > 10000 {
		rto = 10000
	}
	info.rto = rto * 5 / 4
	info.updatedTimestamp = current
}

func (info *RoundTripInfo) Timeout() uint32 {
	info.RLock()
	defer info.RUnlock()

	return info.rto
}

func (info *RoundTripInfo) SmoothedTime() uint32 {
	info.RLock()
	defer info.RUnlock()

	return info.srtt
}

type Updater struct {
	interval        int64
	shouldContinue  func() bool
	shouldTerminate func() bool
	updateFunc      func()
	notifier        *semaphore.Instance
}

func NewUpdater(interval uint32, shouldContinue func() bool, shouldTerminate func() bool, updateFunc func()) *Updater {
	u := &Updater{
		interval:        int64(time.Duration(interval) * time.Millisecond),
		shouldContinue:  shouldContinue,
		shouldTerminate: shouldTerminate,
		updateFunc:      updateFunc,
		notifier:        semaphore.New(1),
	}
	return u
}

func (u *Updater) WakeUp() {
	select {
	case <-u.notifier.Wait():
		go u.run()
	default:
	}
}

func (u *Updater) run() {
	defer u.notifier.Signal()

	if u.shouldTerminate() {
		return
	}
	ticker := time.NewTicker(u.Interval())
	for u.shouldContinue() {
		u.updateFunc()
		<-ticker.C
	}
	ticker.Stop()
}

func (u *Updater) Interval() time.Duration {
	return time.Duration(atomic.LoadInt64(&u.interval))
}

func (u *Updater) SetInterval(d time.Duration) {
	atomic.StoreInt64(&u.interval, int64(d))
}

type ConnMetadata struct {
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
	Conversation uint16
}

// Connection is a KCP connection over UDP.
type Connection struct {
	meta       ConnMetadata
	closer     io.Closer
	rd         time.Time
	wd         time.Time // write deadline
	since      int64
	dataInput  *signal.Notifier
	dataOutput *signal.Notifier
	Config     *Config

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
func NewConnection(meta ConnMetadata, writer PacketWriter, closer io.Closer, config *Config) *Connection {
	newError("#", meta.Conversation, " creating connection to ", meta.RemoteAddr).WriteToLog()

	conn := &Connection{
		meta:       meta,
		closer:     closer,
		since:      nowMillisec(),
		dataInput:  signal.NewNotifier(),
		dataOutput: signal.NewNotifier(),
		Config:     config,
		output:     NewRetryableWriter(NewSegmentWriter(writer)),
		mss:        config.GetMTUValue() - uint32(writer.Overhead()) - DataSegmentOverhead,
		roundTrip: &RoundTripInfo{
			rto:    100,
			minRtt: config.GetTTIValue(),
		},
	}

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
		func() bool {
			return !isTerminating() && (conn.sendingWorker.UpdateNecessary() || conn.receivingWorker.UpdateNecessary())
		},
		isTerminating,
		conn.updateTask)
	conn.pingUpdater = NewUpdater(
		5000, // 5 seconds
		func() bool { return !isTerminated() },
		isTerminated,
		conn.updateTask)
	conn.pingUpdater.WakeUp()

	return conn
}

func (c *Connection) Elapsed() uint32 {
	return uint32(nowMillisec() - c.since)
}

// ReadMultiBuffer implements buf.Reader.
func (c *Connection) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if c == nil {
		return nil, io.EOF
	}

	for {
		if c.State().Is(StateReadyToClose, StateTerminating, StateTerminated) {
			return nil, io.EOF
		}
		mb := c.receivingWorker.ReadMultiBuffer()
		if !mb.IsEmpty() {
			c.dataUpdater.WakeUp()
			return mb, nil
		}

		if c.State() == StatePeerTerminating {
			return nil, io.EOF
		}

		if err := c.waitForDataInput(); err != nil {
			return nil, err
		}
	}
}

func (c *Connection) waitForDataInput() error {
	if c.State() == StatePeerTerminating {
		return io.EOF
	}

	duration := time.Minute
	if !c.rd.IsZero() {
		duration = time.Until(c.rd)
		if duration < 0 {
			return ErrIOTimeout
		}
	}

	select {
	case <-c.dataInput.Wait():
	case <-time.After(duration):
		if !c.rd.IsZero() && c.rd.Before(time.Now()) {
			return ErrIOTimeout
		}
	}

	return nil
}

// Read implements the Conn Read method.
func (c *Connection) Read(b []byte) (int, error) {
	if c == nil {
		return 0, io.EOF
	}

	for {
		if c.State().Is(StateReadyToClose, StateTerminating, StateTerminated) {
			return 0, io.EOF
		}
		nBytes := c.receivingWorker.Read(b)
		if nBytes > 0 {
			c.dataUpdater.WakeUp()
			return nBytes, nil
		}

		if err := c.waitForDataInput(); err != nil {
			return 0, err
		}
	}
}

func (c *Connection) waitForDataOutput() error {
	duration := time.Minute
	if !c.wd.IsZero() {
		duration = time.Until(c.wd)
		if duration < 0 {
			return ErrIOTimeout
		}
	}

	select {
	case <-c.dataOutput.Wait():
	case <-time.After(duration):
		if !c.wd.IsZero() && c.wd.Before(time.Now()) {
			return ErrIOTimeout
		}
	}

	return nil
}

// Write implements io.Writer.
func (c *Connection) Write(b []byte) (int, error) {
	updatePending := false
	defer func() {
		if updatePending {
			c.dataUpdater.WakeUp()
		}
	}()

	for {
		totalWritten := 0
		for {
			if c == nil || c.State() != StateActive {
				return totalWritten, io.ErrClosedPipe
			}
			if !c.sendingWorker.Push(func(bb []byte) (int, error) {
				n := copy(bb[:c.mss], b[totalWritten:])
				totalWritten += n
				return n, nil
			}) {
				break
			}

			updatePending = true

			if totalWritten == len(b) {
				return totalWritten, nil
			}
		}

		if updatePending {
			c.dataUpdater.WakeUp()
			updatePending = false
		}

		if err := c.waitForDataOutput(); err != nil {
			return totalWritten, err
		}
	}
}

// WriteMultiBuffer implements buf.Writer.
func (c *Connection) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer mb.Release()

	updatePending := false
	defer func() {
		if updatePending {
			c.dataUpdater.WakeUp()
		}
	}()

	for {
		for {
			if c == nil || c.State() != StateActive {
				return io.ErrClosedPipe
			}

			mbPtr := &mb

			if !c.sendingWorker.Push(func(bb []byte) (int, error) {
				return mbPtr.Read(bb[:c.mss])
			}) {
				break
			}
			updatePending = true
			if mb.IsEmpty() {
				return nil
			}
		}

		if updatePending {
			c.dataUpdater.WakeUp()
			updatePending = false
		}

		if err := c.waitForDataOutput(); err != nil {
			return err
		}
	}
}

func (c *Connection) SetState(state State) {
	current := c.Elapsed()
	atomic.StoreInt32((*int32)(&c.state), int32(state))
	atomic.StoreUint32(&c.stateBeginTime, current)
	newError("#", c.meta.Conversation, " entering state ", state, " at ", current).AtDebug().WriteToLog()

	switch state {
	case StateReadyToClose:
		c.receivingWorker.CloseRead()
	case StatePeerClosed:
		c.sendingWorker.CloseWrite()
	case StateTerminating:
		c.receivingWorker.CloseRead()
		c.sendingWorker.CloseWrite()
		c.pingUpdater.SetInterval(time.Second)
	case StatePeerTerminating:
		c.sendingWorker.CloseWrite()
		c.pingUpdater.SetInterval(time.Second)
	case StateTerminated:
		c.receivingWorker.CloseRead()
		c.sendingWorker.CloseWrite()
		c.pingUpdater.SetInterval(time.Second)
		c.dataUpdater.WakeUp()
		c.pingUpdater.WakeUp()
		go c.Terminate()
	}
}

// Close closes the connection.
func (c *Connection) Close() error {
	if c == nil {
		return ErrClosedConnection
	}

	c.dataInput.Signal()
	c.dataOutput.Signal()

	switch c.State() {
	case StateReadyToClose, StateTerminating, StateTerminated:
		return ErrClosedConnection
	case StateActive:
		c.SetState(StateReadyToClose)
	case StatePeerClosed:
		c.SetState(StateTerminating)
	case StatePeerTerminating:
		c.SetState(StateTerminated)
	}

	newError("#", c.meta.Conversation, " closing connection to ", c.meta.RemoteAddr).WriteToLog()

	return nil
}

// LocalAddr returns the local network address. The Addr returned is shared by all invocations of LocalAddr, so do not modify it.
func (c *Connection) LocalAddr() net.Addr {
	if c == nil {
		return nil
	}
	return c.meta.LocalAddr
}

// RemoteAddr returns the remote network address. The Addr returned is shared by all invocations of RemoteAddr, so do not modify it.
func (c *Connection) RemoteAddr() net.Addr {
	if c == nil {
		return nil
	}
	return c.meta.RemoteAddr
}

// SetDeadline sets the deadline associated with the listener. A zero time value disables the deadline.
func (c *Connection) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *Connection) SetReadDeadline(t time.Time) error {
	if c == nil || c.State() != StateActive {
		return ErrClosedConnection
	}
	c.rd = t
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *Connection) SetWriteDeadline(t time.Time) error {
	if c == nil || c.State() != StateActive {
		return ErrClosedConnection
	}
	c.wd = t
	return nil
}

// kcp update, input loop
func (c *Connection) updateTask() {
	c.flush()
}

func (c *Connection) Terminate() {
	if c == nil {
		return
	}
	newError("#", c.meta.Conversation, " terminating connection to ", c.RemoteAddr()).WriteToLog()

	//v.SetState(StateTerminated)
	c.dataInput.Signal()
	c.dataOutput.Signal()

	c.closer.Close()
	c.sendingWorker.Release()
	c.receivingWorker.Release()
}

func (c *Connection) HandleOption(opt SegmentOption) {
	if (opt & SegmentOptionClose) == SegmentOptionClose {
		c.OnPeerClosed()
	}
}

func (c *Connection) OnPeerClosed() {
	switch c.State() {
	case StateReadyToClose:
		c.SetState(StateTerminating)
	case StateActive:
		c.SetState(StatePeerClosed)
	}
}

// Input when you received a low level packet (eg. UDP packet), call it
func (c *Connection) Input(segments []Segment) {
	current := c.Elapsed()
	atomic.StoreUint32(&c.lastIncomingTime, current)

	for _, seg := range segments {
		if seg.Conversation() != c.meta.Conversation {
			break
		}

		switch seg := seg.(type) {
		case *DataSegment:
			c.HandleOption(seg.Option)
			c.receivingWorker.ProcessSegment(seg)
			if c.receivingWorker.IsDataAvailable() {
				c.dataInput.Signal()
			}
			c.dataUpdater.WakeUp()
		case *AckSegment:
			c.HandleOption(seg.Option)
			c.sendingWorker.ProcessSegment(current, seg, c.roundTrip.Timeout())
			c.dataOutput.Signal()
			c.dataUpdater.WakeUp()
		case *CmdOnlySegment:
			c.HandleOption(seg.Option)
			if seg.Command() == CommandTerminate {
				switch c.State() {
				case StateActive, StatePeerClosed:
					c.SetState(StatePeerTerminating)
				case StateReadyToClose:
					c.SetState(StateTerminating)
				case StateTerminating:
					c.SetState(StateTerminated)
				}
			}
			if seg.Option == SegmentOptionClose || seg.Command() == CommandTerminate {
				c.dataInput.Signal()
				c.dataOutput.Signal()
			}
			c.sendingWorker.ProcessReceivingNext(seg.ReceivingNext)
			c.receivingWorker.ProcessSendingNext(seg.SendingNext)
			c.roundTrip.UpdatePeerRTO(seg.PeerRTO, current)
			seg.Release()
		default:
		}
	}
}

func (c *Connection) flush() {
	current := c.Elapsed()

	if c.State() == StateTerminated {
		return
	}
	if c.State() == StateActive && current-atomic.LoadUint32(&c.lastIncomingTime) >= 30000 {
		c.Close()
	}
	if c.State() == StateReadyToClose && c.sendingWorker.IsEmpty() {
		c.SetState(StateTerminating)
	}

	if c.State() == StateTerminating {
		newError("#", c.meta.Conversation, " sending terminating cmd.").AtDebug().WriteToLog()
		c.Ping(current, CommandTerminate)

		if current-atomic.LoadUint32(&c.stateBeginTime) > 8000 {
			c.SetState(StateTerminated)
		}
		return
	}
	if c.State() == StatePeerTerminating && current-atomic.LoadUint32(&c.stateBeginTime) > 4000 {
		c.SetState(StateTerminating)
	}

	if c.State() == StateReadyToClose && current-atomic.LoadUint32(&c.stateBeginTime) > 15000 {
		c.SetState(StateTerminating)
	}

	// flush acknowledges
	c.receivingWorker.Flush(current)
	c.sendingWorker.Flush(current)

	if current-atomic.LoadUint32(&c.lastPingTime) >= 3000 {
		c.Ping(current, CommandPing)
	}
}

func (c *Connection) State() State {
	return State(atomic.LoadInt32((*int32)(&c.state)))
}

func (c *Connection) Ping(current uint32, cmd Command) {
	seg := NewCmdOnlySegment()
	seg.Conv = c.meta.Conversation
	seg.Cmd = cmd
	seg.ReceivingNext = c.receivingWorker.NextNumber()
	seg.SendingNext = c.sendingWorker.FirstUnacknowledged()
	seg.PeerRTO = c.roundTrip.Timeout()
	if c.State() == StateReadyToClose {
		seg.Option = SegmentOptionClose
	}
	c.output.Write(seg)
	atomic.StoreUint32(&c.lastPingTime, current)
	seg.Release()
}
