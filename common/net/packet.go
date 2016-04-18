package net

import (
	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
)

// Packet is a network packet to be sent to destination.
type Packet interface {
	common.Releasable

	Destination() Destination
	Chunk() *alloc.Buffer // First chunk of this commnunication
	MoreChunks() bool
}

// NewPacket creates a new Packet with given destination and payload.
func NewPacket(dest Destination, firstChunk *alloc.Buffer, moreChunks bool) Packet {
	return &packetImpl{
		dest:     dest,
		data:     firstChunk,
		moreData: moreChunks,
	}
}

type packetImpl struct {
	dest     Destination
	data     *alloc.Buffer
	moreData bool
}

func (packet *packetImpl) Destination() Destination {
	return packet.dest
}

func (packet *packetImpl) Chunk() *alloc.Buffer {
	return packet.data
}

func (packet *packetImpl) MoreChunks() bool {
	return packet.moreData
}

func (packet *packetImpl) Release() {
	packet.data.Release()
	packet.data = nil
}
