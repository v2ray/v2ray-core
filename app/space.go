package app

type Space struct {
	packetDispatcher PacketDispatcher
}

func NewSpace() *Space {
	return new(Space)
}

func (this *Space) HasPacketDispatcher() bool {
	return this.packetDispatcher != nil
}

func (this *Space) PacketDispatcher() PacketDispatcher {
	return this.packetDispatcher
}

func (this *Space) Bind(object interface{}) {
	if packetDispatcher, ok := object.(PacketDispatcher); ok {
		this.packetDispatcher = packetDispatcher
	}
}
