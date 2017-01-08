package router

import (
	"v2ray.com/core/app"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

var (
	ErrInvalidRule      = errors.New("Invalid Rule")
	ErrNoRuleApplicable = errors.New("No rule applicable")
)

type Router struct {
	domainStrategy Config_DomainStrategy
	rules          []Rule
	//	cache          *RoutingTable
	dnsServer dns.Server
}

func NewRouter(config *Config, space app.Space) *Router {
	r := &Router{
		domainStrategy: config.DomainStrategy,
		//cache:          NewRoutingTable(),
		rules: make([]Rule, len(config.Rule)),
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
	return r
}

// Private: Visible for testing.
func (v *Router) ResolveIP(dest v2net.Destination) []v2net.Destination {
	ips := v.dnsServer.Get(dest.Address.Domain())
	if len(ips) == 0 {
		return nil
	}
	dests := make([]v2net.Destination, len(ips))
	for idx, ip := range ips {
		if dest.Network == v2net.Network_TCP {
			dests[idx] = v2net.TCPDestination(v2net.IPAddress(ip), dest.Port)
		} else {
			dests[idx] = v2net.UDPDestination(v2net.IPAddress(ip), dest.Port)
		}
	}
	return dests
}

func (v *Router) takeDetourWithoutCache(session *proxy.SessionInfo) (string, error) {
	for _, rule := range v.rules {
		if rule.Apply(session) {
			return rule.Tag, nil
		}
	}
	dest := session.Destination
	if v.domainStrategy == Config_IpIfNonMatch && dest.Address.Family().IsDomain() {
		log.Info("Router: Looking up IP for ", dest)
		ipDests := v.ResolveIP(dest)
		if ipDests != nil {
			for _, ipDest := range ipDests {
				log.Info("Router: Trying IP ", ipDest)
				for _, rule := range v.rules {
					if rule.Apply(&proxy.SessionInfo{
						Source:      session.Source,
						Destination: ipDest,
						User:        session.User,
					}) {
						return rule.Tag, nil
					}
				}
			}
		}
	}

	return "", ErrNoRuleApplicable
}

func (v *Router) TakeDetour(session *proxy.SessionInfo) (string, error) {
	//destStr := dest.String()
	//found, tag, err := v.cache.Get(destStr)
	//if !found {
	tag, err := v.takeDetourWithoutCache(session)
	//v.cache.Set(destStr, tag, err)
	return tag, err
	//}
	//return tag, err
}

type RouterFactory struct{}

func (RouterFactory) Create(space app.Space, config interface{}) (app.Application, error) {
	router := NewRouter(config.(*Config), space)
	return router, nil
}

func FromSpace(space app.Space) *Router {
	app := space.(app.AppGetter).GetApp(serial.GetMessageType((*Config)(nil)))
	if app == nil {
		return nil
	}
	return app.(*Router)
}

func init() {
	common.Must(app.RegisterApplicationFactory((*Config)(nil), RouterFactory{}))
}
