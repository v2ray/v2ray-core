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
	b := *buf
	if cap(b) < 2048 {
		return
	}
	bufferPool.Put(b[:cap(b)])
}

func init() {
	bufferPool = bytespool.GetPool(int32(protocol.MaxReceivePacketSize))
}
