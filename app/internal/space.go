package internal

import (
	"github.com/v2ray/v2ray-core/app"
)

type Space struct {
	packetDispatcher      PacketDispatcherWithContext
	dnsCache              DnsCacheWithContext
	pubsub                PubsubWithContext
	inboundHandlerManager InboundHandlerManagerWithContext
	tag                   string
}

func NewSpace(tag string, packetDispatcher PacketDispatcherWithContext, dnsCache DnsCacheWithContext, pubsub PubsubWithContext, inboundHandlerManager InboundHandlerManagerWithContext) *Space {
	return &Space{
		tag:                   tag,
		packetDispatcher:      packetDispatcher,
		dnsCache:              dnsCache,
		pubsub:                pubsub,
		inboundHandlerManager: inboundHandlerManager,
	}
}

func (this *Space) HasPacketDispatcher() bool {
	return this.packetDispatcher != nil
}

func (this *Space) PacketDispatcher() app.PacketDispatcher {
	return &contextedPacketDispatcher{
		packetDispatcher: this.packetDispatcher,
		context: &contextImpl{
			callerTag: this.tag,
		},
	}
}

func (this *Space) HasDnsCache() bool {
	return this.dnsCache != nil
}

func (this *Space) DnsCache() app.DnsCache {
	return &contextedDnsCache{
		dnsCache: this.dnsCache,
		context: &contextImpl{
			callerTag: this.tag,
		},
	}
}

func (this *Space) HasPubsub() bool {
	return this.pubsub != nil
}

func (this *Space) Pubsub() app.Pubsub {
	return &contextedPubsub{
		pubsub: this.pubsub,
		context: &contextImpl{
			callerTag: this.tag,
		},
	}
}

func (this *Space) HasInboundHandlerManager() bool {
	return this.inboundHandlerManager != nil
}

func (this *Space) InboundHandlerManager() app.InboundHandlerManager {
	return &inboundHandlerManagerWithContext{
		manager: this.inboundHandlerManager,
		context: &contextImpl{
			callerTag: this.tag,
		},
	}
}
