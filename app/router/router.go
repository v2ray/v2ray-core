// +build !confonly

package router

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/session"
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

func (r *Router) PickRoute(ctx context.Context) (string, error) {
	rule, err := r.pickRouteInternal(ctx)
	if err != nil {
		return "", err
	}
	return rule.GetTag()
}

func isDomainOutbound(outbound *session.Outbound) bool {
	return outbound != nil && outbound.Target.IsValid() && outbound.Target.Address.Family().IsDomain()
}

func (r *Router) resolveIP(outbound *session.Outbound) error {
	domain := outbound.Target.Address.Domain()
	ips, err := r.dns.LookupIP(domain)
	if err != nil {
		return err
	}

	outbound.ResolvedIPs = ips
	return nil
}

// PickRoute implements routing.Router.
func (r *Router) pickRouteInternal(ctx context.Context) (*Rule, error) {
	outbound := session.OutboundFromContext(ctx)
	if r.domainStrategy == Config_IpOnDemand && isDomainOutbound(outbound) {
		if err := r.resolveIP(outbound); err != nil {
			newError("failed to resolve IP for domain").Base(err).WriteToLog(session.ExportIDToError(ctx))
		}
	}

	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule, nil
		}
	}

	if r.domainStrategy != Config_IpIfNonMatch || !isDomainOutbound(outbound) {
		return nil, common.ErrNoClue
	}

	if err := r.resolveIP(outbound); err != nil {
		newError("failed to resolve IP for domain").Base(err).WriteToLog(session.ExportIDToError(ctx))
		return nil, common.ErrNoClue
	}

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
