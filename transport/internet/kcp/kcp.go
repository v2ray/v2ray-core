// Package kcp - A Fast and Reliable ARQ Protocol
//
// Acknowledgement:
//    skywind3000@github for inventing the KCP protocol
//    xtaci@github for translating to Golang
package kcp

import (
	"github.com/v2ray/v2ray-core/common/log"
)

func _itimediff(later, earlier uint32) int32 {
	return (int32)(later - earlier)
}

type State int

const (
	StateActive       State = 0
	StateReadyToClose State = 1
	StatePeerClosed   State = 2
	StateTerminating  State = 3
	StateTerminated   State = 4
)

// KCP defines a single KCP connection
type KCP struct {
	conv             uint16
	state            State
	stateBeginTime   uint32
	lastIncomingTime uint32
	lastPayloadTime  uint32
	sendingUpdated   bool
	lastPingTime     uint32

	mss                        uint32
	rx_rttvar, rx_srtt, rx_rto uint32
	current, interval          uint32

	receivingWorker *ReceivingWorker
	sendingWorker   *SendingWorker

	fastresend        uint32
	congestionControl bool
	output            *BufferedSegmentWriter
}

// NewKCP create a new kcp control object, 'conv' must equal in two endpoint
// from the same connection.
func NewKCP(conv uint16, output *AuthenticationWriter) *KCP {
	log.Debug("KCP|Core: creating KCP ", conv)
	kcp := new(KCP)
	kcp.conv = conv
	kcp.mss = output.Mtu() - DataSegmentOverhead
	kcp.rx_rto = 100
	kcp.interval = effectiveConfig.Tti
	kcp.output = NewSegmentWriter(output)
	kcp.receivingWorker = NewReceivingWorker(kcp)
	kcp.fastresend = 2
	kcp.congestionControl = effectiveConfig.Congestion
	kcp.sendingWorker = NewSendingWorker(kcp)
	return kcp
}

func (kcp *KCP) SetState(state State) {
	kcp.state = state
	kcp.stateBeginTime = kcp.current

	switch state {
	case StateReadyToClose:
		kcp.receivingWorker.CloseRead()
	case StatePeerClosed:
		kcp.sendingWorker.CloseWrite()
	case StateTerminating:
		kcp.receivingWorker.CloseRead()
		kcp.sendingWorker.CloseWrite()
	case StateTerminated:
		kcp.receivingWorker.CloseRead()
		kcp.sendingWorker.CloseWrite()
	}
}

func (kcp *KCP) HandleOption(opt SegmentOption) {
	if (opt & SegmentOptionClose) == SegmentOptionClose {
		kcp.OnPeerClosed()
	}
}

func (kcp *KCP) OnPeerClosed() {
	if kcp.state == StateReadyToClose {
		kcp.SetState(StateTerminating)
	}
	if kcp.state == StateActive {
		kcp.SetState(StatePeerClosed)
	}
}

func (kcp *KCP) OnClose() {
	if kcp.state == StateActive {
		kcp.SetState(StateReadyToClose)
	}
	if kcp.state == StatePeerClosed {
		kcp.SetState(StateTerminating)
	}
}

// https://tools.ietf.org/html/rfc6298
func (kcp *KCP) update_ack(rtt int32) {
	if kcp.rx_srtt == 0 {
		kcp.rx_srtt = uint32(rtt)
		kcp.rx_rttvar = uint32(rtt) / 2
	} else {
		delta := rtt - int32(kcp.rx_srtt)
		if delta < 0 {
			delta = -delta
		}
		kcp.rx_rttvar = (3*kcp.rx_rttvar + uint32(delta)) / 4
		kcp.rx_srtt = (7*kcp.rx_srtt + uint32(rtt)) / 8
		if kcp.rx_srtt < kcp.interval {
			kcp.rx_srtt = kcp.interval
		}
	}
	var rto uint32
	if kcp.interval < 4*kcp.rx_rttvar {
		rto = kcp.rx_srtt + 4*kcp.rx_rttvar
	} else {
		rto = kcp.rx_srtt + kcp.interval
	}

	if rto > 10000 {
		rto = 10000
	}
	kcp.rx_rto = rto * 3 / 2
}

// Input when you received a low level packet (eg. UDP packet), call it
func (kcp *KCP) Input(data []byte) int {
	kcp.lastIncomingTime = kcp.current

	var seg Segment
	for {
		seg, data = ReadSegment(data)
		if seg == nil {
			break
		}

		switch seg := seg.(type) {
		case *DataSegment:
			kcp.HandleOption(seg.Opt)
			kcp.receivingWorker.ProcessSegment(seg)
			kcp.lastPayloadTime = kcp.current
		case *AckSegment:
			kcp.HandleOption(seg.Opt)
			kcp.sendingWorker.ProcessSegment(seg)
			kcp.lastPayloadTime = kcp.current
		case *CmdOnlySegment:
			kcp.HandleOption(seg.Opt)
			if seg.Cmd == SegmentCommandTerminated {
				if kcp.state == StateActive ||
					kcp.state == StateReadyToClose ||
					kcp.state == StatePeerClosed {
					kcp.SetState(StateTerminating)
				} else if kcp.state == StateTerminating {
					kcp.SetState(StateTerminated)
				}
			}
			kcp.sendingWorker.ProcessReceivingNext(seg.ReceivinNext)
			kcp.receivingWorker.ProcessSendingNext(seg.SendingNext)
		default:
		}
	}

	return 0
}

// flush pending data
func (kcp *KCP) flush() {
	if kcp.state == StateTerminated {
		return
	}
	if kcp.state == StateActive && _itimediff(kcp.current, kcp.lastPayloadTime) >= 30000 {
		kcp.OnClose()
	}

	if kcp.state == StateTerminating {
		kcp.output.Write(&CmdOnlySegment{
			Conv: kcp.conv,
			Cmd:  SegmentCommandTerminated,
		})
		kcp.output.Flush()

		if _itimediff(kcp.current, kcp.stateBeginTime) > 8000 {
			kcp.SetState(StateTerminated)
		}
		return
	}

	if kcp.state == StateReadyToClose && _itimediff(kcp.current, kcp.stateBeginTime) > 15000 {
		kcp.SetState(StateTerminating)
	}

	// flush acknowledges
	kcp.receivingWorker.Flush()
	kcp.sendingWorker.Flush()

	if kcp.sendingWorker.PingNecessary() || kcp.receivingWorker.PingNecessary() || _itimediff(kcp.current, kcp.lastPingTime) >= 5000 {
		seg := NewCmdOnlySegment()
		seg.Conv = kcp.conv
		seg.Cmd = SegmentCommandPing
		seg.ReceivinNext = kcp.receivingWorker.nextNumber
		seg.SendingNext = kcp.sendingWorker.firstUnacknowledged
		if kcp.state == StateReadyToClose {
			seg.Opt = SegmentOptionClose
		}
		kcp.output.Write(seg)
		kcp.lastPingTime = kcp.current
		kcp.sendingUpdated = false
		seg.Release()
	}

	// flash remain segments
	kcp.output.Flush()

}

// Update updates state (call it repeatedly, every 10ms-100ms), or you can ask
// ikcp_check when to call it again (without ikcp_input/_send calling).
// 'current' - current timestamp in millisec.
func (kcp *KCP) Update(current uint32) {
	kcp.current = current
	kcp.flush()
}
