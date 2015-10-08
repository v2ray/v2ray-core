package net

import (
	"github.com/v2ray/v2ray-core/common/alloc"
)

type Packet interface {
	Destination() Destination
	Chunk() *alloc.Buffer // First chunk of this commnunication
	MoreChunks() bool
}

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
