// +build !confonly

package router

//go:generate go run v2ray.com/core/common/errors/errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/routing"
	routing_dns "v2ray.com/core/features/routing/dns"
)

// Router is an implementation of routing.Router.
type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []*Rule
	balancers      map[string]*Balancer
	dns            dns.Client
}

// Route is an implementation of routing.Route.
type Route struct {
	routing.Context
	outboundGroupTags []string
	outboundTag       string
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
func (r *Router) PickRoute(ctx routing.Context) (routing.Route, error) {
	rule, ctx, err := r.pickRouteInternal(ctx)
	if err != nil {
		return nil, err
	}
	tag, err := rule.GetTag()
	if err != nil {
		return nil, err
	}
	return &Route{Context: ctx, outboundTag: tag}, nil
}

func (r *Router) pickRouteInternal(ctx routing.Context) (*Rule, routing.Context, error) {
	if r.domainStrategy == Config_IpOnDemand {
		ctx = routing_dns.ContextWithDNSClient(ctx, r.dns)
	}

	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule, ctx, nil
		}
	}

	if r.domainStrategy != Config_IpIfNonMatch || len(ctx.GetTargetDomain()) == 0 {
		return nil, ctx, common.ErrNoClue
	}

	ctx = routing_dns.ContextWithDNSClient(ctx, r.dns)

	// Try applying rules again if we have IPs.
	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule, ctx, nil
		}
	}

	return nil, ctx, common.ErrNoClue
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

// GetOutboundGroupTags implements routing.Route.
func (r *Route) GetOutboundGroupTags() []string {
	return r.outboundGroupTags
}

// GetOutboundTag implements routing.Route.
func (r *Route) GetOutboundTag() string {
	return r.outboundTag
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
