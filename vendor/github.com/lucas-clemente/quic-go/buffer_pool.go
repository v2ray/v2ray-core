package quic

import (
	"sync"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/common/bytespool"
)

type packetBuffer struct {
	Slice []byte

	// refCount counts how many packets the Slice is used in.
	// It doesn't support concurrent use.
	// It is > 1 when used for coalesced packet.
	refCount int
}

// Split increases the refCount.
// It must be called when a packet buffer is used for more than one packet,
// e.g. when splitting coalesced packets.
func (b *packetBuffer) Split() {
	b.refCount++
}

// Release decreases the refCount.
// It should be called when processing the packet is finished.
// When the refCount reaches 0, the packet buffer is put back into the pool.
func (b *packetBuffer) Release() {
	if cap(b.Slice) < 2048 {
		return
	}
	b.refCount--
	if b.refCount < 0 {
		panic("negative packetBuffer refCount")
	}
	// only put the packetBuffer back if it's not used any more
	if b.refCount == 0 {
		buffer := b.Slice[0:cap(b.Slice)]
		bufferPool.Put(buffer)
	}
}

var bufferPool *sync.Pool

func getPacketBuffer() *packetBuffer {
	buffer := bufferPool.Get().([]byte)
	return &packetBuffer{
		refCount: 1,
		Slice:    buffer[:protocol.MaxReceivePacketSize],
	}
}

func init() {
	bufferPool = bytespool.GetPool(int32(protocol.MaxReceivePacketSize))
}
