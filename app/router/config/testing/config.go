package testing

type RouterConfig struct {
	StrategyValue string
	SettingsValue interface{}
}

func (this *RouterConfig) Strategy() string {
	return this.StrategyValue
}

func (this *RouterConfig) Settings() interface{} {
	return this.SettingsValue
}
