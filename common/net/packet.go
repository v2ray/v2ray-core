package net

type Packet interface {
	Destination() Destination
	Chunk() []byte // First chunk of this commnunication
	MoreChunks() bool
}

func NewTCPPacket(dest Destination) *TCPPacket {
	return &TCPPacket{
		basePacket: basePacket{destination: dest},
	}
}

func NewUDPPacket(dest Destination, data []byte, token uint16) *UDPPacket {
	return &UDPPacket{
		basePacket: basePacket{destination: dest},
		data:       data,
		token:      token,
	}
}

type basePacket struct {
	destination Destination
}

func (base basePacket) Destination() Destination {
	return base.destination
}

type TCPPacket struct {
	basePacket
}

func (packet *TCPPacket) Chunk() []byte {
	return nil
}

func (packet *TCPPacket) MoreChunks() bool {
	return true
}

type UDPPacket struct {
	basePacket
	data  []byte
	token uint16
}

func (packet *UDPPacket) Token() uint16 {
	return packet.token
}

func (packet *UDPPacket) Chunk() []byte {
	return packet.data
}

func (packet *UDPPacket) MoreChunks() bool {
	return false
}
