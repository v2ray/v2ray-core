package rules

import (
	"errors"

	"github.com/v2ray/v2ray-core/app/router"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	InvalidRule      = errors.New("Invalid Rule")
	NoRuleApplicable = errors.New("No rule applicable")
)

type Router struct {
	rules []Rule
}

func NewRouter() *Router {
	return &Router{
		rules: make([]Rule, 0, 16),
	}
}

func (this *Router) AddRule(rule Rule) *Router {
	this.rules = append(this.rules, rule)
	return this
}

func (this *Router) TakeDetour(dest v2net.Destination) (string, error) {
	for _, rule := range this.rules {
		if rule.Apply(dest) {
			return rule.Tag(), nil
		}
	}
	return "", NoRuleApplicable
}

type RouterFactory struct {
}

func (this *RouterFactory) Create(rawConfig interface{}) (router.Router, error) {
	config := rawConfig.(RouterRuleConfig)
	rules := config.Rules()
	router := NewRouter()
	for _, rule := range rules {
		if rule == nil {
			return nil, InvalidRule
		}
		router.AddRule(rule)
	}
	return router, nil
}

func init() {
	router.RegisterRouter("rules", &RouterFactory{})
}
