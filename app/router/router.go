package router

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

var (
	ErrNoRuleApplicable = errors.New("No rule applicable")
)

type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []Rule
	dnsServer      dns.Server
}

func NewRouter(ctx context.Context, config *Config) (*Router, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("Router: No space in context.")
	}
	r := &Router{
		domainStrategy: config.DomainStrategy,
		rules:          make([]Rule, len(config.Rule)),
	}

	space.OnInitialize(func() error {
		for idx, rule := range config.Rule {
			r.rules[idx].Tag = rule.Tag
			cond, err := rule.BuildCondition()
			if err != nil {
				return err
			}
			r.rules[idx].Condition = cond
		}

		r.dnsServer = dns.FromSpace(space)
		if r.dnsServer == nil {
			return errors.New("Router: DNS is not found in the space.")
		}
		return nil
	})
	return r, nil
}

func (v *Router) resolveIP(dest net.Destination) []net.Address {
	ips := v.dnsServer.Get(dest.Address.Domain())
	if len(ips) == 0 {
		return nil
	}
	dests := make([]net.Address, len(ips))
	for idx, ip := range ips {
		dests[idx] = net.IPAddress(ip)
	}
	return dests
}

func (v *Router) TakeDetour(ctx context.Context) (string, error) {
	for _, rule := range v.rules {
		if rule.Apply(ctx) {
			return rule.Tag, nil
		}
	}

	dest, ok := proxy.TargetFromContext(ctx)
	if !ok {
		return "", ErrNoRuleApplicable
	}

	if v.domainStrategy == Config_IpIfNonMatch && dest.Address.Family().IsDomain() {
		log.Trace(errors.New("looking up IP for ", dest).Path("App", "Router"))
		ipDests := v.resolveIP(dest)
		if ipDests != nil {
			ctx = proxy.ContextWithResolveIPs(ctx, ipDests)
			for _, rule := range v.rules {
				if rule.Apply(ctx) {
					return rule.Tag, nil
				}
			}
		}
	}

	return "", ErrNoRuleApplicable
}

func (Router) Interface() interface{} {
	return (*Router)(nil)
}

func (Router) Start() error {
	return nil
}

func (Router) Close() {}

func FromSpace(space app.Space) *Router {
	app := space.GetApplication((*Router)(nil))
	if app == nil {
		return nil
	}
	return app.(*Router)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewRouter(ctx, config.(*Config))
	}))
}
