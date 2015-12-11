package controller

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/internal"
)

// A SpaceController is supposed to be used by a shell to create Spaces. It should not be used
// directly by proxies.
type SpaceController struct {
	packetDispatcher internal.PacketDispatcherWithContext
	dnsCache         internal.DnsCacheWithContext
}

func New() *SpaceController {
	return new(SpaceController)
}

func (this *SpaceController) Bind(object interface{}) {
	if packetDispatcher, ok := object.(internal.PacketDispatcherWithContext); ok {
		this.packetDispatcher = packetDispatcher
	}

	if dnsCache, ok := object.(internal.DnsCacheWithContext); ok {
		this.dnsCache = dnsCache
	}
}

func (this *SpaceController) ForContext(tag string) app.Space {
	return internal.NewSpace(tag, this.packetDispatcher, this.dnsCache)
}
