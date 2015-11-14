package config

type RouterConfig interface {
	Strategy() string
	Settings() interface{}
}
