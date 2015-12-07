package router

type RouterConfig interface {
	Strategy() string
	Settings() interface{}
}
