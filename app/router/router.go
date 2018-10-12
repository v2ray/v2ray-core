package router

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/proxy"
)

// Router is an implementation of routing.Router.
type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []Rule
	dns            dns.Client
}

// NewRouter creates a new Router based on the given config.
func NewRouter(ctx context.Context, config *Config) (*Router, error) {
	v := core.MustFromContext(ctx)
	r := &Router{
		domainStrategy: config.DomainStrategy,
		rules:          make([]Rule, len(config.Rule)),
		dns:            v.DNSClient(),
	}

	for idx, rule := range config.Rule {
		r.rules[idx].Tag = rule.Tag
		cond, err := rule.BuildCondition()
		if err != nil {
			return nil, err
		}
		r.rules[idx].Condition = cond
	}

	if err := v.RegisterFeature(r); err != nil {
		return nil, newError("unable to register Router").Base(err)
	}
	return r, nil
}

type ipResolver struct {
	dns      dns.Client
	ip       []net.Address
	domain   string
	resolved bool
}

func (r *ipResolver) Resolve() []net.Address {
	if r.resolved {
		return r.ip
	}

	newError("looking for IP for domain: ", r.domain).WriteToLog()
	r.resolved = true
	ips, err := r.dns.LookupIP(r.domain)
	if err != nil {
		newError("failed to get IP address").Base(err).WriteToLog()
	}
	if len(ips) == 0 {
		return nil
	}
	r.ip = make([]net.Address, len(ips))
	for i, ip := range ips {
		r.ip[i] = net.IPAddress(ip)
	}
	return r.ip
}

// PickRoute implements routing.Router.
func (r *Router) PickRoute(ctx context.Context) (string, error) {
	resolver := &ipResolver{
		dns: r.dns,
	}

	outbound := session.OutboundFromContext(ctx)
	if r.domainStrategy == Config_IpOnDemand {
		if outbound != nil && outbound.Target.IsValid() && outbound.Target.Address.Family().IsDomain() {
			resolver.domain = outbound.Target.Address.Domain()
			ctx = proxy.ContextWithResolveIPs(ctx, resolver)
		}
	}

	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule.Tag, nil
		}
	}

	if outbound == nil || !outbound.Target.IsValid() {
		return "", common.ErrNoClue
	}

	dest := outbound.Target
	if r.domainStrategy == Config_IpIfNonMatch && dest.Address.Family().IsDomain() {
		resolver.domain = dest.Address.Domain()
		ips := resolver.Resolve()
		if len(ips) > 0 {
			ctx = proxy.ContextWithResolveIPs(ctx, resolver)
			for _, rule := range r.rules {
				if rule.Apply(ctx) {
					return rule.Tag, nil
				}
			}
		}
	}

	return "", common.ErrNoClue
}

// Start implements common.Runnable.
func (*Router) Start() error {
	return nil
}

// Close implements common.Closable.
func (*Router) Close() error {
	return nil
}

// Type implement common.HasType.
func (*Router) Type() interface{} {
	return routing.RouterType()
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewRouter(ctx, config.(*Config))
	}))
}
