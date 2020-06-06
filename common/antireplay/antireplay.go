package antireplay

import (
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"sync"
	"time"
)

func NewAntiReplayWindow(AntiReplayTime int64) *AntiReplayWindow {
	arw := &AntiReplayWindow{}
	arw.AntiReplayTime = AntiReplayTime
	return arw
}

type AntiReplayWindow struct {
	lock           sync.Mutex
	poolA          *cuckoo.Filter
	poolB          *cuckoo.Filter
	lastSwapTime   int64
	PoolSwap       bool
	AntiReplayTime int64
}

func (aw *AntiReplayWindow) Check(sum []byte) bool {
	aw.lock.Lock()

	if aw.lastSwapTime == 0 {
		aw.lastSwapTime = time.Now().Unix()
		aw.poolA = cuckoo.NewFilter(100000)
		aw.poolB = cuckoo.NewFilter(100000)
	}

	tnow := time.Now().Unix()
	timediff := tnow - aw.lastSwapTime

	if timediff >= aw.AntiReplayTime {
		if aw.PoolSwap {
			aw.PoolSwap = false
			aw.poolA.Reset()
		} else {
			aw.PoolSwap = true
			aw.poolB.Reset()
		}
		aw.lastSwapTime = tnow
	}

	ret := aw.poolA.InsertUnique(sum) && aw.poolB.InsertUnique(sum)
	aw.lock.Unlock()
	return ret
}
