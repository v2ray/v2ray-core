package app

type Context interface {
	CallerTag() string
}

type contextImpl struct {
	callerTag string
}

func (this *contextImpl) CallerTag() string {
	return this.callerTag
}

type SpaceController struct {
	packetDispatcher PacketDispatcherWithContext
	dnsCache         DnsCacheWithContext
}

func NewSpaceController() *SpaceController {
	return new(SpaceController)
}

func (this *SpaceController) Bind(object interface{}) {
	if packetDispatcher, ok := object.(PacketDispatcherWithContext); ok {
		this.packetDispatcher = packetDispatcher
	}

	if dnsCache, ok := object.(DnsCacheWithContext); ok {
		this.dnsCache = dnsCache
	}
}

func (this *SpaceController) ForContext(tag string) *Space {
	return newSpace(this, &contextImpl{callerTag: tag})
}

type Space struct {
	packetDispatcher PacketDispatcher
	dnsCache         DnsCache
}

func newSpace(controller *SpaceController, context Context) *Space {
	space := new(Space)
	if controller.packetDispatcher != nil {
		space.packetDispatcher = &contextedPacketDispatcher{
			context:          context,
			packetDispatcher: controller.packetDispatcher,
		}
	}
	if controller.dnsCache != nil {
		space.dnsCache = &contextedDnsCache{
			context:  context,
			dnsCache: controller.dnsCache,
		}
	}
	return space
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
