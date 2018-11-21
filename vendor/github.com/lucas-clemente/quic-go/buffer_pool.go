package quic

import (
	"sync"

	"v2ray.com/core/common/bytespool"

	"github.com/lucas-clemente/quic-go/internal/protocol"
)

var bufferPool *sync.Pool

func getPacketBuffer() *[]byte {
	b := bufferPool.Get().([]byte)
	return &b
}

func putPacketBuffer(buf *[]byte) {
	if cap(*buf) < int(protocol.MaxReceivePacketSize) {
		panic("putPacketBuffer called with packet of wrong size!")
	}
	bufferPool.Put(*buf)
}

func init() {
	bufferPool = bytespool.GetPool(protocol.MaxReceivePacketSize)
}
