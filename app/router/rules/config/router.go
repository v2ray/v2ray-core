package config

type RouterRuleConfig interface {
	Rules() []Rule
}
