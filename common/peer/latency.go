package peer

import (
	"sync"
)

type Latency interface {
	Value() uint64
}

type HasLatency interface {
	ConnectionLatency() Latency
	HandshakeLatency() Latency
}

type AverageLatency struct {
	access sync.Mutex
	value  uint64
}

func (al *AverageLatency) Update(newValue uint64) {
	al.access.Lock()
	defer al.access.Unlock()

	al.value = (al.value + newValue*2) / 3
}

func (al *AverageLatency) Value() uint64 {
	return al.value
}
