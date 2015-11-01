package testing

import (
	"sync/atomic"
)

var (
	port = int32(30000)
)

func PickPort() uint16 {
	return uint16(atomic.AddInt32(&port, 1))
}
