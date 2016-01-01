package app

// Context of a function call from proxy to app.
type Context interface {
	CallerTag() string
}

// A Space contains all apps that may be available in a V2Ray runtime.
// Caller must check the availability of an app by calling HasXXX before getting its instance.
type Space interface {
	HasPacketDispatcher() bool
	PacketDispatcher() PacketDispatcher

	HasDnsCache() bool
	DnsCache() DnsCache

	HasPubsub() bool
	Pubsub() Pubsub

	HasInboundHandlerManager() bool
	InboundHandlerManager() InboundHandlerManager
}
