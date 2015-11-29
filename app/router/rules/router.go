package rules

import (
	"errors"

	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/app/router/rules/config"
	"github.com/v2ray/v2ray-core/app/router/rules/config/json"
	v2net "github.com/v2ray/v2ray-core/common/net"
	pointconfig "github.com/v2ray/v2ray-core/shell/point/config"
)

var (
	InvalidRule      = errors.New("Invalid Rule")
	NoRuleApplicable = errors.New("No rule applicable")

	EmptyTag = pointconfig.DetourTag("")
)

type Router struct {
	rules []config.Rule
}

func (this *Router) TakeDetour(dest v2net.Destination) (pointconfig.DetourTag, error) {
	for _, rule := range this.rules {
		if rule.Apply(dest) {
			return rule.Tag(), nil
		}
	}
	return EmptyTag, NoRuleApplicable
}

type RouterFactory struct {
}

func (this *RouterFactory) Create(rawConfig interface{}) (router.Router, error) {
	config := rawConfig.(*json.RouterRuleConfig)
	rules := config.Rules()
	for _, rule := range rules {
		if rule == nil {
			return nil, InvalidRule
		}
	}
	return &Router{
		rules: rules,
	}, nil
}

func init() {
	router.RegisterRouter("rules", &RouterFactory{})
}
