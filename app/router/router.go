// +build !confonly

package router

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/routing"
)

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

// PickRoute implements routing.Router.
func (r *Router) PickRoute(ctx routing.Context) (string, error) {
	rule, err := r.pickRouteInternal(ctx)
	if err != nil {
		return "", err
	}
	return rule.GetTag()
}

func (r *Router) pickRouteInternal(ctx routing.Context) (*Rule, error) {
	if r.domainStrategy == Config_IpOnDemand {
		ctx = ContextWithDNSClient(ctx, r.dns)
	}

	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule, nil
		}
	}

	if r.domainStrategy != Config_IpIfNonMatch || len(ctx.GetTargetDomain()) == 0 {
		return nil, common.ErrNoClue
	}

	ctx = ContextWithDNSClient(ctx, r.dns)

	// Try applying rules again if we have IPs.
	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule, nil
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

// ContextWithDNSClient creates a new routing context with domain resolving capability. Resolved domain IPs can be retrieved by GetTargetIPs().
func ContextWithDNSClient(ctx routing.Context, client dns.Client) routing.Context {
	return &resolvableContext{Context: ctx, dnsClient: client}
}

type resolvableContext struct {
	routing.Context
	dnsClient   dns.Client
	resolvedIPs []net.IP
}

func (ctx *resolvableContext) GetTargetIPs() []net.IP {
	if ips := ctx.Context.GetTargetIPs(); len(ips) != 0 {
		return ips
	}

	if len(ctx.resolvedIPs) > 0 {
		return ctx.resolvedIPs
	}

	if domain := ctx.GetTargetDomain(); len(domain) != 0 {
		ips, err := ctx.dnsClient.LookupIP(domain)
		if err == nil {
			ctx.resolvedIPs = ips
			return ips
		}
		newError("resolve ip for ", domain).Base(err).WriteToLog()
	}

	return nil
}
