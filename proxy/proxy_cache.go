package proxy

var (
	inboundFactories  = make(map[string]InboundConnectionHandlerFactory)
	outboundFactories = make(map[string]OutboundConnectionHandlerFactory)
)

func RegisterInboundConnectionHandlerFactory(name string, factory InboundConnectionHandlerFactory) error {
	// TODO check name
	inboundFactories[name] = factory
	return nil
}

func RegisterOutboundConnectionHandlerFactory(name string, factory OutboundConnectionHandlerFactory) error {
	// TODO check name
	outboundFactories[name] = factory
	return nil
}

func GetInboundConnectionHandlerFactory(name string) InboundConnectionHandlerFactory {
	factory, found := inboundFactories[name]
	if !found {
		return nil
	}
	return factory
}

func GetOutboundConnectionHandlerFactory(name string) OutboundConnectionHandlerFactory {
	factory, found := outboundFactories[name]
	if !found {
		return nil
	}
	return factory
}
