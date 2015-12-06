package app

type Space struct {
	packetDispatcher PacketDispatcher
	dnsCache         DnsCache
}

func NewSpace() *Space {
	return new(Space)
}

func (this *Space) Bind(object interface{}) {
	if packetDispatcher, ok := object.(PacketDispatcher); ok {
		this.packetDispatcher = packetDispatcher
	}

	if dnsCache, ok := object.(DnsCache); ok {
		this.dnsCache = dnsCache
	}
}

func (this *Space) HasPacketDispatcher() bool {
	return this.packetDispatcher != nil
}

func (this *Space) PacketDispatcher() PacketDispatcher {
	return this.packetDispatcher
}

func (this *Space) HasDnsCache() bool {
	return this.dnsCache != nil
}

func (this *Space) DnsCache() DnsCache {
	return this.dnsCache
}
