package router

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/routing"
)

type key uint32

const (
	resolvedIPsKey key = iota
)

type IPResolver interface {
	Resolve() []net.Address
}

func ContextWithResolveIPs(ctx context.Context, f IPResolver) context.Context {
	return context.WithValue(ctx, resolvedIPsKey, f)
}

func ResolvedIPsFromContext(ctx context.Context) (IPResolver, bool) {
	ips, ok := ctx.Value(resolvedIPsKey).(IPResolver)
	return ips, ok
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		r := new(Router)
		if err := core.RequireFeatures(ctx, func(d dns.Client, ohm outbound.Manager) error {
			return r.Init(config.(*Config), d, ohm)
		}); err != nil {
			return nil, err
		}
		return r, nil
	}))
}

// Router is an implementation of routing.Router.
type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []*Rule
	balancers      map[string]*Balancer
	dns            dns.Client
}

// Init initializes the Router.
func (r *Router) Init(config *Config, d dns.Client, ohm outbound.Manager) error {
	r.domainStrategy = config.DomainStrategy
	r.dns = d

	r.balancers = make(map[string]*Balancer, len(config.BalancingRule))
	for _, rule := range config.BalancingRule {
		balancer, err := rule.Build(ohm)
		if err != nil {
			return err
		}
		r.balancers[rule.Tag] = balancer
	}

	r.rules = make([]*Rule, 0, len(config.Rule))
	for _, rule := range config.Rule {
		cond, err := rule.BuildCondition()
		if err != nil {
			return err
		}
		rr := &Rule{
			Condition: cond,
			Tag:       rule.GetTag(),
		}
		btag := rule.GetBalancingTag()
		if len(btag) > 0 {
			brule, found := r.balancers[btag]
			if !found {
				return newError("balancer ", btag, " not found")
			}
			rr.Balancer = brule
		}
		r.rules = append(r.rules, rr)
	}

	return nil
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

func (r *Router) PickRoute(ctx context.Context) (string, error) {
	rule, err := r.pickRouteInternal(ctx)
	if err != nil {
		return "", err
	}
	return rule.GetTag()
}

// PickRoute implements routing.Router.
func (r *Router) pickRouteInternal(ctx context.Context) (*Rule, error) {
	resolver := &ipResolver{
		dns: r.dns,
	}

	outbound := session.OutboundFromContext(ctx)
	if r.domainStrategy == Config_IpOnDemand {
		if outbound != nil && outbound.Target.IsValid() && outbound.Target.Address.Family().IsDomain() {
			resolver.domain = outbound.Target.Address.Domain()
			ctx = ContextWithResolveIPs(ctx, resolver)
		}
	}

	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule, nil
		}
	}

	if outbound == nil || !outbound.Target.IsValid() {
		return nil, common.ErrNoClue
	}

	dest := outbound.Target
	if r.domainStrategy == Config_IpIfNonMatch && dest.Address.Family().IsDomain() {
		resolver.domain = dest.Address.Domain()
		ips := resolver.Resolve()
		if len(ips) > 0 {
			ctx = ContextWithResolveIPs(ctx, resolver)
			for _, rule := range r.rules {
				if rule.Apply(ctx) {
					return rule, nil
				}
			}
		}
	}

	return nil, common.ErrNoClue
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
