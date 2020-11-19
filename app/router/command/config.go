package command

import (
	"strings"

	"v2ray.com/core/common/net"
	"v2ray.com/core/features/routing"
)

// routingContext is an wrapper of protobuf RoutingContext as implementation of routing.Context and routing.Route.
type routingContext struct {
	*RoutingContext
}

func (c routingContext) GetSourceIPs() []net.IP {
	return mapBytesToIPs(c.RoutingContext.GetSourceIPs())
}

func (c routingContext) GetSourcePort() net.Port {
	return net.Port(c.RoutingContext.GetSourcePort())
}

func (c routingContext) GetTargetIPs() []net.IP {
	return mapBytesToIPs(c.RoutingContext.GetTargetIPs())
}

func (c routingContext) GetTargetPort() net.Port {
	return net.Port(c.RoutingContext.GetTargetPort())
}

// AsRoutingContext converts a protobuf RoutingContext into an implementation of routing.Context.
func AsRoutingContext(r *RoutingContext) routing.Context {
	return routingContext{r}
}

// AsRoutingRoute converts a protobuf RoutingContext into an implementation of routing.Route.
func AsRoutingRoute(r *RoutingContext) routing.Route {
	return routingContext{r}
}

var fieldMap = map[string]func(*RoutingContext, routing.Route){
	"inbound":        func(s *RoutingContext, r routing.Route) { s.InboundTag = r.GetInboundTag() },
	"network":        func(s *RoutingContext, r routing.Route) { s.Network = r.GetNetwork() },
	"ip_source":      func(s *RoutingContext, r routing.Route) { s.SourceIPs = mapIPsToBytes(r.GetSourceIPs()) },
	"ip_target":      func(s *RoutingContext, r routing.Route) { s.TargetIPs = mapIPsToBytes(r.GetTargetIPs()) },
	"port_source":    func(s *RoutingContext, r routing.Route) { s.SourcePort = uint32(r.GetSourcePort()) },
	"port_target":    func(s *RoutingContext, r routing.Route) { s.TargetPort = uint32(r.GetTargetPort()) },
	"domain":         func(s *RoutingContext, r routing.Route) { s.TargetDomain = r.GetTargetDomain() },
	"protocol":       func(s *RoutingContext, r routing.Route) { s.Protocol = r.GetProtocol() },
	"user":           func(s *RoutingContext, r routing.Route) { s.User = r.GetUser() },
	"attributes":     func(s *RoutingContext, r routing.Route) { s.Attributes = r.GetAttributes() },
	"outbound_group": func(s *RoutingContext, r routing.Route) { s.OutboundGroupTags = r.GetOutboundGroupTags() },
	"outbound":       func(s *RoutingContext, r routing.Route) { s.OutboundTag = r.GetOutboundTag() },
}

// AsProtobufMessage takes selectors of fields and returns a function to convert routing.Route to protobuf RoutingContext.
func AsProtobufMessage(fieldSelectors []string) func(routing.Route) *RoutingContext {
	initializers := []func(*RoutingContext, routing.Route){}
	for field, init := range fieldMap {
		if len(fieldSelectors) == 0 { // If selectors not set, retrieve all fields
			initializers = append(initializers, init)
			continue
		}
		for _, selector := range fieldSelectors {
			if strings.HasPrefix(field, selector) {
				initializers = append(initializers, init)
				break
			}
		}
	}
	return func(ctx routing.Route) *RoutingContext {
		message := new(RoutingContext)
		for _, init := range initializers {
			init(message, ctx)
		}
		return message
	}
}

func mapBytesToIPs(bytes [][]byte) []net.IP {
	var ips []net.IP
	for _, rawIP := range bytes {
		ips = append(ips, net.IP(rawIP))
	}
	return ips
}

func mapIPsToBytes(ips []net.IP) [][]byte {
	var bytes [][]byte
	for _, ip := range ips {
		bytes = append(bytes, []byte(ip))
	}
	return bytes
}
