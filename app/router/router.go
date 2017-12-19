package router

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg router -path App,Router

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

var (
	ErrNoRuleApplicable = newError("No rule applicable")
)

type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []Rule
}

func NewRouter(ctx context.Context, config *Config) (*Router, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}
	r := &Router{
		domainStrategy: config.DomainStrategy,
		rules:          make([]Rule, len(config.Rule)),
	}

	space.On(app.SpaceInitializing, func(interface{}) error {
		for idx, rule := range config.Rule {
			r.rules[idx].Tag = rule.Tag
			cond, err := rule.BuildCondition()
			if err != nil {
				return err
			}
			r.rules[idx].Condition = cond
		}
		return nil
	})
	return r, nil
}

type ipResolver struct {
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
	ips, err := net.LookupIP(r.domain)
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

func (r *Router) TakeDetour(ctx context.Context) (string, error) {
	resolver := &ipResolver{}
	if r.domainStrategy == Config_IpOnDemand {
		if dest, ok := proxy.TargetFromContext(ctx); ok && dest.Address.Family().IsDomain() {
			resolver.domain = dest.Address.Domain()
			ctx = proxy.ContextWithResolveIPs(ctx, resolver)
		}
	}

	for _, rule := range r.rules {
		if rule.Apply(ctx) {
			return rule.Tag, nil
		}
	}

	dest, ok := proxy.TargetFromContext(ctx)
	if !ok {
		return "", ErrNoRuleApplicable
	}

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

	return "", ErrNoRuleApplicable
}

func (*Router) Interface() interface{} {
	return (*Router)(nil)
}

func (*Router) Start() error {
	return nil
}

func (*Router) Close() {}

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
