package kcp

import (
	"sync"
)

const (
	defaultRTT = 100
	queueSize  = 10
)

type Queue struct {
	value  [queueSize]uint32
	start  uint32
	length uint32
}

func (v *Queue) Push(value uint32) {
	if v.length < queueSize {
		v.value[v.length] = value
		v.length++
		return
	}
	v.value[v.start] = value
	v.start++
	if v.start == queueSize {
		v.start = 0
	}
}

func (v *Queue) Max() uint32 {
	max := v.value[0]
	for i := 1; i < queueSize; i++ {
		if v.value[i] > max {
			max = v.value[i]
		}
	}
	return max
}

func (v *Queue) Min() uint32 {
	max := v.value[0]
	for i := 1; i < queueSize; i++ {
		if v.value[i] < max {
			max = v.value[i]
		}
	}
	return max
}

type CongestionState byte

const (
	CongestionStateRTTProbe CongestionState = iota
	CongestionStateBandwidthProbe
	CongestionStateTransfer
)

type Congestion struct {
	sync.RWMutex

	state      CongestionState
	stateSince uint32
	limit      uint32 // bytes per 1000 seconds

	rtt           uint32 // millisec
	rttHistory    Queue
	rttUpdateTime uint32

	initialThroughput uint32 // bytes per 1000 seconds

	cycleStartTime      uint32
	cycleBytesConfirmed uint32
	cycleBytesSent      uint32
	cycleBytesLimit     uint32

	cycle                   uint32
	bestCycleBytesConfirmed uint32
	bestCycleBytesSent      uint32
}

func (v *Congestion) SetState(current uint32, state CongestionState) {
	v.state = state
	v.stateSince = current
}

func (v *Congestion) Update(current uint32) {
	switch v.state {
	case CongestionStateRTTProbe:
		if v.rtt > 0 {
			v.SetState(current, CongestionStateBandwidthProbe)
			v.cycleStartTime = current
			v.cycleBytesConfirmed = 0
			v.cycleBytesSent = 0
			v.cycleBytesLimit = v.initialThroughput * v.rtt / 1000 / 1000
		}
	case CongestionStateBandwidthProbe:
		if current-v.cycleStartTime >= v.rtt {

		}
	}
}

func (v *Congestion) AddBytesConfirmed(current uint32, bytesConfirmed uint32) {
	v.Lock()
	defer v.Unlock()

	v.cycleBytesConfirmed += bytesConfirmed
	v.Update(current)
}

func (v *Congestion) UpdateRTT(current uint32, rtt uint32) {
	v.Lock()
	defer v.Unlock()

	if v.state == CongestionStateRTTProbe || rtt < v.rtt {
		v.rtt = rtt
		v.rttUpdateTime = current
	}

	v.Update(current)
}

func (v *Congestion) GetBytesLimit() uint32 {
	v.RLock()
	defer v.RUnlock()

	if v.state == CongestionStateRTTProbe {
		return v.initialThroughput/1000/(1000/defaultRTT) - v.cycleBytesSent
	}

	return v.cycleBytesLimit
}

func (v *Congestion) RoundTripTime() uint32 {
	v.RLock()
	defer v.RUnlock()

	if v.state == CongestionStateRTTProbe {
		return defaultRTT
	}

	return v.rtt
}

func (v *Congestion) Timeout() uint32 {
	v.RLock()
	defer v.RUnlock()

	if v.state == CongestionStateRTTProbe {
		return defaultRTT * 3 / 2
	}

	return v.rtt * 3 / 2
}
