package net

type Packet interface {
	Destination() Destination
	Chunk() []byte // First chunk of this commnunication
	MoreChunks() bool
}

func NewPacket(dest Destination, firstChunk []byte, moreChunks bool) Packet {
	return &packetImpl{
		dest:     dest,
		data:     firstChunk,
		moreData: moreChunks,
	}
}

type packetImpl struct {
	dest     Destination
	data     []byte
	moreData bool
}

func (packet *packetImpl) Destination() Destination {
	return packet.dest
}

func (packet *packetImpl) Chunk() []byte {
	return packet.data
}

func (packet *packetImpl) MoreChunks() bool {
	return packet.moreData
}
